package notifications

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// newTestStore spins up a fresh, migrated in-memory sqlite store, mirroring
// the pattern used by internal/store's own tests (see doc_test.go).
func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := store.RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	return store.New(db, "sqlite")
}

func seedDelivery(t *testing.T, s *store.Store, status store.DeliveryStatus) *store.DeliveryRecord {
	t.Helper()
	d := &store.DeliveryRecord{
		NotificationID: "notif-1",
		Channel:        "smtp",
		Status:         status,
		Attempts:       1,
	}
	if err := s.CreateDelivery(context.Background(), d); err != nil {
		t.Fatalf("seed delivery: %v", err)
	}
	return d
}

func decodePaginated(t *testing.T, rec *httptest.ResponseRecorder) httputil.PaginatedResponse {
	t.Helper()
	var out httputil.PaginatedResponse
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return out
}

func TestListDeliveriesNoFilter(t *testing.T) {
	s := newTestStore(t)
	h := NewHandler(s)

	seedDelivery(t, s, store.DeliveryStatusSent)
	seedDelivery(t, s, store.DeliveryStatusFailed)

	r := httptest.NewRequest(http.MethodGet, "/api/notifications/deliveries", nil)
	rec := httptest.NewRecorder()
	h.ListDeliveries(rec, r)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	out := decodePaginated(t, rec)
	if out.Total != 2 {
		t.Fatalf("total = %d, want 2", out.Total)
	}
	items, ok := out.Items.([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("items = %#v, want 2 items", out.Items)
	}
	if out.Page != httputil.DefaultPage || out.PageSize != httputil.DefaultPageSize {
		t.Fatalf("page/pageSize = %d/%d, want defaults", out.Page, out.PageSize)
	}
}

func TestListDeliveriesStatusFilter(t *testing.T) {
	s := newTestStore(t)
	h := NewHandler(s)

	seedDelivery(t, s, store.DeliveryStatusSent)
	seedDelivery(t, s, store.DeliveryStatusFailed)

	r := httptest.NewRequest(http.MethodGet, "/api/notifications/deliveries?status=failed", nil)
	rec := httptest.NewRecorder()
	h.ListDeliveries(rec, r)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	out := decodePaginated(t, rec)
	if out.Total != 1 {
		t.Fatalf("total = %d, want 1", out.Total)
	}
	items, ok := out.Items.([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v, want 1 item", out.Items)
	}
	item, ok := items[0].(map[string]any)
	if !ok || item["status"] != "failed" {
		t.Fatalf("items[0] = %#v, want status=failed", items[0])
	}
}

func TestListDeliveriesEmpty(t *testing.T) {
	s := newTestStore(t)
	h := NewHandler(s)

	r := httptest.NewRequest(http.MethodGet, "/api/notifications/deliveries", nil)
	rec := httptest.NewRecorder()
	h.ListDeliveries(rec, r)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	out := decodePaginated(t, rec)
	if out.Total != 0 {
		t.Fatalf("total = %d, want 0", out.Total)
	}
	items, ok := out.Items.([]any)
	if !ok || len(items) != 0 {
		t.Fatalf("items = %#v, want empty array", out.Items)
	}
	if !jsonHasEmptyArrayItems(body) {
		t.Fatalf("body = %s, want items to serialize as [] not null", body)
	}
}

// jsonHasEmptyArrayItems is a cheap guard against a nil slice serializing as
// `"items":null` instead of `"items":[]` — the two decode identically via
// out.Items above, so this checks the raw wire format.
func jsonHasEmptyArrayItems(body string) bool {
	return !strings.Contains(body, `"items":null`)
}

func TestListSuccess(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	userID, token, wrapped := apitest.AuthedUser(t, s, "viewer", http.HandlerFunc(h.List))

	// Seed a notification for this user
	n := &store.NotificationRecord{
		UserID:    userID,
		EventType: "test",
		Title:     "Test",
		Message:   "Test notification",
		Read:      false,
	}
	if err := s.CreateNotification(context.Background(), n); err != nil {
		t.Fatalf("seed notification: %v", err)
	}

	// Make a request with auth
	r := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, r)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	out := decodePaginated(t, rec)
	if out.Total != 1 {
		t.Fatalf("total = %d, want 1", out.Total)
	}
}

func TestMarkReadSuccess(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	userID, token, wrapped := apitest.AuthedUser(t, s, "viewer", http.HandlerFunc(h.MarkRead))

	// Seed a notification
	n := &store.NotificationRecord{
		UserID:    userID,
		EventType: "test",
		Title:     "Test",
		Message:   "Test notification",
		Read:      false,
	}
	if err := s.CreateNotification(context.Background(), n); err != nil {
		t.Fatalf("seed notification: %v", err)
	}

	// Mark it as read
	r := httptest.NewRequest(http.MethodPost, "/api/notifications/"+n.ID+"/read", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	r.SetPathValue("id", n.ID)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, r)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestMarkReadNotFound(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	_, token, wrapped := apitest.AuthedUser(t, s, "viewer", http.HandlerFunc(h.MarkRead))

	r := httptest.NewRequest(http.MethodPost, "/api/notifications/missing/read", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	r.SetPathValue("id", "missing")
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, r)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}

func TestReadAllSuccess(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	userID, token, wrapped := apitest.AuthedUser(t, s, "viewer", http.HandlerFunc(h.ReadAll))

	// Seed multiple notifications
	for i := 0; i < 3; i++ {
		n := &store.NotificationRecord{
			UserID:    userID,
			EventType: "test",
			Title:     "Test",
			Message:   "Test notification",
			Read:      false,
		}
		if err := s.CreateNotification(context.Background(), n); err != nil {
			t.Fatalf("seed notification: %v", err)
		}
	}

	// Mark all as read
	r := httptest.NewRequest(http.MethodPost, "/api/notifications/read-all", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, r)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body)
	}
}
