package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/WiseLabz/wiselabz/internal/auth"
)

// mountWorkflowRoutes registers the day-to-day operational resources —
// changes, alerts, attention, runbooks, findings, notifications and saved
// views. It must be called on an already-authenticated group.
func mountWorkflowRoutes(r chi.Router, d routerDeps) {
	// changes/alerts/findings ARE connector-scoped (each row carries a
	// NOT NULL connector FK), unlike templates/runbooks. Their
	// mutations resolve the record's connector in-handler and check
	// store.UserHasConnectorRole themselves, same reasoning as docs.
	r.Route("/changes", func(r chi.Router) {
		r.Get("/", d.changeH.List)
		r.Get("/{id}", d.changeH.Get)
		r.Post("/{id}/ack", d.changeH.Acknowledge)
		r.Post("/{id}/dismiss", d.changeH.Dismiss)
		r.Post("/{id}/ai-update", d.changeH.AIUpdate)
		r.Post("/{id}/explain", d.changeH.Explain)
		r.Post("/bulk-resolve", d.changeH.BulkResolve)
	})

	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", d.alertH.List)
		r.Get("/{id}", d.alertH.Get)
		r.Post("/{id}/resolve", d.alertH.Resolve)
		r.Post("/{id}/dismiss", d.alertH.Dismiss)
		r.Post("/{id}/snooze", d.alertH.Snooze)
		r.Post("/bulk-snooze", d.alertH.BulkSnooze)
	})

	r.Route("/attention", func(r chi.Router) {
		r.Get("/", d.attentionH.List)
	})

	r.Route("/runbooks", func(r chi.Router) {
		r.Get("/", d.runbookH.List)
		r.Get("/{id}", d.runbookH.Get)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireInstanceAdmin)
			r.Post("/", d.runbookH.Create)
			r.Put("/{id}", d.runbookH.Update)
			r.Delete("/{id}", d.runbookH.Delete)
		})
	})

	r.Route("/findings", func(r chi.Router) {
		r.Get("/", d.findingH.List)
		r.Get("/{id}", d.findingH.Get)
		r.Post("/{id}/resolve", d.findingH.Resolve)
	})

	r.Route("/notifications", func(r chi.Router) {
		r.Get("/", d.notifH.List)
		r.Post("/read-all", d.notifH.ReadAll)
		r.Post("/{id}/read", d.notifH.MarkRead)
	})

	// Saved views are a personal resource: any authenticated role (viewer
	// or operator) may list/create/delete their own, scoped by user_id in
	// the handler — no operatorOnly gate here.
	r.Route("/saved-views", func(r chi.Router) {
		r.Get("/", d.savedViewH.List)
		r.Post("/", d.savedViewH.Create)
		r.Delete("/{id}", d.savedViewH.Delete)
	})
}

// mountDashboardRoutes registers the /dashboard tree. It must be called on an
// already-authenticated group.
func mountDashboardRoutes(r chi.Router, d routerDeps) {
	dashboardDefaultOnly := auth.RequirePermission(d.cfg.Store, "can_manage_dashboard_defaults")

	r.Route("/dashboard", func(r chi.Router) {
		r.Get("/overview", d.dashH.Overview)
		// Dashboard layout is a personal resource: any authenticated role
		// (viewer or operator) may read/save/reset their own, scoped by
		// user_id in the handler — no operatorOnly gate here.
		r.Get("/layout", d.dashH.GetLayout)
		r.Put("/layout", d.dashH.SaveLayout)
		r.Post("/layout/reset", d.dashH.ResetLayout)

		r.Group(func(r chi.Router) {
			r.Use(dashboardDefaultOnly)
			r.Get("/layout/admin-default", d.dashH.GetAdminDefault)
			r.Put("/layout/admin-default", d.dashH.PutAdminDefault)
		})
	})
}
