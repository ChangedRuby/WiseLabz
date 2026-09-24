package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// mountWSRoutes registers the WebSocket endpoint.
//
// The upgrade is authorized by a one-time ticket minted from an
// authenticated request (POST /api/ws/ticket), not by the long-lived
// refresh cookie. Open connections are re-validated on every ping.
func mountWSRoutes(r chi.Router, d routerDeps) {
	cfg := d.cfg
	if cfg.WSHub == nil {
		return
	}

	r.With(cfg.AuthMiddleware()).Post("/ws/ticket", func(w http.ResponseWriter, r *http.Request) {
		// The WebSocket stream isn't filtered per API key, so a restricted
		// key could read events for connectors outside its restriction.
		if auth.RejectRestrictedAPIKey(w, r) {
			return
		}
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
