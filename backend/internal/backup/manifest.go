package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

// ManifestVersion is the manifest format version, independent of
// BundleVersion (the manifest describes a bundle; it doesn't replace it).
const ManifestVersion = 1

// Manifest describes a backup bundle's contents and provenance. It is
// written as a sidecar file next to the bundle (see ManifestPath) so a
// checksum/row-count check can confirm the bundle restores cleanly without
// needing the original database.
type Manifest struct {
	Version       int            `json:"version"`
	BundleVersion int            `json:"bundleVersion"`
	AppVersion    string         `json:"appVersion"`
	SchemaVersion uint           `json:"schemaVersion"`
	CreatedAt     string         `json:"createdAt"`
	Checksum      string         `json:"checksumSha256"`
	Counts        map[string]int `json:"counts"`
}

// BundleCounts returns the per-entity row counts recorded in a bundle, used
// both to build a Manifest at export time and as a fallback expectation when
// verifying a bundle that predates the manifest sidecar.
func BundleCounts(b *Bundle) map[string]int {
	return map[string]int{
		"connectors":       len(b.Connectors),
		"docs":             len(b.Docs),
		"docVersions":      len(b.DocVersions),
		"templates":        len(b.Templates),
		"templateSections": len(b.TemplateSections),
	}
}

// ChecksumBytes returns the lowercase hex-encoded sha256 digest of data.
func ChecksumBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// BuildManifest builds the manifest for bundleData, the exact bytes written
// to the bundle file (so Checksum matches what a later read-back computes).
func BuildManifest(b *Bundle, bundleData []byte, appVersion string, schemaVersion uint) Manifest {
	return Manifest{
		Version:       ManifestVersion,
		BundleVersion: b.Version,
		AppVersion:    appVersion,
		SchemaVersion: schemaVersion,
		CreatedAt:     b.ExportedAt,
		Checksum:      ChecksumBytes(bundleData),
		Counts:        BundleCounts(b),
	}
}

// ManifestPath returns the manifest sidecar path for a bundle file path,
// e.g. ".../wiselabz-backup-20260101-000000.json" ->
// ".../wiselabz-backup-20260101-000000.manifest.json".
func ManifestPath(bundlePath string) string {
	return strings.TrimSuffix(bundlePath, ".json") + ".manifest.json"
}

// WriteManifest writes m as indented JSON to path (0o600: same privacy
// rationale as the bundle file it describes — row counts and schema version
// are still infrastructure detail).
func WriteManifest(path string, m Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

// ReadManifest reads and parses a manifest sidecar file. Returns a wrapped
// os.ErrNotExist-compatible error (checkable with errors.Is) when the
// sidecar doesn't exist, e.g. for backups created before this feature.
func ReadManifest(path string) (Manifest, error) {
	var m Manifest
	data, err := os.ReadFile(path)
	if err != nil {
		return m, fmt.Errorf("read manifest: %w", err)
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return m, fmt.Errorf("parse manifest: %w", err)
	}
	return m, nil
}

// AppVersion returns the running binary's module version, following the
// same fallback as api/system.Handler.Info ("dev" when build info isn't
// embedded, e.g. under `go run`).
func AppVersion() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok && buildInfo.Main.Version != "" {
		return buildInfo.Main.Version
	}
	return "dev"
}
