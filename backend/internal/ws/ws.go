// Package ws provides the WebSocket hub, client management, and event types.
package ws

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/WiseLabz/wiselabz/internal/logsafe"
)

// Event types sent over WebSocket.
const (
	EventServiceStatus          = "service.status"
	EventSyncProgress           = "sync.progress"
	EventSyncComplete           = "sync.complete"
	EventChangeDetected         = "change.detected"
	EventAlertCreated           = "alert.created"
	EventAlertResolved          = "alert.resolved"
	EventQualityFindingCreated  = "quality.finding.created"
	EventQualityFindingsChanged = "quality.findings.changed"
	EventDocGenerated           = "doc.generated"
	EventDocAISuggestion        = "doc.ai_suggestion"
	EventDocLockAcquired        = "doc.lock.acquired"
	EventDocLockReleased        = "doc.lock.released"
	EventDocLockExpired         = "doc.lock.expired"
	EventSystemHealth           = "system.health"
	EventSystemNotice           = "system.notice"
)

// Envelope wraps all WebSocket messages.
type Envelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// Client represents a single WebSocket connection.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	role   string
	// sessionHash is the hash of the refresh token the connection was ticketed
	// under ("" when it was not issued from a cookie session, e.g. API key).
	sessionHash      string
	consecutiveDrops int // Protected by hub.mu.
}

// Revalidator reports whether a connection's identity is still acceptable:
// the user still exists, is enabled, still holds the same role, and (when
// sessionHash is non-empty) the session is still active.
type Revalidator func(ctx context.Context, userID, role, sessionHash string) bool

// ticketTTL bounds how long an issued ticket can wait to be redeemed.
const ticketTTL = 30 * time.Second

type ticket struct {
	userID, role, sessionHash string
	expires                   time.Time
}

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	upgrader   websocket.Upgrader
	mu         sync.RWMutex
	clients    map[*Client]bool
	broadcast  chan broadcastMsg
	register   chan *Client
	unregister chan *Client

	revalidate   Revalidator
	pingInterval time.Duration

	ticketMu sync.Mutex
	tickets  map[string]ticket
}

type broadcastMsg struct {
	data   []byte
	userID string // empty = all clients
}

// NewHub creates a new WebSocket hub and starts its run loop.
func NewHub(origins ...string) *Hub {
	h := &Hub{
		upgrader:   websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		clients:    make(map[*Client]bool),
		broadcast:  make(chan broadcastMsg, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		pingInterval: 25 * time.Second,
		tickets:      make(map[string]ticket),
	}
	// Each argument may itself be a comma-separated list. With none configured
	// gorilla's default same-origin check applies; if some were configured but
	// none are usable, every cross-origin upgrade is refused.
	var configured, allowed []string
	for _, arg := range origins {
		for _, o := range strings.Split(arg, ",") {
			if o = strings.TrimSpace(o); o == "" {
				continue
			}
			configured = append(configured, o)
			if o = normalizeOrigin(o); o != "" {
				allowed = append(allowed, o)
			} else {
				slog.Warn("ignoring malformed WebSocket origin", "origin", logsafe.Sanitize(o))
			}
		}
	}
	if len(configured) > 0 {
		h.upgrader.CheckOrigin = func(r *http.Request) bool {
			origin := normalizeOrigin(r.Header.Get("Origin"))
			if origin == "" {
				return false
			}
			for _, a := range allowed {
				if a == origin {
					return true
				}
			}
			return false
		}
	}
	return h
}

// normalizeOrigin lowercases scheme://host[:port] and returns "" for anything
// that is not a bare http(s) origin.
func normalizeOrigin(o string) string {
	u, err := url.Parse(strings.TrimRight(strings.TrimSpace(o), "/"))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.Path != "" || u.RawQuery != "" || u.User != nil {
		return ""
	}
	return strings.ToLower(u.Scheme + "://" + u.Host)
}

// SetRevalidator installs the periodic identity check run for every open
// connection on each ping interval. Call before serving connections.
func (h *Hub) SetRevalidator(fn Revalidator) { h.revalidate = fn }

// Revalidate runs the installed Revalidator (true when none is installed).
func (h *Hub) Revalidate(ctx context.Context, userID, role, sessionHash string) bool {
	return h.revalidate == nil || h.revalidate(ctx, userID, role, sessionHash)
}

// IssueTicket mints a one-time, short-lived ticket that authorizes a single
// WebSocket upgrade for the given identity. The caller must have authenticated
// the user with a normal access token.
func (h *Hub) IssueTicket(userID, role, sessionHash string) (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	id := base64.RawURLEncoding.EncodeToString(buf)
	now := time.Now()

	h.ticketMu.Lock()
	defer h.ticketMu.Unlock()
	for k, t := range h.tickets {
		if now.After(t.expires) {
			delete(h.tickets, k)
		}
	}
	h.tickets[id] = ticket{userID: userID, role: role, sessionHash: sessionHash, expires: now.Add(ticketTTL)}
	return id, nil
}

// RedeemTicket consumes a ticket, returning its identity. A ticket works once.
func (h *Hub) RedeemTicket(id string) (userID, role, sessionHash string, ok bool) {
	h.ticketMu.Lock()
	t, found := h.tickets[id]
	delete(h.tickets, id)
	h.ticketMu.Unlock()
	if !found || time.Now().After(t.expires) {
		return "", "", "", false
	}
	return t.userID, t.role, t.sessionHash, true
}

// Run starts the hub's event loop. Should be run in a goroutine.
func (h *Hub) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			count := len(h.clients)
			h.mu.Unlock()
			slog.Info("WebSocket client connected", "user_id", logsafe.Sanitize(client.userID), "total_clients", count)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			count := len(h.clients)
			h.mu.Unlock()
			slog.Info("WebSocket client disconnected", "user_id", logsafe.Sanitize(client.userID), "total_clients", count)

		case msg := <-h.broadcast:
			h.deliver(msg)

		case <-ticker.C:
			// Heartbeat
			heartbeat, _ := json.Marshal(Envelope{
				Type: EventSystemHealth,
				Payload: map[string]any{
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				},
			})
			h.deliver(broadcastMsg{data: heartbeat})
		}
	}
}

// maxConsecutiveDrops allows short bursts before evicting a stalled client.
const maxConsecutiveDrops = 10

// deliver applies the same backpressure policy to events and heartbeats.
func (h *Hub) deliver(msg broadcastMsg) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		if msg.userID != "" && msg.userID != client.userID {
			continue
		}
		select {
		case client.send <- msg.data:
			client.consecutiveDrops = 0
		default:
			client.consecutiveDrops++
			if client.consecutiveDrops >= maxConsecutiveDrops {
				delete(h.clients, client)
				close(client.send)
				// Closing the socket also interrupts a blocked write immediately.
				if client.conn != nil {
					client.conn.Close() //nolint:errcheck
				}
				slog.Warn("disconnecting slow WebSocket client", "user_id", logsafe.Sanitize(client.userID))
			}
		}
	}
}

// Broadcast queues a message for all clients, dropping it if the hub queue is full.
func (h *Hub) Broadcast(eventType string, payload any) {
	data, err := json.Marshal(Envelope{Type: eventType, Payload: payload})
	if err != nil {
		slog.Error("failed to marshal WS broadcast", "error", err)
		return
	}
	select {
	case h.broadcast <- broadcastMsg{data: data}:
	default:
		slog.Warn("WebSocket broadcast queue full, dropping message")
	}
}

// BroadcastToUser queues a message for a user's connections, dropping it if the hub queue is full.
func (h *Hub) BroadcastToUser(userID, eventType string, payload any) {
	data, err := json.Marshal(Envelope{Type: eventType, Payload: payload})
	if err != nil {
		slog.Error("failed to marshal WS broadcast", "error", err)
		return
	}
	select {
	case h.broadcast <- broadcastMsg{data: data, userID: userID}:
	default:
		slog.Warn("WebSocket broadcast queue full, dropping message", "user_id", logsafe.Sanitize(userID))
	}
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// UpgradeHandler upgrades an HTTP connection to WebSocket.
// Caller must authenticate before upgrading.
func (h *Hub) UpgradeHandler(w http.ResponseWriter, r *http.Request, userID, role, sessionHash string) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		role:   role,

		sessionHash: sessionHash,
	}

	h.register <- client

	go client.writePump()
	go client.readPump()

	return nil
}

// readPump reads messages from the WebSocket connection (pings/pongs only).
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close() //nolint:errcheck
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)) //nolint:errcheck
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)) //nolint:errcheck
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// writePump writes messages from the send channel to the WebSocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(c.hub.pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close() //nolint:errcheck
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) //nolint:errcheck
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) //nolint:errcheck
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) //nolint:errcheck
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
			if !c.hub.Revalidate(context.Background(), c.userID, c.role, c.sessionHash) {
				slog.Info("closing WebSocket: identity no longer valid", "user_id", logsafe.Sanitize(c.userID))
				c.conn.WriteControl(websocket.CloseMessage, //nolint:errcheck
					websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "session no longer valid"), time.Now().Add(time.Second))
				return
			}
		}
	}
}
