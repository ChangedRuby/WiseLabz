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
