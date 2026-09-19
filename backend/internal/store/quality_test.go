package store

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func newConcurrentQualityTestStore(t *testing.T) *Store {
	t.Helper()
	dsn := "file:" + t.TempDir() + "/quality.db?cache=shared&_pragma=busy_timeout(5000)"
	db, err := OpenDB("sqlite", dsn)
	if err != nil {
		t.Fatalf("OpenDB() error: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}
	db.SetMaxOpenConns(8)
	return New(db, "sqlite")
}

func seedQualityConnector(t *testing.T, s *Store, name string) *ConnectorRecord {
	t.Helper()
	c := &ConnectorRecord{Name: name, Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return c
}

func TestUpsertQualityFindingDedup(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	c := seedQualityConnector(t, s, "dedup")
	doc := &DocRecord{Title: "Runbook", Kind: "service", ServiceID: c.ID}
	if err := s.CreateDoc(ctx, doc); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}

	first := &QualityFindingRecord{ConnectorID: c.ID, CheckType: "empty", Severity: "warning", Title: "Empty", Description: "first"}
	if err := s.UpsertQualityFinding(ctx, first); err != nil {
		t.Fatalf("UpsertQualityFinding(first) error: %v", err)
	}
	second := &QualityFindingRecord{ConnectorID: c.ID, DocID: doc.ID, CheckType: "empty", Severity: "critical", Title: "Still empty", Description: "second", RemediationLink: "/docs/" + doc.ID + "/edit"}
	if err := s.UpsertQualityFinding(ctx, second); err != nil {
		t.Fatalf("UpsertQualityFinding(second) error: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("conflict returned ID = %q, want persisted ID %q", second.ID, first.ID)
	}

	findings, total, err := s.ListQualityFindings(ctx, c.ID, "empty", "open", "", 0, 10)
	if err != nil {
		t.Fatalf("ListQualityFindings() error: %v", err)
	}
	if total != 1 || len(findings) != 1 {
		t.Fatalf("findings = %d, total = %d; want one", len(findings), total)
	}
	got := findings[0]
	if got.DetectedCount != 2 || got.DocID != doc.ID || got.Severity != "critical" || got.Description != "second" || got.RemediationLink != second.RemediationLink {
		t.Fatalf("deduplicated finding = %+v", got)
	}
}

func TestResolveThenReopenCreatesFreshRow(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	c := seedQualityConnector(t, s, "reopen")

	first := &QualityFindingRecord{ConnectorID: c.ID, CheckType: "failing", Severity: "critical", Title: "Failing"}
	if err := s.UpsertQualityFinding(ctx, first); err != nil {
		t.Fatalf("UpsertQualityFinding(first) error: %v", err)
	}
	if err := s.ResolveQualityFinding(ctx, c.ID, "failing"); err != nil {
		t.Fatalf("ResolveQualityFinding() error: %v", err)
	}
	second := &QualityFindingRecord{ConnectorID: c.ID, CheckType: "failing", Severity: "warning", Title: "Failing again"}
	if err := s.UpsertQualityFinding(ctx, second); err != nil {
		t.Fatalf("UpsertQualityFinding(second) error: %v", err)
	}
	if second.ID == first.ID {
		t.Fatal("reopened finding reused resolved row ID")
	}

	all, total, err := s.ListQualityFindings(ctx, c.ID, "failing", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListQualityFindings() error: %v", err)
	}
	if total != 2 || len(all) != 2 {
		t.Fatalf("findings = %d, total = %d; want two", len(all), total)
	}
	open, err := s.GetQualityFinding(ctx, second.ID)
	if err != nil {
		t.Fatalf("GetQualityFinding() error: %v", err)
	}
	if open.Status != "open" || open.DetectedCount != 1 || open.ResolvedAt != "" {
		t.Fatalf("reopened finding = %+v", open)
	}
}

func TestListQualityFindingsFilters(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	one := seedQualityConnector(t, s, "one")
	two := seedQualityConnector(t, s, "two")
	base := time.Now().UTC()
	findings := []*QualityFindingRecord{
		{ConnectorID: one.ID, CheckType: "stale", Severity: "warning", Title: "old", LastSeenAt: base.Add(-time.Hour).Format(time.RFC3339Nano)},
		{ConnectorID: one.ID, CheckType: "empty", Severity: "warning", Title: "new", LastSeenAt: base.Format(time.RFC3339Nano)},
		{ConnectorID: two.ID, CheckType: "stale", Severity: "warning", Title: "other", LastSeenAt: base.Add(time.Hour).Format(time.RFC3339Nano)},
	}
	for _, f := range findings {
		if err := s.UpsertQualityFinding(ctx, f); err != nil {
			t.Fatalf("UpsertQualityFinding(%s) error: %v", f.Title, err)
		}
	}
	if err := s.UpdateQualityFindingStatus(ctx, findings[0].ID, "resolved"); err != nil {
		t.Fatalf("UpdateQualityFindingStatus() error: %v", err)
	}

	got, total, err := s.ListQualityFindings(ctx, one.ID, "", "open", "", 0, 10)
	if err != nil {
		t.Fatalf("ListQualityFindings() error: %v", err)
	}
	if total != 1 || len(got) != 1 || got[0].Title != "new" {
		t.Fatalf("connector/status filter = %+v, total %d", got, total)
	}
	got, total, err = s.ListQualityFindings(ctx, "", "stale", "", "", 0, 1)
	if err != nil {
		t.Fatalf("ListQualityFindings(stale) error: %v", err)
	}
	if total != 2 || len(got) != 1 || got[0].Title != "other" {
		t.Fatalf("check/pagination filter = %+v, total %d", got, total)
	}
	openCount, err := s.CountQualityFindingsOpen(ctx)
	if err != nil || openCount != 2 {
		t.Fatalf("CountQualityFindingsOpen() = %d, %v; want 2, nil", openCount, err)
	}
}

func TestUpsertQualityFindingConcurrentDedup(t *testing.T) {
	ctx := context.Background()
	s := newConcurrentQualityTestStore(t)
	c := seedQualityConnector(t, s, "concurrent")

	const workers = 16
	start := make(chan struct{})
	errs := make(chan error, workers)
	ids := make(chan string, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			<-start
			f := &QualityFindingRecord{ConnectorID: c.ID, CheckType: "ownership_incomplete", Severity: "warning", Title: "Owner missing"}
			errs <- s.UpsertQualityFinding(ctx, f)
			ids <- f.ID
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	close(ids)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent UpsertQualityFinding() error: %v", err)
		}
	}

	findings, total, err := s.ListQualityFindings(ctx, c.ID, "ownership_incomplete", "open", "", 0, workers)
	if err != nil {
		t.Fatalf("ListQualityFindings() error: %v", err)
	}
	if total != 1 || len(findings) != 1 || findings[0].DetectedCount != workers {
		t.Fatalf("findings = %+v, total = %d; want one with count %d", findings, total, workers)
	}
	for id := range ids {
		if id != findings[0].ID {
			t.Fatalf("returned ID = %q, want persisted ID %q", id, findings[0].ID)
		}
	}
}

func TestUpsertQualityFindingPostgresDedup(t *testing.T) {
	dsn := os.Getenv("WISELABZ_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("WISELABZ_TEST_POSTGRES_DSN not set; skipping postgres upsert test")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer db.Close() //nolint:errcheck
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := RunMigrations(db, "postgres", logger); err != nil {
		t.Fatalf("RunMigrations(postgres) error: %v", err)
	}
	s := New(db, "postgres")
	connector := &ConnectorRecord{
		ID: uuid.New().String(), Name: "postgres dedup", Category: "virtualization", Type: "proxmox", URL: "https://example.com",
	}
	if err := s.CreateConnector(context.Background(), connector); err != nil {
		t.Fatalf("CreateConnector(postgres) error: %v", err)
	}
	first := &QualityFindingRecord{ConnectorID: connector.ID, CheckType: "stale", Severity: "warning", Title: "Stale"}
	second := &QualityFindingRecord{ConnectorID: connector.ID, CheckType: "stale", Severity: "critical", Title: "Still stale"}
	if err := s.UpsertQualityFinding(context.Background(), first); err != nil {
		t.Fatalf("first postgres upsert error: %v", err)
	}
	if err := s.UpsertQualityFinding(context.Background(), second); err != nil {
		t.Fatalf("second postgres upsert error: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("postgres conflict returned ID = %q, want %q", second.ID, first.ID)
	}
	got, err := s.GetQualityFinding(context.Background(), first.ID)
	if err != nil {
		t.Fatalf("GetQualityFinding(postgres) error: %v", err)
	}
	if got.DetectedCount != 2 || got.Title != "Still stale" {
		t.Fatalf("postgres deduplicated finding = %+v", got)
	}
}

func TestQualityFindingStatusIndexPlan(t *testing.T) {
	skipOnPostgres(t, "EXPLAIN QUERY PLAN is SQLite-only")
	s := newDocTestStore(t)
	rows, err := s.DB().QueryContext(context.Background(), `EXPLAIN QUERY PLAN
		SELECT `+qualityFindingColumns+` FROM quality_findings
		WHERE status = ? ORDER BY last_seen_at DESC LIMIT ? OFFSET ?`, "open", 10, 0)
	if err != nil {
		t.Fatalf("EXPLAIN QUERY PLAN error: %v", err)
	}
	defer rows.Close() //nolint:errcheck
	var id, parent, unused int
	var detail string
	found := false
	for rows.Next() {
		if err := rows.Scan(&id, &parent, &unused, &detail); err != nil {
			t.Fatalf("scan query plan: %v", err)
		}
		if strings.Contains(detail, "idx_quality_findings_status") {
			found = true
		}
	}
	if !found {
		t.Fatal("status-filtered list query did not use idx_quality_findings_status")
	}
}
