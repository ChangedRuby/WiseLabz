// Package connectors provides connector management API handlers.
package connectors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// Handler holds dependencies for connector endpoints.
type Handler struct {
	Store      *store.Store
	SyncEngine *sync.Engine
	Config     *config.Config
	JWT        *auth.Service
	WSHub      *ws.Hub
}

// NewHandler creates a new connector handler.
func NewHandler(s *store.Store, e *sync.Engine, cfg *config.Config, jwtSvc *auth.Service, hub *ws.Hub) *Handler {
	return &Handler{Store: s, SyncEngine: e, Config: cfg, JWT: jwtSvc, WSHub: hub}
}

// List handles GET /api/connectors. Default deny: only connectors the
// caller holds at least a viewer grant on are returned. The grant filter is
// applied in the paginated query so pages and X-Total-Count are accurate.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	_, pageSize, offset := httputil.Paginate(r)
	category := r.URL.Query().Get("category")

	rows, total, err := h.Store.ListConnectorsForUser(r.Context(), auth.UserIDFromContext(r.Context()), category, offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	out := make([]connectorWithRole, 0, len(rows))
	for _, c := range rows {
		out = append(out, connectorWithRole{ConnectorRecord: c.ConnectorRecord, MyRole: c.Role})
	}

	// Spec: GET /connectors returns a bare Connector[] (see openapi.yaml).
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	httputil.JSON(w, http.StatusOK, out)
}

// connectorWithRole embeds a connector record with the requesting user's
// role on it, flattened into the JSON response via embedding.
type connectorWithRole struct {
	store.ConnectorRecord
	MyRole string `json:"myRole"`
}

// Create handles POST /api/connectors.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		Name               string         `json:"name"`
		Category           string         `json:"category"`
		Type               string         `json:"type"`
		URL                string         `json:"url"`
		Owner              string         `json:"owner"`
		VerifyTLS          *bool          `json:"verifyTls"`
		Config             map[string]any `json:"config"`
		Enabled            *bool          `json:"enabled"`
		UserExpiresAt      string         `json:"userExpiresAt"`
		RotationMaxAgeDays *int           `json:"rotationMaxAgeDays"`
	}](w, r)
	if !ok {
		return
	}
	if req.Name == "" || req.Category == "" || req.Type == "" || req.URL == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "name, category, type, and url are required")
		return
	}
	if err := validateRotationFields(req.UserExpiresAt, req.RotationMaxAgeDays); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	verifyTLS := true
	if req.VerifyTLS != nil {
		verifyTLS = *req.VerifyTLS
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	configData := "{}"
	if req.Config != nil {
		data, err := store.MarshalConnectorConfig(req.Type, req.Config, h.Config.Encryption.Key)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		configData = data
	}

	if err := validateConnectorConfig(req.Type, req.URL, verifyTLS, req.Config); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	c := &store.ConnectorRecord{
		Name:               req.Name,
		Category:           req.Category,
		Type:               req.Type,
		URL:                req.URL,
		Owner:              req.Owner,
		VerifyTLS:          verifyTLS,
		ConfigData:         configData,
		Enabled:            enabled,
		UserExpiresAt:      req.UserExpiresAt,
		RotationMaxAgeDays: req.RotationMaxAgeDays,
	}

	if err := h.Store.CreateConnector(r.Context(), c); err != nil {
		httputil.Errorf(w, err)
		return
	}

	// ponytail: auto-grant the creator operator access so an instance admin
	// isn't locked out of the connector they just made (grants are the only
	// source of connector access now — even for admins).
	if _, err := h.Store.UpsertConnectorGrant(r.Context(), auth.UserIDFromContext(r.Context()), c.ID, "operator"); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.create", "connector", c.ID, map[string]any{
		"name": c.Name, "category": c.Category, "type": c.Type,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.create", "error", err)
	}

	httputil.JSON(w, http.StatusCreated, connectorWithRole{ConnectorRecord: *c, MyRole: "operator"})
}

// Get handles GET /api/connectors/{id}. Default deny: 404s (not 403, to
// avoid confirming the connector's existence) if the caller has no grant.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := h.Store.GetConnector(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	role, err := h.Store.GetUserConnectorRole(r.Context(), auth.UserIDFromContext(r.Context()), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if role == "" {
		httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
		return
	}
	httputil.JSON(w, http.StatusOK, connectorWithRole{ConnectorRecord: *c, MyRole: role})
}

// Update handles PUT or PATCH /api/connectors/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	req, ok := httputil.DecodeJSON[struct {
		Name      *string        `json:"name"`
		URL       *string        `json:"url"`
		Owner     *string        `json:"owner"`
		VerifyTLS *bool          `json:"verifyTls"`
		Config    map[string]any `json:"config"`
		Enabled   *bool          `json:"enabled"`
		Category  *string        `json:"category"`
		Type      *string        `json:"type"`
		// ScheduleSeconds, UserExpiresAt and RotationMaxAgeDays are raw JSON so
		// absence (leave unchanged) can be told apart from an explicit `null`
		// (clear): a plain *T/**T decodes both to a nil pointer, losing that
		// distinction.
		ScheduleSeconds    json.RawMessage `json:"scheduleSeconds"`
		UserExpiresAt      json.RawMessage `json:"userExpiresAt"`
		RotationMaxAgeDays json.RawMessage `json:"rotationMaxAgeDays"`
	}](w, r)
	if !ok {
		return
	}

	updates, errMsg := parseScheduleUpdates(req.ScheduleSeconds, req.UserExpiresAt, req.RotationMaxAgeDays)
	if errMsg != "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", errMsg)
		return
	}
	// Repointing a connector makes the server send its stored credentials to
	// the new endpoint, so changing where or how it connects is an
	// instance-admin action; operators may still send unchanged values.
	if (req.URL != nil || req.Type != nil || req.VerifyTLS != nil) && !auth.InstanceAdminFromContext(r.Context()) {
		current, err := h.Store.GetConnector(r.Context(), id)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
				return
			}
			httputil.Errorf(w, err)
			return
		}
		if (req.URL != nil && *req.URL != current.URL) ||
			(req.Type != nil && *req.Type != current.Type) ||
			(req.VerifyTLS != nil && *req.VerifyTLS != current.VerifyTLS) {
			httputil.Error(w, http.StatusForbidden, "forbidden", "Changing url, type or verifyTls requires an instance admin")
			return
		}
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.Owner != nil {
		updates["owner"] = *req.Owner
	}
	if req.VerifyTLS != nil {
		updates["verify_tls"] = *req.VerifyTLS
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}

	if req.Config != nil {
		rec, err := h.Store.GetConnector(r.Context(), id)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
				return
			}
			httputil.Errorf(w, err)
			return
		}
		typ := rec.Type
		if req.Type != nil {
			typ = *req.Type
		}
		url := rec.URL
		if req.URL != nil {
			url = *req.URL
		}
		verifyTLS := rec.VerifyTLS
		if req.VerifyTLS != nil {
			verifyTLS = *req.VerifyTLS
		}
		if err := validateConnectorConfig(typ, url, verifyTLS, req.Config); err != nil {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		changed, err := store.SecretFieldsChanged(typ, rec.ConfigData, req.Config, h.Config.Encryption.Key)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		data, err := store.MarshalConnectorConfig(typ, req.Config, h.Config.Encryption.Key)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		updates["config_data"] = data
		if changed {
			updates["secret_rotated_at"] = time.Now().UTC().Format(time.RFC3339)
		}
	}

	if len(updates) == 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "No fields to update")
		return
	}

	if err := h.Store.UpdateConnector(r.Context(), id, updates); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	// Record only which fields changed, not their values — config_data can
	// carry connector credentials and must never land in the audit log.
	fields := make([]string, 0, len(updates))
	for k := range updates {
		fields = append(fields, k)
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.update", "connector", id, map[string]any{
		"fields": fields,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.update", "error", err)
	}

	c, _ := h.Store.GetConnector(r.Context(), id)
	httputil.JSON(w, http.StatusOK, c)
}

// Delete handles DELETE /api/connectors/{id}.
// If auth config requires step-up, X-Elevation-Token must be validated.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Check for elevation token (step-up auth)
	elevationToken := r.Header.Get("X-Elevation-Token")
	if elevationToken == "" {
		httputil.Error(w, http.StatusBadRequest, "elevation_required", "X-Elevation-Token header required for destructive action")
		return
	}

	if err := h.Store.DeleteConnector(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.delete", "connector", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "connector.delete", "error", err)
	}

	httputil.NoContent(w)
}

// Test handles POST /api/connectors/{id}/test.
// Validates the connector's saved configuration by attempting a connection.
func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rec, err := h.Store.GetConnector(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	cfg, err := store.ParseConnectorConfig(rec.Type, rec.ConfigData, h.Config.Encryption.Key)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	cfg["url"] = rec.URL
	cfg["verify_tls"] = rec.VerifyTLS

	c, err := connector.Get(rec.Type, cfg)
	if err != nil {
		httputil.JSON(w, http.StatusOK, map[string]any{
			"ok":      false,
			"message": fmt.Sprintf("Failed to create connector: %v", err),
		})
		return
	}

	start := time.Now()
	validateErr := c.Validate(r.Context(), cfg)
	latencyMs := int(time.Since(start).Milliseconds())

	if validateErr != nil {
		httputil.JSON(w, http.StatusOK, map[string]any{
			"ok":        false,
			"message":   validateErr.Error(),
			"latencyMs": latencyMs,
		})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"message":   "Connection successful",
		"latencyMs": latencyMs,
	})
}

// Health handles POST /api/connectors/{id}/health.
// Runs a cheap connectivity check (Validate only — no Fetch, no snapshot, no
// docs) and persists the resulting status, so availability can be tracked
// independently of, and between, full sync runs.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rec, err := h.Store.GetConnector(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	cfg, err := store.ParseConnectorConfig(rec.Type, rec.ConfigData, h.Config.Encryption.Key)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	cfg["url"] = rec.URL
	cfg["verify_tls"] = rec.VerifyTLS

	start := time.Now()
	c, connErr := connector.Get(rec.Type, cfg)
	validateErr := connErr
	if connErr == nil {
		validateErr = c.Validate(r.Context(), cfg)
	}
	latency := time.Since(start)

	threshold := connector.DegradedLatencyThreshold
	if schema, err := connector.GetTypeSchema(rec.Type); err == nil {
		threshold = schema.DegradedLatencyThreshold()
	}
	status, message := connector.ClassifyHealth(validateErr, latency, threshold)
	if err := h.Store.UpdateConnector(r.Context(), id, map[string]any{
		"status":         status,
		"status_message": message,
	}); err != nil {
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"status":    status,
		"message":   message,
		"latencyMs": int(latency.Milliseconds()),
	})
}

// Data handles GET /api/connectors/{id}/data.
// Returns the latest fetched service snapshot for the connector.
func (h *Handler) Data(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := h.Store.GetConnector(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	sn, err := h.Store.GetLatestSnapshot(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// No snapshot yet — return empty data, not an error
			httputil.JSON(w, http.StatusOK, map[string]any{
				"serviceName": "",
				"type":        "",
				"sections":    []map[string]any{},
				"metadata":    map[string]string{},
				"fetchedAt":   "",
			})
			return
		}
		httputil.Errorf(w, err)
		return
	}

	var snap connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(sn.Data), &snap); err != nil {
		httputil.Errorf(w, err)
		return
	}

	sections := make([]map[string]any, len(snap.Sections))
	for i, sec := range snap.Sections {
		sections[i] = map[string]any{
			"title":   sec.Title,
			"content": sec.Content,
			"order":   i,
		}
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"serviceName": snap.ServiceName,
		"type":        snap.Type,
		"sections":    sections,
		"metadata":    snap.Metadata,
		"fetchedAt":   snap.FetchedAt.Format(time.RFC3339),
	})
}

// syncsDefaultLimit and syncsMaxLimit mirror httputil.DefaultPageSize/MaxPageSize,
// per the /connectors/{connectorId}/syncs OpenAPI spec (limit default 20, max 100).
const (
	syncsDefaultLimit             = httputil.DefaultPageSize
	syncsMaxLimit                 = httputil.MaxPageSize
	restartPreviewDowntimeSeconds = 30
)

// Syncs handles GET /api/connectors/{id}/syncs (OpenAPI's connectorId path param).
// Returns recent sync run history for a connector, newest first.
func (h *Handler) Syncs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := h.Store.GetConnector(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	limit := syncsDefaultLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > syncsMaxLimit {
		limit = syncsMaxLimit
	}

	runs, err := h.Store.ListSyncRunsByConnector(r.Context(), id, limit)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, runs)
}

// RemovalImpact handles GET /api/connectors/{id}/removal-impact.
// Returns counts of dependent resources before deletion.
func (h *Handler) RemovalImpact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := h.Store.GetConnector(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	snapshots, _ := h.Store.CountSnapshotsByConnector(r.Context(), id)
	docSections, _ := h.Store.CountDocsByConnector(r.Context(), id)

	docs, err := h.Store.ListDocsByService(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	items := make([]map[string]string, 0, len(docs))
	for _, d := range docs {
		items = append(items, map[string]string{"type": "doc", "name": d.Title})
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"trackedServices": 1,
		"docSections":     docSections,
		"snapshots":       snapshots,
		"items":           items,
	})
}

// ConfigFields handles GET /api/connectors/{id}/config-fields. Returns the
// connector's whitelisted config-push fields so the frontend can build a
// field picker without calling a Go interface method directly.
func (h *Handler) ConfigFields(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rec, err := h.Store.GetConnector(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	cfg, err := store.ParseConnectorConfig(rec.Type, rec.ConfigData, h.Config.Encryption.Key)
	if err != nil {
		httputil.Errorf(w, fmt.Errorf("parse config: %w", err))
		return
	}
	cfg["url"] = rec.URL
	cfg["verify_tls"] = rec.VerifyTLS

	conn, err := connector.Get(rec.Type, cfg)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	pusher, ok := conn.(connector.ConfigPusher)
	if !ok {
		httputil.JSON(w, http.StatusOK, []connector.ConfigField{})
		return
	}
	httputil.JSON(w, http.StatusOK, pusher.WritableFields())
}

// parseScheduleUpdates decodes the raw-JSON scheduleSeconds, userExpiresAt and
// rotationMaxAgeDays fields of an Update request into store column updates and
// validates the rotation fields. An absent (nil) field is left out of the
// result; an explicit JSON null is included as a clearing value. The returned message is client-safe for a 400 response and
// empty on success.
func parseScheduleUpdates(scheduleSeconds, userExpiresAtRaw, rotationMaxAgeDaysRaw json.RawMessage) (updates map[string]any, errMsg string) {
	updates = make(map[string]any)
	if scheduleSeconds != nil {
		var v *int
		if err := json.Unmarshal(scheduleSeconds, &v); err != nil {
			return nil, "Invalid scheduleSeconds"
		}
		updates["schedule_seconds"] = v
	}
	var userExpiresAt string
	if userExpiresAtRaw != nil {
		var v *string
		if err := json.Unmarshal(userExpiresAtRaw, &v); err != nil {
			return nil, "Invalid userExpiresAt"
		}
		if v != nil {
			userExpiresAt = *v
		}
		if userExpiresAt == "" {
			updates["user_expires_at"] = nil
		} else {
			updates["user_expires_at"] = userExpiresAt
		}
	}
	var rotationMaxAgeDays *int
	if rotationMaxAgeDaysRaw != nil {
		if err := json.Unmarshal(rotationMaxAgeDaysRaw, &rotationMaxAgeDays); err != nil {
			return nil, "Invalid rotationMaxAgeDays"
		}
		updates["rotation_max_age_days"] = rotationMaxAgeDays
	}
	if err := validateRotationFields(userExpiresAt, rotationMaxAgeDays); err != nil {
		return nil, err.Error()
	}
	return updates, ""
}

// validateConnectorConfig checks a connector config against its type's
// schema (SchemaField pattern/length/enum rules), catching malformed values
// at save time rather than on first fetch. Unknown types are left for
// connector.Get to reject.
// validateRotationFields checks the optional user-set credential rotation
// overrides: userExpiresAt (if non-empty) must be an RFC3339 timestamp,
// rotationMaxAgeDays (if set) must be a positive number of days.
func validateRotationFields(userExpiresAt string, rotationMaxAgeDays *int) error {
	if userExpiresAt != "" {
		if _, err := time.Parse(time.RFC3339, userExpiresAt); err != nil {
			return fmt.Errorf("userExpiresAt must be an RFC3339 timestamp")
		}
	}
	if rotationMaxAgeDays != nil && *rotationMaxAgeDays <= 0 {
		return fmt.Errorf("rotationMaxAgeDays must be a positive number of days")
	}
	return nil
}

func validateConnectorConfig(typ, url string, verifyTLS bool, config map[string]any) error {
	schema, err := connector.GetTypeSchema(typ)
	if err != nil {
		return nil
	}
	cfg := make(map[string]any, len(config)+2)
	for k, v := range config {
		cfg[k] = v
	}
	cfg["url"] = url
	cfg["verify_tls"] = verifyTLS
	return connector.ValidateConfig(*schema, cfg)
}

// Schema handles GET /api/connectors/schema.
func (h *Handler) Schema(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, connector.ListSchemas())
}

// ToggleEnabled handles PUT /api/connectors/{id}/enabled.
func (h *Handler) ToggleEnabled(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req, ok := httputil.DecodeJSON[struct {
		Enabled bool `json:"enabled"`
	}](w, r)
	if !ok {
		return
	}

	if err := h.Store.UpdateConnector(r.Context(), id, map[string]any{
		"enabled": req.Enabled,
	}); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.toggle_enabled", "connector", id, map[string]any{
		"enabled": req.Enabled,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.toggle_enabled", "error", err)
	}

	c, _ := h.Store.GetConnector(r.Context(), id)
	httputil.JSON(w, http.StatusOK, c)
}
