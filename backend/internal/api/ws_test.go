package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/WiseLabz/wiselabz/internal/ws"
)

func TestWebSocketTicketFlow(t *testing.T) {
	t.Parallel()
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

// wsDial mints a ticket with bearer and opens a socket, failing the test on error.
func wsDial(t *testing.T, app *testApp, srv *httptest.Server, bearer string) <-chan ws.Envelope {
	t.Helper()
	rec := app.req(t, http.MethodPost, "/api/ws/ticket", nil, bearer)
	if rec.Code != http.StatusOK {
		t.Fatalf("ticket status = %d: %s", rec.Code, rec.Body)
	}
	var out struct {
		Ticket string `json:"ticket"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil || out.Ticket == "" {
		t.Fatalf("decode ticket: %v %+v", err, out)
	}
	h := http.Header{"Origin": []string{"http://localhost:5173"}}
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http")+"/api/ws?ticket="+out.Ticket, h)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	// A gorilla read timeout is fatal to the connection, so read in a goroutine
	// and let callers time out on the channel instead.
	frames := make(chan ws.Envelope, 16)
	go func() {
		for {
			var env ws.Envelope
			if err := conn.ReadJSON(&env); err != nil {
				return
			}
			if env.Type != ws.EventSystemHealth {
				frames <- env
			}
		}
	}()
	return frames
}

func waitForClients(t *testing.T, hub *ws.Hub, n int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for hub.ClientCount() != n && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := hub.ClientCount(); got != n {
		t.Fatalf("ClientCount() = %d, want %d", got, n)
	}
}

// wsRead returns the next event frame; ok is false when none arrives within wait.
func wsRead(frames <-chan ws.Envelope, wait time.Duration) (ws.Envelope, bool) {
	select {
	case env := <-frames:
		return env, true
	case <-time.After(wait):
		return ws.Envelope{}, false
	}
}

func TestWebSocketConnectorEventsFilteredByGrant(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	srv := httptest.NewServer(app.Router)
	defer srv.Close()

	conn1 := newConnector(t, app, "one")
	conn2 := newConnector(t, app, "two")
	grantedID, granted := app.user(t, "viewer")
	_, ungranted := app.user(t, "viewer")
	// An instance admin gets no bypass: without a grant it receives nothing.
	_, adminTok := app.user(t, "operator")
	app.connectorGrant(t, grantedID, conn1, "viewer")

	restrictedOwner, ownerTok := app.user(t, "operator")
	app.connectorGrant(t, restrictedOwner, conn1, "operator")
	app.connectorGrant(t, restrictedOwner, conn2, "operator")
	key := createKey(t, app, ownerTok, map[string]any{"name": "r", "connectorIds": []string{conn2}})["token"].(string)

	grantedConn := wsDial(t, app, srv, granted)
	ungrantedConn := wsDial(t, app, srv, ungranted)
	adminConn := wsDial(t, app, srv, adminTok)
	// The admin-owned restricted key completes the upgrade (role label is "user").
	keyConn := wsDial(t, app, srv, key)
	waitForClients(t, app.WSHub, 4)

	app.WSHub.BroadcastConnector(conn1, ws.EventSyncProgress, map[string]string{"phase": "x"})
	env, ok := wsRead(grantedConn, time.Second)
	if !ok || env.Type != ws.EventSyncProgress || env.ConnectorID != conn1 {
		t.Fatalf("granted socket got %+v ok=%v", env, ok)
	}
	for name, c := range map[string]<-chan ws.Envelope{"ungranted": ungrantedConn, "admin": adminConn, "restricted key": keyConn} {
		if env, ok := wsRead(c, 50*time.Millisecond); ok {
			t.Errorf("%s socket received connector event %+v", name, env)
		}
	}

	// A global event reaches everyone.
	app.WSHub.Broadcast(ws.EventSystemNotice, map[string]string{"msg": "hi"})
	for name, c := range map[string]<-chan ws.Envelope{"granted": grantedConn, "ungranted": ungrantedConn, "admin": adminConn, "restricted key": keyConn} {
		if env, ok := wsRead(c, time.Second); !ok || env.Type != ws.EventSystemNotice {
			t.Errorf("%s socket missed global event: %+v ok=%v", name, env, ok)
		}
	}

	// The key's own connector reaches it.
	app.WSHub.BroadcastConnector(conn2, ws.EventSyncProgress, map[string]string{"phase": "y"})
	if env, ok := wsRead(keyConn, time.Second); !ok || env.ConnectorID != conn2 {
		t.Errorf("restricted key missed its connector event: %+v ok=%v", env, ok)
	}
}

func TestWebSocketRevokedKeyFailsRevalidation(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	userID, token := app.user(t, "viewer")
	created := createKey(t, app, token, map[string]any{"name": "k"})
	keyID := created["id"].(string)
	ident := ws.Identity{UserID: userID, Role: "user", APIKeyID: keyID}

	if !app.WSHub.Revalidate(context.Background(), ident) {
		t.Fatal("active key failed revalidation")
	}
	if err := app.Store.RevokeAPIKey(context.Background(), keyID); err != nil {
		t.Fatalf("revoke key: %v", err)
	}
	if app.WSHub.Revalidate(context.Background(), ident) {
		t.Error("revoked key passed revalidation")
	}
	other, _ := app.user(t, "viewer")
	if app.WSHub.Revalidate(context.Background(), ws.Identity{UserID: other, Role: "user", APIKeyID: keyID}) {
		t.Error("key owned by another user passed revalidation")
	}
}
