package api

import (
	"github.com/go-chi/chi/v5"
)

// mountSettingsRoutes registers the instance-wide AI and notification settings.
// It must be called on an already-instance-admin group.
func mountSettingsRoutes(r chi.Router, d routerDeps) {
	r.Route("/ai/config", func(r chi.Router) {
		r.Get("/", d.settingH.GetAIConfig)
		r.Put("/", d.settingH.UpdateAIConfig)
		r.Post("/test", d.settingH.TestAIConfig)
		r.Get("/fallback-providers", d.settingH.GetAIFallbackProviders)
		r.Put("/fallback-providers", d.settingH.UpdateAIFallbackProviders)
	})

	r.Route("/notifications/config", func(r chi.Router) {
		r.Get("/", d.settingH.GetNotificationsConfig)
		r.Put("/", d.settingH.UpdateNotificationsConfig)
		r.Post("/test", d.settingH.TestNotificationsConfig)
	})

	r.Get("/notifications/deliveries", d.notifH.ListDeliveries)
}

// mountSystemRoutes registers the /system tree (info, audit, retention, backup,
// diagnostics). It must be called on an already-instance-admin group.
func mountSystemRoutes(r chi.Router, d routerDeps) {
	r.Get("/system/info", d.sysH.Info)
	r.Get("/system/audit", d.sysH.ListAudit)
	r.Get("/system/audit/export", d.sysH.ExportAudit)
	r.Get("/system/jobs", d.sysH.GetJobs)

	r.Get("/system/settings/retention", d.sysH.GetRetentionSettings)
	r.Put("/system/settings/retention", d.sysH.UpdateRetentionSettings)

	r.Route("/system/backup", func(r chi.Router) {
		r.Get("/export", d.sysH.ExportBackup)
		r.Post("/import", d.sysH.ImportBackup)
		r.Get("/schedule", d.sysH.GetBackupSchedule)
		r.Put("/schedule", d.sysH.UpdateBackupSchedule)
		r.Get("/runs", d.sysH.ListBackupRuns)
		r.Post("/run", d.sysH.CreateBackupRun)
	})

	r.Get("/system/diagnostics", d.sysH.Diagnostics)
}

// mountComplianceRoutes registers the /compliance tree. It must be called on an
// already-instance-admin group.
func mountComplianceRoutes(r chi.Router, d routerDeps) {
	r.Route("/compliance", func(r chi.Router) {
		r.Get("/schema", d.complianceH.Schema)
		r.Get("/rules", d.complianceH.List)
		r.Post("/rules", d.complianceH.Create)
		r.Get("/rules/{id}", d.complianceH.Get)
		r.Put("/rules/{id}", d.complianceH.Update)
		r.Delete("/rules/{id}", d.complianceH.Delete)
		r.Post("/rules/test", d.complianceH.Test)
	})
}
