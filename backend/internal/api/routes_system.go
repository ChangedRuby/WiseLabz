package api

import "github.com/go-chi/chi/v5"

func mountRootRoutes(r chi.Router, d routerDeps) {
	r.Get("/healthz", d.sysH.Liveness)
	r.Get("/readyz", d.sysH.Readiness)
}

func mountPublicSystemRoutes(r chi.Router, d routerDeps) {
	r.Get("/health", d.sysH.Health)
}
