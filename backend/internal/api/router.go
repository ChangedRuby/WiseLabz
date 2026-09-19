// Package api provides the HTTP router, middleware chain, and handler wiring.
package api

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/WiseLabz/wiselabz/internal/ai"
	alerthandler "github.com/WiseLabz/wiselabz/internal/api/alerts"
	apikeyhandler "github.com/WiseLabz/wiselabz/internal/api/apikeys"
	attentionhandler "github.com/WiseLabz/wiselabz/internal/api/attention"
	authhandler "github.com/WiseLabz/wiselabz/internal/api/auth"
	changehandler "github.com/WiseLabz/wiselabz/internal/api/changes"
	chathandler "github.com/WiseLabz/wiselabz/internal/api/chat"
	compliancehandler "github.com/WiseLabz/wiselabz/internal/api/compliance"
	connhandler "github.com/WiseLabz/wiselabz/internal/api/connectors"
	dashhandler "github.com/WiseLabz/wiselabz/internal/api/dashboard"
	dochandler "github.com/WiseLabz/wiselabz/internal/api/docs"
	findinghandler "github.com/WiseLabz/wiselabz/internal/api/findings"
	"github.com/WiseLabz/wiselabz/internal/api/middleware"
	notifhandler "github.com/WiseLabz/wiselabz/internal/api/notifications"
	runbookhandler "github.com/WiseLabz/wiselabz/internal/api/runbooks"
	savedviewhandler "github.com/WiseLabz/wiselabz/internal/api/savedviews"
	settinghandler "github.com/WiseLabz/wiselabz/internal/api/settings"
	syshandler "github.com/WiseLabz/wiselabz/internal/api/system"
	tmplhandler "github.com/WiseLabz/wiselabz/internal/api/templates"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/doc"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/quality"
	"github.com/WiseLabz/wiselabz/internal/scheduler"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
	"github.com/WiseLabz/wiselabz/internal/ws"

	// Register connector implementations
	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
)

// Config holds all dependencies needed to construct the router.
type Config struct {
	Store          *store.Store
	JWT            *auth.Service
	Config         *config.Config
	SyncEngine     *sync.Engine
	DocEngine      *doc.Engine
	WSHub          *ws.Hub
	AIRegistry     *ai.Registry
	EmbedRegistry  *ai.EmbedRegistry
	Scheduler      *scheduler.Runner // for backup job scheduling
	QualityChecker *quality.Checker
	BackupDir      string // directory where backups are written
	// SPAFiles serves the embedded frontend build. Only used when Config.Server.Embed is true.
	SPAFiles fs.FS
}

// NewRouter constructs the full chi router with all middleware and route groups.
func NewRouter(cfg Config) chi.Router {
	r := chi.NewRouter()

	// Global middleware chain
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.SecurityHeaders(cfg.Config.Server.TrustedProxies))
	r.Use(middleware.CORS(cfg.Config.Server.Origin))

	// --- Handlers ---
	sysH := syshandler.NewHandler(cfg.Store.DB(), cfg.Config, cfg.Store, cfg.Scheduler, cfg.BackupDir)
	// Register the backup cron job through the handler (not directly against
	// cfg.Scheduler) so its entry ID is tracked and later PUT /schedule calls
	// can remove/replace it instead of stacking duplicate jobs.
	sysH.InitBackupJob(context.Background())
	sysH.InitRetentionJob(context.Background())
	authH := authhandler.NewHandler(cfg.Store, cfg.JWT, cfg.Config)
	apiKeyH := apikeyhandler.NewHandler(cfg.Store)
	settingH := settinghandler.NewHandler(cfg.Store, cfg.Config, cfg.AIRegistry)
	connH := connhandler.NewHandler(cfg.Store, cfg.SyncEngine, cfg.Config, cfg.JWT, cfg.WSHub)
	tmplH := tmplhandler.NewHandler(cfg.Store, cfg.DocEngine)
	changeH := changehandler.NewHandler(cfg.Store, settingH, cfg.AIRegistry, cfg.WSHub)
	alertH := alerthandler.NewHandler(cfg.Store)
	attentionH := attentionhandler.NewHandler(cfg.Store)
	findingH := findinghandler.NewHandler(cfg.Store)
	notifH := notifhandler.NewHandler(cfg.Store)
	runbookH := runbookhandler.NewHandler(cfg.Store)
	dashH := dashhandler.NewHandler(cfg.Store)
	docH := dochandler.NewHandler(cfg.Store, cfg.DocEngine, settingH, cfg.AIRegistry, cfg.EmbedRegistry, cfg.WSHub)
	savedViewH := savedviewhandler.NewHandler(cfg.Store)
	chatH := chathandler.NewHandler(cfg.Store, settingH, cfg.AIRegistry, cfg.EmbedRegistry)
	var ruleEvaluator compliancehandler.RuleEvaluator
	if cfg.QualityChecker != nil {
		ruleEvaluator = cfg.QualityChecker
	}
	complianceH := compliancehandler.NewHandler(cfg.Store, ruleEvaluator)

	if cfg.WSHub != nil {
		cfg.WSHub.SetRevalidator(func(ctx context.Context, userID, role, sessionHash string) bool {
			userRole, disabled, err := cfg.Store.GetUserRoleStatus(ctx, userID)
			if err != nil || disabled || wsRoleLabel(userRole == "admin") != role {
				return false
			}
			if sessionHash == "" {
				return true
			}
			active, err := cfg.Store.HasSessionTokenHash(ctx, userID, sessionHash)
			return err == nil && active
		})
	}

	// The whole API is mounted twice: under the legacy /api prefix and the
	// versioned /api/v1 prefix (an alias, same handlers and middleware).
	mountAPI := func(r chi.Router) {
		// --- System endpoints ---
		r.Get("/health", sysH.Health)

		// Shared across the top-level /api/auth route below and the protected
		// group further down.
		dashboardDefaultOnly := auth.RequirePermission(cfg.Store, "can_manage_dashboard_defaults")
		connOperator := auth.RequireConnectorRole(cfg.Store, "operator", "id")

		// --- Auth routes (mixed public/protected, single mount point) ---
		// chi only allows one Mount per exact pattern, so the protected /api/auth
		// routes (logout, elevate, config) are nested inside this same Route as
		// inner auth-gated groups rather than a second top-level r.Route("/auth", ...).
		// Per-IP throttle on unauthenticated auth endpoints: 5 requests/sec with a
		// burst of 10, so a normal user retrying a typo never trips it but a
		// sustained guessing campaign against one IP does.
		authIPLimit := middleware.RateLimit(5, 10, func(r *http.Request) string { return httputil.ClientIP(r, cfg.Config.Server.TrustedProxies) })
		// Per-user throttle on password re-verification for step-up auth.
		elevateLimit := middleware.RateLimit(1, 5, func(r *http.Request) string { return auth.UserIDFromContext(r.Context()) })

		r.Route("/auth", func(r chi.Router) {
			r.With(authIPLimit).Post("/login", authH.Login)
			r.With(authIPLimit).Post("/oidc/callback", authH.OIDCCallback)
			r.With(authIPLimit).Post("/refresh", authH.Refresh)
			r.Get("/providers", authH.Providers)

			r.Group(func(r chi.Router) {
				r.Use(cfg.AuthMiddleware())
				r.Post("/logout", authH.Logout)
				r.With(elevateLimit).Post("/elevate", authH.Elevate)
				r.Route("/api-keys", func(r chi.Router) {
					r.Get("/", apiKeyH.List)
					r.Post("/", apiKeyH.Create)
					r.Delete("/{id}", apiKeyH.Revoke)
				})

				r.Group(func(r chi.Router) {
					r.Use(auth.RequireInstanceAdmin)
					r.Get("/config", settingH.GetAuthConfig)
					r.Put("/config", settingH.UpdateAuthConfig)
					r.Put("/providers/{providerId}/enabled", settingH.UpdateProviderEnabled)
				})
			})
		})

		// --- Read-only doc share links (unauthenticated, token in the path) ---
		// Deliberately outside cfg.AuthMiddleware(): the whole point is access
		// without an account. docH.ResolveShareLink is the token/expiry/revoke
		// gate; every route in this group is read-only (no save/lock/ai-suggest).
		r.Route("/share/{token}", func(r chi.Router) {
			r.Use(docH.ResolveShareLink)
			r.Get("/tree", docH.ShareLinkTree)
			r.Get("/docs/{docId}", docH.ShareLinkDoc)
		})

		// --- Protected routes (authenticated) ---
		r.Group(func(r chi.Router) {
			r.Use(cfg.AuthMiddleware())

			r.Route("/me", func(r chi.Router) {
				r.Get("/", authH.Me)
				r.Patch("/", authH.UpdateMe)
				r.Post("/password", authH.ChangePassword)
				r.Get("/sessions", authH.ListSessions)
				r.Delete("/sessions/{id}", authH.DeleteSession)
			})

			// Each resource mounts at a single path: GET endpoints are open to any
			// authenticated user, mutating endpoints are nested in an inner
			// operator-only group. chi panics if the same pattern is r.Route()'d
			// twice on one mux, so read/write must share one Route block.

			connViewer := auth.RequireConnectorRole(cfg.Store, "viewer", "id")

			r.Route("/connectors", func(r chi.Router) {
				// List/maintenance-windows are cross-connector and filter inside
				// the handler; Get 404s inside the handler on a missing grant
				// (not 403, to avoid confirming existence) rather than gated
				// here. Every other GET takes a connector {id} and can use the
				// viewer-role middleware directly.
				r.Get("/", connH.List)
				r.Get("/schema", connH.Schema)
				r.Get("/{id}", connH.Get)
				r.Get("/maintenance-windows", connH.ListActiveMaintenance)

				r.Group(func(r chi.Router) {
					r.Use(connViewer)
					r.Get("/{id}/data", connH.Data)
					r.Get("/{id}/syncs", connH.Syncs)
					r.Get("/{id}/removal-impact", connH.RemovalImpact)
					r.Get("/{id}/config-fields", connH.ConfigFields)
					r.Get("/{id}/maintenance-window", connH.GetMaintenanceWindow)
				})

				// Creating a connector has no existing grant to check against, so
				// it's instance-admin only; the creator is auto-granted operator
				// on the new connector (see connH.Create).
				r.With(auth.RequireInstanceAdmin).Post("/", connH.Create)

				// Instance-admin-only grant management for this connector — a
				// deliberate exception to the per-connector-role pattern below,
				// same reasoning as saved-views/dashboard-layout above: granting
				// access is an instance-wide action, not something a connector's
				// own operators can do to each other.
				r.Route("/{id}/permissions", func(r chi.Router) {
					r.Use(auth.RequireInstanceAdmin)
					r.Get("/", connH.ListPermissions)
					r.Put("/{userId}", connH.PutPermission)
					r.Delete("/{userId}", connH.DeletePermission)
				})

				r.Group(func(r chi.Router) {
					r.Use(connOperator)
					r.Post("/{id}/test", connH.Test)
					r.Post("/{id}/health", connH.Health)
					r.Post("/{id}/restart", connH.RestartPreview)
					r.Post("/{id}/start", connH.StartPreview)
					r.Post("/{id}/stop", connH.StopPreview)
					r.Post("/{id}/config-push", connH.ConfigPush)
					r.Patch("/{id}", connH.Update) // compatibility for clients predating the OpenAPI PUT contract
					r.Put("/{id}", connH.Update)
					r.Put("/{id}/enabled", connH.ToggleEnabled)
					r.Post("/{id}/sync", connH.Sync)
					r.Post("/{id}/maintenance-window", connH.OpenMaintenanceWindow) // no elevation: reversible and time-boxed
					r.Delete("/{id}/maintenance-window", connH.CloseMaintenanceWindow)

					r.Group(func(r chi.Router) {
						r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "connector.delete"))
						r.Delete("/{id}", connH.Delete)
					})
				})

				// Bulk routes take a body-supplied ID list, not a path {id}, so
				// they can't use connOperator middleware — each handler checks
				// per-ID via store.UserHasConnectorRole itself.
				r.Group(func(r chi.Router) {
					r.Post("/bulk-sync", connH.BulkSync)
					r.Post("/bulk-reauth", connH.BulkReauth)

					r.Group(func(r chi.Router) {
						// One elevation token covers the whole bulk-restart batch,
						// via its own distinct action string (see BulkRestart).
						r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "connector.bulkRestart"))
						r.Post("/bulk-restart", connH.BulkRestart)
					})
				})
			})

			r.Route("/docs", func(r chi.Router) {
				// GET routes are default-deny inside the handlers (List/Tree
				// filter, Get/ByService 404 on a missing grant); lab-wide docs
				// (no connector) stay visible to any authenticated user.
				r.Get("/", docH.List)
				r.Get("/tree", docH.Tree)
				r.Get("/template-schema", docH.TemplateSchema)
				r.Get("/service/{id}", docH.ByService)
				r.Get("/{id}", docH.Get)
				r.Get("/{id}/versions", docH.Versions)
				r.Get("/{id}/versions/{rev}", docH.Version)
				r.Get("/{id}/lock", docH.GetLock)

				// Mutations resolve their doc/connector ID in-handler (a doc ID in
				// the path, or a connectorId in the body for Generate) and check
				// store.UserHasConnectorRole themselves (docH.requireDocOperator) —
				// they can't use connOperator middleware, which only reads a
				// connector ID directly from the path.
				r.Post("/generate", docH.Generate)
				r.Put("/{id}", docH.Save)
				r.Post("/{id}/versions/{rev}/restore", docH.Restore)
				r.Post("/{id}/ai-suggest", docH.AISuggest)
				r.Post("/{id}/lock", docH.AcquireLock)
				r.Post("/{id}/lock/release", docH.ReleaseLock)

				// Regenerates the single lab-wide Lab Topology doc — instance-admin,
				// same as any other lab-wide (no connector) doc mutation.
				r.With(auth.RequireInstanceAdmin).Post("/topology", docH.GenerateTopology)

				// Share links: creation checks operator on every connector covered
				// by the requested subtree in-handler (docH.requireShareCreateAccess),
				// same resolve-then-check pattern as the rest of this group.
				r.Route("/share-links", func(r chi.Router) {
					r.Get("/", docH.ListShareLinks)
					r.Post("/", docH.CreateShareLink)
					r.Delete("/{id}", docH.RevokeShareLink)
				})
			})

			r.Route("/chat/conversations", func(r chi.Router) {
				r.Post("/", chatH.CreateConversation)
				r.Get("/", chatH.ListConversations)
				r.Get("/{id}", chatH.GetConversation)
				r.Post("/{id}/messages", chatH.PostMessage)
			})

			r.Route("/templates", func(r chi.Router) {
				r.Get("/", tmplH.List)
				r.Get("/{id}", tmplH.Get)
				r.Get("/{id}/versions", tmplH.Versions)
				r.Get("/{id}/versions/{rev}", tmplH.Version)
				r.Post("/{id}/preview", tmplH.Preview)

				r.Group(func(r chi.Router) {
					r.Use(auth.RequireInstanceAdmin)
					r.Post("/", tmplH.Create)
					r.Put("/{id}", tmplH.Update)
					r.Post("/{id}/versions/{rev}/restore", tmplH.Restore)

					r.Group(func(r chi.Router) {
						r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "template.delete"))
						r.Delete("/{id}", tmplH.Delete)
					})
				})
			})

			// changes/alerts/findings ARE connector-scoped (each row carries a
			// NOT NULL connector FK), unlike templates/runbooks above/below. Their
			// mutations resolve the record's connector in-handler and check
			// store.UserHasConnectorRole themselves, same reasoning as docs.
			r.Route("/changes", func(r chi.Router) {
				r.Get("/", changeH.List)
				r.Get("/{id}", changeH.Get)
				r.Post("/{id}/ack", changeH.Acknowledge)
				r.Post("/{id}/dismiss", changeH.Dismiss)
				r.Post("/{id}/ai-update", changeH.AIUpdate)
				r.Post("/{id}/explain", changeH.Explain)
				r.Post("/bulk-resolve", changeH.BulkResolve)
			})

			r.Route("/alerts", func(r chi.Router) {
				r.Get("/", alertH.List)
				r.Get("/{id}", alertH.Get)
				r.Post("/{id}/resolve", alertH.Resolve)
				r.Post("/{id}/dismiss", alertH.Dismiss)
				r.Post("/{id}/snooze", alertH.Snooze)
				r.Post("/bulk-snooze", alertH.BulkSnooze)
			})

			r.Route("/attention", func(r chi.Router) {
				r.Get("/", attentionH.List)
			})

			r.Route("/runbooks", func(r chi.Router) {
				r.Get("/", runbookH.List)
				r.Get("/{id}", runbookH.Get)

				r.Group(func(r chi.Router) {
					r.Use(auth.RequireInstanceAdmin)
					r.Post("/", runbookH.Create)
					r.Put("/{id}", runbookH.Update)
					r.Delete("/{id}", runbookH.Delete)
				})
			})

			r.Route("/findings", func(r chi.Router) {
				r.Get("/", findingH.List)
				r.Get("/{id}", findingH.Get)
				r.Post("/{id}/resolve", findingH.Resolve)
			})

			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", notifH.List)
				r.Post("/read-all", notifH.ReadAll)
				r.Post("/{id}/read", notifH.MarkRead)
			})

			// Saved views are a personal resource: any authenticated role (viewer
			// or operator) may list/create/delete their own, scoped by user_id in
			// the handler — no operatorOnly gate here.
			r.Route("/saved-views", func(r chi.Router) {
				r.Get("/", savedViewH.List)
				r.Post("/", savedViewH.Create)
				r.Delete("/{id}", savedViewH.Delete)
			})

			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/overview", dashH.Overview)
				// Dashboard layout is a personal resource: any authenticated role
				// (viewer or operator) may read/save/reset their own, scoped by
				// user_id in the handler — no operatorOnly gate here.
				r.Get("/layout", dashH.GetLayout)
				r.Put("/layout", dashH.SaveLayout)
				r.Post("/layout/reset", dashH.ResetLayout)

				r.Group(func(r chi.Router) {
					r.Use(dashboardDefaultOnly)
					r.Get("/layout/admin-default", dashH.GetAdminDefault)
					r.Put("/layout/admin-default", dashH.PutAdminDefault)
				})
			})

			// --- Instance-admin-only routes ---
			r.Group(func(r chi.Router) {
				r.Use(auth.RequireInstanceAdmin)

				r.Route("/ai/config", func(r chi.Router) {
					r.Get("/", settingH.GetAIConfig)
					r.Put("/", settingH.UpdateAIConfig)
					r.Post("/test", settingH.TestAIConfig)
					r.Get("/fallback-providers", settingH.GetAIFallbackProviders)
					r.Put("/fallback-providers", settingH.UpdateAIFallbackProviders)
				})

				r.Route("/notifications/config", func(r chi.Router) {
					r.Get("/", settingH.GetNotificationsConfig)
					r.Put("/", settingH.UpdateNotificationsConfig)
					r.Post("/test", settingH.TestNotificationsConfig)
				})

				r.Get("/notifications/deliveries", notifH.ListDeliveries)

				r.Get("/system/info", sysH.Info)
				r.Get("/system/audit", sysH.ListAudit)
				r.Get("/system/audit/export", sysH.ExportAudit)

				r.Get("/system/settings/retention", sysH.GetRetentionSettings)
				r.Put("/system/settings/retention", sysH.UpdateRetentionSettings)

				r.Route("/system/backup", func(r chi.Router) {
					r.Get("/export", sysH.ExportBackup)
					r.Post("/import", sysH.ImportBackup)
					r.Get("/schedule", sysH.GetBackupSchedule)
					r.Put("/schedule", sysH.UpdateBackupSchedule)
					r.Get("/runs", sysH.ListBackupRuns)
					r.Post("/run", sysH.CreateBackupRun)
				})

				r.Get("/system/diagnostics", sysH.Diagnostics)

				r.Route("/compliance", func(r chi.Router) {
					r.Get("/schema", complianceH.Schema)
					r.Get("/rules", complianceH.List)
					r.Post("/rules", complianceH.Create)
					r.Get("/rules/{id}", complianceH.Get)
					r.Put("/rules/{id}", complianceH.Update)
					r.Delete("/rules/{id}", complianceH.Delete)
					r.Post("/rules/test", complianceH.Test)
				})

				r.Route("/users", func(r chi.Router) {
					r.Get("/", authH.ListUsers)
					r.Post("/", authH.CreateUser)
					r.Patch("/{id}", authH.UpdateUser)

					r.Group(func(r chi.Router) {
						r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "user.delete"))
						r.Delete("/{id}", authH.DeleteUser)
					})

					r.Group(func(r chi.Router) {
						r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "user.resetPassword"))
						r.Post("/{id}/reset-password", authH.ResetPassword)
					})
				})

				r.Post("/sync", connH.SyncAll)
			})
		})

		// --- WebSocket endpoint ---
		// The upgrade is authorized by a one-time ticket minted from an
		// authenticated request (POST /api/ws/ticket), not by the long-lived
		// refresh cookie. Open connections are re-validated on every ping.
		if cfg.WSHub != nil {
			r.With(cfg.AuthMiddleware()).Post("/ws/ticket", func(w http.ResponseWriter, r *http.Request) {
				userID := auth.UserIDFromContext(r.Context())
				// Bind the ticket to the caller's refresh session when the cookie is
				// present so logout / password change closes the socket.
				var sessionHash string
				if cookie, err := r.Cookie("refresh_token"); err == nil {
					sessionHash = store.HashToken(cookie.Value)
					if active, err := cfg.Store.HasSessionTokenHash(r.Context(), userID, sessionHash); err != nil || !active {
						sessionHash = ""
					}
				}
				id, err := cfg.WSHub.IssueTicket(userID, wsRoleLabel(auth.InstanceAdminFromContext(r.Context())), sessionHash)
				if err != nil {
					httputil.Errorf(w, err)
					return
				}
				httputil.JSON(w, http.StatusOK, map[string]any{"ticket": id})
			})
			r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
				userID, role, sessionHash, ok := cfg.WSHub.RedeemTicket(r.URL.Query().Get("ticket"))
				if !ok {
					httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Invalid or expired ticket")
					return
				}
				if !cfg.WSHub.Revalidate(r.Context(), userID, role, sessionHash) {
					httputil.Error(w, http.StatusUnauthorized, "unauthorized", "User not found or disabled")
					return
				}
				if err := cfg.WSHub.UpgradeHandler(w, r, userID, role, sessionHash); err != nil {
					slog.Error("WebSocket upgrade failed", "error", err)
				}
			})
		}
	}

	r.Route("/api", mountAPI)
	r.Route("/api/v1", mountAPI)

	if cfg.Config.Server.Embed && cfg.SPAFiles != nil {
		r.NotFound(spaHandler(cfg.SPAFiles))
	}

	return r
}

// wsRoleLabel is a cosmetic presence/broadcast tag, not a security boundary
// (per-connector access is never carried over the WS connection).
func wsRoleLabel(instanceAdmin bool) string {
	if instanceAdmin {
		return "admin"
	}
	return "user"
}

// spaHandler serves the embedded SPA build, falling back to index.html for
// client-side routes. It only handles paths chi's router didn't already match,
// so /api/* routes are never shadowed — unmatched API paths get a JSON 404.
func spaHandler(files fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(files))

	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			httputil.Error(w, http.StatusNotFound, "not_found", "Resource not found")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if info, err := fs.Stat(files, path); err != nil || info.IsDir() {
			r = r.Clone(r.Context())
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	}
}

// AuthMiddleware returns chi-compatible auth middleware from the JWT service.
func (cfg Config) AuthMiddleware() func(http.Handler) http.Handler {
	return auth.AuthMiddleware(cfg.JWT, cfg.Store)
}
