package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/connector"
	syncengine "github.com/WiseLabz/wiselabz/internal/sync"
)

// The engine calls this seam after its final database write, so tests can
// inspect persisted results and clean up without racing a detached sync.
type completedSync chan string

func (done completedSync) RegenerateForConnector(_ context.Context, id string) error {
	done <- id
	return nil
}

type actionConnector struct {
	bulkFakeConnector
	fields  chan []string
	release <-chan struct{}
	fail    bool
}

func (c *actionConnector) Fetch(ctx context.Context, cfg map[string]any) (*connector.ServiceSnapshot, error) {
	<-c.release
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	fields, _ := cfg["fields"].([]string)
	c.fields <- fields
	if c.fail {
		return nil, errors.New("fixture fetch failed")
	}
	return &connector.ServiceSnapshot{ServiceName: "async fixture", FetchedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}, nil
}

func TestActionAsyncSync(t *testing.T) {
	for _, mode := range []string{"single empty", "single fields", "single failure", "bulk", "all"} {
		t.Run(mode, func(t *testing.T) {
			h := newTestHandler(t)
			fields := make(chan []string, 1)
			release := make(chan struct{})
			done := make(completedSync, 1)
			connector.Register(connector.TypeSchema{Type: "action_async", Name: "Async", Category: "networking"}, func(map[string]any) (connector.Connector, error) {
				return &actionConnector{fields: fields, release: release, fail: mode == "single failure"}, nil
			})
			h.SyncEngine = syncengine.NewEngine(h.Store, nil, nil, nil, h.Config.Encryption.Key)
			h.SyncEngine.SetDocRegenerator(done)
			c := seedCoverageConnector(t, h, "async", "networking")
			if err := h.Store.UpdateConnector(context.Background(), c.ID, map[string]any{"type": "action_async", "enabled": true}); err != nil {
				t.Fatal(err)
			}
			user := apitest.NewUser(t, h.Store, "operator")
			apitest.GrantConnectorRole(t, h.Store, user, c.ID, "operator")
			fn := http.HandlerFunc(h.Sync)
			body := ""
			action := "connector.sync"
			status := 202
			switch mode {
			case "single fields":
				body = `{"fields":["vms","storage"]}`
			case "bulk":
				fn = h.BulkSync
				body = `{"ids":["` + c.ID + `"]}`
				action = "connector.bulk_sync"
				status = 200
			case "all":
				fn = h.SyncAll
				action = "connector.sync_all"
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			r := actionRequest(c.ID, body).WithContext(auth.ContextWithUser(ctx, user, false))
			rr := actionResponse(t, fn, r, status)
			cancel() // detached work must outlive the initiating request
			close(release)
			jobID := ""
			if mode == "bulk" {
				items := bulkResults(t, rr)
				if len(items) != 1 || items[0].ID != c.ID || items[0].Status != "success" {
					t.Fatalf("results=%+v", items)
				}
				jobID = items[0].JobID
			} else {
				var result struct {
					JobID     string  `json:"jobId"`
					ServiceID *string `json:"serviceId"`
				}
				if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				jobID = result.JobID
				if mode == "all" {
					if result.ServiceID != nil {
						t.Fatalf("serviceId=%v", result.ServiceID)
					}
				} else if result.ServiceID == nil || *result.ServiceID != c.ID {
					t.Fatalf("response=%s", rr.Body.String())
				}
			}
			if jobID == "" {
				t.Fatal("missing job ID")
			}
			select {
			case id := <-done:
				if id != c.ID {
					t.Fatalf("completed=%s", id)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("sync did not complete")
			}
			var gotFields []string
			select {
			case gotFields = <-fields:
			case <-time.After(time.Second):
				t.Fatal("fetch was not called")
			}
			var wantFields []string
			if mode == "single fields" {
				wantFields = []string{"vms", "storage"}
			}
			if !reflect.DeepEqual(gotFields, wantFields) {
				t.Fatalf("fields=%v, want %v", gotFields, wantFields)
			}
			runs, err := h.Store.ListSyncRunsByConnector(context.Background(), c.ID, 10)
			wantStatus := "success"
			if mode == "single failure" {
				wantStatus = "error"
			}
			if err != nil || len(runs) != 1 || string(runs[0].Status) != wantStatus {
				t.Fatalf("runs=%+v err=%v", runs, err)
			}
			rows, _, err := h.Store.ListAuditRecords(context.Background(), action, "connector", "", "", 0, 10)
			if err != nil || len(rows) != 1 {
				t.Fatalf("audit=%+v err=%v", rows, err)
			}
			var detail struct {
				JobID string `json:"jobId"`
			}
			if err := json.Unmarshal([]byte(rows[0].Detail), &detail); err != nil {
				t.Fatal(err)
			}
			if detail.JobID != jobID {
				t.Fatalf("audit job=%s want %s", detail.JobID, jobID)
			}
		})
	}
}

type failingPushConnector struct {
	bulkFakeConnector
	mode            string
	fetches, pushes int
}

func (c *failingPushConnector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{{Key: "memory"}}
}
func (c *failingPushConnector) Fetch(context.Context, map[string]any) (*connector.ServiceSnapshot, error) {
	c.fetches++
	if c.mode == "pre" || c.mode == "post" && c.fetches == 2 {
		return nil, errors.New("fixture fetch failed")
	}
	return &connector.ServiceSnapshot{ServiceName: "push fixture"}, nil
}
func (c *failingPushConnector) ConfigPush(context.Context, map[string]any, string, string, any) error {
	c.pushes++
	if c.mode == "push" {
		return errors.New("fixture push failed")
	}
	return nil
}

func TestActionConfigPushFailures(t *testing.T) {
	for _, mode := range []string{"pre", "push", "post"} {
		t.Run(mode, func(t *testing.T) {
			h := newTestHandler(t)
			fake := &failingPushConnector{mode: mode}
			connector.Register(connector.TypeSchema{Type: "action_push", Name: "Push", Category: "networking"}, func(map[string]any) (connector.Connector, error) { return fake, nil })
			c := seedCoverageConnector(t, h, "push", "networking")
			if err := h.Store.UpdateConnector(context.Background(), c.ID, map[string]any{"type": "action_push"}); err != nil {
				t.Fatal(err)
			}
			token, err := h.JWT.IssueElevation("", "connector.configPush")
			if err != nil {
				t.Fatal(err)
			}
			r := actionRequest(c.ID, `{"entityRef":"100","fieldKey":"memory","value":4096,"previousValue":2048}`)
			r.Header.Set("X-Elevation-Token", token.Token)
			status := 500
			if mode == "push" {
				status = 502
			}
			rr := actionResponse(t, h.ConfigPush, r, status)
			if mode == "push" && !strings.Contains(rr.Body.String(), "config_push_failed") {
				t.Fatalf("response=%s", rr.Body.String())
			}
			wantFetches, wantPushes := 1, 1
			if mode == "pre" {
				wantPushes = 0
			}
			if mode == "post" {
				wantFetches = 2
			}
			if fake.fetches != wantFetches || fake.pushes != wantPushes {
				t.Fatalf("fetches=%d pushes=%d", fake.fetches, fake.pushes)
			}
			rows, _, err := h.Store.ListAuditRecords(context.Background(), "connector.configPush", "connector", "", "", 0, 10)
			if err != nil || len(rows) != 0 {
				t.Fatalf("failed push audit=%+v err=%v", rows, err)
			}
		})
	}
}
