package docs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/WiseLabz/wiselabz/internal/ws"
)

// TestBroadcastDocEventScoping checks that a connector doc's lock event goes
// to that connector's grant holders only, while a lab-wide doc's goes to all.
func TestBroadcastDocEventScoping(t *testing.T) {
	h := newTestHandler(t)
	hub := ws.NewHub()
	hub.SetConnectorAudience(func(_ context.Context, connectorID string) ([]string, error) {
		if connectorID == "c1" {
			return []string{"granted"}, nil
		}
		return nil, nil
	})
	go hub.Run(t.Context())
	h.WSHub = hub

	dial := func(userID string) chan string {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = hub.UpgradeHandler(w, r, ws.Identity{UserID: userID, Role: "user"})
		}))
		t.Cleanup(srv.Close)
		conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http"), nil)
		if err != nil {
			t.Fatalf("dial: %v", err)
		}
		t.Cleanup(func() { _ = conn.Close() })
		frames := make(chan string, 8)
		go func() {
			for {
				var env ws.Envelope
				if conn.ReadJSON(&env) != nil {
					return
				}
				if env.Type == ws.EventDocLockAcquired {
					frames <- env.ConnectorID
				}
			}
		}()
		return frames
	}
	granted, ungranted := dial("granted"), dial("ungranted")
	deadline := time.Now().Add(2 * time.Second)
	for hub.ClientCount() != 2 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	h.broadcastDocEvent("c1", ws.EventDocLockAcquired, map[string]any{"docId": "d1"})
	select {
	case got := <-granted:
		if got != "c1" {
			t.Errorf("connectorId = %q, want c1", got)
		}
	case <-time.After(time.Second):
		t.Fatal("granted user missed connector doc lock")
	}
	select {
	case <-ungranted:
		t.Fatal("ungranted user received connector doc lock")
	case <-time.After(25 * time.Millisecond):
	}

	h.broadcastDocEvent("", ws.EventDocLockAcquired, map[string]any{"docId": "lab"})
	for name, ch := range map[string]chan string{"granted": granted, "ungranted": ungranted} {
		select {
		case got := <-ch:
			if got != "" {
				t.Errorf("%s: lab-wide connectorId = %q, want empty", name, got)
			}
		case <-time.After(time.Second):
			t.Errorf("%s user missed lab-wide doc lock", name)
		}
	}
}
