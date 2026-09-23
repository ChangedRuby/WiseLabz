package api

import (
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/go-chi/chi/v5"
)

func mountReportRoutes(r chi.Router, d routerDeps) {
	r.Route("/reports", func(r chi.Router) {
		r.Use(auth.RequireInstanceAdmin)
		r.Get("/", d.reportH.List)
		r.Get("/{id}", d.reportH.Get)
		r.Get("/{id}/download", d.reportH.Download)
		r.Get("/definitions", d.reportH.ListDefinitions)
		r.Post("/definitions", d.reportH.Create)
		r.Put("/definitions/{id}", d.reportH.Update)
		r.Delete("/definitions/{id}", d.reportH.Delete)
		r.Post("/definitions/{id}/run", d.reportH.Run)
	})
}
