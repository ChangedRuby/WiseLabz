package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedDoc(t *testing.T, app *testApp) *store.DocRecord {
	t.Helper()
	d := &store.DocRecord{Title: "Doc", Kind: "service", ServiceID: "svc-1", Content: "hello"}
	if err := app.Store.CreateDoc(context.Background(), d); err != nil {
		t.Fatalf("seed doc: %v", err)
	}
	return d
}

func TestDocsListAndGetSuccess(t *testing.T) {
	app := newTestApp(t)
	viewerID, viewerToken := app.user(t, "viewer")
	d := seedDoc(t, app)
	app.connectorGrant(t, viewerID, d.ServiceID, "viewer")

	rec := app.req(t, http.MethodGet, "/api/docs", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodGet, "/api/docs/"+d.ID, nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestDocsSaveRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	d := seedDoc(t, app)

	rec := app.req(t, http.MethodPut, "/api/docs/"+d.ID, map[string]any{"content": "x"}, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestDocsSaveValidation(t *testing.T) {
	app := newTestApp(t)
	opID, opToken := app.user(t, "operator")
	d := seedDoc(t, app)
	app.connectorGrant(t, opID, d.ServiceID, "operator")

	rec := app.req(t, http.MethodPut, "/api/docs/"+d.ID, "not-json", opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestDocsSaveSuccess(t *testing.T) {
	app := newTestApp(t)
	opID, opToken := app.user(t, "operator")
	d := seedDoc(t, app)
	app.connectorGrant(t, opID, d.ServiceID, "operator")

	rec := app.req(t, http.MethodPut, "/api/docs/"+d.ID, map[string]any{"content": "updated content"}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	got, err := app.Store.GetDoc(context.Background(), d.ID)
	if err != nil {
		t.Fatalf("GetDoc: %v", err)
	}
	if got.Content != "updated content" {
		t.Errorf("Content = %q, want %q", got.Content, "updated content")
	}
}

func TestDocLockRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	d := seedDoc(t, app)

	rec := app.req(t, http.MethodPost, "/api/docs/"+d.ID+"/lock", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("acquire status = %d, want 403; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodPost, "/api/docs/"+d.ID+"/lock/release", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("release status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestDocLockHappyPath(t *testing.T) {
	app := newTestApp(t)
	opID, opToken := app.user(t, "operator")
	d := seedDoc(t, app)
	app.connectorGrant(t, opID, d.ServiceID, "operator")

	rec := app.req(t, http.MethodGet, "/api/docs/"+d.ID+"/lock", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("get (unlocked) status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var emptyLock map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &emptyLock); err != nil {
		t.Fatalf("unmarshal empty lock: %v", err)
	}
	if len(emptyLock) != 0 {
		t.Fatalf("empty lock = %v, want {}", emptyLock)
	}

	rec = app.req(t, http.MethodPost, "/api/docs/"+d.ID+"/lock", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("acquire status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var acquired store.DocLockRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &acquired); err != nil {
		t.Fatalf("unmarshal acquired lock: %v", err)
	}
	if acquired.DocID != d.ID {
		t.Errorf("acquired.DocID = %q, want %q", acquired.DocID, d.ID)
	}

	rec = app.req(t, http.MethodGet, "/api/docs/"+d.ID+"/lock", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("get (locked) status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodPost, "/api/docs/"+d.ID+"/lock/release", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("release status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodGet, "/api/docs/"+d.ID+"/lock", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("get (after release) status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var releasedLock map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &releasedLock); err != nil {
		t.Fatalf("unmarshal released lock: %v", err)
	}
	if len(releasedLock) != 0 {
		t.Fatalf("released lock = %v, want {}", releasedLock)
	}
}

func TestDocLockConflict(t *testing.T) {
	app := newTestApp(t)
	user1, op1Token := app.user(t, "operator")
	user2, op2Token := app.user(t, "operator")
	d := seedDoc(t, app)
	app.connectorGrant(t, user1, d.ServiceID, "operator")
	app.connectorGrant(t, user2, d.ServiceID, "operator")

	rec := app.req(t, http.MethodPost, "/api/docs/"+d.ID+"/lock", nil, op1Token)
	if rec.Code != http.StatusOK {
		t.Fatalf("user1 acquire status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodPost, "/api/docs/"+d.ID+"/lock", nil, op2Token)
	if rec.Code != http.StatusConflict {
		t.Fatalf("user2 acquire status = %d, want 409; body = %s", rec.Code, rec.Body)
	}
	var conflict store.DocLockRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &conflict); err != nil {
		t.Fatalf("unmarshal conflict response: %v", err)
	}
	if conflict.UserID != user1 {
		t.Errorf("conflict holder = %q, want %q", conflict.UserID, user1)
	}
}
