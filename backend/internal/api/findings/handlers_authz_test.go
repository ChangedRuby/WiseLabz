package findings

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

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

// finding seeds an open finding; ruleID keeps otherwise-identical findings
// distinct under the open-finding uniqueness key.
func (f *fixture) finding(t *testing.T, connectorID, ruleID string) string {
	return f.findingOfType(t, connectorID, ruleID, "stale")
}

func (f *fixture) findingOfType(t *testing.T, connectorID, ruleID, checkType string) string {
	t.Helper()
	rec := &store.QualityFindingRecord{ConnectorID: connectorID, RuleID: ruleID, CheckType: checkType, Severity: "warning", Title: ruleID}
	if err := f.s.UpsertQualityFinding(context.Background(), rec); err != nil {
		t.Fatalf("upsert finding: %v", err)
	}
	return rec.ID
}

// call issues a request through the real auth middleware as a fresh user with
// the given per-connector grants (connectorID -> role).
func (f *fixture) call(t *testing.T, method, target, id string, fn http.HandlerFunc, grants map[string]string) *httptest.ResponseRecorder {
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
	req := httptest.NewRequest(method, target, nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestGetAuthz(t *testing.T) {
	f := newFixture(t)
	cid := f.connector(t, "wiki")
	id := f.finding(t, cid, "r0")

	tests := []struct {
		name  string
		grant map[string]string
		want  int
	}{
		{"no grant hides finding", nil, http.StatusNotFound},
		{"viewer", map[string]string{cid: "viewer"}, http.StatusOK},
		{"operator", map[string]string{cid: "operator"}, http.StatusOK},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := f.call(t, http.MethodGet, "/api/findings/"+id, id, f.h.Get, tc.grant)
			if rr.Code != tc.want {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tc.want, rr.Body.String())
			}
			if tc.want == http.StatusOK {
				var body map[string]any
				if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body["connectorName"] != "wiki" || body["checkType"] != "stale" {
					t.Errorf("unexpected body: %v", body)
				}
				if body["docId"] != nil || body["resolvedAt"] != nil {
					t.Errorf("empty optional fields should be null: %v", body)
				}
			}
		})
	}
}

func TestResolveAuthz(t *testing.T) {
	f := newFixture(t)
	cid := f.connector(t, "wiki")

	tests := []struct {
		name  string
		grant map[string]string
		want  int
	}{
		{"no grant forbidden", nil, http.StatusForbidden},
		{"viewer forbidden", map[string]string{cid: "viewer"}, http.StatusForbidden},
		{"operator resolves", map[string]string{cid: "operator"}, http.StatusOK},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := f.finding(t, cid, "check-"+strconv.Itoa(i))
			rr := f.call(t, http.MethodPost, "/api/findings/"+id+"/resolve", id, f.h.Resolve, tc.grant)
			if rr.Code != tc.want {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tc.want, rr.Body.String())
			}
			got, err := f.s.GetQualityFinding(context.Background(), id)
			if err != nil {
				t.Fatal(err)
			}
			wantStatus := "open"
			if tc.want == http.StatusOK {
				wantStatus = "resolved"
			}
			if got.Status != wantStatus {
				t.Errorf("stored status = %q, want %q", got.Status, wantStatus)
			}
		})
	}
}

func TestListFiltersByGrantAndPaginates(t *testing.T) {
	f := newFixture(t)
	visible := f.connector(t, "visible")
	hidden := f.connector(t, "hidden")
	for i := 0; i < 3; i++ {
		f.finding(t, visible, "v"+strconv.Itoa(i))
	}
	f.finding(t, hidden, "h0")
	f.findingOfType(t, visible, "e0", "empty")

	type listResp struct {
		Items    []map[string]any `json:"items"`
		Total    int              `json:"total"`
		Page     int              `json:"page"`
		PageSize int              `json:"pageSize"`
	}
	list := func(query string, grants map[string]string) listResp {
		t.Helper()
		rr := f.call(t, http.MethodGet, "/api/findings"+query, "", f.h.List, grants)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d; body=%s", rr.Code, rr.Body.String())
		}
		var out listResp
		if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
		return out
	}

	got := list("", map[string]string{visible: "viewer"})
	if len(got.Items) != 4 {
		t.Fatalf("items = %d, want 4 (hidden connector must be filtered)", len(got.Items))
	}
	for _, it := range got.Items {
		if it["connectorId"] != visible {
			t.Errorf("leaked finding from connector %v", it["connectorId"])
		}
	}
	if got.Page != 1 || got.PageSize != 20 {
		t.Errorf("page=%d pageSize=%d, want 1/20", got.Page, got.PageSize)
	}

	if got := list("", nil); len(got.Items) != 0 {
		t.Errorf("ungranted user saw %d items", len(got.Items))
	}

	// Pagination: page size 2 splits the 5 stored findings across pages, and
	// the visible connector filter is applied to each page.
	p1 := list("?pageSize=2&page=1&connectorId="+visible, map[string]string{visible: "viewer"})
	p2 := list("?pageSize=2&page=2&connectorId="+visible, map[string]string{visible: "viewer"})
	if len(p1.Items) != 2 || len(p2.Items) != 2 || p1.Total != 4 {
		t.Fatalf("p1=%d p2=%d total=%d, want 2/2/4", len(p1.Items), len(p2.Items), p1.Total)
	}
	if p1.Items[0]["id"] == p2.Items[0]["id"] {
		t.Error("pages overlap")
	}

	// Oversized page size is clamped.
	if got := list("?pageSize=100000", map[string]string{visible: "viewer"}); got.PageSize != 100 {
		t.Errorf("pageSize = %d, want clamp to 100", got.PageSize)
	}

	// Filters.
	if got := list("?status=resolved&connectorId="+visible, map[string]string{visible: "viewer"}); len(got.Items) != 0 {
		t.Errorf("status filter returned %d items", len(got.Items))
	}
	if got := list("?checkType=empty&connectorId="+visible, map[string]string{visible: "viewer"}); len(got.Items) != 1 {
		t.Errorf("checkType filter returned %d items, want 1", len(got.Items))
	}
}
