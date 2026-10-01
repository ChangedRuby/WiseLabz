package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/WiseLabz/wiselabz/internal/api/middleware"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
)

// mountAuthRoutes registers the /auth tree: the unauthenticated login/refresh
// endpoints plus the auth-gated session, API-key and auth-config groups.
//
// chi only allows one Mount per exact pattern, so the protected /api/auth
// routes (logout, elevate, config) are nested inside this same Route as
// inner auth-gated groups rather than a second top-level r.Route("/auth", ...).
func mountAuthRoutes(r chi.Router, d routerDeps) {
	cfg := d.cfg

	// Per-IP throttle on unauthenticated auth endpoints: 5 requests/sec with a
	// burst of 10, so a normal user retrying a typo never trips it but a
	// sustained guessing campaign against one IP does.
	authIPLimit := middleware.RateLimit(5, 10, func(r *http.Request) string { return httputil.ClientIP(r, cfg.Config.Server.TrustedProxies) })
	// Per-user throttle on password re-verification for step-up auth.
	elevateLimit := middleware.RateLimit(1, 5, func(r *http.Request) string { return auth.UserIDFromContext(r.Context()) })

	r.Route("/auth", func(r chi.Router) {
		r.With(authIPLimit).Post("/login", d.authH.Login)
		r.With(authIPLimit).Post("/login/mfa", d.authH.LoginMFA)
		r.With(authIPLimit).Post("/login/mfa/webauthn/begin", d.authH.PostWebAuthnLoginBegin)
		r.With(authIPLimit).Post("/oidc/callback", d.authH.OIDCCallback)
		r.With(authIPLimit).Post("/refresh", d.authH.Refresh)
		r.Get("/providers", d.authH.Providers)

		r.Group(func(r chi.Router) {
			r.Use(cfg.AuthMiddleware())
			r.Post("/logout", d.authH.Logout)
			// API keys must not mint elevation tokens.
			r.Group(func(r chi.Router) {
				r.Use(auth.RejectAPIKey)
				r.With(elevateLimit).Post("/elevate", d.authH.Elevate)
				r.With(elevateLimit).Post("/elevate/webauthn/begin", d.authH.PostWebAuthnElevateBegin)
				r.Get("/elevate/methods", d.authH.ElevateMethods)
				r.With(elevateLimit).Post("/elevate/oidc/begin", d.authH.ElevateOIDCBegin)
				r.With(elevateLimit).Post("/elevate/oidc/complete", d.authH.ElevateOIDCComplete)
			})
			r.Route("/api-keys", func(r chi.Router) {
				r.Get("/", d.apiKeyH.List)
				r.Post("/", d.apiKeyH.Create)
				r.Delete("/{id}", d.apiKeyH.Revoke)
			})

			r.Group(func(r chi.Router) {
				r.Use(auth.RequireInstanceAdmin)
				r.Get("/config", d.settingH.GetAuthConfig)
				r.Put("/config", d.settingH.UpdateAuthConfig)
				r.Put("/providers/{providerId}/enabled", d.settingH.UpdateProviderEnabled)
			})
		})
	})
}

// mountMeRoutes registers the caller's own profile and session routes. It must
// be called on an already-authenticated group.
func mountMeRoutes(r chi.Router, d routerDeps) {
	cfg := d.cfg

	r.Route("/me", func(r chi.Router) {
		r.Get("/", d.authH.Me)
		r.Patch("/", d.authH.UpdateMe)
		r.With(auth.RejectAPIKey).Post("/password", d.authH.ChangePassword)
		r.Get("/sessions", d.authH.ListSessions)
		r.Delete("/sessions/{id}", d.authH.DeleteSession)

		r.Route("/mfa", func(r chi.Router) {
			r.Use(auth.RejectAPIKey)
			r.Get("/", d.authH.GetMFA)
			// Starting enrollment needs step-up; confirm/finish depend on state the start step creates.
			enrollElevation := auth.RequireElevationUnlessEnrollOnly(cfg.JWT, cfg.Store, cfg.Store, "mfa.manage")
			r.With(enrollElevation).Post("/totp", d.authH.PostMFATOTP)
			r.Post("/totp/{id}/confirm", d.authH.PostMFATOTPConfirm)
			r.With(enrollElevation).Post("/webauthn/register/begin", d.authH.PostWebAuthnRegisterBegin)
			r.Post("/webauthn/register/finish", d.authH.PostWebAuthnRegisterFinish)

			r.Group(func(r chi.Router) {
				r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "mfa.manage"))
				r.Post("/recovery-codes", d.authH.PostMFARecoveryCodes)
				r.Delete("/factors/{id}", d.authH.DeleteMFAFactor)
			})
		})
	})
}

// mountUserRoutes registers instance-admin user management. It must be called
// on an already-instance-admin group.
func mountUserRoutes(r chi.Router, d routerDeps) {
	cfg := d.cfg

	r.Route("/users", func(r chi.Router) {
		r.Get("/", d.authH.ListUsers)
		r.Post("/", d.authH.CreateUser)
		r.Patch("/{id}", d.authH.UpdateUser)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "user.delete"))
			r.Delete("/{id}", d.authH.DeleteUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "user.resetPassword"))
			r.Post("/{id}/reset-password", d.authH.ResetPassword)
		})

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireElevation(cfg.JWT, cfg.Store, "user.resetMfa"))
			r.Post("/{id}/reset-mfa", d.authH.ResetMFA)
		})
	})
}
