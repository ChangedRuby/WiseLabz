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
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
)

// Snapshots handles GET /api/connectors/{id}/snapshots.
func (h *Handler) Snapshots(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := h.Store.GetConnector(r.Context(), id); err != nil {
		snapshotStoreError(w, err)
		return
	}
	limit := httputil.DefaultPageSize
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > httputil.MaxPageSize {
		limit = httputil.MaxPageSize
	}
	keyset, sort, cursorID, ok := httputil.Cursor(w, r)
	if !ok {
		return
	}
	items, _, err := h.Store.ListSnapshotsByConnectorKeyset(r.Context(), id, store.Keyset{Sort: sort, ID: cursorID}, limit)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if next := httputil.NextCursor(items, limit, func(s store.SnapshotSummary) (string, string) {
		return s.FetchedAt, s.ID
	}); keyset && next != "" {
		w.Header().Set(httputil.NextCursorHeader, next)
	}
	httputil.JSON(w, http.StatusOK, items)
}

// Snapshot handles GET /api/connectors/{id}/snapshots/{snapshotId}.
func (h *Handler) Snapshot(w http.ResponseWriter, r *http.Request) {
	record, err := h.Store.GetSnapshotForConnector(r.Context(), r.PathValue("id"), r.PathValue("snapshotId"))
	if err != nil {
		snapshotStoreError(w, err)
		return
	}
	snap, err := decodeStoredSnapshot(record)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, struct {
		ID           string                        `json:"id"`
		ConnectorID  string                        `json:"connectorId"`
		ServiceName  string                        `json:"serviceName"`
		Type         string                        `json:"type"`
		Sections     []connector.SnapshotSection   `json:"sections"`
		Entities     []connector.SnapshotEntity    `json:"entities"`
		Dependencies []connector.ServiceDependency `json:"dependencies"`
		Metadata     map[string]string             `json:"metadata"`
		FetchedAt    string                        `json:"fetchedAt"`
	}{record.ID, record.ConnectorID, snap.ServiceName, snap.Type, snap.Sections, snap.Entities, snap.Dependencies, snap.Metadata, record.FetchedAt})
}

// SnapshotDiff handles GET /api/connectors/{id}/snapshots/diff.
func (h *Handler) SnapshotDiff(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fromID, toID, format := query.Get("from"), query.Get("to"), query.Get("format")
	if fromID == "" || toID == "" {
		fields := make([]httputil.FieldError, 0, 2)
		if fromID == "" {
			fields = append(fields, httputil.FieldError{Field: "from", Msg: "is required"})
		}
		if toID == "" {
			fields = append(fields, httputil.FieldError{Field: "to", Msg: "is required"})
		}
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "from and to are required", fields)
		return
	}
	if query.Has("format") && format != "json" && format != "csv" && format != "md" && format != "html" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_format", "format must be json, csv, md, or html", []httputil.FieldError{{Field: "format", Msg: "must be json, csv, md, or html"}})
		return
	}
	id := r.PathValue("id")
	conn, err := h.Store.GetConnector(r.Context(), id)
	if err != nil {
		snapshotStoreError(w, err)
		return
	}
	from, err := h.Store.GetSnapshotForConnector(r.Context(), id, fromID)
	if err != nil {
		snapshotStoreError(w, err)
		return
	}
	to, err := h.Store.GetSnapshotForConnector(r.Context(), id, toID)
	if err != nil {
		snapshotStoreError(w, err)
		return
	}
	previous, err := decodeStoredSnapshot(from)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	current, err := decodeStoredSnapshot(to)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	diff := sync.BuildSnapshotDiff(previous, current)
	diff.Provenance = sync.SnapshotDiffProvenance{
		ConnectorID: id, ConnectorName: conn.Name,
		From:        sync.SnapshotSourceFromRecord(from.ID, from.FetchedAt, from.Data),
		To:          sync.SnapshotSourceFromRecord(to.ID, to.FetchedAt, to.Data),
		GeneratedAt: time.Now().UTC(), GeneratedBy: auth.UserIDFromContext(r.Context()),
	}
	if format == "" {
		httputil.JSON(w, http.StatusOK, diff)
		return
	}
	data, contentType, err := sync.RenderSnapshotDiff(diff, format)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "snapshot.diff.export", "connector", id, map[string]any{"from": fromID, "to": toID, "format": format}); err != nil {
		slog.Error("failed to record audit", "action", "snapshot.diff.export", "error", err)
	}
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="wiselabz-snapshot-diff-%s-%s.%s"`, id, time.Now().UTC().Format("20060102-150405"), format))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func decodeStoredSnapshot(record *store.SnapshotRecord) (*connector.ServiceSnapshot, error) {
	var snap connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(record.Data), &snap); err != nil {
		return nil, fmt.Errorf("decode stored snapshot %s: %w", record.ID, err)
	}
	if snap.Sections == nil {
		snap.Sections = []connector.SnapshotSection{}
	}
	if snap.Entities == nil {
		snap.Entities = []connector.SnapshotEntity{}
	}
	if snap.Dependencies == nil {
		snap.Dependencies = []connector.ServiceDependency{}
	}
	if snap.Metadata == nil {
		snap.Metadata = map[string]string{}
	}
	return &snap, nil
}

func snapshotStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Snapshot or connector not found")
		return
	}
	httputil.Errorf(w, err)
}
