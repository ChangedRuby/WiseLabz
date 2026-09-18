package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebSocketTicketFlow(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "viewer")
	srv := httptest.NewServer(app.Router)
	defer srv.Close()

	dial := func(query string) (*websocket.Conn, *http.Response, error) {
		h := http.Header{"Origin": []string{"http://localhost:5173"}}
		return websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http")+"/api/ws"+query, h)
	}
	mint := func() string {
		rec := app.req(t, http.MethodPost, "/api/ws/ticket", nil, token)
		if rec.Code != http.StatusOK {
			t.Fatalf("ticket status = %d: %s", rec.Code, rec.Body)
		}
		var out struct {
			Ticket string `json:"ticket"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil || out.Ticket == "" {
			t.Fatalf("decode ticket: %v %+v", err, out)
		}
		return out.Ticket
	}

	if rec := app.req(t, http.MethodPost, "/api/ws/ticket", nil, ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("unauthenticated ticket status = %d, want 401", rec.Code)
	}
	if _, resp, err := dial(""); err == nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("dial without ticket should be 401, err=%v", err)
	}

	ticket := mint()
	conn, _, err := dial("?ticket=" + ticket)
	if err != nil {
		t.Fatalf("dial with ticket: %v", err)
	}
	_ = conn.Close()
	if _, resp, err := dial("?ticket=" + ticket); err == nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("ticket reuse should be 401, err=%v", err)
	}

	// A ticket minted for a user who is then disabled is refused at upgrade.
	ticket = mint()
	if err := app.Store.UpdateUser(t.Context(), userID, map[string]any{"disabled": true}); err != nil {
		t.Fatalf("disable user: %v", err)
	}
	if _, resp, err := dial("?ticket=" + ticket); err == nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("disabled user's ticket should be 401, err=%v", err)
	}
}
