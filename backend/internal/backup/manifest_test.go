package backup_test

import (
	"context"
	"os"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// TestManifestChecksumRoundTrip is the manifest/checksum round-trip test:
// ExportToFile must write a manifest sidecar whose checksum matches the
// bundle file's actual sha256, and whose counts match the exported records.
func TestManifestChecksumRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	tmpDir := t.TempDir()

	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "doc-1", Title: "Runbook", Content: "steps"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	run, err := backup.ExportToFile(ctx, s, tmpDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	if run.ManifestPath == "" {
		t.Fatal("Run.ManifestPath is empty")
	}
	if run.Checksum == "" {
		t.Fatal("Run.Checksum is empty")
	}
	wantManifestPath := backup.ManifestPath(run.FilePath)
	if run.ManifestPath != wantManifestPath {
		t.Errorf("Run.ManifestPath = %q, want %q", run.ManifestPath, wantManifestPath)
	}

	bundleData, err := os.ReadFile(run.FilePath)
	if err != nil {
		t.Fatalf("read bundle file: %v", err)
	}
	gotChecksum := backup.ChecksumBytes(bundleData)
	if gotChecksum != run.Checksum {
		t.Errorf("recomputed checksum = %q, want %q", gotChecksum, run.Checksum)
	}

	manifest, err := backup.ReadManifest(run.ManifestPath)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	if manifest.Checksum != gotChecksum {
		t.Errorf("manifest.Checksum = %q, want %q", manifest.Checksum, gotChecksum)
	}
	if manifest.Version != backup.ManifestVersion {
		t.Errorf("manifest.Version = %d, want %d", manifest.Version, backup.ManifestVersion)
	}
	if manifest.BundleVersion != backup.BundleVersion {
		t.Errorf("manifest.BundleVersion = %d, want %d", manifest.BundleVersion, backup.BundleVersion)
	}
	if manifest.AppVersion == "" {
		t.Error("manifest.AppVersion is empty")
	}
	if manifest.CreatedAt == "" {
		t.Error("manifest.CreatedAt is empty")
	}
	if got := manifest.Counts["docs"]; got != 1 {
		t.Errorf("manifest.Counts[docs] = %d, want 1", got)
	}
}

// TestImportFromFileVerifiesChecksum exercises the "verify the checksum on
// import" requirement: ImportFromFile must reject a bundle whose bytes no
// longer match its manifest checksum, and must accept an intact one.
func TestImportFromFileVerifiesChecksum(t *testing.T) {
	ctx := context.Background()
	src := newTestStore(t)
	if err := src.CreateDoc(ctx, &store.DocRecord{ID: "doc-1", Title: "Runbook", Content: "steps"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	tmpDir := t.TempDir()
	run, err := backup.ExportToFile(ctx, src, tmpDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	t.Run("intact bundle imports cleanly", func(t *testing.T) {
		dst := newTestStore(t)
		result, err := backup.ImportFromFile(ctx, dst, run.FilePath)
		if err != nil {
			t.Fatalf("ImportFromFile: %v", err)
		}
		if result.Docs.Imported != 1 {
			t.Errorf("Docs.Imported = %d, want 1", result.Docs.Imported)
		}
	})

	t.Run("corrupted bundle is rejected", func(t *testing.T) {
		corrupted := corruptFile(t, run.FilePath)
		dst := newTestStore(t)
		if _, err := backup.ImportFromFile(ctx, dst, corrupted); err == nil {
			t.Fatal("ImportFromFile on a corrupted bundle: got nil error, want a checksum mismatch error")
		}
	})

	t.Run("bundle without a manifest sidecar still imports (backward compatible)", func(t *testing.T) {
		noManifestDir := t.TempDir()
		bundleData, err := os.ReadFile(run.FilePath)
		if err != nil {
			t.Fatalf("read bundle: %v", err)
		}
		path := noManifestDir + "/wiselabz-backup-legacy.json"
		if err := os.WriteFile(path, bundleData, 0o600); err != nil {
			t.Fatalf("write legacy bundle: %v", err)
		}
		dst := newTestStore(t)
		result, err := backup.ImportFromFile(ctx, dst, path)
		if err != nil {
			t.Fatalf("ImportFromFile without manifest: %v", err)
		}
		if result.Docs.Imported != 1 {
			t.Errorf("Docs.Imported = %d, want 1", result.Docs.Imported)
		}
	})
}

// corruptFile copies path to a sibling file with one byte flipped, and
// returns the new path, leaving the original untouched.
func corruptFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if len(data) == 0 {
		t.Fatalf("%s is empty", path)
	}
	corrupted := append([]byte(nil), data...)
	corrupted[len(corrupted)/2] ^= 0xFF

	dst := path + ".corrupt.json"
	if err := os.WriteFile(dst, corrupted, 0o600); err != nil {
		t.Fatalf("write corrupted bundle: %v", err)
	}
	manifestData, err := os.ReadFile(backup.ManifestPath(path))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if err := os.WriteFile(backup.ManifestPath(dst), manifestData, 0o600); err != nil {
		t.Fatalf("write manifest copy: %v", err)
	}
	return dst
}
