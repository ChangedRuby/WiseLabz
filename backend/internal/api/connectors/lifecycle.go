package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// RestartPreview handles POST /api/connectors/{id}/restart. dryRun=true
// previews the restart from the latest stored snapshot without touching the
// connector; dryRun absent/false performs the real, elevation-gated restart.
func (h *Handler) RestartPreview(w http.ResponseWriter, r *http.Request) {
	h.serveLifecycleFromRequest(w, r, "restart")
}

// StartPreview handles POST /api/connectors/{id}/start. Same dry-run/mutate
// split as RestartPreview.
func (h *Handler) StartPreview(w http.ResponseWriter, r *http.Request) {
	h.serveLifecycleFromRequest(w, r, "start")
}

// StopPreview handles POST /api/connectors/{id}/stop. Same dry-run/mutate
// split as RestartPreview. estimatedDowntimeSeconds is 0 in the preview
// since a stop's downtime is indefinite until an explicit start, not a
// bounded window.
func (h *Handler) StopPreview(w http.ResponseWriter, r *http.Request) {
	h.serveLifecycleFromRequest(w, r, "stop")
}

// serveLifecycleFromRequest resolves connectorID from the path and
// entityRef from the JSON body — read for both the dry-run and real-mutate
// cases, so a dry-run can preview against a specific entity too — and
// delegates to ServeLifecycleOp with no extra audit detail. This is the
// path used by the connector-scoped restart/start/stop endpoints; the
// runbooks handler calls ServeLifecycleOp directly with a stored
// connector/verb/entityRef and its own extraAudit instead.
func (h *Handler) serveLifecycleFromRequest(w http.ResponseWriter, r *http.Request, verb string) {
	id := r.PathValue("id")
	var body struct {
		EntityRef string `json:"entityRef"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body) // ponytail: absent/empty body means entityRef == "", matches default
	}
	h.ServeLifecycleOp(w, r, id, verb, body.EntityRef, nil)
}

// ServeLifecycleOp implements the shared dry-run/mutate split for the
// lifecycle verbs (restart/start/stop), driven by an already-resolved
// connectorID/verb/entityRef rather than reading them itself — so both the
// connector-scoped handlers above (id from the path, entityRef from the
// body) and the runbooks ExecuteStep handler (id/verb/entityRef from a
// stored RunbookStepRecord, body ignored) can share it.
//
// dryRun=true previews from the latest stored snapshot without touching the
// connector, targeting entityRef's matching entity when one is given.
// dryRun absent/false performs the real, elevation-gated mutation per ADR
// 0001/0002: config parse, verb-support check ("unsupported_operation"),
// X-Elevation-Token validation against "connector.<verb>", entityRef
// validation, a critical AlertRecord + ws broadcast on failure (no
// rollback), and an audit row on success whose detail is {entityRef} merged
// with extraAudit (e.g. {runbookId, stepId} from the runbooks handler).
func (h *Handler) ServeLifecycleOp(w http.ResponseWriter, r *http.Request, connectorID, verb, entityRef string, extraAudit map[string]any) {
	dryRun := len(r.URL.Query()["dryRun"]) == 1 && r.URL.Query()["dryRun"][0] == "true"
	if dryRun {
		h.lifecycleOpPreview(w, r, connectorID, verb, entityRef)
		return
	}
	h.lifecycleOpMutate(w, r, connectorID, verb, entityRef, extraAudit)
}

// lifecycleOpPreview builds the dry-run preview response for connectorID's
// latest snapshot. When entityRef matches a snapshot entity by ExternalID,
// that entity's name is used as targetService instead of the top-level
// service name, so a per-entity restart/start/stop previews the entity
// actually being targeted.
func (h *Handler) lifecycleOpPreview(w http.ResponseWriter, r *http.Request, connectorID, verb, entityRef string) {
	if _, err := h.Store.GetConnector(r.Context(), connectorID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	sn, err := h.Store.GetLatestSnapshot(r.Context(), connectorID)
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

	targetService := snap.ServiceName
	if entityRef != "" {
		for _, e := range snap.Entities {
			if e.ExternalID == entityRef {
				targetService = e.Name
				break
			}
		}
	}

	dependencies := snap.Dependencies
	if dependencies == nil {
		dependencies = []connector.ServiceDependency{}
	}
	downtime := restartPreviewDowntimeSeconds
	if verb == "stop" {
		// A stop's downtime is indefinite (no scheduled restart), not a
		// bounded estimate — 0 rather than inventing a new field.
		downtime = 0
	}
	httputil.JSON(w, http.StatusOK, map[string]any{
		"targetService":            targetService,
		"estimatedDowntimeSeconds": downtime,
		"dependentServices":        dependencies,
	})
}

// lifecycleOpMutate handles the real, mutating side of restart/start/stop
// (dryRun absent/false). Per ADR 0001/0002: gated by the verb's elevation
// action, no rollback on failure (an AlertRecord is raised instead), and
// only successes are audited.
func (h *Handler) lifecycleOpMutate(w http.ResponseWriter, r *http.Request, connectorID, verb, entityRef string, extraAudit map[string]any) {
	rec, err := h.Store.GetConnector(r.Context(), connectorID)
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

	fn, ok := connector.LifecycleOp(conn, verb)
	if !ok {
		httputil.Error(w, http.StatusBadRequest, "unsupported_operation", "connector does not support "+verb)
		return
	}

	elevateAction := "connector." + verb
	if err := auth.ValidateElevationHeader(h.JWT, h.Store, elevateAction, r); err != nil {
		auth.WriteElevationError(w, err)
		return
	}

	if err := connector.ValidateCompositeRef(entityRef); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "invalid entityRef")
		return
	}

	if err := fn(r.Context(), cfg, entityRef); err != nil {
		slog.Error("connector "+verb+" failed", "connector", connectorID, "error", err)
		alert := &store.AlertRecord{
			ServiceID:   connectorID,
			Severity:    "critical",
			Title:       fmt.Sprintf("%s failed for %s", capitalize(verb), rec.Name),
			Description: err.Error(),
		}
		if createErr := h.Store.CreateAlert(r.Context(), alert); createErr != nil {
			slog.Error("failed to create "+verb+" failure alert", "error", createErr)
		} else if h.WSHub != nil {
			h.WSHub.BroadcastConnector(connectorID, ws.EventAlertCreated, map[string]any{
				"alertId":   alert.ID,
				"serviceId": connectorID,
				"severity":  alert.Severity,
				"title":     alert.Title,
			})
		}
		httputil.Error(w, http.StatusBadGateway, verb+"_failed", err.Error())
		return
	}

	detail := map[string]any{"entityRef": entityRef}
	for k, v := range extraAudit {
		detail[k] = v
	}
	auditAction := "connector." + verb
	if err := h.Store.RecordAuditFromContext(r.Context(), auditAction, "connector", connectorID, detail); err != nil {
		slog.Error("failed to record audit", "action", auditAction, "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"status": verb + "ed"})
}

// capitalize upper-cases a word's first byte (ASCII verbs only: "restart",
// "start", "stop") for a display-friendly alert title.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
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

// restartConnector performs one connector's restart (config load,
// LifecycleOp lookup, call, failure-alert) without the per-call elevation
// check or per-call audit write — those are handled once, batch-wide, by
// BulkRestart. The single-connector RestartPreview/ServeLifecycleOp path
// still does its own per-call elevation + audit, since it isn't part of a
// batch.
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
	restart, ok := connector.LifecycleOp(conn, "restart")
	if !ok {
		return fmt.Errorf("connector does not support restart")
	}

	if err := restart(ctx, cfg, ""); err != nil {
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
			h.WSHub.BroadcastConnector(rec.ID, ws.EventAlertCreated, map[string]any{
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
