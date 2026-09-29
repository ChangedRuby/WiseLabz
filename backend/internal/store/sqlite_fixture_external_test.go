package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/store/storetest"
)

func TestMigratedSQLiteCopiesAreIsolated(t *testing.T) {
	for _, fixture := range []struct {
		name string
		copy func(testing.TB) string
	}{
		{"store", store.MigratedSQLiteForTest},
		{"storetest", storetest.MigratedSQLite},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			openCopy := func() *sql.DB {
				path := fixture.copy(t)
				info, err := os.Stat(path)
				if err != nil {
					t.Fatal(err)
				}
				if info.Mode().Perm() != 0o600 {
					t.Fatalf("copy permissions = %o, want 600", info.Mode().Perm())
				}
				db, err := sql.Open("sqlite", "file:"+path+"?cache=shared")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = db.Close() })
				return db
			}
			checkClean := func(db *sql.DB) error {
				status, err := store.GetMigrationStatus(db, "sqlite")
				if err != nil {
					return err
				}
				if status.Dirty || status.Pending() || status.Current != status.Latest {
					return fmt.Errorf("copy migration status = %+v", status)
				}
				for _, table := range []string{"docs", "users"} {
					var count int
					if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
						return err
					}
					if count != 0 {
						return fmt.Errorf("copy has %d unexpected %s rows", count, table)
					}
				}
				return nil
			}

			// Exercise concurrent initialization/copying without parallel tests.
			var copies sync.WaitGroup
			for range 8 {
				copies.Add(1)
				go func() {
					defer copies.Done()
					if err := checkClean(openCopy()); err != nil {
						t.Error(err)
					}
				}()
			}
			copies.Wait()

			first, second := openCopy(), openCopy()
			if err := checkClean(first); err != nil {
				t.Fatal(err)
			}
			s := store.New(first, "sqlite")
			if err := s.Init(context.Background(), "fixture-admin-password"); err != nil {
				t.Fatal(err)
			}
			if err := s.CreateDoc(context.Background(), &store.DocRecord{ID: "only-first", Title: "private", Content: "first copy"}); err != nil {
				t.Fatal(err)
			}
			if err := checkClean(second); err != nil {
				t.Fatalf("existing copy leaked seed data: %v", err)
			}
			if err := checkClean(openCopy()); err != nil {
				t.Fatalf("cached template leaked seed data: %v", err)
			}
		})
	}
}
