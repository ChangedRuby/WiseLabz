package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedCoverageConnector(t *testing.T, h *Handler, name, category string) *store.ConnectorRecord {
	t.Helper()
	c := &store.ConnectorRecord{Name: name, Category: category, Type: "custom", URL: "https://example.com", ConfigData: "{}"}
	if err := h.Store.CreateConnector(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestListFiltersGrantsBeforePagination(t *testing.T) {
	h := newTestHandler(t)
	user := apitest.NewUser(t, h.Store, "viewer")
	for i, role := range []string{"", "viewer", "operator", "viewer"} {
		category := "networking"
		if i == 3 {
			category = "dns"
		}
		c := &store.ConnectorRecord{Name: fmt.Sprintf("connector-%d", i), Category: category, Type: "custom", URL: "https://example.com", ConfigData: "{}", CreatedAt: fmt.Sprintf("2025-01-%02dT00:00:00Z", 10-i)}
		if err := h.Store.CreateConnector(context.Background(), c); err != nil {
			t.Fatal(err)
		}
		if role != "" {
			apitest.GrantConnectorRole(t, h.Store, user, c.ID, role)
		}
	}
	for _, tc := range []struct {
		query        string
		count        int
		names, roles []string
	}{
		{"?category=networking&pageSize=1&page=1", 2, []string{"connector-1"}, []string{"viewer"}},
		{"?category=networking&pageSize=1&page=2", 2, []string{"connector-2"}, []string{"operator"}},
		{"?category=networking&pageSize=1&page=3", 2, []string{}, []string{}},
		{"?category=dns", 1, []string{"connector-3"}, []string{"viewer"}},
		{"?page=bad&pageSize=-1", 3, []string{"connector-1", "connector-2", "connector-3"}, []string{"viewer", "operator", "viewer"}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+tc.query, nil)
			req = req.WithContext(auth.ContextWithUser(req.Context(), user, false))
			rr := httptest.NewRecorder()
			h.List(rr, req)
			if rr.Code != 200 || rr.Header().Get("X-Total-Count") != fmt.Sprint(tc.count) {
				t.Fatalf("status=%d total=%s body=%s", rr.Code, rr.Header().Get("X-Total-Count"), rr.Body.String())
			}
			var rows []connectorWithRole
			if err := json.Unmarshal(rr.Body.Bytes(), &rows); err != nil {
				t.Fatal(err)
			}
			if rows == nil || len(rows) != len(tc.names) {
				t.Fatalf("rows=%s", rr.Body.String())
			}
			for i, name := range tc.names {
				if rows[i].Name != name || rows[i].MyRole != tc.roles[i] {
					t.Fatalf("row=%+v", rows[i])
				}
			}
		})
	}
}

func TestGetRequiresGrantEvenForInstanceAdmin(t *testing.T) {
	h := newTestHandler(t)
	c := seedCoverageConnector(t, h, "Private", "networking")
	for _, tc := range []struct {
		role  string
		admin bool
		want  int
	}{{"", false, 404}, {"", true, 404}, {"viewer", false, 200}, {"operator", false, 200}, {"viewer", true, 200}} {
		t.Run(fmt.Sprintf("%s/admin=%v", tc.role, tc.admin), func(t *testing.T) {
			userRole := "viewer"
			if tc.admin {
				userRole = "operator"
			}
			user := apitest.NewUser(t, h.Store, userRole)
			if tc.role != "" {
				apitest.GrantConnectorRole(t, h.Store, user, c.ID, tc.role)
			}
			req := httptest.NewRequest("GET", "/", nil)
			req.SetPathValue("id", c.ID)
			req = req.WithContext(auth.ContextWithUser(req.Context(), user, tc.admin))
			rr := httptest.NewRecorder()
			h.Get(rr, req)
			if rr.Code != tc.want {
				t.Fatalf("status=%d: %s", rr.Code, rr.Body.String())
			}
			if tc.want == 200 {
				var row connectorWithRole
				if err := json.Unmarshal(rr.Body.Bytes(), &row); err != nil {
					t.Fatal(err)
				}
				if row.ID != c.ID || row.MyRole != tc.role {
					t.Fatalf("row=%+v", row)
				}
			}
		})
	}
}

func TestSyncsLimitAndIsolation(t *testing.T) {
	h := newTestHandler(t)
	c := seedCoverageConnector(t, h, "History", "networking")
	other := seedCoverageConnector(t, h, "Other", "networking")
	ctx := context.Background()
	for i := 0; i < 105; i++ {
		run := &store.SyncRunRecord{ConnectorID: c.ID, Status: store.SyncRunStatusSuccess, StartedAt: time.Date(2025, 1, 1, 0, i, 0, 0, time.UTC).Format(time.RFC3339)}
		if err := h.Store.CreateSyncRun(ctx, run); err != nil {
			t.Fatal(err)
		}
	}
	if err := h.Store.CreateSyncRun(ctx, &store.SyncRunRecord{ConnectorID: other.ID, Status: store.SyncRunStatusError, StartedAt: "2026-01-01T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		limit string
		want  int
	}{{"", 20}, {"2", 2}, {"bad", 20}, {"-1", 20}, {"0", 20}, {"1000", 100}} {
		t.Run(tc.limit, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/?limit="+tc.limit, nil)
			req.SetPathValue("id", c.ID)
			rr := httptest.NewRecorder()
			h.Syncs(rr, req)
			if rr.Code != 200 {
				t.Fatalf("status=%d: %s", rr.Code, rr.Body.String())
			}
			var runs []store.SyncRunRecord
			if err := json.Unmarshal(rr.Body.Bytes(), &runs); err != nil {
				t.Fatal(err)
			}
			if len(runs) != tc.want {
				t.Fatalf("count=%d, want %d", len(runs), tc.want)
			}
			for i, run := range runs {
				want := time.Date(2025, 1, 1, 0, 104-i, 0, 0, time.UTC).Format(time.RFC3339)
				if run.ConnectorID != c.ID || run.StartedAt != want {
					t.Fatalf("run %d=%+v, want time %s", i, run, want)
				}
			}
		})
	}
}

func TestDataSnapshotPaths(t *testing.T) {
	for _, tc := range []struct {
		name, data string
		status     int
	}{
		{"no snapshot", "", 200},
		{"snapshot", `{"serviceName":"Fixture","type":"custom","sections":[{"title":"Overview","content":"body"}],"metadata":{"region":"test"},"fetchedAt":"2025-01-01T00:00:00Z"}`, 200},
		{"corrupt snapshot", "{", 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t)
			c := seedCoverageConnector(t, h, "Fixture", "networking")
			if tc.data != "" {
				if err := h.Store.CreateSnapshot(context.Background(), &store.SnapshotRecord{ConnectorID: c.ID, Data: tc.data, FetchedAt: "2025-01-01T00:00:00Z"}); err != nil {
					t.Fatal(err)
				}
			}
			req := httptest.NewRequest("GET", "/", nil)
			req.SetPathValue("id", c.ID)
			rr := httptest.NewRecorder()
			h.Data(rr, req)
			if rr.Code != tc.status {
				t.Fatalf("status=%d: %s", rr.Code, rr.Body.String())
			}
			if tc.status != 200 {
				return
			}
			var result struct {
				ServiceName string `json:"serviceName"`
				Sections    []struct {
					Title, Content string
					Order          int
				}
				Metadata  map[string]string
				FetchedAt string `json:"fetchedAt"`
			}
			if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if tc.data == "" {
				if result.Sections == nil || len(result.Sections) != 0 || result.Metadata == nil || result.ServiceName != "" {
					t.Fatalf("empty data=%s", rr.Body.String())
				}
				return
			}
			if result.ServiceName != "Fixture" || len(result.Sections) != 1 || result.Sections[0].Content != "body" || result.Sections[0].Title != "Overview" || result.Sections[0].Order != 0 || result.Metadata["region"] != "test" || result.FetchedAt != "2025-01-01T00:00:00Z" {
				t.Fatalf("data=%s", rr.Body.String())
			}
		})
	}
}

func TestConnectorStoreErrorPaths(t *testing.T) {
	h := newTestHandler(t)
	for name, fn := range map[string]http.HandlerFunc{"list": h.List, "get": h.Get, "create": h.Create, "update": h.Update, "delete": h.Delete, "test": h.Test, "health": h.Health, "data": h.Data, "syncs": h.Syncs, "removal": h.RemovalImpact, "config fields": h.ConfigFields, "toggle": h.ToggleEnabled} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Test","category":"networking","type":"custom","url":"https://example.com","enabled":false}`)).WithContext(ctx)
			req.SetPathValue("id", "missing")
			req.Header.Set("X-Elevation-Token", "present")
			req = req.WithContext(auth.ContextWithUser(req.Context(), "", true))
			rr := httptest.NewRecorder()
			fn(rr, req)
			if rr.Code != 500 {
				t.Fatalf("status=%d: %s", rr.Code, rr.Body.String())
			}
		})
	}
	for name, fn := range map[string]http.HandlerFunc{"health": h.Health, "data": h.Data, "syncs": h.Syncs, "config fields": h.ConfigFields} {
		t.Run(name+" missing", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.SetPathValue("id", "missing")
			rr := httptest.NewRecorder()
			fn(rr, req)
			if rr.Code != 404 {
				t.Fatalf("status=%d: %s", rr.Code, rr.Body.String())
			}
		})
	}
}
