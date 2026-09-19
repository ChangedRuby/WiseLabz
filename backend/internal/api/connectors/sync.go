package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/logsafe"
	"github.com/WiseLabz/wiselabz/internal/store"
)

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
