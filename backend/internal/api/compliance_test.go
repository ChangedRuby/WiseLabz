package api_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestComplianceSchemaRequiresInstanceAdmin(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/compliance/schema", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestComplianceSchemaReturnsConnectorAttributeCatalog(t *testing.T) {
	app := newTestApp(t)
	_, adminToken := app.user(t, "operator")

	rec := app.req(t, http.MethodGet, "/api/compliance/schema", nil, adminToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var schema map[string]map[string][]struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &schema); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	pfsense, ok := schema["pfsense"]
	if !ok {
		t.Fatalf("schema missing pfsense connector type: %+v", schema)
	}
	ruleAttrs, ok := pfsense["rule"]
	if !ok || len(ruleAttrs) == 0 {
		t.Fatalf("schema missing pfsense rule attributes: %+v", pfsense)
	}
	found := false
	for _, a := range ruleAttrs {
		if a.Name == "enabled" && a.Type == "boolean" {
			found = true
		}
	}
	if !found {
		t.Errorf("pfsense rule attributes missing boolean \"enabled\": %+v", ruleAttrs)
	}
}
