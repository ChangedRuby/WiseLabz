// Package runbooks provides API handlers for actionable runbooks attached to
// change types, alert severities, and finding check types.
package runbooks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/api/connectors"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for runbook endpoints. ConnH is the connectors
// handler; ExecuteStep delegates to its ServeLifecycleOp so a step's
// execution shares the exact same dry-run/elevation/audit/alert path as a
// direct connector restart/start/stop.
type Handler struct {
	Store *store.Store
	ConnH *connectors.Handler
}

// NewHandler creates a new runbook handler.
func NewHandler(s *store.Store, connH *connectors.Handler) *Handler {
	return &Handler{Store: s, ConnH: connH}
}

func validTargetType(t string) bool {
	return t == "change_type" || t == "alert_severity" || t == "finding_check_type"
}

const targetTypeErrorMsg = "targetType must be change_type, alert_severity, or finding_check_type"

// maxRunbookSteps caps how many steps one runbook may hold — enough for any
// realistic remediation sequence without letting authoring turn into an
// unbounded list.
const maxRunbookSteps = 20

func validVerb(v string) bool {
	return v == "restart" || v == "start" || v == "stop"
}

// stepInput is one element of the "steps" array in RunbookCreate/RunbookUpdate.
type stepInput struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ConnectorID string `json:"connectorId"`
	Verb        string `json:"verb"`
	EntityRef   string `json:"entityRef"`
}

// stepResponse is one element of the "steps" array in a runbook response.
// CanExecute/ExecuteBlockedReason are computed per calling user, never
// stored.
type stepResponse struct {
	ID                   string `json:"id"`
	Position             int    `json:"position"`
	Title                string `json:"title"`
	ConnectorID          string `json:"connectorId"`
	ConnectorName        string `json:"connectorName"`
	Verb                 string `json:"verb"`
	EntityRef            string `json:"entityRef"`
	CanExecute           bool   `json:"canExecute"`
	ExecuteBlockedReason string `json:"executeBlockedReason"`
}

// runbookResponse is a RunbookRecord plus its steps, as returned by
// Get/List/Create/Update.
type runbookResponse struct {
	store.RunbookRecord
	Steps []stepResponse `json:"steps"`
}

// validateSteps checks a Create/Update steps[] payload and, for each valid
// step, returns a store.RunbookStepRecord (position is assigned later, by
// array index, in ReplaceRunbookSteps) ready to persist. Field errors are
// keyed "steps[i].<field>" so the client can point at the offending row.
// Authoring only checks that the connector exists and its type supports
// the verb — no connector-grant check at authoring time (linking grants
// nothing; see the Handler doc comment).
func (h *Handler) validateSteps(ctx context.Context, inputs []stepInput) ([]*store.RunbookStepRecord, []httputil.FieldError) {
	if len(inputs) > maxRunbookSteps {
		return nil, []httputil.FieldError{{Field: "steps", Msg: fmt.Sprintf("must not exceed %d steps", maxRunbookSteps)}}
	}

	var fieldErrs []httputil.FieldError
	steps := make([]*store.RunbookStepRecord, 0, len(inputs))
	for i, in := range inputs {
		prefix := fmt.Sprintf("steps[%d]", i)

		if strings.TrimSpace(in.Title) == "" {
			fieldErrs = append(fieldErrs, httputil.FieldError{Field: prefix + ".title", Msg: "is required"})
		}
		if !validVerb(in.Verb) {
			fieldErrs = append(fieldErrs, httputil.FieldError{Field: prefix + ".verb", Msg: "must be restart, start, or stop"})
		}

		conn, err := h.Store.GetConnector(ctx, in.ConnectorID)
		switch {
		case err != nil:
			fieldErrs = append(fieldErrs, httputil.FieldError{Field: prefix + ".connectorId", Msg: "connector not found"})
		case validVerb(in.Verb) && !connector.SupportsLifecycleVerb(conn.Type, in.Verb):
			fieldErrs = append(fieldErrs, httputil.FieldError{Field: prefix + ".verb", Msg: "connector does not support this verb"})
		}

		if err := connector.ValidateCompositeRef(in.EntityRef); err != nil {
			fieldErrs = append(fieldErrs, httputil.FieldError{Field: prefix + ".entityRef", Msg: "invalid entityRef"})
		}

		steps = append(steps, &store.RunbookStepRecord{
			ID:          in.ID,
			Title:       in.Title,
			ConnectorID: in.ConnectorID,
			Verb:        in.Verb,
			EntityRef:   in.EntityRef,
		})
	}
	if len(fieldErrs) > 0 {
		return nil, fieldErrs
	}
	return steps, nil
}

// toStepResponses builds the response steps for one runbook's steps,
// computing canExecute per userID. connectorNames/canExecute are caller-
// provided caches so List (many runbooks, possibly sharing connectors) does
// one lookup per connector instead of one per step.
func (h *Handler) toStepResponses(ctx context.Context, userID string, steps []*store.RunbookStepRecord, connectorNames map[string]string, canExecute map[string]bool) ([]stepResponse, error) {
	out := make([]stepResponse, 0, len(steps))
	for _, st := range steps {
		name, ok := connectorNames[st.ConnectorID]
		if !ok {
			conn, err := h.Store.GetConnector(ctx, st.ConnectorID)
			if err != nil && !errors.Is(err, store.ErrNotFound) {
				return nil, err
			}
			if conn != nil {
				name = conn.Name
			}
			connectorNames[st.ConnectorID] = name
		}

		can, ok := canExecute[st.ConnectorID]
		if !ok {
			var err error
			can, err = h.Store.UserHasConnectorRole(ctx, userID, st.ConnectorID, "operator")
			if err != nil {
				return nil, err
			}
			canExecute[st.ConnectorID] = can
		}
		reason := ""
		if !can {
			reason = "no_operator_grant"
		}

		out = append(out, stepResponse{
			ID:                   st.ID,
			Position:             st.Position,
			Title:                st.Title,
			ConnectorID:          st.ConnectorID,
			ConnectorName:        name,
			Verb:                 st.Verb,
			EntityRef:            st.EntityRef,
			CanExecute:           can,
			ExecuteBlockedReason: reason,
		})
	}
	return out, nil
}

// toRunbookResponse builds the response for a single runbook (Get/Create/Update).
func (h *Handler) toRunbookResponse(ctx context.Context, rb *store.RunbookRecord, steps []*store.RunbookStepRecord) (runbookResponse, error) {
	stepResps, err := h.toStepResponses(ctx, auth.UserIDFromContext(ctx), steps, map[string]string{}, map[string]bool{})
	if err != nil {
		return runbookResponse{}, err
	}
	return runbookResponse{RunbookRecord: *rb, Steps: stepResps}, nil
}

// stepAuditDetail renders steps for the runbook.create/update audit detail.
func stepAuditDetail(steps []*store.RunbookStepRecord) []map[string]any {
	out := make([]map[string]any, 0, len(steps))
	for _, st := range steps {
		out = append(out, map[string]any{
			"id":          st.ID,
			"connectorId": st.ConnectorID,
			"verb":        st.Verb,
			"entityRef":   st.EntityRef,
		})
	}
	return out
}

// List handles GET /api/runbooks. changeType, alertSeverity, and
// findingCheckType are mutually exclusive filters on
// target_type/target_value; the table is expected to stay small, so
// filtering and pagination happen in Go rather than SQL.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	changeType := r.URL.Query().Get("changeType")
	alertSeverity := r.URL.Query().Get("alertSeverity")
	findingCheckType := r.URL.Query().Get("findingCheckType")

	set := 0
	for _, v := range []string{changeType, alertSeverity, findingCheckType} {
		if v != "" {
			set++
		}
	}
	if set > 1 {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "changeType, alertSeverity, and findingCheckType are mutually exclusive", []httputil.FieldError{
			{Field: "changeType", Msg: "is mutually exclusive with alertSeverity and findingCheckType"},
			{Field: "alertSeverity", Msg: "is mutually exclusive with changeType and findingCheckType"},
			{Field: "findingCheckType", Msg: "is mutually exclusive with changeType and alertSeverity"},
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
	case findingCheckType != "":
		targetType, targetValue = "finding_check_type", findingCheckType
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
	pageItems := filtered[start:end]

	ids := make([]string, len(pageItems))
	for i, rb := range pageItems {
		ids[i] = rb.ID
	}
	stepsByRunbook, err := h.Store.ListRunbookSteps(r.Context(), ids)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	connectorNames := map[string]string{}
	canExecute := map[string]bool{}
	responses := make([]runbookResponse, 0, len(pageItems))
	for _, rb := range pageItems {
		stepResps, err := h.toStepResponses(r.Context(), userID, stepsByRunbook[rb.ID], connectorNames, canExecute)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		responses = append(responses, runbookResponse{RunbookRecord: *rb, Steps: stepResps})
	}

	httputil.WritePaginated(w, responses, page, pageSize, total)
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
	steps, err := h.Store.ListRunbookStepsFor(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	resp, err := h.toRunbookResponse(r.Context(), rb, steps)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, resp)
}

// createRequest is the body of POST /api/runbooks (RunbookCreate).
type createRequest struct {
	Title       string      `json:"title"`
	Body        string      `json:"body"`
	TargetType  string      `json:"targetType"`
	TargetValue string      `json:"targetValue"`
	SnapshotID  *string     `json:"snapshotId"`
	DocID       *string     `json:"docId"`
	Steps       []stepInput `json:"steps"`
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
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", targetTypeErrorMsg, []httputil.FieldError{{Field: "targetType", Msg: targetTypeErrorMsg}})
		return
	}
	steps, fieldErrs := h.validateSteps(r.Context(), req.Steps)
	if len(fieldErrs) > 0 {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "invalid steps", fieldErrs)
		return
	}

	created, savedSteps, err := h.Store.CreateRunbookWithSteps(r.Context(), &store.RunbookRecord{
		Title:       req.Title,
		Body:        req.Body,
		TargetType:  req.TargetType,
		TargetValue: req.TargetValue,
		SnapshotID:  req.SnapshotID,
		DocID:       req.DocID,
	}, steps)
	if errors.Is(err, store.ErrConflict) {
		httputil.Error(w, http.StatusConflict, "conflict", "A runbook already exists for this target")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "runbook.create", "runbook", created.ID, map[string]any{
		"title": created.Title,
		"steps": stepAuditDetail(savedSteps),
	}); err != nil {
		slog.Error("failed to record audit", "action", "runbook.create", "error", err)
	}

	resp, err := h.toRunbookResponse(r.Context(), created, savedSteps)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, resp)
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
// from "explicitly null" on the nullable fields — and, for "steps", "absent"
// (leave the runbook's steps unchanged) from "present" (replace all of
// them, even with an empty array).
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	raw, ok := httputil.DecodeJSON[map[string]json.RawMessage](w, r)
	if !ok {
		return
	}

	updates := map[string]any{}
	var changedFields []string
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
		changedFields = append(changedFields, jsonKey)
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
		changedFields = append(changedFields, jsonKey)
	}

	if tt, ok := updates["target_type"]; ok && !validTargetType(tt.(string)) {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", targetTypeErrorMsg, []httputil.FieldError{{Field: "targetType", Msg: targetTypeErrorMsg}})
		return
	}

	var stepRecords []*store.RunbookStepRecord
	replaceSteps := false
	if v, ok := raw["steps"]; ok {
		replaceSteps = true
		changedFields = append(changedFields, "steps")
		var inputs []stepInput
		if err := json.Unmarshal(v, &inputs); err != nil {
			httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "steps must be an array", []httputil.FieldError{{Field: "steps", Msg: "must be an array"}})
			return
		}
		var fieldErrs []httputil.FieldError
		stepRecords, fieldErrs = h.validateSteps(r.Context(), inputs)
		if len(fieldErrs) > 0 {
			httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "invalid steps", fieldErrs)
			return
		}
	}

	rb, savedSteps, err := h.Store.UpdateRunbookWithSteps(r.Context(), id, updates, stepRecords, replaceSteps)
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

	auditDetail := map[string]any{"changedFields": changedFields}
	if replaceSteps {
		auditDetail["steps"] = stepAuditDetail(savedSteps)
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "runbook.update", "runbook", id, auditDetail); err != nil {
		slog.Error("failed to record audit", "action", "runbook.update", "error", err)
	}

	resp, err := h.toRunbookResponse(r.Context(), rb, savedSteps)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /api/runbooks/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Store.DeleteRunbook(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Runbook not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "runbook.delete", "runbook", id, map[string]any{
		"title": rb.Title,
	}); err != nil {
		slog.Error("failed to record audit", "action", "runbook.delete", "error", err)
	}

	httputil.NoContent(w)
}

// ExecuteStep handles POST /api/runbooks/{id}/steps/{stepId}/execute
// [?dryRun=true]. The target connector/verb/entityRef always come from the
// stored step, never from the request body, and execution is delegated to
// the connectors handler's ServeLifecycleOp — the exact same dry-run
// preview / elevation-gated mutate / failure-alert / audit path as a
// direct connector restart/start/stop, with runbookId/stepId merged into
// the audit detail. Authoring a step grants nothing on its own: the caller
// must additionally hold at least an operator grant on the step's
// connector, checked here (not at authoring time).
func (h *Handler) ExecuteStep(w http.ResponseWriter, r *http.Request) {
	runbookID := r.PathValue("id")
	stepID := r.PathValue("stepId")

	if _, err := h.Store.GetRunbook(r.Context(), runbookID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Runbook not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	step, err := h.Store.GetRunbookStep(r.Context(), runbookID, stepID)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Runbook step not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	ok, err := h.Store.UserHasConnectorRole(r.Context(), auth.UserIDFromContext(r.Context()), step.ConnectorID, "operator")
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !ok {
		httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
		return
	}

	h.ConnH.ServeLifecycleOp(w, r, step.ConnectorID, step.Verb, step.EntityRef, map[string]any{
		"runbookId": runbookID,
		"stepId":    stepID,
	})
}
