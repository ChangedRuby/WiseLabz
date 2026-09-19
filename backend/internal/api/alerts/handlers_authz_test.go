package alerts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/store"
)

type fixture struct {
	s *store.Store
	h *Handler
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	s := apitest.NewStore(t)
	return &fixture{s: s, h: NewHandler(s)}
}

func (f *fixture) connector(t *testing.T, name string) string {
	t.Helper()
	c := &store.ConnectorRecord{Name: name, Category: "networking", Type: "generic", Enabled: true}
	if err := f.s.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	return c.ID
}

func (f *fixture) alert(t *testing.T, serviceID, title string) string {
	t.Helper()
	a := &store.AlertRecord{ServiceID: serviceID, Severity: "warning", Title: title, Status: "pending"}
	if err := f.s.CreateAlert(context.Background(), a); err != nil {
		t.Fatalf("create alert: %v", err)
	}
	return a.ID
}

// call runs fn through the real auth middleware as a fresh user holding the
// given per-connector grants (connectorID -> role).
func (f *fixture) call(t *testing.T, method, target, body, id string, fn http.HandlerFunc, grants map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id != "" {
			r.SetPathValue("id", id)
		}
		fn(w, r)
	})
	uid, tok, h := apitest.AuthedUser(t, f.s, "viewer", inner)
	for cid, role := range grants {
		apitest.GrantConnectorRole(t, f.s, uid, cid, role)
	}
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestGetAuthz(t *testing.T) {
	f := newFixture(t)
	cid := f.connector(t, "svc")
	id := f.alert(t, cid, "a")

	tests := []struct {
		name  string
		grant map[string]string
		want  int
	}{
		{"no grant hides alert", nil, http.StatusNotFound},
		{"viewer", map[string]string{cid: "viewer"}, http.StatusOK},
		{"operator", map[string]string{cid: "operator"}, http.StatusOK},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := f.call(t, http.MethodGet, "/x", "", id, f.h.Get, tc.grant)
			if rr.Code != tc.want {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tc.want, rr.Body.String())
			}
		})
	}
}

func TestMutationAuthz(t *testing.T) {
	f := newFixture(t)
	cid := f.connector(t, "svc")
	until := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)

	actions := []struct {
		name       string
		fn         http.HandlerFunc
		body       string
		wantStatus string
	}{
		{"resolve", f.h.Resolve, "", "resolved"},
		{"dismiss", f.h.Dismiss, "", "dismissed"},
		{"snooze", f.h.Snooze, `{"until":"` + until + `"}`, "snoozed"},
	}
	roles := []struct {
		name  string
		grant map[string]string
		want  int
	}{
		{"no grant", nil, http.StatusForbidden},
		{"viewer", map[string]string{cid: "viewer"}, http.StatusForbidden},
		{"operator", map[string]string{cid: "operator"}, http.StatusOK},
	}
	for _, act := range actions {
		for i, role := range roles {
			t.Run(act.name+"/"+role.name, func(t *testing.T) {
				id := f.alert(t, cid, act.name+strconv.Itoa(i))
				rr := f.call(t, http.MethodPost, "/x", act.body, id, act.fn, role.grant)
				if rr.Code != role.want {
					t.Fatalf("status = %d, want %d; body=%s", rr.Code, role.want, rr.Body.String())
				}
				got, err := f.s.GetAlert(context.Background(), id)
				if err != nil {
					t.Fatal(err)
				}
				want := "pending"
				if role.want == http.StatusOK {
					want = act.wantStatus
				}
				if got.Status != want {
					t.Errorf("stored status = %q, want %q", got.Status, want)
				}
			})
		}
	}
}

func TestListFiltersByGrantAndPaginates(t *testing.T) {
	f := newFixture(t)
	visible := f.connector(t, "visible")
	hidden := f.connector(t, "hidden")
	for i := 0; i < 3; i++ {
		f.alert(t, visible, "v"+strconv.Itoa(i))
	}
	f.alert(t, hidden, "h")

	type listResp struct {
		Items    []store.AlertRecord `json:"items"`
		Total    int                 `json:"total"`
		Page     int                 `json:"page"`
		PageSize int                 `json:"pageSize"`
	}
	list := func(query string, grants map[string]string) listResp {
		t.Helper()
		rr := f.call(t, http.MethodGet, "/api/alerts"+query, "", "", f.h.List, grants)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d; body=%s", rr.Code, rr.Body.String())
		}
		var out listResp
		if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
		return out
	}
	g := map[string]string{visible: "viewer"}

	got := list("", g)
	if len(got.Items) != 3 {
		t.Fatalf("items = %d, want 3 (hidden connector filtered)", len(got.Items))
	}
	for _, a := range got.Items {
		if a.ServiceID != visible {
			t.Errorf("leaked alert from connector %s", a.ServiceID)
		}
	}
	if got := list("", nil); len(got.Items) != 0 {
		t.Errorf("ungranted user saw %d alerts", len(got.Items))
	}

	p1 := list("?serviceId="+visible+"&pageSize=2&page=1", g)
	p2 := list("?serviceId="+visible+"&pageSize=2&page=2", g)
	if len(p1.Items) != 2 || len(p2.Items) != 1 || p1.Total != 3 || p2.Page != 2 {
		t.Fatalf("p1=%d p2=%d total=%d page=%d, want 2/1/3/2", len(p1.Items), len(p2.Items), p1.Total, p2.Page)
	}
	if p1.Items[0].ID == p2.Items[0].ID || p1.Items[1].ID == p2.Items[0].ID {
		t.Error("pages overlap")
	}
	if got := list("?pageSize=9999", g); got.PageSize != 100 {
		t.Errorf("pageSize = %d, want clamp to 100", got.PageSize)
	}
	if got := list("?serviceId="+visible+"&status=resolved", g); len(got.Items) != 0 {
		t.Errorf("status filter returned %d", len(got.Items))
	}
	if got := list("?serviceId="+visible+"&severity=warning", g); len(got.Items) != 3 {
		t.Errorf("severity filter returned %d, want 3", len(got.Items))
	}
}

func TestBulkSnoozeAuthzPerItem(t *testing.T) {
	f := newFixture(t)
	opConn := f.connector(t, "op")
	viewConn := f.connector(t, "view")
	okID := f.alert(t, opConn, "ok")
	deniedID := f.alert(t, viewConn, "denied")
	until := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	body := `{"ids":["` + okID + `","` + deniedID + `","missing"],"until":"` + until + `"}`

	rr := f.call(t, http.MethodPost, "/x", body, "", f.h.BulkSnooze, map[string]string{opConn: "operator", viewConn: "viewer"})
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rr.Code, rr.Body.String())
	}
	var resp struct {
		Results []bulkSnoozeItemResult `json:"results"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	got := map[string]bulkSnoozeItemResult{}
	for _, r := range resp.Results {
		got[r.ID] = r
	}
	if got[okID].Status != "success" {
		t.Errorf("operator item = %+v", got[okID])
	}
	if got[deniedID].Reason != "forbidden" {
		t.Errorf("viewer item = %+v, want forbidden", got[deniedID])
	}
	if got["missing"].Reason != "not_found" {
		t.Errorf("missing item = %+v, want not_found", got["missing"])
	}

	for id, want := range map[string]string{okID: "snoozed", deniedID: "pending"} {
		a, err := f.s.GetAlert(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}
		if a.Status != want {
			t.Errorf("alert %s status = %q, want %q", id, a.Status, want)
		}
	}
}
