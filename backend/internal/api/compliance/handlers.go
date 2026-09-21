// Package compliance provides API handlers for compliance rules and schema.
package compliance

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/compliance"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// RuleEvaluator evaluates or resolves a persisted rule after an admin change.
type RuleEvaluator interface {
	EvaluateRule(context.Context, string) error
	ResolveRule(context.Context, string) error
}

// Handler serves admin compliance schema and rule endpoints.
type Handler struct {
	Store     *store.Store
	Evaluator RuleEvaluator
}

// NewHandler creates a compliance handler.
func NewHandler(s *store.Store, evaluator RuleEvaluator) *Handler {
	return &Handler{Store: s, Evaluator: evaluator}
}

// Schema handles GET /api/compliance/schema.
func (h *Handler) Schema(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, connector.AttributeCatalog())
}

type ruleRequest struct {
	Name            string                 `json:"name"`
	ConnectorType   string                 `json:"connectorType"`
	EntityKind      string                 `json:"entityKind"`
	Conditions      []compliance.Condition `json:"conditions"`
	Severity        string                 `json:"severity"`
	Title           string                 `json:"title"`
	RemediationLink string                 `json:"remediationLink"`
	Enabled         bool                   `json:"enabled"`
}

func catalog() compliance.Catalog {
	source := connector.AttributeCatalog()
	out := make(compliance.Catalog, len(source))
	for typ, kinds := range source {
		out[typ] = make(map[string][]compliance.AttributeSpec, len(kinds))
		for kind, attributes := range kinds {
			out[typ][kind] = make([]compliance.AttributeSpec, len(attributes))
			for i, attribute := range attributes {
				out[typ][kind][i] = compliance.AttributeSpec{Name: attribute.Name, Type: attribute.Type}
			}
		}
	}
	return out
}

func (req ruleRequest) record(id string) (store.ComplianceRuleRecord, error) {
	conditions, err := json.Marshal(req.Conditions)
	if err != nil {
		return store.ComplianceRuleRecord{}, err
	}
	return store.ComplianceRuleRecord{ID: id, Name: req.Name, ConnectorType: req.ConnectorType, EntityKind: req.EntityKind, Conditions: string(conditions), Severity: req.Severity, Title: req.Title, RemediationLink: req.RemediationLink, Enabled: req.Enabled}, nil
}

func toRule(record store.ComplianceRuleRecord) (compliance.Rule, error) {
	var conditions []compliance.Condition
	if err := json.Unmarshal([]byte(record.Conditions), &conditions); err != nil {
		return compliance.Rule{}, err
	}
	return compliance.Rule{ID: record.ID, Name: record.Name, ConnectorType: record.ConnectorType, EntityKind: record.EntityKind, Conditions: conditions, Severity: record.Severity, Title: record.Title, RemediationLink: record.RemediationLink, Enabled: record.Enabled}, nil
}

func validRecord(record store.ComplianceRuleRecord) error {
	rule, err := toRule(record)
	if err != nil {
		return err
	}
	if rule.Name == "" || rule.Title == "" {
		var fields []compliance.FieldError
		for _, f := range []struct{ name, value string }{{"name", rule.Name}, {"title", rule.Title}} {
			if f.value == "" {
				fields = append(fields, compliance.FieldError{Field: f.name, Msg: "is required"})
			}
		}
		return &compliance.ValidationError{Message: "name and title are required", Fields: fields}
	}
	if rule.Severity != "info" && rule.Severity != "warning" && rule.Severity != "critical" {
		return &compliance.ValidationError{
			Message: "severity must be info, warning, or critical",
			Fields:  []compliance.FieldError{{Field: "severity", Msg: "must be info, warning, or critical"}},
		}
	}
	return compliance.Validate(rule, catalog())
}

// writeRuleRejection writes the 400 for a rule that failed record() or
// validRecord(), attaching the offending rule fields as details whenever the
// failure named them. A rejection that names none — a conditions payload that
// will not round-trip through JSON — keeps the plain error envelope.
func writeRuleRejection(w http.ResponseWriter, err error) {
	var invalid *compliance.ValidationError
	if !errors.As(err, &invalid) || len(invalid.Fields) == 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	fields := make([]httputil.FieldError, len(invalid.Fields))
	for i, f := range invalid.Fields {
		fields[i] = httputil.FieldError{Field: f.Field, Msg: f.Msg}
	}
	httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", err.Error(), fields)
}

func response(record store.ComplianceRuleRecord) (map[string]any, error) {
	rule, err := toRule(record)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id": rule.ID, "name": rule.Name, "connectorType": rule.ConnectorType,
		"entityKind": rule.EntityKind, "conditions": rule.Conditions, "severity": rule.Severity,
		"title": rule.Title, "remediationLink": rule.RemediationLink, "enabled": rule.Enabled,
		"createdAt": record.CreatedAt, "updatedAt": record.UpdatedAt,
	}, nil
}

// List handles GET /api/compliance/rules.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	rules, err := h.Store.ListComplianceRules(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	items := make([]map[string]any, 0, len(rules))
	for _, rule := range rules {
		item, err := response(rule)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		items = append(items, item)
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"items": items})
}

// Get handles GET /api/compliance/rules/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	rule, err := h.Store.GetComplianceRule(r.Context(), r.PathValue("id"))
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Compliance rule not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	response, err := response(*rule)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, response)
}

// Create handles POST /api/compliance/rules.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[ruleRequest](w, r)
	if !ok {
		return
	}
	rule, err := req.record("")
	if err == nil {
		err = validRecord(rule)
	}
	if err != nil {
		writeRuleRejection(w, err)
		return
	}
	if err := h.Store.CreateComplianceRule(r.Context(), &rule); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if rule.Enabled && h.Evaluator != nil {
		if err := h.Evaluator.EvaluateRule(r.Context(), rule.ID); err != nil {
			httputil.Errorf(w, err)
			return
		}
	}
	h.audit(r, "compliance_rule.create", rule.ID, []string{"name", "connectorType", "entityKind", "conditions", "severity", "title", "remediationLink", "enabled"})
	response, err := response(rule)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, response)
}

// Update handles PUT /api/compliance/rules/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	old, err := h.Store.GetComplianceRule(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Compliance rule not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	req, ok := httputil.DecodeJSON[ruleRequest](w, r)
	if !ok {
		return
	}
	rule, err := req.record(id)
	if err == nil {
		err = validRecord(rule)
	}
	if err != nil {
		writeRuleRejection(w, err)
		return
	}
	rule.CreatedAt = old.CreatedAt
	changed := changedFields(*old, rule)
	if err := h.Store.UpdateComplianceRule(r.Context(), &rule); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if h.Evaluator != nil {
		if old.ConnectorType != rule.ConnectorType {
			connectors, err := h.Store.ListAllConnectors(r.Context())
			if err != nil {
				httputil.Errorf(w, err)
				return
			}
			for _, connector := range connectors {
				if connector.Type == old.ConnectorType {
					if err := h.Store.ResolveQualityFindingForRule(r.Context(), connector.ID, id); err != nil {
						httputil.Errorf(w, err)
						return
					}
				}
			}
		}
		if rule.Enabled {
			err = h.Evaluator.EvaluateRule(r.Context(), id)
		} else {
			err = h.Evaluator.ResolveRule(r.Context(), id)
		}
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
	}
	action := "compliance_rule.update"
	if old.Enabled != rule.Enabled {
		if rule.Enabled {
			action = "compliance_rule.enable"
		} else {
			action = "compliance_rule.disable"
		}
	}
	h.audit(r, action, id, changed)
	response, err := response(rule)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, response)
}

func changedFields(old, current store.ComplianceRuleRecord) []string {
	fields := make([]string, 0, 8)
	for _, field := range []struct{ name, old, current string }{
		{"name", old.Name, current.Name}, {"connectorType", old.ConnectorType, current.ConnectorType}, {"entityKind", old.EntityKind, current.EntityKind}, {"conditions", old.Conditions, current.Conditions}, {"severity", old.Severity, current.Severity}, {"title", old.Title, current.Title}, {"remediationLink", old.RemediationLink, current.RemediationLink},
	} {
		if field.old != field.current {
			fields = append(fields, field.name)
		}
	}
	if old.Enabled != current.Enabled {
		fields = append(fields, "enabled")
	}
	return fields
}

// Delete handles DELETE /api/compliance/rules/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if h.Evaluator != nil {
		if err := h.Evaluator.ResolveRule(r.Context(), id); err != nil && !errors.Is(err, store.ErrNotFound) {
			httputil.Errorf(w, err)
			return
		}
	}
	if err := h.Store.DeleteComplianceRule(r.Context(), id); errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Compliance rule not found")
		return
	} else if err != nil {
		httputil.Errorf(w, err)
		return
	}
	h.audit(r, "compliance_rule.delete", id, nil)
	httputil.NoContent(w)
}

// Test evaluates the supplied (unsaved) rule against current snapshots only.
func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[ruleRequest](w, r)
	if !ok {
		return
	}
	record, err := req.record("")
	if err == nil {
		err = validRecord(record)
	}
	if err != nil {
		writeRuleRejection(w, err)
		return
	}
	rule, _ := toRule(record)
	connectors, err := h.Store.ListAllConnectors(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	items := make([]map[string]any, 0)
	for _, c := range connectors {
		if c.Type != rule.ConnectorType {
			continue
		}
		snapshot, err := h.Store.GetLatestSnapshot(r.Context(), c.ID)
		if errors.Is(err, store.ErrNotFound) {
			continue
		}
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		var source connector.ServiceSnapshot
		if err := json.Unmarshal([]byte(snapshot.Data), &source); err != nil {
			httputil.Errorf(w, err)
			return
		}
		entities := make([]compliance.Entity, len(source.Entities))
		for i, entity := range source.Entities {
			entities[i] = compliance.Entity{Kind: entity.Kind, Name: entity.Name, Attributes: entity.Attributes}
		}
		matches := compliance.Evaluate(rule, compliance.Snapshot{Entities: entities})
		if len(matches) > 0 {
			items = append(items, map[string]any{"connectorId": c.ID, "connectorName": c.Name, "entities": matches})
		}
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) audit(r *http.Request, action, id string, changed []string) {
	if err := h.Store.RecordAuditFromContext(r.Context(), action, "compliance_rule", id, map[string]any{"changedFields": changed}); err != nil {
		slog.Error("failed to record audit", "action", action, "error", err)
	}
}
