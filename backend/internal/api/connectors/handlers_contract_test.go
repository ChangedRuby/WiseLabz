package connectors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
)

// errorEnvelope mirrors httputil.ErrorResponse with details typed as the
// array docs/openapi.yaml now declares, so the assertions below fail loudly
// if the envelope ever regresses to a free-form object.
type errorEnvelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details []struct {
		Field string `json:"field"`
		Msg   string `json:"msg"`
	} `json:"details"`
}

func decodeEnvelope(t *testing.T, rr *httptest.ResponseRecorder) errorEnvelope {
	t.Helper()
	var got errorEnvelope
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error envelope: %v; body = %s", err, rr.Body.String())
	}
	return got
}

// fieldMsgs flattens the details array into field -> msg for order-independent
// assertions.
func fieldMsgs(e errorEnvelope) map[string]string {
	m := make(map[string]string, len(e.Details))
	for _, d := range e.Details {
		m[d.Field] = d.Msg
	}
	return m
}

func TestCreateValidationDetails(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(`{"name":"svc"}`))
	rr := httptest.NewRecorder()
	h.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
	got := decodeEnvelope(t, rr)
	if got.Code != "invalid_request" {
		t.Errorf("code = %q, want invalid_request", got.Code)
	}
	msgs := fieldMsgs(got)
	// name was supplied; the other three required fields were not.
	for _, field := range []string{"category", "type", "url"} {
		if msgs[field] == "" {
			t.Errorf("details missing an entry for %q; got %+v", field, got.Details)
		}
	}
	if _, ok := msgs["name"]; ok {
		t.Errorf("details reports name, which was supplied; got %+v", got.Details)
	}

	apitest.AssertMatchesSpec(t, req, rr.Result())
}

func TestCreateRotationValidationDetails(t *testing.T) {
	h := newTestHandler(t)

	body := `{"name":"svc","category":"virtualization","type":"custom","url":"https://svc.example.com",` +
		`"userExpiresAt":"not-a-timestamp","rotationMaxAgeDays":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(body))
	rr := httptest.NewRecorder()
	h.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
	msgs := fieldMsgs(decodeEnvelope(t, rr))
	if len(msgs) != 2 || msgs["userExpiresAt"] == "" || msgs["rotationMaxAgeDays"] == "" {
		t.Errorf("details = %+v, want one entry per bad rotation field", msgs)
	}

	apitest.AssertMatchesSpec(t, req, rr.Result())
}

func TestUpdateValidationDetails(t *testing.T) {
	h := newTestHandler(t)

	createRR := httptest.NewRecorder()
	h.Create(createRR, httptest.NewRequest(http.MethodPost, "/api/connectors",
		strings.NewReader(`{"name":"Original","category":"virtualization","type":"custom","url":"https://a.example.com"}`)))
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("no id in create response: %s", createRR.Body.String())
	}

	req := httptest.NewRequest(http.MethodPut, "/api/connectors/"+id,
		strings.NewReader(`{"userExpiresAt":"nope"}`))
	req.SetPathValue("id", id)
	rr := httptest.NewRecorder()
	h.Update(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
	if msgs := fieldMsgs(decodeEnvelope(t, rr)); msgs["userExpiresAt"] == "" {
		t.Errorf("details = %+v, want an entry for userExpiresAt", msgs)
	}

	apitest.AssertMatchesSpec(t, req, rr.Result())
}

// TestGetNotFoundMatchesSpec covers a non-validation response through the same
// helper, so the 404 envelope is held to the spec too.
func TestGetNotFoundMatchesSpec(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/connectors/missing", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Get(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	apitest.AssertMatchesSpec(t, req, rr.Result())
}
