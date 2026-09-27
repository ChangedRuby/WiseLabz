package leader

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func postgresDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("WISELABZ_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("WISELABZ_TEST_POSTGRES_DSN not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(4)
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSecondElectorWaitsThenTakesOver(t *testing.T) {
	db := postgresDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	first := New(db, 20*time.Millisecond)
	second := New(db, 20*time.Millisecond)
	defer first.Close()  //nolint:errcheck
	defer second.Close() //nolint:errcheck
	if err := first.Campaign(ctx); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- second.Campaign(ctx) }()
	select {
	case err := <-done:
		t.Fatalf("second acquired held lock: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("second did not take over")
	}
}

func TestWatchReportsTerminatedSession(t *testing.T) {
	db := postgresDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	e := New(db, 20*time.Millisecond)
	defer e.Close() //nolint:errcheck
	if err := e.Campaign(ctx); err != nil {
		t.Fatal(err)
	}
	lost := e.Watch(ctx)
	var pid int
	if err := e.conn.QueryRowContext(ctx, "SELECT pg_backend_pid()").Scan(&pid); err != nil {
		t.Fatal(err)
	}
	var killed bool
	if err := db.QueryRowContext(ctx, "SELECT pg_terminate_backend($1)", pid).Scan(&killed); err != nil || !killed {
		t.Fatalf("terminate backend: killed=%v err=%v", killed, err)
	}
	select {
	case err := <-lost:
		if err == nil {
			t.Fatal("watch did not report loss")
		}
	case <-ctx.Done():
		t.Fatal("watch did not detect loss")
	}
}
