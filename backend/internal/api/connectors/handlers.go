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
	"strings"
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

	updates := make(map[string]any)
	if req.ScheduleSeconds != nil {
		var scheduleSeconds *int
		if err := json.Unmarshal(req.ScheduleSeconds, &scheduleSeconds); err != nil {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid scheduleSeconds")
			return
		}
		updates["schedule_seconds"] = scheduleSeconds
	}
	var userExpiresAt string
	if req.UserExpiresAt != nil {
		var v *string
		if err := json.Unmarshal(req.UserExpiresAt, &v); err != nil {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid userExpiresAt")
			return
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
	if req.RotationMaxAgeDays != nil {
		if err := json.Unmarshal(req.RotationMaxAgeDays, &rotationMaxAgeDays); err != nil {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid rotationMaxAgeDays")
			return
		}
		updates["rotation_max_age_days"] = rotationMaxAgeDays
	}
	if err := validateRotationFields(userExpiresAt, rotationMaxAgeDays); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
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

// RestartPreview handles POST /api/connectors/{id}/restart. dryRun=true
// previews the restart from the latest stored snapshot without touching the
// connector; dryRun absent/false performs the real, elevation-gated restart.
func (h *Handler) RestartPreview(w http.ResponseWriter, r *http.Request) {
	h.mutatingOpPreview(w, r, mutatingOp{
		verb:          "restart",
		auditAction:   "connector.restart",
		elevateAction: "connector.restart",
		call: func(conn connector.Connector) (func(ctx context.Context, cfg map[string]any, entityRef string) error, bool) {
			restarter, ok := conn.(connector.Restarter)
			if !ok {
				return nil, false
			}
			return restarter.Restart, true
		},
	})
}

// StartPreview handles POST /api/connectors/{id}/start. Same dry-run/mutate
// split as RestartPreview.
func (h *Handler) StartPreview(w http.ResponseWriter, r *http.Request) {
	h.mutatingOpPreview(w, r, mutatingOp{
		verb:          "start",
		auditAction:   "connector.start",
		elevateAction: "connector.start",
		call: func(conn connector.Connector) (func(ctx context.Context, cfg map[string]any, entityRef string) error, bool) {
			starter, ok := conn.(connector.Starter)
			if !ok {
				return nil, false
			}
			return starter.Start, true
		},
	})
}

// StopPreview handles POST /api/connectors/{id}/stop. Same dry-run/mutate
// split as RestartPreview. estimatedDowntimeSeconds is 0 in the preview
// since a stop's downtime is indefinite until an explicit start, not a
// bounded window.
func (h *Handler) StopPreview(w http.ResponseWriter, r *http.Request) {
	h.mutatingOpPreview(w, r, mutatingOp{
		verb:          "stop",
		auditAction:   "connector.stop",
		elevateAction: "connector.stop",
		call: func(conn connector.Connector) (func(ctx context.Context, cfg map[string]any, entityRef string) error, bool) {
			stopper, ok := conn.(connector.Stopper)
			if !ok {
				return nil, false
			}
			return stopper.Stop, true
		},
	})
}

// mutatingOp describes one lab-mutating verb (restart/start/stop) sharing
// the ADR 0001/0002 dry-run-preview + elevation-gated-mutate shape.
type mutatingOp struct {
	verb          string // "restart", "start", "stop" — used in error/audit text
	auditAction   string
	elevateAction string
	// call type-asserts conn against the verb's interface (Restarter/
	// Starter/Stopper) and returns its method, or ok=false if unsupported.
	call func(conn connector.Connector) (fn func(ctx context.Context, cfg map[string]any, entityRef string) error, ok bool)
}

// mutatingOpPreview implements the shared dry-run/mutate split for restart,
// start, and stop: dryRun=true returns a preview built from the latest
// stored snapshot without touching the connector; dryRun absent/false
// performs the real, elevation-gated mutation.
func (h *Handler) mutatingOpPreview(w http.ResponseWriter, r *http.Request, op mutatingOp) {
	dryRun := len(r.URL.Query()["dryRun"]) == 1 && r.URL.Query()["dryRun"][0] == "true"
	if !dryRun {
		h.mutateOp(w, r, op)
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
	downtime := restartPreviewDowntimeSeconds
	if op.verb == "stop" {
		// A stop's downtime is indefinite (no scheduled restart), not a
		// bounded estimate — 0 rather than inventing a new field.
		downtime = 0
	}
	httputil.JSON(w, http.StatusOK, map[string]any{
		"targetService":            snap.ServiceName,
		"estimatedDowntimeSeconds": downtime,
		"dependentServices":        dependencies,
	})
}

// mutateOp handles the real, mutating side of restart/start/stop (dryRun
// absent/false). Per ADR 0001/0002: gated by the verb's elevation action,
// no rollback on failure (an AlertRecord is raised instead), and only
// successes are audited.
func (h *Handler) mutateOp(w http.ResponseWriter, r *http.Request, op mutatingOp) {
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

	fn, ok := op.call(conn)
	if !ok {
		httputil.Error(w, http.StatusBadRequest, "unsupported_operation", "connector does not support "+op.verb)
		return
	}

	if err := auth.ValidateElevationHeader(h.JWT, h.Store, op.elevateAction, r); err != nil {
		auth.WriteElevationError(w, err)
		return
	}

	var body struct {
		EntityRef string `json:"entityRef"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body) // ponytail: absent/empty body means entityRef == "", matches default
	}

	if err := connector.ValidateCompositeRef(body.EntityRef); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "invalid entityRef")
		return
	}

	if err := fn(r.Context(), cfg, body.EntityRef); err != nil {
		slog.Error("connector "+op.verb+" failed", "connector", id, "error", err)
		alert := &store.AlertRecord{
			ServiceID:   id,
			Severity:    "critical",
			Title:       fmt.Sprintf("%s failed for %s", capitalize(op.verb), rec.Name),
			Description: err.Error(),
		}
		if createErr := h.Store.CreateAlert(r.Context(), alert); createErr != nil {
			slog.Error("failed to create "+op.verb+" failure alert", "error", createErr)
		} else if h.WSHub != nil {
			h.WSHub.Broadcast(ws.EventAlertCreated, map[string]any{
				"alertId":   alert.ID,
				"serviceId": id,
				"severity":  alert.Severity,
				"title":     alert.Title,
			})
		}
		httputil.Error(w, http.StatusBadGateway, op.verb+"_failed", err.Error())
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), op.auditAction, "connector", id, map[string]any{
		"entityRef": body.EntityRef,
	}); err != nil {
		slog.Error("failed to record audit", "action", op.auditAction, "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"status": op.verb + "ed"})
}

// capitalize upper-cases a word's first byte (ASCII verbs only: "restart",
// "start", "stop") for a display-friendly alert title.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
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

// ConfigPush handles POST /api/connectors/{id}/config-push. Per ADR 0003:
// field-level partial update against a per-connector whitelist, snapshot
// before write, verify-diff after, auto-revert-then-alert on mismatch.
func (h *Handler) ConfigPush(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	req, ok := httputil.DecodeJSON[struct {
		EntityRef string `json:"entityRef"`
		FieldKey  string `json:"fieldKey"`
		Value     any    `json:"value"`
		// PreviousValue is the field's currently-displayed value, as the
		// frontend read it from GET /{id}/data before the user edited it.
		// The handler has no generic way to re-derive "the old value of
		// fieldKey" from a rendered ServiceSnapshot, so the revert path
		// (ADR 0003) re-pushes this rather than something inferred.
		PreviousValue any `json:"previousValue"`
	}](w, r)
	if !ok {
		return
	}
	if err := connector.ValidateCompositeRef(req.EntityRef); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "invalid entityRef")
		return
	}
	if req.FieldKey == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "fieldKey is required")
		return
	}

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
		httputil.Error(w, http.StatusBadRequest, "unsupported_operation", "connector does not support config-push")
		return
	}
	fieldOK := false
	for _, f := range pusher.WritableFields() {
		if f.Key == req.FieldKey {
			fieldOK = true
			break
		}
	}
	if !fieldOK {
		httputil.Error(w, http.StatusBadRequest, "unsupported_field", fmt.Sprintf("field %q is not writable for this connector", req.FieldKey))
		return
	}

	if err := auth.ValidateElevationHeader(h.JWT, h.Store, "connector.configPush", r); err != nil {
		auth.WriteElevationError(w, err)
		return
	}

	pre, err := conn.Fetch(r.Context(), cfg)
	if err != nil {
		httputil.Errorf(w, fmt.Errorf("pre-push fetch: %w", err))
		return
	}

	if err := pusher.ConfigPush(r.Context(), cfg, req.EntityRef, req.FieldKey, req.Value); err != nil {
		slog.Error("connector config-push failed", "connector", id, "field", req.FieldKey, "error", err)
		httputil.Error(w, http.StatusBadGateway, "config_push_failed", err.Error())
		return
	}

	post, err := conn.Fetch(r.Context(), cfg)
	if err != nil {
		httputil.Errorf(w, fmt.Errorf("post-push fetch: %w", err))
		return
	}

	if !configPushLanded(pre, post) {
		h.revertConfigPush(w, r, pusher, cfg, rec, req.EntityRef, req.FieldKey, req.PreviousValue)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.configPush", "connector", id, map[string]any{
		"fieldKey":  req.FieldKey,
		"entityRef": req.EntityRef,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.configPush", "error", err)
	}

	httputil.JSON(w, http.StatusOK, post)
}

// configPushLanded reports whether pre→post shows any change at all.
// sync.Compare's diff results, non-empty, are treated as "the write landed
// as some change" — the field-scoped check ADR 0003 calls for; a
// per-field-key section filter would need a fixed field→section mapping
// per connector, which the curated whitelist doesn't carry today.
// ponytail: any detected change counts as "landed" rather than verifying
// the new value matches exactly — tighten if a connector's diff proves too
// coarse to catch a landed-but-wrong-value push.
func configPushLanded(pre, post *connector.ServiceSnapshot) bool {
	return len(sync.Compare(pre, post)) > 0
}

// revertConfigPush handles ADR 0003's mismatch path: re-push the field's
// pre-push value (the revert), then raise a critical alert describing the
// mismatch — escalating further if the revert call itself fails, since
// that's the one hard stop in this feature set requiring manual
// intervention. No further automation is attempted either way.
// ponytail: reverts exactly one field per call by design (matches
// ConfigPush's own one-field-per-call contract), not a multi-field batch.
func (h *Handler) revertConfigPush(w http.ResponseWriter, r *http.Request, pusher connector.ConfigPusher, cfg map[string]any, rec *store.ConnectorRecord, entityRef, fieldKey string, previousValue any) {
	id := rec.ID
	revertErr := pusher.ConfigPush(r.Context(), cfg, entityRef, fieldKey, previousValue)

	desc := fmt.Sprintf("Pushed field %q did not verify after write; auto-revert to the pre-push value was attempted.", fieldKey)
	if revertErr != nil {
		desc = fmt.Sprintf("Pushed field %q did not verify after write; auto-revert FAILED (%v) — manual intervention required.", fieldKey, revertErr)
		slog.Error("config-push auto-revert failed", "connector", id, "field", fieldKey, "error", revertErr)
	} else {
		slog.Warn("config-push mismatch, auto-reverted", "connector", id, "field", fieldKey)
	}

	alert := &store.AlertRecord{
		ServiceID:   id,
		Severity:    "critical",
		Title:       fmt.Sprintf("Config push mismatch for %s", rec.Name),
		Description: desc,
	}
	if createErr := h.Store.CreateAlert(r.Context(), alert); createErr != nil {
		slog.Error("failed to create config-push mismatch alert", "error", createErr)
	} else if h.WSHub != nil {
		h.WSHub.Broadcast(ws.EventAlertCreated, map[string]any{
			"alertId":   alert.ID,
			"serviceId": id,
			"severity":  alert.Severity,
			"title":     alert.Title,
		})
	}

	httputil.Error(w, http.StatusConflict, "config_push_mismatch", "pushed value did not verify; auto-revert attempted, see the new alert")
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
		if _, err := h.SyncEngine.RunSyncFields(h.SyncEngine.BaseContext(), id, jobID, req.Fields); err != nil {
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
		if _, err := h.SyncEngine.RunSyncAll(h.SyncEngine.BaseContext(), jobID); err != nil {
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

// OpenMaintenanceWindow handles POST /api/connectors/{id}/maintenance-window.
// Opens a time-boxed maintenance window on a connector: while active, the
// scheduler skips it and drift alerting/change-record creation from any sync
// (scheduled or manual) is suppressed, though snapshots keep being saved.
// Operator-only; unlike the other mutating connector actions, this one is
// reversible and time-boxed so it doesn't require step-up elevation.
func (h *Handler) OpenMaintenanceWindow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := h.Store.GetConnector(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	req, ok := httputil.DecodeJSON[struct {
		DurationMinutes int `json:"durationMinutes"`
	}](w, r)
	if !ok {
		return
	}
	if req.DurationMinutes <= 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "durationMinutes must be positive")
		return
	}

	now := time.Now().UTC()
	m := &store.MaintenanceWindowRecord{
		ConnectorID: id,
		StartsAt:    now.Format(time.RFC3339),
		EndsAt:      now.Add(time.Duration(req.DurationMinutes) * time.Minute).Format(time.RFC3339),
		CreatedBy:   auth.UserIDFromContext(r.Context()),
	}
	if err := h.Store.CreateMaintenanceWindow(r.Context(), m); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.maintenanceWindow.open", "connector", id, map[string]any{
		"durationMinutes": req.DurationMinutes, "endsAt": m.EndsAt,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.maintenanceWindow.open", "error", err)
	}

	httputil.JSON(w, http.StatusCreated, m)
}

// CloseMaintenanceWindow handles DELETE /api/connectors/{id}/maintenance-window.
// Ends the connector's active maintenance window early. A no-op (200) when
// there is none.
func (h *Handler) CloseMaintenanceWindow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	active, err := h.Store.GetActiveMaintenanceWindow(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if active == nil {
		httputil.JSON(w, http.StatusOK, map[string]any{"closed": false})
		return
	}

	if err := h.Store.CloseMaintenanceWindow(r.Context(), active.ID); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.maintenanceWindow.close", "connector", id, map[string]any{
		"windowId": active.ID,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.maintenanceWindow.close", "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"closed": true})
}

// GetMaintenanceWindow handles GET /api/connectors/{id}/maintenance-window.
// Returns the connector's active maintenance window, or null if none.
func (h *Handler) GetMaintenanceWindow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	m, err := h.Store.GetActiveMaintenanceWindow(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, m)
}

// ListActiveMaintenance handles GET /api/connectors/maintenance-windows.
// Returns every currently active maintenance window, across all connectors —
// used by the ServicesPage badge to avoid an N+1 fetch per row.
func (h *Handler) ListActiveMaintenance(w http.ResponseWriter, r *http.Request) {
	windows, err := h.Store.ListActiveMaintenanceWindows(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	ids := make([]string, len(windows))
	for i, m := range windows {
		ids[i] = m.ConnectorID
	}
	allowed, err := h.Store.FilterConnectorIDsByGrant(r.Context(), auth.UserIDFromContext(r.Context()), ids, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	isAllowed := make(map[string]bool, len(allowed))
	for _, id := range allowed {
		isAllowed[id] = true
	}
	filtered := make([]store.MaintenanceWindowRecord, 0, len(windows))
	for _, m := range windows {
		if isAllowed[m.ConnectorID] {
			filtered = append(filtered, m)
		}
	}
	httputil.JSON(w, http.StatusOK, filtered)
}

// bulkRequest is the shared body shape for bulk-sync/bulk-reauth/bulk-restart:
// an explicit, caller-supplied list of connector IDs.
type bulkRequest struct {
	IDs []string `json:"ids"`
}

// bulkItemResult is the per-item outcome in a bulk connector-action response.
type bulkItemResult struct {
	ID     string `json:"id"`
	Status string `json:"status"` // "success" | "error"
	JobID  string `json:"jobId,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// decodeBulkRequest validates the common ids shape shared by all three bulk
// connector endpoints. Returns ok=false after writing the error response.
func decodeBulkRequest(w http.ResponseWriter, r *http.Request) (bulkRequest, bool) {
	req, ok := httputil.DecodeJSON[bulkRequest](w, r)
	if !ok {
		return req, false
	}
	if len(req.IDs) == 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "ids must be a non-empty array")
		return req, false
	}
	if len(req.IDs) > httputil.MaxBulkIDs {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "ids must contain at most 500 items")
		return req, false
	}
	return req, true
}

// splitByConnectorGrant looks up which of the requested ids exist, then
// partitions the existing ones into ones the caller has at least operator on
// and the rest. Nonexistent ids are surfaced as "not_found" and unauthorized
// existing ids as "forbidden", both pre-rendered as bulkItemResults so
// callers can append them straight into their results slice alongside the
// outcomes for the ids they go on to process.
func (h *Handler) splitByConnectorGrant(ctx context.Context, ids []string) (allowed []string, found map[string]store.ConnectorRecord, results []bulkItemResult, err error) {
	found, err = h.Store.ListConnectorsByID(ctx, ids)
	if err != nil {
		return nil, nil, nil, err
	}
	existing := make([]string, 0, len(found))
	for _, id := range ids {
		if _, ok := found[id]; ok {
			existing = append(existing, id)
		} else {
			results = append(results, bulkItemResult{ID: id, Status: "error", Reason: "not_found"})
		}
	}
	allowed, err = h.Store.FilterConnectorIDsByGrant(ctx, auth.UserIDFromContext(ctx), existing, "operator")
	if err != nil {
		return nil, nil, nil, err
	}
	isAllowed := make(map[string]bool, len(allowed))
	for _, id := range allowed {
		isAllowed[id] = true
	}
	for _, id := range existing {
		if !isAllowed[id] {
			results = append(results, bulkItemResult{ID: id, Status: "error", Reason: "forbidden"})
		}
	}
	return allowed, found, results, nil
}

// BulkSync handles POST /api/connectors/bulk-sync. Fans out an async
// RunSyncFields per connector (same as the single-connector Sync handler),
// auditing and reporting per-item results. One bad ID never aborts the
// batch; one audit record is written per resolved item.
func (h *Handler) BulkSync(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBulkRequest(w, r)
	if !ok {
		return
	}

	allowedIDs, _, results, err := h.splitByConnectorGrant(r.Context(), req.IDs)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	auditRecords := make([]store.AuditRecord, 0, len(allowedIDs))
	for _, id := range allowedIDs {
		jobID := uuid.New().String()
		go func(connectorID, jobID string) {
			if _, err := h.SyncEngine.RunSyncFields(h.SyncEngine.BaseContext(), connectorID, jobID, nil); err != nil {
				slog.Error("bulk sync failed", "connector", logsafe.Sanitize(connectorID), "job", jobID, "error", logsafe.Sanitize(err.Error()))
			}
		}(id, jobID)
		auditRecords = append(auditRecords, store.AuditRecord{TargetID: id, Detail: fmt.Sprintf(`{"jobId":%q}`, jobID)})
		results = append(results, bulkItemResult{ID: id, Status: "success", JobID: jobID})
	}

	if err := h.Store.RecordAuditBatchFromContext(r.Context(), "connector.bulk_sync", "connector", auditRecords); err != nil {
		slog.Error("failed to record audit", "action", "connector.bulk_sync", "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"results": results})
}

// BulkReauth handles POST /api/connectors/bulk-reauth. For each connector
// implementing connector.CredentialRefresher, refreshes and persists its
// credentials via sync.Engine.RefreshCredentials. Not elevation-gated
// (lower-risk, matches how single sync/re-auth aren't gated today either).
func (h *Handler) BulkReauth(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBulkRequest(w, r)
	if !ok {
		return
	}

	allowedIDs, _, results, err := h.splitByConnectorGrant(r.Context(), req.IDs)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	auditRecords := make([]store.AuditRecord, 0, len(allowedIDs))
	for _, id := range allowedIDs {
		if err := h.SyncEngine.RefreshCredentials(r.Context(), id); err != nil {
			results = append(results, bulkItemResult{ID: id, Status: "error", Reason: err.Error()})
			continue
		}
		auditRecords = append(auditRecords, store.AuditRecord{TargetID: id})
		results = append(results, bulkItemResult{ID: id, Status: "success"})
	}

	if err := h.Store.RecordAuditBatchFromContext(r.Context(), "connector.bulk_reauth", "connector", auditRecords); err != nil {
		slog.Error("failed to record audit", "action", "connector.bulk_reauth", "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"results": results})
}

// BulkRestart handles POST /api/connectors/bulk-restart. The router gates
// this route with one elevation check for the whole batch (action
// "connector.bulkRestart" — a distinct action from "connector.restart"
// since RequireElevation validates one token against one exact action
// string, not per item). Per ID, reuses the same restart logic as the
// single-connector Restart handler.
func (h *Handler) BulkRestart(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBulkRequest(w, r)
	if !ok {
		return
	}

	allowedIDs, found, results, err := h.splitByConnectorGrant(r.Context(), req.IDs)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	auditRecords := make([]store.AuditRecord, 0, len(allowedIDs))
	for _, id := range allowedIDs {
		rec := found[id]
		if err := h.restartConnector(r.Context(), &rec); err != nil {
			results = append(results, bulkItemResult{ID: id, Status: "error", Reason: err.Error()})
			continue
		}
		auditRecords = append(auditRecords, store.AuditRecord{TargetID: id})
		results = append(results, bulkItemResult{ID: id, Status: "success"})
	}

	if err := h.Store.RecordAuditBatchFromContext(r.Context(), "connector.bulk_restart", "connector", auditRecords); err != nil {
		slog.Error("failed to record audit", "action", "connector.bulk_restart", "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"results": results})
}

// restartConnector performs one connector's restart (config load, Restarter
// type-assert, call, failure-alert) without the per-call elevation check or
// per-call audit write — those are handled once, batch-wide, by BulkRestart.
// The single-connector RestartPreview/mutateOp path still does its own
// per-call elevation + audit, since it isn't part of a batch.
func (h *Handler) restartConnector(ctx context.Context, rec *store.ConnectorRecord) error {
	cfg, err := store.ParseConnectorConfig(rec.Type, rec.ConfigData, h.Config.Encryption.Key)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	cfg["url"] = rec.URL
	cfg["verify_tls"] = rec.VerifyTLS

	conn, err := connector.Get(rec.Type, cfg)
	if err != nil {
		return err
	}
	restarter, ok := conn.(connector.Restarter)
	if !ok {
		return fmt.Errorf("connector does not support restart")
	}

	if err := restarter.Restart(ctx, cfg, ""); err != nil {
		slog.Error("connector restart failed", "connector", rec.ID, "error", err)
		alert := &store.AlertRecord{
			ServiceID:   rec.ID,
			Severity:    "critical",
			Title:       fmt.Sprintf("Restart failed for %s", rec.Name),
			Description: err.Error(),
		}
		if createErr := h.Store.CreateAlert(ctx, alert); createErr != nil {
			slog.Error("failed to create restart failure alert", "error", createErr)
		} else if h.WSHub != nil {
			h.WSHub.Broadcast(ws.EventAlertCreated, map[string]any{
				"alertId":   alert.ID,
				"serviceId": rec.ID,
				"severity":  alert.Severity,
				"title":     alert.Title,
			})
		}
		return err
	}
	return nil
}
