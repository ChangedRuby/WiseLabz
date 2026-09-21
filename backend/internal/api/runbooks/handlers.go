// Package runbooks provides API handlers for actionable runbooks attached to
// change types and alert severities.
package runbooks

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for runbook endpoints.
type Handler struct {
	Store *store.Store
}

// NewHandler creates a new runbook handler.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s}
}

func validTargetType(t string) bool {
	return t == "change_type" || t == "alert_severity"
}

// List handles GET /api/runbooks. changeType and alertSeverity are mutually
// exclusive filters on target_type/target_value; the table is expected to
// stay small, so filtering and pagination happen in Go rather than SQL.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	changeType := r.URL.Query().Get("changeType")
	alertSeverity := r.URL.Query().Get("alertSeverity")
	if changeType != "" && alertSeverity != "" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "changeType and alertSeverity are mutually exclusive", []httputil.FieldError{
			{Field: "changeType", Msg: "is mutually exclusive with alertSeverity"},
			{Field: "alertSeverity", Msg: "is mutually exclusive with changeType"},
		})
		return
	}

	page, pageSize, offset := httputil.Paginate(r)

	all, err := h.Store.ListRunbooks(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	var targetType, targetValue string
	switch {
	case changeType != "":
		targetType, targetValue = "change_type", changeType
	case alertSeverity != "":
		targetType, targetValue = "alert_severity", alertSeverity
	}

	filtered := all
	if targetType != "" {
		filtered = make([]*store.RunbookRecord, 0, len(all))
		for _, rb := range all {
			if rb.TargetType == targetType && rb.TargetValue == targetValue {
				filtered = append(filtered, rb)
			}
		}
	}

	total := len(filtered)
	start := offset
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	httputil.WritePaginated(w, filtered[start:end], page, pageSize, total)
}

// Get handles GET /api/runbooks/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rb, err := h.Store.GetRunbook(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Runbook not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, rb)
}

// createRequest is the body of POST /api/runbooks (RunbookCreate).
type createRequest struct {
	Title       string  `json:"title"`
	Body        string  `json:"body"`
	TargetType  string  `json:"targetType"`
	TargetValue string  `json:"targetValue"`
	SnapshotID  *string `json:"snapshotId"`
	DocID       *string `json:"docId"`
}

// Create handles POST /api/runbooks.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[createRequest](w, r)
	if !ok {
		return
	}
	if fieldErrs := httputil.MissingFields("title", req.Title, "targetValue", req.TargetValue); len(fieldErrs) > 0 {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "title and targetValue are required", fieldErrs)
		return
	}
	if !validTargetType(req.TargetType) {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "targetType must be change_type or alert_severity", []httputil.FieldError{{Field: "targetType", Msg: "must be change_type or alert_severity"}})
		return
	}

	created, err := h.Store.CreateRunbook(r.Context(), &store.RunbookRecord{
		Title:       req.Title,
		Body:        req.Body,
		TargetType:  req.TargetType,
		TargetValue: req.TargetValue,
		SnapshotID:  req.SnapshotID,
		DocID:       req.DocID,
	})
	if errors.Is(err, store.ErrConflict) {
		httputil.Error(w, http.StatusConflict, "conflict", "A runbook already exists for this target")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, created)
}

// updateStringFields maps RunbookUpdate JSON keys to their store column names
// for plain string fields.
var updateStringFields = map[string]string{
	"title":       "title",
	"body":        "body",
	"targetType":  "target_type",
	"targetValue": "target_value",
}

// updateNullableFields maps RunbookUpdate JSON keys to their store column
// names for nullable (*string) fields.
var updateNullableFields = map[string]string{
	"snapshotId": "snapshot_id",
	"docId":      "doc_id",
}

// Update handles PUT /api/runbooks/{id}. Only fields present in the request
// body are applied, so a raw map is decoded first to distinguish "absent"
// from "explicitly null" on the nullable fields.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	raw, ok := httputil.DecodeJSON[map[string]json.RawMessage](w, r)
	if !ok {
		return
	}

	updates := map[string]any{}
	for jsonKey, col := range updateStringFields {
		v, ok := raw[jsonKey]
		if !ok {
			continue
		}
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", jsonKey+" must be a string", []httputil.FieldError{{Field: jsonKey, Msg: "must be a string"}})
			return
		}
		updates[col] = s
	}
	for jsonKey, col := range updateNullableFields {
		v, ok := raw[jsonKey]
		if !ok {
			continue
		}
		var s *string
		if err := json.Unmarshal(v, &s); err != nil {
			httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", jsonKey+" must be a string or null", []httputil.FieldError{{Field: jsonKey, Msg: "must be a string or null"}})
			return
		}
		updates[col] = s
	}

	if tt, ok := updates["target_type"]; ok && !validTargetType(tt.(string)) {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "targetType must be change_type or alert_severity", []httputil.FieldError{{Field: "targetType", Msg: "must be change_type or alert_severity"}})
		return
	}

	rb, err := h.Store.UpdateRunbook(r.Context(), id, updates)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Runbook not found")
		return
	}
	if errors.Is(err, store.ErrConflict) {
		httputil.Error(w, http.StatusConflict, "conflict", "A runbook already exists for this target")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, rb)
}

// Delete handles DELETE /api/runbooks/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.Store.DeleteRunbook(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Runbook not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.NoContent(w)
}
