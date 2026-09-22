package connectors

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// PinGoldenSnapshot handles POST /api/connectors/{id}/golden-snapshot.
// Pins a snapshot as the connector's known-good configuration baseline for
// drift detection (#275). Body may omit snapshotId to pin the connector's
// current latest snapshot.
func (h *Handler) PinGoldenSnapshot(w http.ResponseWriter, r *http.Request) {
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
		SnapshotID string `json:"snapshotId"`
	}](w, r)
	if !ok {
		return
	}

	snapshotID := req.SnapshotID
	if snapshotID == "" {
		latest, err := h.Store.GetLatestSnapshot(r.Context(), id)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				httputil.Error(w, http.StatusBadRequest, "invalid_request", "Connector has no snapshot to pin")
				return
			}
			httputil.Errorf(w, err)
			return
		}
		snapshotID = latest.ID
	} else {
		snap, err := h.Store.GetSnapshotByID(r.Context(), snapshotID)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				httputil.Error(w, http.StatusBadRequest, "invalid_request", "Snapshot not found")
				return
			}
			httputil.Errorf(w, err)
			return
		}
		if snap.ConnectorID != id {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "Snapshot does not belong to this connector")
			return
		}
	}

	g := &store.GoldenSnapshotRecord{
		ConnectorID: id,
		SnapshotID:  snapshotID,
		PinnedBy:    auth.UserIDFromContext(r.Context()),
	}
	if err := h.Store.PinGoldenSnapshot(r.Context(), g); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.goldenSnapshot.pin", "connector", id, map[string]any{
		"snapshotId": snapshotID,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.goldenSnapshot.pin", "error", err)
	}

	httputil.JSON(w, http.StatusOK, g)
}

// UnpinGoldenSnapshot handles DELETE /api/connectors/{id}/golden-snapshot.
// Removes the connector's golden snapshot, if any. A no-op (200) when none
// is pinned.
func (h *Handler) UnpinGoldenSnapshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.Store.GetGoldenSnapshot(r.Context(), id)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		httputil.Errorf(w, err)
		return
	}
	if existing == nil {
		httputil.JSON(w, http.StatusOK, map[string]any{"unpinned": false})
		return
	}

	if err := h.Store.UnpinGoldenSnapshot(r.Context(), id); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.goldenSnapshot.unpin", "connector", id, map[string]any{
		"snapshotId": existing.SnapshotID,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.goldenSnapshot.unpin", "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"unpinned": true})
}

// GetGoldenSnapshot handles GET /api/connectors/{id}/golden-snapshot.
// Returns the connector's pinned golden snapshot, or null if none.
func (h *Handler) GetGoldenSnapshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	g, err := h.Store.GetGoldenSnapshot(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.JSON(w, http.StatusOK, nil)
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, g)
}
