package backup_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedOneDoc(ctx context.Context, t *testing.T, s *store.Store) {
	t.Helper()
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "doc-1", Title: "Runbook", Content: "steps"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}
}

func TestVerifyBundleFilePass(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	seedOneDoc(ctx, t, s)

	tmpDir := t.TempDir()
	run, err := backup.ExportToFile(ctx, s, tmpDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	res := backup.VerifyBundleFile(ctx, run.FilePath)
	if res.Status != "pass" {
		t.Fatalf("VerifyBundleFile status = %q, want pass (error: %s)", res.Status, res.Error)
	}
	if res.Checksum != run.Checksum {
		t.Errorf("res.Checksum = %q, want %q", res.Checksum, run.Checksum)
	}
	if res.ActualCounts["docs"] != 1 {
		t.Errorf("ActualCounts[docs] = %d, want 1", res.ActualCounts["docs"])
	}
	if res.ExpectedCounts["docs"] != res.ActualCounts["docs"] {
		t.Errorf("ExpectedCounts[docs]=%d != ActualCounts[docs]=%d", res.ExpectedCounts["docs"], res.ActualCounts["docs"])
	}
}

// TestVerifyBundleFileDetectsChecksumCorruption is the corruption-detection
// case: a bit-flipped bundle file must fail verification via the checksum
// check, before any restore is attempted.
func TestVerifyBundleFileDetectsChecksumCorruption(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	seedOneDoc(ctx, t, s)

	tmpDir := t.TempDir()
	run, err := backup.ExportToFile(ctx, s, tmpDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	corrupted := corruptFile(t, run.FilePath)
	res := backup.VerifyBundleFile(ctx, corrupted)
	if res.Status != "fail" {
		t.Fatalf("VerifyBundleFile status = %q, want fail", res.Status)
	}
	if res.Error == "" {
		t.Error("res.Error is empty for a failed verification")
	}
}

// TestVerifyBundleFileDetectsRowCountMismatch simulates a bundle whose
// content diverges from what was recorded at export time (e.g. a bug that
// truncates rows without touching the file's checksum-relevant bytes): the
// manifest's counts are tampered independently of the checksum, so the
// checksum check still passes but the restore-and-diff step must catch it.
func TestVerifyBundleFileDetectsRowCountMismatch(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	seedOneDoc(ctx, t, s)

	tmpDir := t.TempDir()
	run, err := backup.ExportToFile(ctx, s, tmpDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	manifest, err := backup.ReadManifest(run.ManifestPath)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	manifest.Counts["docs"] = manifest.Counts["docs"] + 1 // claim one more doc than the bundle actually has
	if err := backup.WriteManifest(run.ManifestPath, manifest); err != nil {
		t.Fatalf("WriteManifest: %v", err)
	}

	res := backup.VerifyBundleFile(ctx, run.FilePath)
	if res.Status != "fail" {
		t.Fatalf("VerifyBundleFile status = %q, want fail", res.Status)
	}
}

func TestVerifyBundleFileMissingFile(t *testing.T) {
	res := backup.VerifyBundleFile(context.Background(), "/nonexistent/wiselabz-backup-x.json")
	if res.Status != "fail" {
		t.Fatalf("VerifyBundleFile status = %q, want fail", res.Status)
	}
}

func TestRunVerifyOnceRecordsResult(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	seedOneDoc(ctx, t, s)

	tmpDir := t.TempDir()
	if _, err := backup.ExportToFile(ctx, s, tmpDir); err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	backup.RunVerifyOnce(ctx, tmpDir, logger)

	results, err := backup.ListVerifications(tmpDir, 0)
	if err != nil {
		t.Fatalf("ListVerifications: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Status != "pass" {
		t.Errorf("results[0].Status = %q, want pass", results[0].Status)
	}
}

func TestRunVerifyOnceNoBundles(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	// Must not panic when the directory has no bundles yet.
	backup.RunVerifyOnce(context.Background(), tmpDir, logger)

	results, err := backup.ListVerifications(tmpDir, 0)
	if err != nil {
		t.Fatalf("ListVerifications: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

func TestListVerificationsNewestFirstAndLimit(t *testing.T) {
	tmpDir := t.TempDir()
	for i := 0; i < 3; i++ {
		res := backup.VerificationResult{ID: string(rune('a' + i)), Status: "pass"}
		if err := backup.RecordVerification(tmpDir, res); err != nil {
			t.Fatalf("RecordVerification: %v", err)
		}
	}

	results, err := backup.ListVerifications(tmpDir, 2)
	if err != nil {
		t.Fatalf("ListVerifications: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	if results[0].ID != "c" || results[1].ID != "b" {
		t.Errorf("results = %+v, want newest-first [c, b]", results)
	}
}

func TestLatestBundlePicksMostRecent(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	tmpDir := t.TempDir()

	var last backup.Run
	for i := 0; i < 3; i++ {
		run, err := backup.ExportToFile(ctx, s, tmpDir)
		if err != nil {
			t.Fatalf("ExportToFile: %v", err)
		}
		last = run
	}

	got, err := backup.LatestBundle(tmpDir)
	if err != nil {
		t.Fatalf("LatestBundle: %v", err)
	}
	if got != last.FilePath {
		t.Errorf("LatestBundle = %q, want %q", got, last.FilePath)
	}
}

func TestLatestBundleNoBundles(t *testing.T) {
	tmpDir := t.TempDir()
	if _, err := backup.LatestBundle(tmpDir); err == nil {
		t.Fatal("LatestBundle on an empty directory: got nil error, want an error")
	}
}
