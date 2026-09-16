// Package connectors provides connector management API handlers.
package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/logsafe"
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

// List handles GET /api/connectors.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	_, pageSize, offset := httputil.Paginate(r)
	category := r.URL.Query().Get("category")

	var connectors []store.ConnectorRecord
	var err error

	if category != "" {
		connectors, _, err = h.Store.ListConnectorsByCategory(r.Context(), category, offset, pageSize)
	} else {
		connectors, _, err = h.Store.ListConnectors(r.Context(), offset, pageSize)
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if connectors == nil {
		connectors = []store.ConnectorRecord{}
	}

	// Spec: GET /connectors returns a bare Connector[] (see openapi.yaml).
	httputil.JSON(w, http.StatusOK, connectors)
}

// Create handles POST /api/connectors.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		Name      string         `json:"name"`
		Category  string         `json:"category"`
		Type      string         `json:"type"`
		URL       string         `json:"url"`
		Owner     string         `json:"owner"`
		VerifyTLS *bool          `json:"verifyTls"`
		Config    map[string]any `json:"config"`
		Enabled   *bool          `json:"enabled"`
	}](w, r)
	if !ok {
		return
	}
	if req.Name == "" || req.Category == "" || req.Type == "" || req.URL == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "name, category, type, and url are required")
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
		Name:       req.Name,
		Category:   req.Category,
		Type:       req.Type,
		URL:        req.URL,
		Owner:      req.Owner,
		VerifyTLS:  verifyTLS,
		ConfigData: configData,
		Enabled:    enabled,
	}

	if err := h.Store.CreateConnector(r.Context(), c); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.create", "connector", c.ID, map[string]any{
		"name": c.Name, "category": c.Category, "type": c.Type,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.create", "error", err)
	}

	httputil.JSON(w, http.StatusCreated, c)
}

// Get handles GET /api/connectors/{id}.
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
	httputil.JSON(w, http.StatusOK, c)
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
		// ScheduleSeconds is raw JSON so absence (leave unchanged) can be told
		// apart from an explicit `null` (clear to manual-only): a **int decodes
		// both to a nil outer pointer, losing that distinction.
		ScheduleSeconds json.RawMessage `json:"scheduleSeconds"`
	}](w, r)
	if !ok {
		return
	}

	updates := make(map[string]any)
	if req.ScheduleSeconds != nil {
		var scheduleSeconds *int
		if err := json.Unmarshal(req.ScheduleSeconds, &scheduleSeconds); err != nil {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid scheduleSeconds")
			return
		}
		updates["schedule_seconds"] = scheduleSeconds
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
		data, err := store.MarshalConnectorConfig(typ, req.Config, h.Config.Encryption.Key)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		updates["config_data"] = data
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

	// Elevation validation will be done in Phase 6 when sync engine is built
	// For now, accept the token as a placeholder
	_ = elevationToken

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

// RestartPreview handles POST /api/connectors/{id}/restart. dryRun=true
// previews the restart from the latest stored snapshot without touching the
// connector; dryRun absent/false performs the real, elevation-gated restart.
func (h *Handler) RestartPreview(w http.ResponseWriter, r *http.Request) {
	dryRun := len(r.URL.Query()["dryRun"]) == 1 && r.URL.Query()["dryRun"][0] == "true"
	if !dryRun {
		h.restart(w, r)
		return
	}

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
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "No snapshot available for connector")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	var snap connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(sn.Data), &snap); err != nil {
		httputil.Errorf(w, fmt.Errorf("decode service snapshot: %w", err))
		return
	}

	dependencies := snap.Dependencies
	if dependencies == nil {
		dependencies = []connector.ServiceDependency{}
	}
	httputil.JSON(w, http.StatusOK, map[string]any{
		"targetService":            snap.ServiceName,
		"estimatedDowntimeSeconds": restartPreviewDowntimeSeconds,
		"dependentServices":        dependencies,
	})
}

// restart handles the real, mutating side of POST /api/connectors/{id}/restart
// (dryRun absent/false). Per ADR 0001: gated by the "connector.restart"
// elevation action, no rollback on failure (an AlertRecord is raised
// instead), and only successes are audited.
func (h *Handler) restart(w http.ResponseWriter, r *http.Request) {
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

	restarter, ok := conn.(connector.Restarter)
	if !ok {
		httputil.Error(w, http.StatusBadRequest, "unsupported_operation", "connector does not support restart")
		return
	}

	if err := auth.ValidateElevationHeader(h.JWT, h.Store, "connector.restart", r); err != nil {
		auth.WriteElevationError(w, err)
		return
	}

	var body struct {
		EntityRef string `json:"entityRef"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body) // ponytail: absent/empty body means entityRef == "", matches default
	}

	if err := restarter.Restart(r.Context(), cfg, body.EntityRef); err != nil {
		slog.Error("connector restart failed", "connector", id, "error", err)
		alert := &store.AlertRecord{
			ServiceID:   id,
			Severity:    "critical",
			Title:       fmt.Sprintf("Restart failed for %s", rec.Name),
			Description: err.Error(),
		}
		if createErr := h.Store.CreateAlert(r.Context(), alert); createErr != nil {
			slog.Error("failed to create restart failure alert", "error", createErr)
		} else if h.WSHub != nil {
			h.WSHub.Broadcast(ws.EventAlertCreated, map[string]any{
				"alertId":   alert.ID,
				"serviceId": id,
				"severity":  alert.Severity,
				"title":     alert.Title,
			})
		}
		httputil.Error(w, http.StatusBadGateway, "restart_failed", err.Error())
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.restart", "connector", id, map[string]any{
		"entityRef": body.EntityRef,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.restart", "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"status": "restarted"})
}

// validateConnectorConfig checks a connector config against its type's
// schema (SchemaField pattern/length/enum rules), catching malformed values
// at save time rather than on first fetch. Unknown types are left for
// connector.Get to reject.
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

// Sync handles POST /api/connectors/{id}/sync.
// Triggers a sync for a single connector. Returns 202 with job info; the sync
// itself runs asynchronously and its progress streams over /ws. An optional
// JSON body {"fields": ["vms","storage"]} requests a partial fetch — a
// dashboard quick-check doesn't need to force a full node/VM/container fetch.
func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := h.Store.GetConnector(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	var req struct {
		Fields []string `json:"fields"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	jobID := uuid.New().String()
	go func() {
		if _, err := h.SyncEngine.RunSyncFields(context.Background(), id, jobID, req.Fields); err != nil {
			slog.Error("sync failed", "connector", logsafe.Sanitize(id), "job", jobID, "error", logsafe.Sanitize(err.Error()))
		}
	}()

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.sync", "connector", id, map[string]any{
		"jobId": jobID,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.sync", "error", err)
	}

	httputil.JSON(w, http.StatusAccepted, map[string]any{
		"jobId":     jobID,
		"serviceId": id,
	})
}

// SyncAll handles POST /api/sync.
// Triggers a global sync of all enabled connectors. Returns 202 with job info;
// the sync itself runs asynchronously.
func (h *Handler) SyncAll(w http.ResponseWriter, r *http.Request) {
	jobID := uuid.New().String()
	go func() {
		if _, err := h.SyncEngine.RunSyncAll(context.Background(), jobID); err != nil {
			slog.Error("global sync failed", "job", jobID, "error", err)
		}
	}()

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.sync_all", "connector", "", map[string]any{
		"jobId": jobID,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.sync_all", "error", err)
	}

	httputil.JSON(w, http.StatusAccepted, map[string]any{
		"jobId":     jobID,
		"serviceId": nil,
	})
}
