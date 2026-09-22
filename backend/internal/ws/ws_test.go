package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHubConfiguredOrigin(t *testing.T) {
	hub := NewHub("https://app.example.com")
	allowed := httptest.NewRequest("GET", "/api/ws", nil)
	allowed.Header.Set("Origin", "https://app.example.com")
	if !hub.upgrader.CheckOrigin(allowed) {
		t.Fatal("configured origin was rejected")
	}
	rejected := httptest.NewRequest("GET", "/api/ws", nil)
	rejected.Header.Set("Origin", "https://evil.example.com")
	if hub.upgrader.CheckOrigin(rejected) {
		t.Fatal("cross-origin request was accepted")
	}
}

func TestHubBroadcastRouting(t *testing.T) {
	hub := NewHub()
	go hub.Run(context.Background())

	first := &Client{hub: hub, send: make(chan []byte, 1), userID: "first"}
	second := &Client{hub: hub, send: make(chan []byte, 1), userID: "second"}
	hub.register <- first
	hub.register <- second
	t.Cleanup(func() {
		hub.unregister <- first
		hub.unregister <- second
	})

	deadline := time.Now().Add(time.Second)
	for hub.ClientCount() != 2 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := hub.ClientCount(); got != 2 {
		t.Fatalf("ClientCount() = %d, want 2", got)
	}
	hub.Broadcast(EventSyncComplete, map[string]string{"job": "all"})
	assertEnvelope(t, <-first.send, EventSyncComplete, "all")
	assertEnvelope(t, <-second.send, EventSyncComplete, "all")

	hub.BroadcastToUser("first", EventSystemNotice, map[string]string{"job": "one"})
	assertEnvelope(t, <-first.send, EventSystemNotice, "one")
	select {
	case data := <-second.send:
		t.Fatalf("second user received %s", data)
	case <-time.After(25 * time.Millisecond):
	}
}

func assertEnvelope(t *testing.T, data []byte, wantType, wantJob string) {
	t.Helper()
	var got struct {
		Type    string            `json:"type"`
		Payload map[string]string `json:"payload"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if got.Type != wantType || got.Payload["job"] != wantJob {
		t.Fatalf("envelope = %+v", got)
	}
}

func TestDocLockEventBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run(context.Background())

	client := &Client{hub: hub, send: make(chan []byte, 3), userID: "test"}
	hub.register <- client
	t.Cleanup(func() {
		hub.unregister <- client
	})

	deadline := time.Now().Add(time.Second)
	for hub.ClientCount() != 1 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := hub.ClientCount(); got != 1 {
		t.Fatalf("ClientCount() = %d, want 1", got)
	}

	hub.Broadcast(EventDocLockAcquired, map[string]string{"docId": "doc-1", "userId": "user-1"})
	hub.Broadcast(EventDocLockReleased, map[string]string{"docId": "doc-1", "userId": "user-1"})
	hub.Broadcast(EventDocLockExpired, map[string]string{"docId": "doc-1", "userId": "user-1"})

	var env struct {
		Type string `json:"type"`
	}

	data := <-client.send
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("unmarshal event 1: %v", err)
	}
	if env.Type != EventDocLockAcquired {
		t.Fatalf("event 1: got type %q, want %q", env.Type, EventDocLockAcquired)
	}

	data = <-client.send
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("unmarshal event 2: %v", err)
	}
	if env.Type != EventDocLockReleased {
		t.Fatalf("event 2: got type %q, want %q", env.Type, EventDocLockReleased)
	}

	data = <-client.send
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("unmarshal event 3: %v", err)
	}
	if env.Type != EventDocLockExpired {
		t.Fatalf("event 3: got type %q, want %q", env.Type, EventDocLockExpired)
	}
}

// TestUpgradeHandlerAndWritePump tests the full round trip: UpgradeHandler upgrades
// the connection, and writePump sends messages to the client.
func TestUpgradeHandlerAndWritePump(t *testing.T) {
	hub := NewHub()
	go hub.Run(context.Background())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hub.UpgradeHandler(w, r, "user-test-123", "admin", "")
		if err != nil {
			t.Logf("UpgradeHandler error: %v", err)
		}
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Wait for client to register
	deadline := time.Now().Add(time.Second)
	for hub.ClientCount() != 1 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	if got := hub.ClientCount(); got != 1 {
		t.Fatalf("ClientCount() = %d, want 1", got)
	}

	// Broadcast a message and verify it reaches the client
	hub.Broadcast(EventSyncComplete, map[string]string{"job": "test-job"})

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage error: %v", err)
	}

	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("Unmarshal envelope error: %v", err)
	}

	if env.Type != EventSyncComplete {
		t.Errorf("envelope type = %q, want %q", env.Type, EventSyncComplete)
	}

	payload, ok := env.Payload.(map[string]interface{})
	if !ok {
		t.Errorf("payload type = %T, want map[string]interface{}", env.Payload)
	}
	if payload["job"] != "test-job" {
		t.Errorf("payload job = %v, want test-job", payload["job"])
	}
}

// TestReadPumpGarbageInput tests that readPump handles non-JSON garbage input
// without crashing the server.
func TestReadPumpGarbageInput(t *testing.T) {
	hub := NewHub()
	go hub.Run(context.Background())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hub.UpgradeHandler(w, r, "user-garbage-456", "viewer", "")
		if err != nil {
			t.Logf("UpgradeHandler error: %v", err)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Wait for client to register
	deadline := time.Now().Add(time.Second)
	for hub.ClientCount() != 1 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	// Send garbage (non-JSON) data to the server
	err = conn.WriteMessage(websocket.TextMessage, []byte("this is not json {malformed"))
	if err != nil {
		t.Fatalf("WriteMessage error: %v", err)
	}

	// Server should not crash. Connection may stay alive or close gracefully.
	// Set a read deadline to detect if connection is responsive.
	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, _ = conn.ReadMessage()
	// Error is acceptable (timeout or connection closed); no error also OK.
	// The important thing is that the hub and server don't crash.

	// Verify hub is still responsive by broadcasting to a different client
	if hub.ClientCount() == 0 {
		// Original client disconnected, which is fine
		return
	}

	// If client is still connected, verify hub still processes broadcasts
	hub.Broadcast(EventSystemHealth, map[string]string{"status": "ok"})
	// If we get here without crash, test passes
}

// TestClientCloseDisconnect tests that readPump and writePump properly
// handle client disconnection and remove the client from the hub.
func TestClientCloseDisconnect(t *testing.T) {
	hub := NewHub()
	go hub.Run(context.Background())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hub.UpgradeHandler(w, r, "user-close-789", "operator", "")
		if err != nil {
			t.Logf("UpgradeHandler error: %v", err)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}

	// Wait for client to register
	deadline := time.Now().Add(time.Second)
	for hub.ClientCount() != 1 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	if got := hub.ClientCount(); got != 1 {
		t.Fatalf("ClientCount() = %d, want 1 after connect", got)
	}

	// Close the client connection
	_ = conn.Close()

	// Wait for the hub to unregister the client
	deadline = time.Now().Add(time.Second)
	for hub.ClientCount() != 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	if got := hub.ClientCount(); got != 0 {
		t.Fatalf("ClientCount() = %d after close, want 0", got)
	}
}

// TestBroadcastToUserAfterUpgrade tests BroadcastToUser with a real WebSocket
// connection, verifying that UpgradeHandler and writePump work together.
func TestBroadcastToUserAfterUpgrade(t *testing.T) {
	hub := NewHub()
	go hub.Run(context.Background())

	// Create two separate servers/connections for two users
	conn1 := setupWSConnection(t, hub, "alice")
	defer func() { _ = conn1.Close() }()

	conn2 := setupWSConnection(t, hub, "bob")
	defer func() { _ = conn2.Close() }()

	// Wait for both to register
	deadline := time.Now().Add(time.Second)
	for hub.ClientCount() != 2 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	// Broadcast to alice only
	hub.BroadcastToUser("alice", EventAlertCreated, map[string]string{"alertId": "alert-1"})

	// Alice should receive it
	_ = conn1.SetReadDeadline(time.Now().Add(time.Second))
	_, data, err := conn1.ReadMessage()
	if err != nil {
		t.Fatalf("conn1.ReadMessage error: %v", err)
	}

	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if env.Type != EventAlertCreated {
		t.Errorf("alice got type %q, want %q", env.Type, EventAlertCreated)
	}

	// Bob should not receive it (with timeout)
	_ = conn2.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err = conn2.ReadMessage()
	if err == nil {
		t.Error("bob received a message meant for alice")
	}
}

// setupWSConnection is a helper that creates a WebSocket connection to a test server
// backed by the given hub, for a specific user.
func setupWSConnection(t *testing.T, hub *Hub, userID string) *websocket.Conn {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hub.UpgradeHandler(w, r, userID, "user", "")
		if err != nil {
			t.Logf("UpgradeHandler error: %v", err)
		}
	}))
	t.Cleanup(func() { server.Close() })

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial error for %s: %v", userID, err)
	}

	return conn
}

func TestHubMultipleConfiguredOrigins(t *testing.T) {
	hub := NewHub("https://a.example.com, https://b.example.com:8443/", "http://c.example.com")
	check := func(origin string) bool {
		r := httptest.NewRequest("GET", "/api/ws", nil)
		if origin != "" {
			r.Header.Set("Origin", origin)
		}
		return hub.upgrader.CheckOrigin(r)
	}
	for _, ok := range []string{"https://a.example.com", "https://b.example.com:8443", "http://c.example.com", "HTTPS://A.example.com"} {
		if !check(ok) {
			t.Errorf("origin %q rejected", ok)
		}
	}
	for _, bad := range []string{"", "https://evil.example.com", "https://a.example.com.evil.com", "https://b.example.com", "https://a.example.com,https://b.example.com:8443"} {
		if check(bad) {
			t.Errorf("origin %q accepted", bad)
		}
	}
}

func TestHubMalformedOriginsFailClosed(t *testing.T) {
	hub := NewHub("not-a-url", "*", "ftp://x")
	r := httptest.NewRequest("GET", "/api/ws", nil)
	r.Header.Set("Origin", "https://app.example.com")
	r.Host = "app.example.com"
	if hub.upgrader.CheckOrigin(r) {
		t.Fatal("malformed configuration must not fall back to accepting origins")
	}
}

func TestTicketsAreSingleUseAndExpire(t *testing.T) {
	hub := NewHub()
	id, err := hub.IssueTicket("u1", "admin", "sess")
	if err != nil {
		t.Fatal(err)
	}
	uid, role, sess, ok := hub.RedeemTicket(id)
	if !ok || uid != "u1" || role != "admin" || sess != "sess" {
		t.Fatalf("redeem = %q %q %q %v", uid, role, sess, ok)
	}
	if _, _, _, ok := hub.RedeemTicket(id); ok {
		t.Error("ticket redeemed twice")
	}
	if _, _, _, ok := hub.RedeemTicket("bogus"); ok {
		t.Error("unknown ticket accepted")
	}

	id, _ = hub.IssueTicket("u1", "admin", "")
	hub.ticketMu.Lock()
	tk := hub.tickets[id]
	tk.expires = time.Now().Add(-time.Second)
	hub.tickets[id] = tk
	hub.ticketMu.Unlock()
	if _, _, _, ok := hub.RedeemTicket(id); ok {
		t.Error("expired ticket accepted")
	}
}

// TestRevalidationClosesConnection verifies that once the revalidator reports the
// identity as invalid (disabled/revoked/role change), the open socket is closed on the next tick.
func TestRevalidationClosesConnection(t *testing.T) {
	hub := NewHub()
	hub.pingInterval = 20 * time.Millisecond
	var valid atomic.Bool
	valid.Store(true)
	hub.SetRevalidator(func(_ context.Context, userID, _, _ string) bool {
		return userID == "u1" && valid.Load()
	})
	go hub.Run(context.Background())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = hub.UpgradeHandler(w, r, "u1", "viewer", "sess")
	}))
	defer server.Close()
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	errCh := make(chan error, 1)
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Still valid across several ticks: the connection must stay open.
	select {
	case err := <-errCh:
		t.Fatalf("connection closed while identity was valid: %v", err)
	case <-time.After(150 * time.Millisecond):
	}

	valid.Store(false)
	select {
	case err := <-errCh:
		if !websocket.IsCloseError(err, websocket.ClosePolicyViolation) {
			t.Fatalf("expected policy-violation close, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("connection not closed after revalidation failure")
	}
}

func TestBroadcastFullQueueDoesNotBlock(t *testing.T) {
	for _, targeted := range []bool{false, true} {
		hub := NewHub()
		for i := 0; i < cap(hub.broadcast); i++ {
			hub.broadcast <- broadcastMsg{}
		}
		done := make(chan struct{})
		go func() {
			if targeted {
				hub.BroadcastToUser("user", EventAlertCreated, nil)
			} else {
				hub.Broadcast(EventSyncComplete, nil)
			}
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(time.Second):
			<-hub.broadcast // Release the blocked call on the old implementation.
			<-done
			t.Errorf("broadcast blocked with targeted=%v", targeted)
		}
	}
}

func TestSlowClientEviction(t *testing.T) {
	hub := NewHub()
	connections := make(chan *websocket.Conn, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := hub.upgrader.Upgrade(w, r, nil)
		if err == nil {
			connections <- conn
		}
	}))
	defer server.Close()
	peer, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close() //nolint:errcheck
	conn := <-connections
	defer conn.Close() //nolint:errcheck
	// No write pump: deterministically model a client whose queue cannot drain.
	slow := &Client{conn: conn, send: make(chan []byte, 1), userID: "slow"}
	healthy := &Client{send: make(chan []byte, 1), userID: "healthy"}
	hub.clients[slow] = true
	hub.clients[healthy] = true
	slow.send <- []byte("queued")
	msg := broadcastMsg{data: []byte("event")}
	for i := 0; i < maxConsecutiveDrops-1; i++ {
		hub.deliver(msg)
		<-healthy.send
	}
	if hub.ClientCount() != 2 {
		t.Fatal("client evicted before drop threshold")
	}
	// A successful delivery resets the streak.
	<-slow.send
	hub.deliver(msg)
	<-healthy.send
	for i := 0; i < maxConsecutiveDrops-1; i++ {
		hub.deliver(msg)
		<-healthy.send
	}
	if hub.ClientCount() != 2 {
		t.Fatal("successful send did not reset drop streak")
	}
	// Traffic for another user neither increments nor resets the streak.
	hub.deliver(broadcastMsg{data: msg.data, userID: "healthy"})
	<-healthy.send
	hub.deliver(msg)
	<-healthy.send
	if hub.ClientCount() != 1 {
		t.Fatal("slow client was not evicted")
	}
	<-slow.send
	if _, ok := <-slow.send; ok {
		t.Fatal("evicted client's queue remains open")
	}
	_ = peer.SetReadDeadline(time.Now().Add(time.Second))
	if _, _, err := peer.ReadMessage(); err == nil {
		t.Fatal("evicted client's socket remains open")
	} else if timeout, ok := err.(interface{ Timeout() bool }); ok && timeout.Timeout() {
		t.Fatal("socket timed out instead of closing")
	}
	// The read pump may unregister an already evicted client; it must be safe.
	go hub.Run(context.Background())
	hub.unregister <- slow
	hub.BroadcastToUser("healthy", EventSystemNotice, nil)
	select {
	case <-healthy.send:
	case <-time.After(time.Second):
		t.Fatal("hub stopped delivering after eviction and unregister")
	}
	hub.unregister <- healthy
}
