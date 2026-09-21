package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/httputil"
)

// complianceRule returns validComplianceRule with mut applied, for cases that
// reject exactly one field of an otherwise-valid rule.
func complianceRule(mut func(map[string]any)) map[string]any {
	rule := validComplianceRule()
	mut(rule)
	return rule
}

// TestValidationErrorDetails pins the {field, msg} details envelope (issue
// #350) for one representative validation failure per handler package, so a
// handler cannot quietly drop back to a detail-less 400. It also asserts the
// response still matches the Error schema in docs/openapi.yaml.
func TestValidationErrorDetails(t *testing.T) {
	cases := []struct {
		name    string
		role    string
		method  string
		path    string
		body    any
		details []httputil.FieldError
	}{
		{
			name: "savedviews list surface", role: "viewer",
			method: "GET", path: "/api/saved-views?surface=bogus",
			details: []httputil.FieldError{{Field: "surface", Msg: "must be one of: services, changes, alerts"}},
		},
		{
			name: "apikeys create name", role: "viewer",
			method: "POST", path: "/api/auth/api-keys", body: map[string]any{},
			details: []httputil.FieldError{{Field: "name", Msg: "is required"}},
		},
		{
			name: "chat create conversation scopeType", role: "viewer",
			method: "POST", path: "/api/chat/conversations", body: map[string]any{"scopeType": "nope"},
			details: []httputil.FieldError{{Field: "scopeType", Msg: `must be "doc" or "lab"`}},
		},
		{
			name: "system audit export format", role: "operator",
			method: "GET", path: "/api/system/audit/export?format=xml",
			details: []httputil.FieldError{{Field: "format", Msg: `must be "json" or "csv"`}},
		},
		{
			name: "system retention cronExpr", role: "operator",
			method: "PUT", path: "/api/system/settings/retention", body: map[string]any{"cronExpr": ""},
			details: []httputil.FieldError{{Field: "cronExpr", Msg: "must not be empty"}},
		},
		{
			name: "docs generate required fields", role: "operator",
			method: "POST", path: "/api/docs/generate", body: map[string]any{},
			details: []httputil.FieldError{
				{Field: "templateId", Msg: "is required"},
				{Field: "connectorId", Msg: "is required"},
			},
		},
		{
			name: "templates create name", role: "operator",
			method: "POST", path: "/api/templates", body: map[string]any{},
			details: []httputil.FieldError{{Field: "name", Msg: "is required"}},
		},
		{
			name: "alerts bulk-snooze until", role: "viewer",
			method: "POST", path: "/api/alerts/bulk-snooze", body: map[string]any{"ids": []string{"a"}},
			details: []httputil.FieldError{{Field: "until", Msg: "is required"}},
		},
		{
			name: "changes bulk-resolve status", role: "viewer",
			method: "POST", path: "/api/changes/bulk-resolve", body: map[string]any{"status": "nope", "ids": []string{"a"}},
			details: []httputil.FieldError{{Field: "status", Msg: "must be acknowledged or dismissed"}},
		},
		{
			name: "runbooks list mutually exclusive filters", role: "viewer",
			method: "GET", path: "/api/runbooks?changeType=upgrade&alertSeverity=critical",
			details: []httputil.FieldError{
				{Field: "changeType", Msg: "is mutually exclusive with alertSeverity"},
				{Field: "alertSeverity", Msg: "is mutually exclusive with changeType"},
			},
		},
		{
			name: "users create required fields", role: "operator",
			method: "POST", path: "/api/users", body: map[string]any{},
			details: []httputil.FieldError{
				{Field: "username", Msg: "is required"},
				{Field: "password", Msg: "is required"},
			},
		},
		{
			name: "auth login required fields", role: "",
			method: "POST", path: "/api/auth/login", body: map[string]any{},
			details: []httputil.FieldError{
				{Field: "username", Msg: "is required"},
				{Field: "password", Msg: "is required"},
			},
		},
		{
			name: "compliance rule connectorType", role: "operator",
			method: "POST", path: "/api/compliance/rules",
			body:    complianceRule(func(r map[string]any) { r["connectorType"] = "unknown" }),
			details: []httputil.FieldError{{Field: "connectorType", Msg: "is not a known connector type"}},
		},
		{
			// The offending condition is the second one: details must point at
			// that row, not just at "conditions".
			name: "compliance rule indexed condition", role: "operator",
			method: "POST", path: "/api/compliance/rules",
			body: complianceRule(func(r map[string]any) {
				r["conditions"] = []map[string]any{
					{"attribute": "action", "op": "eq", "value": "pass"},
					{"attribute": "nope", "op": "eq", "value": "x"},
				}
			}),
			details: []httputil.FieldError{{Field: "conditions[1].attribute", Msg: "is not a known attribute for this entity kind"}},
		},
		{
			name: "compliance rule name and title", role: "operator",
			method: "POST", path: "/api/compliance/rules",
			body: complianceRule(func(r map[string]any) { r["name"] = ""; r["title"] = "" }),
			details: []httputil.FieldError{
				{Field: "name", Msg: "is required"},
				{Field: "title", Msg: "is required"},
			},
		},
		{
			name: "settings fallback provider", role: "operator",
			method: "PUT", path: "/api/ai/config/fallback-providers", body: []map[string]any{{"model": "m"}},
			details: []httputil.FieldError{{Field: "[0].provider", Msg: "is required"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp(t)
			var token string
			if tc.role != "" {
				_, token = app.user(t, tc.role)
			}

			req := app.newRequest(t, tc.method, tc.path, tc.body, token)
			rec := app.serve(req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("%s %s: got %d, want 400 (body: %s)", tc.method, tc.path, rec.Code, rec.Body.String())
			}

			var got struct {
				Code    string                `json:"code"`
				Message string                `json:"message"`
				Details []httputil.FieldError `json:"details"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode error body %q: %v", rec.Body.String(), err)
			}
			if got.Code == "" || got.Message == "" {
				t.Errorf("error envelope missing code/message: %+v", got)
			}
			if len(got.Details) != len(tc.details) {
				t.Fatalf("details = %+v, want %+v", got.Details, tc.details)
			}
			for i, want := range tc.details {
				if got.Details[i] != want {
					t.Errorf("details[%d] = %+v, want %+v", i, got.Details[i], want)
				}
			}

			apitest.AssertMatchesSpec(t, app.newRequest(t, tc.method, tc.path, tc.body, token), rec.Result())
		})
	}
}
