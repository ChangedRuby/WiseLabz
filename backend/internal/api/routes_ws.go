package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
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

	// A read-only key may mint a ticket: it only opens a read stream, so the
	// POST counts as safe. Connector events are filtered per connection.
	r.With(auth.TreatAsSafeMethod, cfg.AuthMiddleware()).Post("/ws/ticket", func(w http.ResponseWriter, r *http.Request) {
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
		ident := ws.Identity{
			UserID:       userID,
			Role:         wsRoleLabel(auth.InstanceAdminFromContext(r.Context())),
			SessionHash:  sessionHash,
			APIKeyID:     auth.APIKeyIDFromContext(r.Context()),
			ConnectorIDs: auth.APIKeyRestrictionFromContext(r.Context()).ConnectorIDs,
		}
		id, err := cfg.WSHub.IssueTicket(ident)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		httputil.JSON(w, http.StatusOK, map[string]any{"ticket": id})
	})
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ident, ok := cfg.WSHub.RedeemTicket(r.URL.Query().Get("ticket"))
		if !ok {
			httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Invalid or expired ticket")
			return
		}
		if !cfg.WSHub.Revalidate(r.Context(), ident) {
			httputil.Error(w, http.StatusUnauthorized, "unauthorized", "User not found or disabled")
			return
		}
		if err := cfg.WSHub.UpgradeHandler(w, r, ident); err != nil {
			slog.Error("WebSocket upgrade failed", "error", err)
		}
	})
}
