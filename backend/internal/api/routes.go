package api

import (
	"github.com/go-chi/chi/v5"

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
	notifhandler "github.com/WiseLabz/wiselabz/internal/api/notifications"
	runbookhandler "github.com/WiseLabz/wiselabz/internal/api/runbooks"
	savedviewhandler "github.com/WiseLabz/wiselabz/internal/api/savedviews"
	settinghandler "github.com/WiseLabz/wiselabz/internal/api/settings"
	syshandler "github.com/WiseLabz/wiselabz/internal/api/system"
	tmplhandler "github.com/WiseLabz/wiselabz/internal/api/templates"
	"github.com/WiseLabz/wiselabz/internal/auth"
)

// routerDeps carries the Config plus the handler instances built once in
// NewRouter, so every mount*Routes helper below takes a single value instead of
// a long parameter list. It exists only to let the route tree be split across
// files; it is not a dependency-injection layer.
type routerDeps struct {
	cfg         Config
	sysH        *syshandler.Handler
	authH       *authhandler.Handler
	apiKeyH     *apikeyhandler.Handler
	settingH    *settinghandler.Handler
	connH       *connhandler.Handler
	tmplH       *tmplhandler.Handler
	changeH     *changehandler.Handler
	alertH      *alerthandler.Handler
	attentionH  *attentionhandler.Handler
	findingH    *findinghandler.Handler
	notifH      *notifhandler.Handler
	runbookH    *runbookhandler.Handler
	dashH       *dashhandler.Handler
	docH        *dochandler.Handler
	savedViewH  *savedviewhandler.Handler
	chatH       *chathandler.Handler
	complianceH *compliancehandler.Handler
}

// mountAPIRoutes registers the whole API surface on r. It is mounted twice by
// NewRouter — under the legacy /api prefix and the versioned /api/v1 prefix (an
// alias, same handlers and middleware).
func mountAPIRoutes(r chi.Router, d routerDeps) {
	mountPublicSystemRoutes(r, d)

	mountAuthRoutes(r, d)
	mountShareRoutes(r, d)

	// --- Protected routes (authenticated) ---
	r.Group(func(r chi.Router) {
		r.Use(d.cfg.AuthMiddleware())

		mountMeRoutes(r, d)

		// Each resource mounts at a single path: GET endpoints are open to any
		// authenticated user, mutating endpoints are nested in an inner
		// operator-only group. chi panics if the same pattern is r.Route()'d
		// twice on one mux, so read/write must share one Route block.

		mountConnectorRoutes(r, d)
		mountDocRoutes(r, d)
		mountChatRoutes(r, d)
		mountTemplateRoutes(r, d)
		mountWorkflowRoutes(r, d)
		mountDashboardRoutes(r, d)

		// --- Instance-admin-only routes ---
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireInstanceAdmin)

			mountSettingsRoutes(r, d)
			mountSystemRoutes(r, d)
			mountComplianceRoutes(r, d)
			mountUserRoutes(r, d)

			r.Post("/sync", d.connH.SyncAll)
		})
	})

	mountWSRoutes(r, d)
}
