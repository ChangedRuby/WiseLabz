package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/httputil"
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

// complianceCondition is one condition row of a rule request body.
func complianceCondition(attribute, op string, value any) map[string]any {
	return map[string]any{"attribute": attribute, "op": op, "value": value}
}

// badRegexMessage is the message validRecord produces for an uncompilable
// pattern, built from regexp's own error rather than pinned to its wording.
func badRegexMessage(t *testing.T, pattern string) string {
	t.Helper()
	_, err := regexp.Compile(pattern)
	if err == nil {
		t.Fatalf("regexp.Compile(%q) unexpectedly succeeded", pattern)
	}
	return fmt.Sprintf("invalid regex for %q: %v", "action", err)
}

// TestComplianceRuleValidation pins every compliance rule rejection: the 400
// envelope's code and message stay exactly what the check produced, and the
// details array names the offending field — with an indexed path
// (conditions[1].attribute) for the per-condition checks, so a client can
// point at the exact row of a multi-condition rule.
//
// Each case runs against all three routes that validate a rule body, since
// Create, Update and Test share one rejection path and must not drift apart.
func TestComplianceRuleValidation(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")

	created := app.req(t, http.MethodPost, "/api/compliance/rules", validComplianceRule(), token)
	if created.Code != http.StatusCreated {
		t.Fatalf("seed create = %d: %s", created.Code, created.Body.String())
	}
	var seeded struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &seeded); err != nil {
		t.Fatal(err)
	}

	manyConditions := make([]map[string]any, 21)
	for i := range manyConditions {
		manyConditions[i] = complianceCondition("action", "eq", "pass")
	}

	cases := []struct {
		name    string
		mutate  func(map[string]any)
		message string
		details []httputil.FieldError
	}{
		{
			name:    "name and title empty",
			mutate:  func(r map[string]any) { r["name"] = ""; r["title"] = "" },
			message: "name and title are required",
			details: []httputil.FieldError{
				{Field: "name", Msg: "is required"},
				{Field: "title", Msg: "is required"},
			},
		},
		{
			name:    "title empty",
			mutate:  func(r map[string]any) { r["title"] = "" },
			message: "name and title are required",
			details: []httputil.FieldError{{Field: "title", Msg: "is required"}},
		},
		{
			name:    "severity not a known level",
			mutate:  func(r map[string]any) { r["severity"] = "fatal" },
			message: "severity must be info, warning, or critical",
			details: []httputil.FieldError{{Field: "severity", Msg: "must be info, warning, or critical"}},
		},
		{
			name:    "connectorType unknown",
			mutate:  func(r map[string]any) { r["connectorType"] = "unknown" },
			message: `unknown connector type "unknown"`,
			details: []httputil.FieldError{{Field: "connectorType", Msg: "is not a known connector type"}},
		},
		{
			name:    "entityKind unknown for connector type",
			mutate:  func(r map[string]any) { r["entityKind"] = "nope" },
			message: `unknown entity kind "nope" for connector type "pfsense"`,
			details: []httputil.FieldError{{Field: "entityKind", Msg: "is not a known entity kind for this connector type"}},
		},
		{
			name:    "conditions empty",
			mutate:  func(r map[string]any) { r["conditions"] = []map[string]any{} },
			message: "at least one condition is required",
			details: []httputil.FieldError{{Field: "conditions", Msg: "must contain at least one condition"}},
		},
		{
			name:    "conditions over the limit",
			mutate:  func(r map[string]any) { r["conditions"] = manyConditions },
			message: "at most 20 conditions are allowed",
			details: []httputil.FieldError{{Field: "conditions", Msg: "must contain at most 20 conditions"}},
		},
		{
			// The offending condition is the second one: details must point at
			// that row, not just at "conditions".
			name: "condition attribute unknown",
			mutate: func(r map[string]any) {
				r["conditions"] = []map[string]any{
					complianceCondition("action", "eq", "pass"),
					complianceCondition("nope", "eq", "x"),
				}
			},
			message: `unknown attribute "nope"`,
			details: []httputil.FieldError{{Field: "conditions[1].attribute", Msg: "is not a known attribute for this entity kind"}},
		},
		{
			name: "condition op invalid for attribute type",
			mutate: func(r map[string]any) {
				r["conditions"] = []map[string]any{complianceCondition("log", "contains", "x")}
			},
			message: `operator "contains" is not valid for boolean attribute "log"`,
			details: []httputil.FieldError{{Field: "conditions[0].op", Msg: "is not a valid operator for boolean attributes"}},
		},
		{
			name: "regex value not a string",
			mutate: func(r map[string]any) {
				r["conditions"] = []map[string]any{complianceCondition("action", "regex", 42)}
			},
			message: `regex value for "action" must be a string`,
			details: []httputil.FieldError{{Field: "conditions[0].value", Msg: "must be a string when op is regex"}},
		},
		{
			name: "regex value over the length limit",
			mutate: func(r map[string]any) {
				r["conditions"] = []map[string]any{complianceCondition("action", "regex", strings.Repeat("a", 257))}
			},
			message: `regex for "action" exceeds 256 characters`,
			details: []httputil.FieldError{{Field: "conditions[0].value", Msg: "must be at most 256 characters"}},
		},
		{
			name: "regex value does not compile",
			mutate: func(r map[string]any) {
				r["conditions"] = []map[string]any{complianceCondition("action", "regex", "(")}
			},
			message: badRegexMessage(t, "("),
			details: []httputil.FieldError{{Field: "conditions[0].value", Msg: "is not a valid regular expression"}},
		},
	}

	routes := []struct{ method, path string }{
		{http.MethodPost, "/api/compliance/rules"},
		{http.MethodPut, "/api/compliance/rules/" + seeded.ID},
		{http.MethodPost, "/api/compliance/rules/test"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, route := range routes {
				body := validComplianceRule()
				tc.mutate(body)

				rec := app.req(t, route.method, route.path, body, token)
				if rec.Code != http.StatusBadRequest {
					t.Fatalf("%s %s = %d, want 400 (body: %s)", route.method, route.path, rec.Code, rec.Body.String())
				}

				var got struct {
					Code    string                `json:"code"`
					Message string                `json:"message"`
					Details []httputil.FieldError `json:"details"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("%s %s: decode %q: %v", route.method, route.path, rec.Body.String(), err)
				}
				if got.Code != "invalid_request" {
					t.Errorf("%s %s: code = %q, want invalid_request", route.method, route.path, got.Code)
				}
				if got.Message != tc.message {
					t.Errorf("%s %s: message = %q, want %q", route.method, route.path, got.Message, tc.message)
				}
				if !slices.Equal(got.Details, tc.details) {
					t.Errorf("%s %s: details = %+v, want %+v", route.method, route.path, got.Details, tc.details)
				}

				apitest.AssertMatchesSpec(t, app.newRequest(t, route.method, route.path, body, token), rec.Result())
			}
		})
	}
}
