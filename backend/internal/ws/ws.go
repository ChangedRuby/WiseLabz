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
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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
	// EventSystemResync tells clients to refetch volatile state. The relay emits
	// it after its listener reconnects following a gap (ADR 0005).
	EventSystemResync = "system.resync"
)

// Envelope wraps all WebSocket messages. ID is unique per emitted event (clients
// use it to drop duplicates); TS is the emit time, UTC, millisecond RFC 3339.
// ConnectorID is set on connector-scoped events (see BroadcastConnector).
type Envelope struct {
	ID          string `json:"id"`
	TS          string `json:"ts"`
	Type        string `json:"type"`
	ConnectorID string `json:"connectorId,omitempty"`
	Payload     any    `json:"payload"`
}

// Audience kinds. An event goes to every client, to one user's clients, or to
// the clients allowed to read one connector.
const (
	AudienceAll       = "all"
	AudienceUser      = "user"
	AudienceConnector = "connector"
)

// Audience is who an event is for. It is the `audience` field of the ADR 0005
// relay payload: a connector audience travels as the connector ID, and each
// replica resolves its readers locally when it delivers.
type Audience struct {
	Kind        string `json:"kind"`
	UserID      string `json:"userID,omitempty"`
	ConnectorID string `json:"connectorID,omitempty"`
}

// Identity is what a ticket authorizes and what a connection carries.
type Identity struct {
	UserID string
	Role   string
	// SessionHash is the hash of the refresh token the connection was ticketed
	// under ("" when it was not issued from a cookie session, e.g. API key).
	SessionHash string
	// APIKeyID is the key the ticket was minted with ("" for a session).
	APIKeyID string
	// ConnectorIDs is the API key's connector restriction; empty means the
	// connection may see every connector its user can read.
	ConnectorIDs []string
}

// ConnectorAudience returns the IDs of the users allowed to read connectorID.
type ConnectorAudience func(ctx context.Context, connectorID string) ([]string, error)

// audienceTimeout bounds a connector audience lookup so a stuck query drops
// the event instead of blocking the emitter.
const audienceTimeout = 5 * time.Second

// newEnvelope stamps a fresh id and timestamp on an event.
func newEnvelope(eventType string, payload any) Envelope {
	return Envelope{
		ID:      uuid.NewString(),
		TS:      time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		Type:    eventType,
		Payload: payload,
	}
}

// newHeartbeat builds the periodic system.health frame.
func newHeartbeat() Envelope {
	return newEnvelope(EventSystemHealth, map[string]any{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Client represents a single WebSocket connection.
type Client struct {
	hub              *Hub
	conn             *websocket.Conn
	send             chan []byte
	id               Identity
	consecutiveDrops int // Protected by hub.mu.
}

// Revalidator reports whether a connection's identity is still acceptable:
// the user still exists, is enabled, still holds the same role, (when
// SessionHash is set) the session is still active, and (when APIKeyID is set)
// the key is still valid.
type Revalidator func(ctx context.Context, id Identity) bool

// ticketTTL bounds how long an issued ticket can wait to be redeemed.
const ticketTTL = 30 * time.Second

type ticket struct {
	id      Identity
	expires time.Time
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
	audience     ConnectorAudience
	pingInterval time.Duration
	// heartbeatInterval paces the system.health broadcast.
	heartbeatInterval time.Duration

	ticketMu sync.Mutex
	tickets  map[string]ticket
}

type broadcastMsg struct {
	data []byte
	aud  Audience // zero value = all clients
	// readers is the resolved user set of a connector audience.
	readers map[string]struct{}
}

// reaches reports whether c is in the message's audience. A connector event
// needs the user to be a reader and, for a connector-restricted API key, the
// connector to be in the key's list.
func (m broadcastMsg) reaches(c *Client) bool {
	switch m.aud.Kind {
	case AudienceUser:
		return c.id.UserID == m.aud.UserID
	case AudienceConnector:
		_, ok := m.readers[c.id.UserID]
		return ok && (len(c.id.ConnectorIDs) == 0 || slices.Contains(c.id.ConnectorIDs, m.aud.ConnectorID))
	}
	return true
}

// NewHub creates a new WebSocket hub and starts its run loop.
func NewHub(origins ...string) *Hub {
	h := &Hub{
		upgrader:   websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		clients:    make(map[*Client]bool),
		broadcast:  make(chan broadcastMsg, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		pingInterval:      25 * time.Second,
		heartbeatInterval: 30 * time.Second,
		tickets:           make(map[string]ticket),
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
func (h *Hub) Revalidate(ctx context.Context, id Identity) bool {
	return h.revalidate == nil || h.revalidate(ctx, id)
}

// SetConnectorAudience installs the lookup that resolves a connector-scoped
// event to the users allowed to read it. Without one, connector events are
// dropped. Call before serving connections.
func (h *Hub) SetConnectorAudience(fn ConnectorAudience) { h.audience = fn }

// IssueTicket mints a one-time, short-lived ticket that authorizes a single
// WebSocket upgrade for the given identity. The caller must have authenticated
// the user with a normal access token.
func (h *Hub) IssueTicket(ident Identity) (string, error) {
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
	h.tickets[id] = ticket{id: ident, expires: now.Add(ticketTTL)}
	return id, nil
}

// RedeemTicket consumes a ticket, returning its identity. A ticket works once.
func (h *Hub) RedeemTicket(id string) (Identity, bool) {
	h.ticketMu.Lock()
	t, found := h.tickets[id]
	delete(h.tickets, id)
	h.ticketMu.Unlock()
	if !found || time.Now().After(t.expires) {
		return Identity{}, false
	}
	return t.id, true
}

// Run starts the hub's event loop. Should be run in a goroutine. Returns
// when ctx is canceled.
func (h *Hub) Run(ctx context.Context) {
	ticker := time.NewTicker(h.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			count := len(h.clients)
			h.mu.Unlock()
			slog.Info("WebSocket client connected", "user_id", logsafe.Sanitize(client.id.UserID), "total_clients", count)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			count := len(h.clients)
			h.mu.Unlock()
			slog.Info("WebSocket client disconnected", "user_id", logsafe.Sanitize(client.id.UserID), "total_clients", count)

		case msg := <-h.broadcast:
			h.deliver(msg)

		case <-ticker.C:
			heartbeat, _ := json.Marshal(newHeartbeat())
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
		if !msg.reaches(client) {
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
				slog.Warn("disconnecting slow WebSocket client", "user_id", logsafe.Sanitize(client.id.UserID))
			}
		}
	}
}

// Broadcast queues a global event for all clients. Use BroadcastConnector for
// anything tied to a connector.
func (h *Hub) Broadcast(eventType string, payload any) {
	h.publish(newEnvelope(eventType, payload), Audience{Kind: AudienceAll})
}

// BroadcastToUser queues a message for a user's connections.
func (h *Hub) BroadcastToUser(userID, eventType string, payload any) {
	h.publish(newEnvelope(eventType, payload), Audience{Kind: AudienceUser, UserID: userID})
}

// BroadcastConnector queues a connector-scoped event for the connections
// allowed to read connectorID. It resolves the readers in the caller's
// goroutine; the caller must not hold a database transaction or open rows.
func (h *Hub) BroadcastConnector(connectorID, eventType string, payload any) {
	env := newEnvelope(eventType, payload)
	env.ConnectorID = connectorID
	h.publish(env, Audience{Kind: AudienceConnector, ConnectorID: connectorID})
}

// publish resolves aud on this replica and queues env, dropping it if the hub
// queue is full. Connector audiences fail closed: with no lookup installed, no
// connector ID, or a failed lookup, the event is dropped. The ADR 0005 relay
// receiver delivers through here too, so every replica filters the same way.
func (h *Hub) publish(env Envelope, aud Audience) {
	msg := broadcastMsg{aud: aud}
	if aud.Kind == AudienceConnector {
		if h.ClientCount() == 0 {
			return
		}
		if aud.ConnectorID == "" || h.audience == nil {
			slog.Warn("dropping connector WS event: no connector or audience lookup", "type", env.Type)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), audienceTimeout)
		ids, err := h.audience(ctx, aud.ConnectorID)
		cancel()
		if err != nil {
			slog.Warn("dropping connector WS event: audience lookup failed",
				"connector_id", logsafe.Sanitize(aud.ConnectorID), "type", env.Type, "error", err)
			return
		}
		if len(ids) == 0 {
			return
		}
		msg.readers = make(map[string]struct{}, len(ids))
		for _, id := range ids {
			msg.readers[id] = struct{}{}
		}
	}
	data, err := json.Marshal(env)
	if err != nil {
		slog.Error("failed to marshal WS broadcast", "error", err)
		return
	}
	msg.data = data
	select {
	case h.broadcast <- msg:
	default:
		slog.Warn("WebSocket broadcast queue full, dropping message", "type", env.Type)
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
func (h *Hub) UpgradeHandler(w http.ResponseWriter, r *http.Request, id Identity) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
		id:   id,
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
			if !c.hub.Revalidate(context.Background(), c.id) {
				slog.Info("closing WebSocket: identity no longer valid", "user_id", logsafe.Sanitize(c.id.UserID))
				c.conn.WriteControl(websocket.CloseMessage, //nolint:errcheck
					websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "session no longer valid"), time.Now().Add(time.Second))
				return
			}
		}
	}
}
