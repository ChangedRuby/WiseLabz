package runbooks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/api/connectors"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"

	// Register connector implementations (proxmox, etc.), same reasoning as
	// internal/api/connectors' own handlers_test.go: SupportsLifecycleVerb
	// probes the registry, which is empty without this.
	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	s := apitest.NewStore(t)
	cfg := &config.Config{Encryption: config.EncryptionSettings{Key: "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="}}
	jwtSvc := auth.NewService("test-secret-test-secret-test-secret", time.Hour, time.Hour)
	connH := connectors.NewHandler(s, sync.NewEngine(s, nil, nil, nil, cfg.Encryption.Key), cfg, jwtSvc, nil)
	return NewHandler(s, connH)
}

func TestListMutuallyExclusiveFilters(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/runbooks?changeType=deploy&alertSeverity=high", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestGetNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/runbooks/missing", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestCreate(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(`{"title":""}`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid target type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(`{"title":"t","targetType":"bogus","targetValue":"v"}`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("happy path then conflict on duplicate target", func(t *testing.T) {
		body := `{"title":"Deploy runbook","body":"steps","targetType":"change_type","targetValue":"deploy"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}

		dupReq := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(body))
		dupRR := httptest.NewRecorder()
		h.Create(dupRR, dupReq)
		if dupRR.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d; body=%s", dupRR.Code, http.StatusConflict, dupRR.Body.String())
		}
	})
}

func TestUpdateAndDelete(t *testing.T) {
	h := newTestHandler(t)

	createReq := httptest.NewRequest(http.MethodPost, "/api/runbooks",
		strings.NewReader(`{"title":"Original","body":"b","targetType":"alert_severity","targetValue":"high"}`))
	createRR := httptest.NewRecorder()
	h.Create(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("no id in create response: %s", createRR.Body.String())
	}

	t.Run("update invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/"+id, strings.NewReader(`{`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("update invalid field type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/"+id, strings.NewReader(`{"title":123}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("update invalid target type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/"+id, strings.NewReader(`{"targetType":"bogus"}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("update not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/missing", strings.NewReader(`{"title":"x"}`))
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("update happy path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/"+id, strings.NewReader(`{"title":"Updated","docId":null}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
	})

	t.Run("delete not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/runbooks/missing", nil)
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Delete(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("delete happy path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/runbooks/"+id, nil)
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Delete(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
		}
	})
}

func seedProxmoxConnector(t *testing.T, h *Handler) string {
	t.Helper()
	c := &store.ConnectorRecord{Name: "Proxmox", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := h.Store.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return c.ID
}

func TestCreateStepsValidation(t *testing.T) {
	h := newTestHandler(t)
	connID := seedProxmoxConnector(t, h)

	cases := []struct {
		name string
		step string
	}{
		{"missing title", `{"connectorId":"` + connID + `","verb":"restart"}`},
		{"invalid verb", `{"title":"t","connectorId":"` + connID + `","verb":"reboot"}`},
		{"unknown connector", `{"title":"t","connectorId":"missing","verb":"restart"}`},
		{"bad entityRef", `{"title":"t","connectorId":"` + connID + `","verb":"restart","entityRef":"../etc"}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body := `{"title":"t","targetType":"change_type","targetValue":"x","steps":[` + c.step + `]}`
			req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(body))
			rr := httptest.NewRecorder()
			h.Create(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
			}
		})
	}

	t.Run("too many steps", func(t *testing.T) {
		steps := "["
		for i := 0; i < 21; i++ {
			if i > 0 {
				steps += ","
			}
			steps += `{"title":"t","connectorId":"` + connID + `","verb":"restart"}`
		}
		steps += "]"
		body := `{"title":"t","targetType":"change_type","targetValue":"y","steps":` + steps + `}`
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
		}
	})

	t.Run("happy path with steps and canExecute false without grant", func(t *testing.T) {
		body := `{"title":"t","targetType":"change_type","targetValue":"z","steps":[{"title":"Restart it","connectorId":"` + connID + `","verb":"restart","entityRef":"100"}]}`
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(body))
		user := apitest.NewUser(t, h.Store, "viewer")
		req = req.WithContext(auth.ContextWithUser(req.Context(), user, false))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body=%s", rr.Code, rr.Body.String())
		}
		var resp runbookResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(resp.Steps) != 1 {
			t.Fatalf("len(steps) = %d, want 1", len(resp.Steps))
		}
		if resp.Steps[0].CanExecute {
			t.Error("canExecute should be false without an operator grant")
		}
		if resp.Steps[0].ExecuteBlockedReason != "no_operator_grant" {
			t.Errorf("executeBlockedReason = %q, want no_operator_grant", resp.Steps[0].ExecuteBlockedReason)
		}
		if resp.Steps[0].ConnectorName != "Proxmox" {
			t.Errorf("connectorName = %q, want Proxmox", resp.Steps[0].ConnectorName)
		}

		apitest.GrantConnectorRole(t, h.Store, user, connID, "operator")
		getReq := httptest.NewRequest(http.MethodGet, "/api/runbooks/"+resp.ID, nil)
		getReq.SetPathValue("id", resp.ID)
		getReq = getReq.WithContext(auth.ContextWithUser(getReq.Context(), user, false))
		getRR := httptest.NewRecorder()
		h.Get(getRR, getReq)
		var got runbookResponse
		if err := json.Unmarshal(getRR.Body.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal get: %v", err)
		}
		if !got.Steps[0].CanExecute {
			t.Error("canExecute should be true after granting operator")
		}
	})
}

func TestExecuteStepForbiddenWithoutOperatorGrant(t *testing.T) {
	h := newTestHandler(t)
	connID := seedProxmoxConnector(t, h)

	body := `{"title":"t","targetType":"change_type","targetValue":"exec","steps":[{"title":"Restart","connectorId":"` + connID + `","verb":"restart"}]}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(body))
	createRR := httptest.NewRecorder()
	h.Create(createRR, createReq)
	if createRR.Code != http.StatusCreated {
		t.Fatalf("create status = %d; body=%s", createRR.Code, createRR.Body.String())
	}
	var created runbookResponse
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	stepID := created.Steps[0].ID

	user := apitest.NewUser(t, h.Store, "viewer") // no connector grant at all
	req := httptest.NewRequest(http.MethodPost, "/api/runbooks/"+created.ID+"/steps/"+stepID+"/execute", nil)
	req.SetPathValue("id", created.ID)
	req.SetPathValue("stepId", stepID)
	req = req.WithContext(auth.ContextWithUser(req.Context(), user, false))
	rr := httptest.NewRecorder()
	h.ExecuteStep(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body=%s", rr.Code, rr.Body.String())
	}
}

func TestExecuteStepNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/runbooks/missing/steps/missing/execute", nil)
	req.SetPathValue("id", "missing")
	req.SetPathValue("stepId", "missing")
	rr := httptest.NewRecorder()
	h.ExecuteStep(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rr.Code, rr.Body.String())
	}
}
