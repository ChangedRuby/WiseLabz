package store

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
)

type countingAuthDB struct {
	DBTX
	reads  atomic.Int64
	writes atomic.Int64
}

func (db *countingAuthDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	db.reads.Add(1)
	return db.DBTX.QueryRowContext(ctx, query, args...)
}

func (db *countingAuthDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	result, err := db.DBTX.ExecContext(ctx, query, args...)
	if err == nil {
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return result, rowsErr
		}
		db.writes.Add(rows)
	}
	return result, err
}

func TestUserStatusCache(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := &User{Username: "cached", InstanceAdminRole: "admin"}
	if err := s.CreateUser(ctx, u); err != nil {
		t.Fatal(err)
	}
	db := &countingAuthDB{DBTX: s.db}
	s.db = db
	check := func(role string, disabled bool) {
		t.Helper()
		got, off, err := s.GetUserRoleStatus(ctx, u.ID)
		if err != nil || got != role || off != disabled {
			t.Fatalf("status = %q, %v, %v", got, off, err)
		}
	}
	check("admin", false)
	for range 10 {
		check("admin", false)
	}
	if got := db.reads.Load(); got != 1 {
		t.Fatalf("warm cache made %d reads, want 1", got)
	}
	s.userStatusMu.Lock()
	entry := s.userStatuses[u.ID]
	entry.expires = time.Now().Add(-time.Second)
	s.userStatuses[u.ID] = entry
	s.userStatusMu.Unlock()
	check("admin", false)
	if got := db.reads.Load(); got != 2 {
		t.Fatalf("expired cache made %d reads, want 2", got)
	}
	if err := s.UpdateUser(ctx, u.ID, map[string]any{"instance_admin_role": "user"}); err != nil {
		t.Fatal(err)
	}
	check("user", false)
	if err := s.UpdateUser(ctx, u.ID, map[string]any{"disabled": true}); err != nil {
		t.Fatal(err)
	}
	check("user", true)
	if err := s.DeleteUser(ctx, u.ID); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, _, err := s.GetUserRoleStatus(ctx, u.ID); !errors.Is(err, ErrNotFound) {
			t.Fatalf("deleted user: %v", err)
		}
	}
	if got := db.reads.Load(); got != 6 {
		t.Fatalf("reads = %d, want 6 (errors must not be cached)", got)
	}
}

func TestUserStatusCacheTransaction(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := &User{Username: "transaction", InstanceAdminRole: "admin"}
	if err := s.CreateUser(ctx, u); err != nil {
		t.Fatal(err)
	}
	for _, rollback := range []bool{true, false} {
		if _, _, err := s.GetUserRoleStatus(ctx, u.ID); err != nil {
			t.Fatal(err)
		}
		stop := errors.New("rollback")
		err := s.WithinTransaction(ctx, func(tx *Store) error {
			if err := tx.UpdateUser(ctx, u.ID, map[string]any{"disabled": true}); err != nil {
				return err
			}
			if rollback {
				return stop
			}
			return nil
		})
		if rollback && !errors.Is(err, stop) || !rollback && err != nil {
			t.Fatal(err)
		}
		_, disabled, err := s.GetUserRoleStatus(ctx, u.ID)
		if err != nil || disabled == rollback {
			t.Fatalf("disabled after transaction = %v, %v", disabled, err)
		}
	}
}

func TestUserStatusCacheConcurrentMutation(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := &User{Username: "concurrent"}
	if err := s.CreateUser(ctx, u); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 20 {
				if _, _, err := s.GetUserRoleStatus(ctx, u.ID); err != nil {
					t.Error(err)
				}
			}
		})
	}
	if err := s.UpdateUser(ctx, u.ID, map[string]any{"disabled": true}); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	_, disabled, err := s.GetUserRoleStatus(ctx, u.ID)
	if err != nil || !disabled {
		t.Fatalf("stale status after mutation: %v, %v", disabled, err)
	}
}

func TestAPIKeyLastUsedThrottle(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := &User{Username: "throttle"}
	if err := s.CreateUser(ctx, u); err != nil {
		t.Fatal(err)
	}
	key := &APIKey{UserID: u.ID, Name: "key", Role: "viewer", TokenHash: HashToken("throttle")}
	if err := s.CreateAPIKey(ctx, key); err != nil {
		t.Fatal(err)
	}
	// Count affected rows, including writes of an unchanged timestamp.
	db := &countingAuthDB{DBTX: s.db}
	s.db = db
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			if err := s.TouchAPIKeyLastUsed(ctx, key.ID); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if count := db.writes.Load(); count != 1 {
		t.Fatalf("writes = %d, want 1", count)
	}
	claims, err := s.LookupAPIKey(ctx, key.TokenHash)
	if err != nil || claims.LastUsedAt == "" {
		t.Fatalf("lookup timestamp: %#v, %v", claims, err)
	}
	if _, err := s.db.ExecContext(ctx, "UPDATE api_keys SET last_used_at = ?", time.Now().UTC().Add(-2*time.Minute).Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}
	if err := s.TouchAPIKeyLastUsed(ctx, key.ID); err != nil {
		t.Fatal(err)
	}
	if count := db.writes.Load(); count != 3 {
		t.Fatalf("writes after expiry = %d, want 3 including test setup", count)
	}
}

func TestCachedUserStatusRejectsRevokedAccess(t *testing.T) {
	for _, mutation := range []string{"demote", "disable", "delete"} {
		t.Run(mutation, func(t *testing.T) {
			s := newDocTestStore(t)
			ctx := context.Background()
			u := &User{Username: "session", InstanceAdminRole: "admin"}
			if err := s.CreateUser(ctx, u); err != nil {
				t.Fatal(err)
			}
			svc := auth.NewService("test-secret", time.Minute, time.Hour)
			pair, err := svc.IssuePair(u.ID, true)
			if err != nil {
				t.Fatal(err)
			}
			handler := auth.AuthMiddleware(svc, s)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
			request := func() int {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)
				return rec.Code
			}
			for range 2 {
				if got := request(); got != http.StatusOK {
					t.Fatalf("warm request = %d", got)
				}
			}
			switch mutation {
			case "demote":
				err = s.UpdateUser(ctx, u.ID, map[string]any{"instance_admin_role": "user"})
			case "disable":
				err = s.UpdateUser(ctx, u.ID, map[string]any{"disabled": true})
			case "delete":
				err = s.DeleteUser(ctx, u.ID)
			}
			if err != nil {
				t.Fatal(err)
			}
			if got := request(); got != http.StatusUnauthorized {
				t.Fatalf("request after %s = %d, want 401", mutation, got)
			}
		})
	}
}
