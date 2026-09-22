package backup

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// DefaultVerifyCronExpr is the cron expression the scheduled verify job
// registers with by default (daily, after the default backup window).
// Unlike the backup schedule, this isn't currently operator-configurable —
// there's no persisted schedule row for it, matching the "quality"/"digest"
// jobs in cmd/server/main.go.
const DefaultVerifyCronExpr = "0 4 * * *"

// bundleFilePrefix and bundleFileSuffix bracket the timestamp in filenames
// written by ExportToFile ("wiselabz-backup-<ts>.json"). Kept in sync with
// the format string there.
const (
	bundleFilePrefix   = "wiselabz-backup-"
	bundleFileSuffix   = ".json"
	manifestFileSuffix = ".manifest.json"
)

// VerificationResult is the outcome of restoring a bundle into a scratch
// database and comparing its row counts against the bundle's manifest.
type VerificationResult struct {
	ID             string         `json:"id"`
	BundlePath     string         `json:"bundlePath"`
	Status         string         `json:"status"` // "pass" or "fail"
	Error          string         `json:"error,omitempty"`
	Checksum       string         `json:"checksumSha256,omitempty"`
	ExpectedCounts map[string]int `json:"expectedCounts,omitempty"`
	ActualCounts   map[string]int `json:"actualCounts,omitempty"`
	StartedAt      string         `json:"startedAt"`
	FinishedAt     string         `json:"finishedAt"`
}

// LatestBundle returns the most recently created wiselabz-backup-*.json
// file in dir. Bundle filenames embed a sortable "20060102-150405" UTC
// timestamp, so a lexicographic sort orders them chronologically.
func LatestBundle(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read backup directory: %w", err)
	}
	var candidates []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		// Manifest sidecars also end in ".json" (see ManifestPath) —
		// exclude them explicitly rather than just matching the bundle
		// suffix.
		if strings.HasPrefix(name, bundleFilePrefix) && strings.HasSuffix(name, bundleFileSuffix) && !strings.HasSuffix(name, manifestFileSuffix) {
			candidates = append(candidates, name)
		}
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("no backup bundles found in %s", dir)
	}
	sort.Strings(candidates)
	return filepath.Join(dir, candidates[len(candidates)-1]), nil
}

// VerifyBundleFile checksums bundlePath against its manifest sidecar, then
// restores it into a throwaway in-memory SQLite database and compares the
// imported row counts against the manifest's recorded counts. It never
// touches the real database. This is the shared logic behind the scheduled
// verify job and the `backup verify`/`backup restore` CLI subcommands.
func VerifyBundleFile(ctx context.Context, bundlePath string) VerificationResult {
	res := VerificationResult{
		ID:         uuid.New().String(),
		BundlePath: bundlePath,
		StartedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	data, err := os.ReadFile(bundlePath)
	if err != nil {
		return failVerification(res, fmt.Errorf("read bundle file: %w", err))
	}
	res.Checksum = ChecksumBytes(data)

	manifest, mErr := ReadManifest(ManifestPath(bundlePath))
	hasManifest := mErr == nil
	switch {
	case hasManifest:
		if manifest.Checksum != res.Checksum {
			return failVerification(res, fmt.Errorf("checksum mismatch: manifest expects %s, bundle file has %s (corrupted or modified)", manifest.Checksum, res.Checksum))
		}
		res.ExpectedCounts = manifest.Counts
	case errors.Is(mErr, os.ErrNotExist):
		// No manifest sidecar (e.g. a bundle exported before this feature).
		// Nothing to check the checksum against; row counts fall back to the
		// bundle's own contents below.
	default:
		return failVerification(res, fmt.Errorf("read manifest: %w", mErr))
	}

	var bundle Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return failVerification(res, fmt.Errorf("parse bundle: %w", err))
	}
	if err := ValidateBundle(&bundle); err != nil {
		return failVerification(res, fmt.Errorf("validate bundle: %w", err))
	}
	if !hasManifest {
		res.ExpectedCounts = BundleCounts(&bundle)
	}

	scratch, err := newScratchStore(ctx)
	if err != nil {
		return failVerification(res, fmt.Errorf("create scratch database: %w", err))
	}
	defer func() { _ = scratch.Close() }()

	importResult, err := Import(ctx, scratch, &bundle)
	if err != nil {
		return failVerification(res, fmt.Errorf("import into scratch database: %w", err))
	}

	res.ActualCounts = map[string]int{
		"connectors":       importResult.Connectors.Imported,
		"docs":             importResult.Docs.Imported,
		"docVersions":      importResult.DocVersions.Imported,
		"templates":        importResult.Templates.Imported,
		"templateSections": importResult.TemplateSections.Imported,
	}
	for entity, want := range res.ExpectedCounts {
		if got := res.ActualCounts[entity]; got != want {
			return failVerification(res, fmt.Errorf("row count mismatch for %s: expected %d, restored %d", entity, want, got))
		}
	}

	res.Status = "pass"
	res.FinishedAt = time.Now().UTC().Format(time.RFC3339)
	return res
}

func failVerification(res VerificationResult, err error) VerificationResult {
	res.Status = "fail"
	res.Error = err.Error()
	res.FinishedAt = time.Now().UTC().Format(time.RFC3339)
	return res
}

// newScratchStore opens a fresh in-memory SQLite database migrated to the
// current schema, for restore verification only. It never touches disk;
// callers must Close() it when done.
func newScratchStore(ctx context.Context) (*store.Store, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("open scratch database: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping scratch database: %w", err)
	}
	if err := store.RunMigrations(db, "sqlite", slog.New(slog.DiscardHandler)); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate scratch database: %w", err)
	}
	return store.New(db, "sqlite"), nil
}

// RunVerifyOnce verifies the latest backup bundle in dir, records the
// result (see RecordVerification), and logs pass/fail — following the same
// log-on-failure convention as other scheduled jobs (e.g.
// internal/retention.RunCleanupOnce). It's the function the scheduled
// "backup-verify" cron job and the server wiring in cmd/server/main.go call.
func RunVerifyOnce(ctx context.Context, dir string, logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}

	path, err := LatestBundle(dir)
	if err != nil {
		logger.Warn("backup verify: no backup bundle to verify", "dir", dir, "error", err)
		return
	}

	res := VerifyBundleFile(ctx, path)
	if err := RecordVerification(dir, res); err != nil {
		logger.Error("backup verify: record result", "error", err)
	}

	if res.Status == "pass" {
		logger.Info("backup verify: passed", "bundle", path, "checksum", res.Checksum)
		return
	}
	logger.Error("backup verify: failed", "bundle", path, "error", res.Error)
}
