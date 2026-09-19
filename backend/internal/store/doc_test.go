package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// newDocTestStore returns a migrated Store. It uses a per-test SQLite file by
// default; when WISELABZ_TEST_POSTGRES_DSN is set it instead uses a fresh,
// isolated schema in that Postgres database so the same store tests run against
// both dialects.
func newDocTestStore(t *testing.T) *Store {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	if pgDSN := os.Getenv("WISELABZ_TEST_POSTGRES_DSN"); pgDSN != "" {
		return newPostgresTestStore(t, pgDSN, logger)
	}

	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}
	db.SetMaxOpenConns(1)

	return New(db, "sqlite")
}

// newPostgresTestStore creates a uniquely named schema, migrates it, and drops
// it on cleanup, so tests sharing one Postgres database stay isolated.
func newPostgresTestStore(t *testing.T, dsn string, logger *slog.Logger) *Store {
	t.Helper()
	admin, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	schema := "t_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := admin.Exec("CREATE SCHEMA " + schema); err != nil {
		_ = admin.Close()
		t.Fatalf("create schema: %v", err)
	}

	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse postgres dsn: %v", err)
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()

	db, err := sql.Open("pgx", u.String())
	if err != nil {
		t.Fatalf("open postgres schema db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		_, _ = admin.Exec("DROP SCHEMA " + schema + " CASCADE")
		_ = admin.Close()
	})
	if err := RunMigrations(db, "postgres", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}
	return New(db, "postgres")
}

// skipOnPostgres skips tests that depend on SQLite-only features.
func skipOnPostgres(t *testing.T, why string) {
	t.Helper()
	if os.Getenv("WISELABZ_TEST_POSTGRES_DSN") != "" {
		t.Skip(why)
	}
}

// mustCreateUser inserts a user with the given id (if absent) so foreign keys resolve.
func mustCreateUser(t *testing.T, s *Store, id string) {
	t.Helper()
	if u, err := s.GetUserByID(context.Background(), id); err == nil && u != nil {
		return
	}
	if err := s.CreateUser(context.Background(), &User{ID: id, Username: id, DisplayName: id, Email: id + "@example.com"}); err != nil {
		t.Fatalf("CreateUser(%s) error: %v", id, err)
	}
}

func TestUpdateDocOptimisticConcurrency(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	d := &DocRecord{Title: "Test Doc", Content: "v1"}
	if err := s.CreateDoc(ctx, d); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}
	if d.CurrentVersion != 1 {
		t.Fatalf("CreateDoc() CurrentVersion = %d, want 1", d.CurrentVersion)
	}

	// Matching expectedVersion succeeds and bumps the version.
	match := d.CurrentVersion
	if err := s.UpdateDoc(ctx, d.ID, "v2", &match); err != nil {
		t.Fatalf("UpdateDoc() with matching version error: %v", err)
	}
	got, err := s.GetDoc(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetDoc() error: %v", err)
	}
	if got.Content != "v2" || got.CurrentVersion != 2 {
		t.Fatalf("after matching update: content=%q version=%d, want content=v2 version=2", got.Content, got.CurrentVersion)
	}

	// Stale expectedVersion is rejected and leaves the doc untouched.
	stale := 1
	err = s.UpdateDoc(ctx, d.ID, "v3-should-not-land", &stale)
	if !errors.Is(err, ErrVersionConflict) {
		t.Fatalf("UpdateDoc() with stale version error = %v, want ErrVersionConflict", err)
	}
	got, err = s.GetDoc(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetDoc() error: %v", err)
	}
	if got.Content != "v2" || got.CurrentVersion != 2 {
		t.Fatalf("after rejected update: content=%q version=%d, want unchanged content=v2 version=2", got.Content, got.CurrentVersion)
	}

	// nil expectedVersion keeps today's unconditional last-write-wins behavior.
	if err := s.UpdateDoc(ctx, d.ID, "v4", nil); err != nil {
		t.Fatalf("UpdateDoc() with nil version error: %v", err)
	}
	got, err = s.GetDoc(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetDoc() error: %v", err)
	}
	if got.Content != "v4" || got.CurrentVersion != 3 {
		t.Fatalf("after unconditional update: content=%q version=%d, want content=v4 version=3", got.Content, got.CurrentVersion)
	}

	// Non-existent doc: ErrNotFound with or without an expected version.
	if err := s.UpdateDoc(ctx, "missing-id", "x", nil); !errors.Is(err, ErrNotFound) {
		t.Fatalf("UpdateDoc() on missing doc (nil version) error = %v, want ErrNotFound", err)
	}
	someVersion := 1
	if err := s.UpdateDoc(ctx, "missing-id", "x", &someVersion); !errors.Is(err, ErrNotFound) {
		t.Fatalf("UpdateDoc() on missing doc (with version) error = %v, want ErrNotFound", err)
	}
}

func TestListDocsGroupedByService(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	serviceA, serviceB := mustCreateMaintenanceConnector(t, s), mustCreateMaintenanceConnector(t, s)
	for _, serviceID := range []string{serviceA, serviceB} {
		if err := s.CreateDoc(ctx, &DocRecord{Title: serviceID, Kind: "service", ServiceID: serviceID, Content: "content"}); err != nil {
			t.Fatalf("CreateDoc() error: %v", err)
		}
	}

	got, err := s.ListDocsGroupedByService(ctx)
	if err != nil {
		t.Fatalf("ListDocsGroupedByService() error: %v", err)
	}
	if len(got[serviceA]) != 1 || len(got[serviceB]) != 1 {
		t.Fatalf("grouped docs = %+v, want one doc per service", got)
	}
}

func TestTemplateVersions(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	tmpl := &TemplateRecord{
		Name:        "Runbook",
		Description: "Initial",
		AppliesTo:   `["docker"]`,
	}
	if err := s.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatalf("CreateTemplate() error: %v", err)
	}
	if tmpl.CurrentVersion != 1 {
		t.Fatalf("CreateTemplate() CurrentVersion = %d, want 1", tmpl.CurrentVersion)
	}

	gotTemplate, err := s.GetTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("GetTemplate() error: %v", err)
	}
	if gotTemplate.CurrentVersion != 1 {
		t.Fatalf("GetTemplate() CurrentVersion = %d, want 1", gotTemplate.CurrentVersion)
	}
	templates, total, err := s.ListTemplates(ctx, 0, 10)
	if err != nil {
		t.Fatalf("ListTemplates() error: %v", err)
	}
	if total != 1 || len(templates) != 1 || templates[0].CurrentVersion != 1 {
		t.Fatalf("ListTemplates() = total %d, templates %#v; want one template at version 1", total, templates)
	}

	sectionsJSON, err := json.Marshal([]TemplateVersionSection{
		{Title: "Instalação 🚀", Order: 2, Body: "linha 1\nlinha 2"},
		{Title: "Overview", Order: 1, Body: "{{service.name}}"},
	})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	v1 := &TemplateVersionRecord{
		TemplateID:  tmpl.ID,
		Rev:         1,
		Name:        tmpl.Name,
		Description: tmpl.Description,
		AppliesTo:   tmpl.AppliesTo,
		Sections:    string(sectionsJSON),
		Author:      "operator-1",
		Trigger:     "save",
	}
	if err := s.CreateTemplateVersion(ctx, v1); err != nil {
		t.Fatalf("CreateTemplateVersion(v1) error: %v", err)
	}
	if v1.ID == "" || v1.CreatedAt == "" {
		t.Fatalf("CreateTemplateVersion(v1) did not populate identity fields: %#v", v1)
	}

	if err := s.UpdateTemplate(ctx, tmpl.ID, map[string]any{"description": "Restored"}); err != nil {
		t.Fatalf("UpdateTemplate() error: %v", err)
	}
	v2 := &TemplateVersionRecord{
		TemplateID:  tmpl.ID,
		Rev:         2,
		Name:        tmpl.Name,
		Description: "Restored",
		Sections:    "[]",
		Trigger:     "restore",
	}
	if err := s.CreateTemplateVersion(ctx, v2); err != nil {
		t.Fatalf("CreateTemplateVersion(v2) error: %v", err)
	}

	versions, err := s.GetTemplateVersions(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("GetTemplateVersions() error: %v", err)
	}
	if len(versions) != 2 || versions[0].Rev != 2 || versions[1].Rev != 1 {
		t.Fatalf("GetTemplateVersions() revisions = %#v, want [2, 1]", versions)
	}
	if versions[0].AppliesTo != "" || versions[0].Author != "" {
		t.Fatalf("nullable strings = appliesTo %q, author %q; want empty", versions[0].AppliesTo, versions[0].Author)
	}
	var sections []TemplateVersionSection
	if err := json.Unmarshal([]byte(versions[1].Sections), &sections); err != nil {
		t.Fatalf("Unmarshal(stored sections) error: %v", err)
	}
	if len(sections) != 2 || sections[0].Title != "Instalação 🚀" || sections[0].Body != "linha 1\nlinha 2" {
		t.Fatalf("stored sections = %#v, want lossless JSON round trip", sections)
	}

	gotTemplate, err = s.GetTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("GetTemplate() after update error: %v", err)
	}
	if gotTemplate.CurrentVersion != 2 {
		t.Fatalf("GetTemplate() CurrentVersion = %d, want 2", gotTemplate.CurrentVersion)
	}
}

func TestTemplateVersionConstraintsAndRawJSON(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	tmpl := &TemplateRecord{Name: "Runbook"}
	if err := s.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatalf("CreateTemplate() error: %v", err)
	}
	base := TemplateVersionRecord{
		TemplateID: tmpl.ID,
		Rev:        1,
		Name:       tmpl.Name,
		AppliesTo:  "not-json",
		Sections:   "not-json",
		Trigger:    "save",
	}
	if err := s.CreateTemplateVersion(ctx, &base); err != nil {
		t.Fatalf("CreateTemplateVersion(raw JSON) error: %v", err)
	}
	versions, err := s.GetTemplateVersions(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("GetTemplateVersions() error: %v", err)
	}
	if len(versions) != 1 || versions[0].AppliesTo != "not-json" || versions[0].Sections != "not-json" {
		t.Fatalf("GetTemplateVersions() = %#v, want malformed JSON preserved for boundary validation", versions)
	}

	duplicate := base
	duplicate.ID = ""
	if err := s.CreateTemplateVersion(ctx, &duplicate); err == nil {
		t.Fatal("CreateTemplateVersion(duplicate rev) error = nil, want unique constraint error")
	}
	invalidTrigger := base
	invalidTrigger.ID = ""
	invalidTrigger.Rev = 2
	invalidTrigger.Trigger = "manual"
	if err := s.CreateTemplateVersion(ctx, &invalidTrigger); err == nil {
		t.Fatal("CreateTemplateVersion(invalid trigger) error = nil, want check constraint error")
	}

	empty, err := s.GetTemplateVersions(ctx, "missing-template")
	if err != nil {
		t.Fatalf("GetTemplateVersions(missing) error: %v", err)
	}
	if empty == nil || len(empty) != 0 {
		t.Fatalf("GetTemplateVersions(missing) = %#v, want non-nil empty slice", empty)
	}
}

func TestUpdateTemplateConcurrentVersionBumps(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	tmpl := &TemplateRecord{Name: "Runbook"}
	if err := s.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatalf("CreateTemplate() error: %v", err)
	}

	const updates = 8
	start := make(chan struct{})
	errs := make(chan error, updates)
	var wg sync.WaitGroup
	for range updates {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- s.UpdateTemplate(ctx, tmpl.ID, map[string]any{})
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent UpdateTemplate() error: %v", err)
		}
	}

	got, err := s.GetTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("GetTemplate() error: %v", err)
	}
	if got.CurrentVersion != 1+updates {
		t.Fatalf("CurrentVersion = %d, want %d", got.CurrentVersion, 1+updates)
	}
}

func TestTemplateVersionIndexAndCascade(t *testing.T) {
	skipOnPostgres(t, "EXPLAIN QUERY PLAN and PRAGMA are SQLite-only")
	ctx := context.Background()
	s := newDocTestStore(t)
	tmpl := &TemplateRecord{Name: "Runbook"}
	if err := s.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatalf("CreateTemplate() error: %v", err)
	}
	if err := s.CreateTemplateVersion(ctx, &TemplateVersionRecord{
		TemplateID: tmpl.ID, Rev: 1, Name: tmpl.Name, Sections: "[]", Trigger: "save",
	}); err != nil {
		t.Fatalf("CreateTemplateVersion() error: %v", err)
	}

	rows, err := s.db.QueryContext(ctx, `EXPLAIN QUERY PLAN
		SELECT id, template_id, rev, name, description, applies_to, sections, author, trigger, created_at
		FROM template_versions WHERE template_id = ? ORDER BY rev DESC`, tmpl.ID)
	if err != nil {
		t.Fatalf("EXPLAIN QUERY PLAN error: %v", err)
	}
	defer rows.Close() //nolint:errcheck
	var plan strings.Builder
	for rows.Next() {
		var id, parent, unused int
		var detail string
		if err := rows.Scan(&id, &parent, &unused, &detail); err != nil {
			t.Fatalf("scan query plan: %v", err)
		}
		plan.WriteString(detail)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate query plan: %v", err)
	}
	if !strings.Contains(plan.String(), "idx_template_versions_tmpl_rev") {
		t.Fatalf("query plan %q does not use template version index", plan.String())
	}

	if _, err := s.db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
	if err := s.DeleteTemplate(ctx, tmpl.ID); err != nil {
		t.Fatalf("DeleteTemplate() error: %v", err)
	}
	versions, err := s.GetTemplateVersions(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("GetTemplateVersions() after delete error: %v", err)
	}
	if len(versions) != 0 {
		t.Fatalf("versions after template delete = %#v, want empty", versions)
	}
}

func TestDocLockConflict(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	mustCreateUser(t, s, "user-1")
	mustCreateUser(t, s, "user-2")
	d := &DocRecord{Title: "Test Doc", Content: "v1"}
	if err := s.CreateDoc(ctx, d); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}

	lock1, err := s.AcquireDocLock(ctx, d.ID, "user-1")
	if err != nil {
		t.Fatalf("AcquireDocLock(user-1) error: %v", err)
	}
	if lock1.UserID != "user-1" {
		t.Fatalf("lock1.UserID = %q, want user-1", lock1.UserID)
	}

	lock2, err := s.AcquireDocLock(ctx, d.ID, "user-2")
	if !errors.Is(err, ErrLockHeldByOther) {
		t.Fatalf("AcquireDocLock(user-2) error = %v, want ErrLockHeldByOther", err)
	}
	if lock2.UserID != "user-1" {
		t.Fatalf("returned lock.UserID = %q, want user-1", lock2.UserID)
	}
}

func TestDocLockRenewalByHolder(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	mustCreateUser(t, s, "user-1")
	mustCreateUser(t, s, "user-2")
	d := &DocRecord{Title: "Test Doc", Content: "v1"}
	if err := s.CreateDoc(ctx, d); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}

	lock1, err := s.AcquireDocLock(ctx, d.ID, "user-1")
	if err != nil {
		t.Fatalf("first AcquireDocLock() error: %v", err)
	}
	first := lock1.ExpiresAt

	time.Sleep(1100 * time.Millisecond)

	lock2, err := s.AcquireDocLock(ctx, d.ID, "user-1")
	if err != nil {
		t.Fatalf("second AcquireDocLock() (renewal) error: %v", err)
	}
	if lock2.UserID != "user-1" {
		t.Fatalf("renewal: lock.UserID = %q, want user-1", lock2.UserID)
	}
	if lock2.ExpiresAt <= first {
		t.Fatalf("renewal: ExpiresAt did not advance: old %q, new %q", first, lock2.ExpiresAt)
	}
}

func TestDocLockAcquireAfterExpiry(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	mustCreateUser(t, s, "user-1")
	mustCreateUser(t, s, "user-2")
	d := &DocRecord{Title: "Test Doc", Content: "v1"}
	if err := s.CreateDoc(ctx, d); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}

	_, err := s.AcquireDocLock(ctx, d.ID, "user-1")
	if err != nil {
		t.Fatalf("AcquireDocLock(user-1) error: %v", err)
	}

	past := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	if _, err := s.db.ExecContext(ctx, `UPDATE doc_locks SET expires_at = ? WHERE doc_id = ?`, past, d.ID); err != nil {
		t.Fatalf("update expires_at: %v", err)
	}

	lock2, err := s.AcquireDocLock(ctx, d.ID, "user-2")
	if err != nil {
		t.Fatalf("AcquireDocLock(user-2) after expiry error: %v", err)
	}
	if lock2.UserID != "user-2" {
		t.Fatalf("after expiry: lock.UserID = %q, want user-2", lock2.UserID)
	}
}

func TestDocLockReleaseOnlyByHolder(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	mustCreateUser(t, s, "user-1")
	mustCreateUser(t, s, "user-2")
	d := &DocRecord{Title: "Test Doc", Content: "v1"}
	if err := s.CreateDoc(ctx, d); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}

	_, err := s.AcquireDocLock(ctx, d.ID, "user-1")
	if err != nil {
		t.Fatalf("AcquireDocLock() error: %v", err)
	}

	if err := s.ReleaseDocLock(ctx, d.ID, "user-2"); err != nil {
		t.Fatalf("ReleaseDocLock(user-2) error: %v", err)
	}

	got, err := s.GetDocLock(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetDocLock() after non-holder release error: %v", err)
	}
	if got.UserID != "user-1" {
		t.Fatalf("after non-holder release: lock.UserID = %q, want user-1 (unchanged)", got.UserID)
	}

	if err := s.ReleaseDocLock(ctx, d.ID, "user-1"); err != nil {
		t.Fatalf("ReleaseDocLock(user-1) error: %v", err)
	}

	_, err = s.GetDocLock(ctx, d.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetDocLock() after holder release error = %v, want ErrNotFound", err)
	}
}
