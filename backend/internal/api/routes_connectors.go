package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/WiseLabz/wiselabz/internal/auth"
)

// mountConnectorRoutes registers the /connectors tree. It must be called on an
// already-authenticated group.
func mountConnectorRoutes(r chi.Router, d routerDeps) {
	cfg := d.cfg

	connViewer := auth.RequireConnectorRole(cfg.Store, "viewer", "id")
	connOperator := auth.RequireConnectorRole(cfg.Store, "operator", "id")

	r.Route("/connectors", func(r chi.Router) {
		// List/maintenance-windows are cross-connector and filter inside
		// the handler; Get 404s inside the handler on a missing grant
		// (not 403, to avoid confirming existence) rather than gated
		// here. Every other GET takes a connector {id} and can use the
		// viewer-role middleware directly.
		r.Get("/", d.connH.List)
		r.Get("/schema", d.connH.Schema)
		r.Get("/{id}", d.connH.Get)
		r.Get("/maintenance-windows", d.connH.ListActiveMaintenance)

		r.Group(func(r chi.Router) {
			r.Use(connViewer)
			r.Get("/{id}/data", d.connH.Data)
			r.Get("/{id}/uptime", d.connH.Uptime)
			r.Get("/{id}/syncs", d.connH.Syncs)
			r.Get("/{id}/removal-impact", d.connH.RemovalImpact)
			r.Get("/{id}/config-fields", d.connH.ConfigFields)
			r.Get("/{id}/maintenance-window", d.connH.GetMaintenanceWindow)
		})

		// Creating a connector has no existing grant to check against, so
		// it's instance-admin only; the creator is auto-granted operator
		// on the new connector (see connH.Create).
		r.With(auth.RequireInstanceAdmin).Post("/", d.connH.Create)

		// Instance-admin-only grant management for this connector — a
		// deliberate exception to the per-connector-role pattern below,
		// same reasoning as saved-views/dashboard-layout: granting
		// access is an instance-wide action, not something a connector's
		// own operators can do to each other.
		r.Route("/{id}/permissions", func(r chi.Router) {
			r.Use(auth.RequireInstanceAdmin)
			r.Get("/", d.connH.ListPermissions)
			r.Put("/{userId}", d.connH.PutPermission)
			r.Delete("/{userId}", d.connH.DeletePermission)
		})

		r.Group(func(r chi.Router) {
			r.Use(connOperator)
			r.Post("/{id}/test", d.connH.Test)
			r.Post("/{id}/health", d.connH.Health)
			r.Post("/{id}/restart", d.connH.RestartPreview)
			r.Post("/{id}/start", d.connH.StartPreview)
			r.Post("/{id}/stop", d.connH.StopPreview)
			r.Post("/{id}/config-push", d.connH.ConfigPush)
			r.Put("/{id}", d.connH.Update)
			r.Put("/{id}/enabled", d.connH.ToggleEnabled)
			r.Post("/{id}/sync", d.connH.Sync)
			r.Post("/{id}/maintenance-window", d.connH.OpenMaintenanceWindow) // no elevation: reversible and time-boxed
			r.Delete("/{id}/maintenance-window", d.connH.CloseMaintenanceWindow)

			r.Group(func(r chi.Router) {
				r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "connector.delete"))
				r.Delete("/{id}", d.connH.Delete)
			})
		})

		// Bulk routes take a body-supplied ID list, not a path {id}, so
		// they can't use connOperator middleware — each handler checks
		// per-ID via store.UserHasConnectorRole itself.
		r.Group(func(r chi.Router) {
			r.Post("/bulk-sync", d.connH.BulkSync)
			r.Post("/bulk-reauth", d.connH.BulkReauth)

			r.Group(func(r chi.Router) {
				// One elevation token covers the whole bulk-restart batch,
				// via its own distinct action string (see BulkRestart).
				r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "connector.bulkRestart"))
				r.Post("/bulk-restart", d.connH.BulkRestart)
			})
		})
	})
}
