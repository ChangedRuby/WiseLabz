package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/store"
)

type dashboardLayout struct {
	Widgets []map[string]any `json:"widgets"`
}

// operatorWithPermission seeds an operator with can_manage_dashboard_defaults
// set, bypassing testApp.user (which only sets role) since this permission is
// a separate boolean flag layered on top of the operator role.
func (a *testApp) operatorWithPermission(t *testing.T, canManage bool) (userID, accessToken string) {
	t.Helper()

	hash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	u := &store.User{
		Username:                   "operator-" + uuid.New().String(),
		DisplayName:                "Test Operator",
		InstanceAdminRole:          "admin",
		AuthSource:                 "local",
		PasswordHash:               hash,
		CanManageDashboardDefaults: canManage,
	}
	if err := a.Store.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	pair, err := a.JWT.IssuePair(u.ID, u.InstanceAdminRole == "admin")
	if err != nil {
		t.Fatalf("issue pair: %v", err)
	}
	return u.ID, pair.AccessToken
}

func TestDashboardLayoutPerUserIsolation(t *testing.T) {
	app := newTestApp(t)
	_, aliceToken := app.user(t, "viewer")
	_, bobToken := app.user(t, "viewer")

	saveRec := app.req(t, http.MethodPut, "/api/dashboard/layout", map[string]any{
		"widgets": []map[string]any{{"id": "alice-widget", "type": "service_status", "x": 0, "y": 0, "w": 1, "h": 1}},
	}, aliceToken)
	if saveRec.Code != http.StatusOK {
		t.Fatalf("viewer save layout status = %d, want 200; body = %s", saveRec.Code, saveRec.Body)
	}

	aliceGet := app.req(t, http.MethodGet, "/api/dashboard/layout", nil, aliceToken)
	var aliceLayout dashboardLayout
	if err := json.Unmarshal(aliceGet.Body.Bytes(), &aliceLayout); err != nil {
		t.Fatalf("decode alice layout: %v", err)
	}
	if len(aliceLayout.Widgets) != 1 || aliceLayout.Widgets[0]["id"] != "alice-widget" {
		t.Fatalf("alice layout = %+v, want her saved widget", aliceLayout)
	}

	// Bob never saved a layout — he should see the administrator default,
	// not Alice's widget and not an empty layout.
	bobGet := app.req(t, http.MethodGet, "/api/dashboard/layout", nil, bobToken)
	var bobLayout dashboardLayout
	if err := json.Unmarshal(bobGet.Body.Bytes(), &bobLayout); err != nil {
		t.Fatalf("decode bob layout: %v", err)
	}
	if len(bobLayout.Widgets) == 0 {
		t.Fatalf("bob layout = %+v, want the seeded admin default (non-empty)", bobLayout)
	}
	for _, w := range bobLayout.Widgets {
		if w["id"] == "alice-widget" {
			t.Fatalf("bob layout leaked alice's widget: %+v", bobLayout)
		}
	}
}

func TestDashboardAdminDefaultPermissionGate(t *testing.T) {
	app := newTestApp(t)
	_, plainOperatorToken := app.operatorWithPermission(t, false)
	_, permittedOperatorToken := app.operatorWithPermission(t, true)
	_, viewerToken := app.user(t, "viewer")

	for _, tok := range []string{viewerToken, plainOperatorToken} {
		rec := app.req(t, http.MethodGet, "/api/dashboard/layout/admin-default", nil, tok)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("GET admin-default without permission status = %d, want 403; body = %s", rec.Code, rec.Body)
		}
		rec = app.req(t, http.MethodPut, "/api/dashboard/layout/admin-default", map[string]any{"widgets": []map[string]any{}}, tok)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("PUT admin-default without permission status = %d, want 403; body = %s", rec.Code, rec.Body)
		}
	}

	getRec := app.req(t, http.MethodGet, "/api/dashboard/layout/admin-default", nil, permittedOperatorToken)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET admin-default with permission status = %d, want 200; body = %s", getRec.Code, getRec.Body)
	}

	newDefault := []map[string]any{{"id": "org-widget", "type": "alert_summary", "x": 0, "y": 0, "w": 2, "h": 1}}
	putRec := app.req(t, http.MethodPut, "/api/dashboard/layout/admin-default", map[string]any{"widgets": newDefault}, permittedOperatorToken)
	if putRec.Code != http.StatusOK {
		t.Fatalf("PUT admin-default with permission status = %d, want 200; body = %s", putRec.Code, putRec.Body)
	}

	// A brand-new user now gets the updated default.
	_, newUserToken := app.user(t, "viewer")
	newUserGet := app.req(t, http.MethodGet, "/api/dashboard/layout", nil, newUserToken)
	var newUserLayout dashboardLayout
	if err := json.Unmarshal(newUserGet.Body.Bytes(), &newUserLayout); err != nil {
		t.Fatalf("decode new user layout: %v", err)
	}
	if len(newUserLayout.Widgets) != 1 || newUserLayout.Widgets[0]["id"] != "org-widget" {
		t.Fatalf("new user layout = %+v, want the updated admin default", newUserLayout)
	}
}

func TestDashboardResetRestoresAdminDefault(t *testing.T) {
	app := newTestApp(t)
	_, userToken := app.user(t, "viewer")

	// Save a custom layout.
	saveRec := app.req(t, http.MethodPut, "/api/dashboard/layout", map[string]any{
		"widgets": []map[string]any{{"id": "custom-widget", "type": "docs_health", "x": 0, "y": 0, "w": 1, "h": 1}},
	}, userToken)
	if saveRec.Code != http.StatusOK {
		t.Fatalf("save custom layout status = %d, want 200; body = %s", saveRec.Code, saveRec.Body)
	}

	// Capture the current admin default to compare against after reset.
	_, permittedOperatorToken := app.operatorWithPermission(t, true)
	defaultGet := app.req(t, http.MethodGet, "/api/dashboard/layout/admin-default", nil, permittedOperatorToken)
	var adminDefault dashboardLayout
	if err := json.Unmarshal(defaultGet.Body.Bytes(), &adminDefault); err != nil {
		t.Fatalf("decode admin default: %v", err)
	}

	resetRec := app.req(t, http.MethodPost, "/api/dashboard/layout/reset", nil, userToken)
	if resetRec.Code != http.StatusOK {
		t.Fatalf("reset status = %d, want 200; body = %s", resetRec.Code, resetRec.Body)
	}
	var resetLayout dashboardLayout
	if err := json.Unmarshal(resetRec.Body.Bytes(), &resetLayout); err != nil {
		t.Fatalf("decode reset response: %v", err)
	}
	if len(resetLayout.Widgets) != len(adminDefault.Widgets) {
		t.Fatalf("reset layout = %+v, want it to match the admin default %+v", resetLayout, adminDefault)
	}

	// Verify persisted, not just returned.
	afterGet := app.req(t, http.MethodGet, "/api/dashboard/layout", nil, userToken)
	var afterLayout dashboardLayout
	if err := json.Unmarshal(afterGet.Body.Bytes(), &afterLayout); err != nil {
		t.Fatalf("decode layout after reset: %v", err)
	}
	for _, w := range afterLayout.Widgets {
		if w["id"] == "custom-widget" {
			t.Fatalf("layout after reset still has the custom widget: %+v", afterLayout)
		}
	}
}

func TestDashboardOverviewDaysWindow(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	recent := &store.ChangeRecord{
		ServiceID: "svc-1", ChangeType: "config", Severity: "info",
		Summary: "recent change", Diff: "[]", AffectedDocIDs: "[]",
		DetectedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := app.Store.CreateChange(context.Background(), recent); err != nil {
		t.Fatalf("seed recent change: %v", err)
	}

	old := &store.ChangeRecord{
		ServiceID: "svc-1", ChangeType: "config", Severity: "info",
		Summary: "old change", Diff: "[]", AffectedDocIDs: "[]",
		DetectedAt: time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339),
	}
	if err := app.Store.CreateChange(context.Background(), old); err != nil {
		t.Fatalf("seed old change: %v", err)
	}

	type overviewResp struct {
		RecentChanges []struct {
			ID string `json:"id"`
		} `json:"recentChanges"`
	}

	narrowRec := app.req(t, http.MethodGet, "/api/dashboard/overview?days=1", nil, viewerToken)
	if narrowRec.Code != http.StatusOK {
		t.Fatalf("days=1 status = %d, want 200; body = %s", narrowRec.Code, narrowRec.Body)
	}
	var narrow overviewResp
	if err := json.Unmarshal(narrowRec.Body.Bytes(), &narrow); err != nil {
		t.Fatalf("decode days=1 body: %v", err)
	}
	for _, c := range narrow.RecentChanges {
		if c.ID == old.ID {
			t.Errorf("days=1 recentChanges included the old change: %+v", narrow.RecentChanges)
		}
	}
	foundRecent := false
	for _, c := range narrow.RecentChanges {
		if c.ID == recent.ID {
			foundRecent = true
		}
	}
	if !foundRecent {
		t.Errorf("days=1 recentChanges missing the recent change: %+v", narrow.RecentChanges)
	}

	wideRec := app.req(t, http.MethodGet, "/api/dashboard/overview?days=90", nil, viewerToken)
	if wideRec.Code != http.StatusOK {
		t.Fatalf("days=90 status = %d, want 200; body = %s", wideRec.Code, wideRec.Body)
	}
	var wide overviewResp
	if err := json.Unmarshal(wideRec.Body.Bytes(), &wide); err != nil {
		t.Fatalf("decode days=90 body: %v", err)
	}
	foundOld := false
	for _, c := range wide.RecentChanges {
		if c.ID == old.ID {
			foundOld = true
		}
	}
	if !foundOld {
		t.Errorf("days=90 recentChanges missing the old change: %+v", wide.RecentChanges)
	}
}
