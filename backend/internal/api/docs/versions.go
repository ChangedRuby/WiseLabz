package docs

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Versions handles GET /api/docs/{id}/versions.
func (h *Handler) Versions(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	versions, err := h.Store.GetDocVersions(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, versions)
}

// Version handles GET /api/docs/{id}/versions/{rev}.
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	docID := r.PathValue("id")
	rev, err := strconv.Atoi(r.PathValue("rev"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid rev")
		return
	}

	versions, err := h.Store.GetDocVersions(r.Context(), docID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	for _, v := range versions {
		if v.Rev == rev {
			httputil.JSON(w, http.StatusOK, v)
			return
		}
	}
	httputil.Error(w, http.StatusNotFound, "not_found", "Version not found")
}

// Restore handles POST /api/docs/{id}/versions/{rev}/restore.
// Restores a doc to a previous version.
func (h *Handler) Restore(w http.ResponseWriter, r *http.Request) {
	docID := r.PathValue("id")
	rev, err := strconv.Atoi(r.PathValue("rev"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid rev")
		return
	}

	versions, err := h.Store.GetDocVersions(r.Context(), docID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	var target *store.DocVersionRecord
	for _, v := range versions {
		if v.Rev == rev {
			target = &v
			break
		}
	}
	if target == nil {
		httputil.Error(w, http.StatusNotFound, "not_found", "Version not found")
		return
	}

	existingDoc, err := h.Store.GetDoc(r.Context(), docID)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.requireDocOperator(w, r, existingDoc.ServiceID) {
		return
	}

	if err := h.Store.UpdateDoc(r.Context(), docID, target.Content, nil); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "doc.restore", "doc", docID, map[string]any{
		"rev": rev,
	}); err != nil {
		slog.Error("failed to record audit", "action", "doc.restore", "error", err)
	}

	d, _ := h.Store.GetDoc(r.Context(), docID)
	httputil.JSON(w, http.StatusOK, d)
}
