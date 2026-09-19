package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestOpenDBEnablesSQLiteForeignKeys(t *testing.T) {
	db, err := OpenDB("sqlite", "file:"+filepath.Join(t.TempDir(), "test.db")+"?cache=shared")
	if err != nil {
		t.Fatalf("OpenDB() error: %v", err)
	}
	defer db.Close() //nolint:errcheck

	var enabled int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&enabled); err != nil {
		t.Fatalf("query foreign_keys pragma: %v", err)
	}
	if enabled != 1 {
		t.Fatalf("foreign_keys = %d, want 1", enabled)
	}
}

func TestWithinTransactionRollsBack(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	template := &TemplateRecord{Name: "before"}
	if err := s.CreateTemplate(ctx, template); err != nil {
		t.Fatalf("CreateTemplate() error: %v", err)
	}
	wantErr := errors.New("stop")
	err := s.WithinTransaction(ctx, func(tx *Store) error {
		if err := tx.UpdateTemplate(ctx, template.ID, map[string]any{"name": "after"}); err != nil {
			return err
		}
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("WithinTransaction() error = %v, want %v", err, wantErr)
	}
	got, err := s.GetTemplate(ctx, template.ID)
	if err != nil {
		t.Fatalf("GetTemplate() error: %v", err)
	}
	if got.Name != "before" || got.CurrentVersion != 1 {
		t.Fatalf("template after rollback = %#v", got)
	}
}

func TestOpenDBSetsSQLiteDurabilityPragmas(t *testing.T) {
	db, err := OpenDB("sqlite", "file:"+filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("OpenDB() error: %v", err)
	}
	defer db.Close() //nolint:errcheck

	var mode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&mode); err != nil || mode != "wal" {
		t.Fatalf("journal_mode = %q, %v, want wal", mode, err)
	}
	var timeout, sync int
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&timeout); err != nil || timeout != 5000 {
		t.Fatalf("busy_timeout = %d, %v, want 5000", timeout, err)
	}
	if err := db.QueryRow("PRAGMA synchronous").Scan(&sync); err != nil || sync != 1 {
		t.Fatalf("synchronous = %d, %v, want 1 (NORMAL)", sync, err)
	}
}
