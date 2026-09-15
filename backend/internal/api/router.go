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
	"github.com/WiseLabz/wiselabz/internal/scheduler"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
	"github.com/WiseLabz/wiselabz/internal/ws"

	// Register connector implementations
	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
)

// Config holds all dependencies needed to construct the router.
type Config struct {
	Store         *store.Store
	JWT           *auth.Service
	Config        *config.Config
	SyncEngine    *sync.Engine
	DocEngine     *doc.Engine
	WSHub         *ws.Hub
	AIRegistry    *ai.Registry
	EmbedRegistry *ai.EmbedRegistry
	Scheduler     *scheduler.Runner // for backup job scheduling
	BackupDir     string            // directory where backups are written
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
	connH := connhandler.NewHandler(cfg.Store, cfg.SyncEngine, cfg.Config)
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

	// --- System endpoints ---
	r.Get("/api/health", sysH.Health)

	// Shared across the top-level /api/auth route below and the protected
	// group further down.
	operatorOnly := auth.RequireRole("operator")
	dashboardDefaultOnly := auth.RequirePermission(cfg.Store, "can_manage_dashboard_defaults")

	// --- Auth routes (mixed public/protected, single mount point) ---
	// chi only allows one Mount per exact pattern, so the protected /api/auth
	// routes (logout, elevate, config) are nested inside this same Route as
	// inner auth-gated groups rather than a second top-level r.Route("/api/auth", ...).
	// Per-IP throttle on unauthenticated auth endpoints: 5 requests/sec with a
	// burst of 10, so a normal user retrying a typo never trips it but a
	// sustained guessing campaign against one IP does.
	authIPLimit := middleware.RateLimit(5, 10, func(r *http.Request) string { return httputil.ClientIP(r, cfg.Config.Server.TrustedProxies) })
	// Per-user throttle on password re-verification for step-up auth.
	elevateLimit := middleware.RateLimit(1, 5, func(r *http.Request) string { return auth.UserIDFromContext(r.Context()) })

	r.Route("/api/auth", func(r chi.Router) {
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
				r.Use(operatorOnly)
				r.Get("/config", settingH.GetAuthConfig)
				r.Put("/config", settingH.UpdateAuthConfig)
				r.Put("/providers/{providerId}/enabled", settingH.UpdateProviderEnabled)
			})
		})
	})

	// --- Protected routes (authenticated) ---
	r.Group(func(r chi.Router) {
		r.Use(cfg.AuthMiddleware())

		r.Route("/api/me", func(r chi.Router) {
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

		r.Route("/api/connectors", func(r chi.Router) {
			r.Get("/", connH.List)
			r.Get("/schema", connH.Schema)
			r.Get("/{id}", connH.Get)
			r.Get("/{id}/data", connH.Data)
			r.Get("/{id}/syncs", connH.Syncs)
			r.Get("/{id}/removal-impact", connH.RemovalImpact)

			r.Group(func(r chi.Router) {
				r.Use(operatorOnly)
				r.Post("/", connH.Create)
				r.Post("/{id}/test", connH.Test)
				r.Post("/{id}/health", connH.Health)
				r.Post("/{id}/restart", connH.RestartPreview)
				r.Patch("/{id}", connH.Update) // compatibility for clients predating the OpenAPI PUT contract
				r.Put("/{id}", connH.Update)
				r.Put("/{id}/enabled", connH.ToggleEnabled)
				r.Post("/{id}/sync", connH.Sync)

				r.Group(func(r chi.Router) {
					r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "connector.delete"))
					r.Delete("/{id}", connH.Delete)
				})
			})
		})

		r.Route("/api/docs", func(r chi.Router) {
			r.Get("/", docH.List)
			r.Get("/tree", docH.Tree)
			r.Get("/template-schema", docH.TemplateSchema)
			r.Get("/service/{id}", docH.ByService)
			r.Get("/{id}", docH.Get)
			r.Get("/{id}/versions", docH.Versions)
			r.Get("/{id}/versions/{rev}", docH.Version)
			r.Get("/{id}/lock", docH.GetLock)

			r.Group(func(r chi.Router) {
				r.Use(operatorOnly)
				r.Post("/generate", docH.Generate)
				r.Put("/{id}", docH.Save)
				r.Post("/{id}/versions/{rev}/restore", docH.Restore)
				r.Post("/{id}/ai-suggest", docH.AISuggest)
				r.Post("/{id}/lock", docH.AcquireLock)
				r.Post("/{id}/lock/release", docH.ReleaseLock)
			})
		})

		r.Route("/api/chat/conversations", func(r chi.Router) {
			r.Post("/", chatH.CreateConversation)
			r.Get("/", chatH.ListConversations)
			r.Get("/{id}", chatH.GetConversation)
			r.Post("/{id}/messages", chatH.PostMessage)
		})

		r.Route("/api/templates", func(r chi.Router) {
			r.Get("/", tmplH.List)
			r.Get("/{id}", tmplH.Get)
			r.Get("/{id}/versions", tmplH.Versions)
			r.Get("/{id}/versions/{rev}", tmplH.Version)
			r.Post("/{id}/preview", tmplH.Preview)

			r.Group(func(r chi.Router) {
				r.Use(operatorOnly)
				r.Post("/", tmplH.Create)
				r.Put("/{id}", tmplH.Update)
				r.Post("/{id}/versions/{rev}/restore", tmplH.Restore)

				r.Group(func(r chi.Router) {
					r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "template.delete"))
					r.Delete("/{id}", tmplH.Delete)
				})
			})
		})

		r.Route("/api/changes", func(r chi.Router) {
			r.Get("/", changeH.List)
			r.Get("/{id}", changeH.Get)

			r.Group(func(r chi.Router) {
				r.Use(operatorOnly)
				r.Post("/{id}/ack", changeH.Acknowledge)
				r.Post("/{id}/dismiss", changeH.Dismiss)
				r.Post("/{id}/ai-update", changeH.AIUpdate)
				r.Post("/{id}/explain", changeH.Explain)
				r.Post("/bulk-resolve", changeH.BulkResolve)
			})
		})

		r.Route("/api/alerts", func(r chi.Router) {
			r.Get("/", alertH.List)
			r.Get("/{id}", alertH.Get)

			r.Group(func(r chi.Router) {
				r.Use(operatorOnly)
				r.Post("/{id}/resolve", alertH.Resolve)
				r.Post("/{id}/dismiss", alertH.Dismiss)
				r.Post("/{id}/snooze", alertH.Snooze)
				r.Post("/bulk-snooze", alertH.BulkSnooze)
			})
		})

		r.Route("/api/attention", func(r chi.Router) {
			r.Get("/", attentionH.List)
		})

		r.Route("/api/runbooks", func(r chi.Router) {
			r.Get("/", runbookH.List)
			r.Get("/{id}", runbookH.Get)

			r.Group(func(r chi.Router) {
				r.Use(operatorOnly)
				r.Post("/", runbookH.Create)
				r.Put("/{id}", runbookH.Update)
				r.Delete("/{id}", runbookH.Delete)
			})
		})

		r.Route("/api/findings", func(r chi.Router) {
			r.Get("/", findingH.List)
			r.Get("/{id}", findingH.Get)

			r.Group(func(r chi.Router) {
				r.Use(operatorOnly)
				r.Post("/{id}/resolve", findingH.Resolve)
			})
		})

		r.Route("/api/notifications", func(r chi.Router) {
			r.Get("/", notifH.List)
			r.Post("/read-all", notifH.ReadAll)
			r.Post("/{id}/read", notifH.MarkRead)
		})

		// Saved views are a personal resource: any authenticated role (viewer
		// or operator) may list/create/delete their own, scoped by user_id in
		// the handler — no operatorOnly gate here.
		r.Route("/api/saved-views", func(r chi.Router) {
			r.Get("/", savedViewH.List)
			r.Post("/", savedViewH.Create)
			r.Delete("/{id}", savedViewH.Delete)
		})

		r.Route("/api/dashboard", func(r chi.Router) {
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

		// --- Operator-only routes ---
		r.Group(func(r chi.Router) {
			r.Use(operatorOnly)

			r.Route("/api/ai/config", func(r chi.Router) {
				r.Get("/", settingH.GetAIConfig)
				r.Put("/", settingH.UpdateAIConfig)
				r.Post("/test", settingH.TestAIConfig)
			})

			r.Route("/api/notifications/config", func(r chi.Router) {
				r.Get("/", settingH.GetNotificationsConfig)
				r.Put("/", settingH.UpdateNotificationsConfig)
				r.Post("/test", settingH.TestNotificationsConfig)
			})

			r.Get("/api/notifications/deliveries", notifH.ListDeliveries)

			r.Get("/api/system/info", sysH.Info)
			r.Get("/api/system/audit", sysH.ListAudit)
			r.Get("/api/system/audit/export", sysH.ExportAudit)

			r.Get("/api/system/settings/retention", sysH.GetRetentionSettings)
			r.Put("/api/system/settings/retention", sysH.UpdateRetentionSettings)

			r.Route("/api/system/backup", func(r chi.Router) {
				r.Get("/export", sysH.ExportBackup)
				r.Post("/import", sysH.ImportBackup)
				r.Get("/schedule", sysH.GetBackupSchedule)
				r.Put("/schedule", sysH.UpdateBackupSchedule)
				r.Get("/runs", sysH.ListBackupRuns)
				r.Post("/run", sysH.CreateBackupRun)
			})

			r.Get("/api/system/diagnostics", sysH.Diagnostics)

			r.Route("/api/users", func(r chi.Router) {
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

			r.Post("/api/sync", connH.SyncAll)
		})
	})

	// --- WebSocket endpoint ---
	r.Get("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Missing refresh token")
			return
		}
		claims, err := cfg.JWT.ValidateRefresh(cookie.Value)
		if err != nil {
			httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Invalid access token")
			return
		}
		active, err := cfg.Store.HasSessionTokenHash(r.Context(), claims.UserID, store.HashToken(cookie.Value))
		if err != nil || !active {
			httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Refresh token is no longer active")
			return
		}
		user, err := cfg.Store.GetUserByID(r.Context(), claims.UserID)
		if err != nil || user.Disabled {
			httputil.Error(w, http.StatusUnauthorized, "unauthorized", "User not found or disabled")
			return
		}
		if cfg.WSHub != nil {
			if err := cfg.WSHub.UpgradeHandler(w, r, claims.UserID, claims.Role); err != nil {
				slog.Error("WebSocket upgrade failed", "error", err)
			}
		}
	})

	if cfg.Config.Server.Embed && cfg.SPAFiles != nil {
		r.NotFound(spaHandler(cfg.SPAFiles))
	}

	return r
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
