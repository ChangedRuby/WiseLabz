package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestSavedViewsCreateListDelete(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	createRec := app.req(t, http.MethodPost, "/api/saved-views", map[string]any{
		"surface": "services",
		"name":    "Down services",
		"filters": map[string]any{"status": "down"},
	}, viewerToken)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", createRec.Code, createRec.Body)
	}
	var created store.SavedView
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create body: %v", err)
	}
	if created.Name != "Down services" || created.Surface != "services" {
		t.Fatalf("created view = %+v, want name/surface set", created)
	}

	listRec := app.req(t, http.MethodGet, "/api/saved-views?surface=services", nil, viewerToken)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200; body = %s", listRec.Code, listRec.Body)
	}
	var views []store.SavedView
	if err := json.Unmarshal(listRec.Body.Bytes(), &views); err != nil {
		t.Fatalf("decode list body: %v", err)
	}
	if len(views) != 1 || views[0].ID != created.ID {
		t.Fatalf("list = %+v, want just the created view", views)
	}

	deleteRec := app.req(t, http.MethodDelete, "/api/saved-views/"+created.ID, nil, viewerToken)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204; body = %s", deleteRec.Code, deleteRec.Body)
	}

	listAfter := app.req(t, http.MethodGet, "/api/saved-views?surface=services", nil, viewerToken)
	var viewsAfter []store.SavedView
	if err := json.Unmarshal(listAfter.Body.Bytes(), &viewsAfter); err != nil {
		t.Fatalf("decode list-after body: %v", err)
	}
	if len(viewsAfter) != 0 {
		t.Fatalf("list after delete = %+v, want empty", viewsAfter)
	}
}

func TestSavedViewsRequireValidSurface(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, token := app.user(t, "viewer")

	rec := app.req(t, http.MethodPost, "/api/saved-views", map[string]any{
		"surface": "not-a-real-surface",
		"name":    "bad view",
	}, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("create with bad surface status = %d, want 400; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodGet, "/api/saved-views?surface=not-a-real-surface", nil, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("list with bad surface status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

// TestSavedViewsOwnershipBoundary ensures a user can never delete another
// user's saved view — same spirit as the audit trail's role-boundary test.
func TestSavedViewsOwnershipBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, ownerToken := app.user(t, "viewer")
	_, otherToken := app.user(t, "operator")

	createRec := app.req(t, http.MethodPost, "/api/saved-views", map[string]any{
		"surface": "changes",
		"name":    "Critical only",
	}, ownerToken)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", createRec.Code, createRec.Body)
	}
	var created store.SavedView
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create body: %v", err)
	}

	rec := app.req(t, http.MethodDelete, "/api/saved-views/"+created.ID, nil, otherToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("delete by non-owner status = %d, want 403; body = %s", rec.Code, rec.Body)
	}

	// Still there for the owner.
	listRec := app.req(t, http.MethodGet, "/api/saved-views?surface=changes", nil, ownerToken)
	var views []store.SavedView
	if err := json.Unmarshal(listRec.Body.Bytes(), &views); err != nil {
		t.Fatalf("decode list body: %v", err)
	}
	if len(views) != 1 {
		t.Fatalf("owner's list after forbidden delete = %+v, want the view still present", views)
	}
}

func TestSavedViewsDeleteMissingIsNotFound(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, token := app.user(t, "viewer")

	rec := app.req(t, http.MethodDelete, "/api/saved-views/does-not-exist", nil, token)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}
