// Package api provides the HTTP router, middleware chain, and handler wiring.
package api

import (
	"context"
	"io/fs"
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

	d := newRouterDeps(cfg)
	mountRootRoutes(r, d)

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
	mountAPI := func(r chi.Router) { mountAPIRoutes(r, d) }

	r.Route("/api", mountAPI)
	r.Route("/api/v1", mountAPI)

	if cfg.Config.Server.Embed && cfg.SPAFiles != nil {
		r.NotFound(spaHandler(cfg.SPAFiles))
	}

	return r
}

// newRouterDeps constructs every domain handler once, so the /api and /api/v1
// mounts share the same instances.
func newRouterDeps(cfg Config) routerDeps {
	sysH := syshandler.NewHandler(cfg.Store.DB(), cfg.Config, cfg.Store, cfg.Scheduler, cfg.BackupDir)
	// Register the backup cron job through the handler (not directly against
	// cfg.Scheduler) so its entry ID is tracked and later PUT /schedule calls
	// can remove/replace it instead of stacking duplicate jobs.
	sysH.InitBackupJob(context.Background())
	sysH.InitRetentionJob(context.Background())

	settingH := settinghandler.NewHandler(cfg.Store, cfg.Config, cfg.AIRegistry)

	var ruleEvaluator compliancehandler.RuleEvaluator
	if cfg.QualityChecker != nil {
		ruleEvaluator = cfg.QualityChecker
	}

	return routerDeps{
		cfg:         cfg,
		sysH:        sysH,
		authH:       authhandler.NewHandler(cfg.Store, cfg.JWT, cfg.Config),
		apiKeyH:     apikeyhandler.NewHandler(cfg.Store),
		settingH:    settingH,
		connH:       connhandler.NewHandler(cfg.Store, cfg.SyncEngine, cfg.Config, cfg.JWT, cfg.WSHub),
		tmplH:       tmplhandler.NewHandler(cfg.Store, cfg.DocEngine),
		changeH:     changehandler.NewHandler(cfg.Store, settingH, cfg.AIRegistry, cfg.WSHub),
		alertH:      alerthandler.NewHandler(cfg.Store),
		attentionH:  attentionhandler.NewHandler(cfg.Store),
		findingH:    findinghandler.NewHandler(cfg.Store),
		notifH:      notifhandler.NewHandler(cfg.Store),
		runbookH:    runbookhandler.NewHandler(cfg.Store),
		dashH:       dashhandler.NewHandler(cfg.Store),
		docH:        dochandler.NewHandler(cfg.Store, cfg.DocEngine, settingH, cfg.AIRegistry, cfg.EmbedRegistry, cfg.WSHub),
		savedViewH:  savedviewhandler.NewHandler(cfg.Store),
		chatH:       chathandler.NewHandler(cfg.Store, settingH, cfg.AIRegistry, cfg.EmbedRegistry),
		complianceH: compliancehandler.NewHandler(cfg.Store, ruleEvaluator),
	}
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
