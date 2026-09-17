package api_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func validComplianceRule() map[string]any {
	return map[string]any{
		"name": "Open any rule", "connectorType": "pfsense", "entityKind": "rule",
		"conditions": []map[string]any{{"attribute": "action", "op": "eq", "value": "pass"}},
		"severity":   "warning", "title": "Open firewall rule", "remediationLink": "", "enabled": false,
	}
}

func TestComplianceRulesCRUDAndAdminGate(t *testing.T) {
	app := newTestApp(t)
	_, adminToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")

	for _, endpoint := range []struct{ method, path string }{
		{http.MethodGet, "/api/compliance/rules"}, {http.MethodPost, "/api/compliance/rules"},
		{http.MethodGet, "/api/compliance/rules/nope"}, {http.MethodPut, "/api/compliance/rules/nope"},
		{http.MethodDelete, "/api/compliance/rules/nope"}, {http.MethodPost, "/api/compliance/rules/test"},
	} {
		rec := app.req(t, endpoint.method, endpoint.path, validComplianceRule(), viewerToken)
		if rec.Code != http.StatusForbidden {
			t.Errorf("%s %s = %d, want 403", endpoint.method, endpoint.path, rec.Code)
		}
	}

	created := app.req(t, http.MethodPost, "/api/compliance/rules", validComplianceRule(), adminToken)
	if created.Code != http.StatusCreated {
		t.Fatalf("create = %d: %s", created.Code, created.Body.String())
	}
	var rule struct {
		ID         string `json:"id"`
		Conditions []struct {
			Attribute string `json:"attribute"`
		} `json:"conditions"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &rule); err != nil {
		t.Fatal(err)
	}
	if rule.ID == "" || len(rule.Conditions) != 1 {
		t.Fatalf("unexpected create response: %s", created.Body.String())
	}

	got := app.req(t, http.MethodGet, "/api/compliance/rules/"+rule.ID, nil, adminToken)
	if got.Code != http.StatusOK {
		t.Fatalf("get = %d: %s", got.Code, got.Body.String())
	}
	updated := validComplianceRule()
	updated["enabled"] = true
	put := app.req(t, http.MethodPut, "/api/compliance/rules/"+rule.ID, updated, adminToken)
	if put.Code != http.StatusOK {
		t.Fatalf("update = %d: %s", put.Code, put.Body.String())
	}
	preview := app.req(t, http.MethodPost, "/api/compliance/rules/test", updated, adminToken)
	if preview.Code != http.StatusOK {
		t.Fatalf("test = %d: %s", preview.Code, preview.Body.String())
	}
	deleted := app.req(t, http.MethodDelete, "/api/compliance/rules/"+rule.ID, nil, adminToken)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete = %d: %s", deleted.Code, deleted.Body.String())
	}
	if rec := app.req(t, http.MethodGet, "/api/compliance/rules/"+rule.ID, nil, adminToken); rec.Code != http.StatusNotFound {
		t.Fatalf("get deleted = %d", rec.Code)
	}
}

func TestComplianceRuleValidation(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")
	rule := validComplianceRule()
	rule["conditions"] = []map[string]any{}
	if rec := app.req(t, http.MethodPost, "/api/compliance/rules", rule, token); rec.Code != http.StatusBadRequest {
		t.Fatalf("empty conditions = %d", rec.Code)
	}
	rule = validComplianceRule()
	rule["connectorType"] = "unknown"
	if rec := app.req(t, http.MethodPost, "/api/compliance/rules", rule, token); rec.Code != http.StatusBadRequest {
		t.Fatalf("unknown type = %d", rec.Code)
	}
}
