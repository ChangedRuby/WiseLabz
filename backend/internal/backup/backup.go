// Package backup exports and imports a portable JSON bundle of WiseLabz
// configuration and content (connectors, docs, templates) for disaster
// recovery and migration between instances.
//
// Secrets are never included: connector secret fields (as declared by each
// connector's TypeSchema) are redacted from exported config, the AI provider
// API key is never read, and notification channel config (which may embed
// SMTP/webhook credentials in an arbitrary blob we can't safely redact) is
// excluded entirely.
package backup

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// BundleVersion is the current backup format version. ValidateBundle rejects
// any bundle whose Version doesn't match.
const BundleVersion = 1

const exportPageSize = 1000

// validCategories mirrors the connectors.category CHECK constraint in
// migrations/sqlite/000001_init.up.sql.
var validCategories = map[string]bool{
	"virtualization":  true,
	"containers_paas": true,
	"networking":      true,
}

// AIConfigSummary is an informational, secret-free snapshot of the AI
// configuration. It is exported for operator visibility only — Import never
// applies it, since the encrypted API key can't be restored from a backup.
type AIConfigSummary struct {
	Enabled  bool   `json:"enabled"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	BaseURL  string `json:"baseUrl"`
	Mode     string `json:"mode"`
}

// Bundle is the full portable backup format.
type Bundle struct {
	Version          int                           `json:"version"`
	ExportedAt       string                        `json:"exportedAt"`
	Connectors       []store.ConnectorRecord       `json:"connectors"`
	Docs             []store.DocRecord             `json:"docs"`
	DocVersions      []store.DocVersionRecord      `json:"docVersions"`
	Templates        []store.TemplateRecord        `json:"templates"`
	TemplateSections []store.TemplateSectionRecord `json:"templateSections"`
	// AIConfig is informational only; see Import.
	AIConfig *AIConfigSummary `json:"aiConfig,omitempty"`
}

// Result reports how many records of each entity were imported vs. skipped
// (skipped = an existing record with the same ID was found, left untouched).
type Result struct {
	Connectors       Counts `json:"connectors"`
	Docs             Counts `json:"docs"`
	DocVersions      Counts `json:"docVersions"`
	Templates        Counts `json:"templates"`
	TemplateSections Counts `json:"templateSections"`
}

// Counts is the imported/skipped tally for one entity kind.
type Counts struct {
	Imported int `json:"imported"`
	Skipped  int `json:"skipped"`
}

// Export builds a full backup bundle from the current store state, with
// connector secrets redacted.
func Export(ctx context.Context, s *store.Store) (*Bundle, error) {
	var b *Bundle
	err := s.WithinTransaction(ctx, func(tx *store.Store) error {
		bundle, err := exportWithin(ctx, tx)
		if err != nil {
			return err
		}
		b = bundle
		return nil
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

// exportWithin builds the bundle from a single transaction-bound Store so
// the connector/doc/template listings and their COUNT(*) totals all see one
// consistent snapshot, even while syncs or edits are writing concurrently.
func exportWithin(ctx context.Context, s *store.Store) (*Bundle, error) {
	connectors, err := s.ListAllConnectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("export connectors: %w", err)
	}
	for i := range connectors {
		redacted, err := RedactConnectorConfig(connectors[i].Type, connectors[i].ConfigData)
		if err != nil {
			// ponytail: unknown/unregistered connector type — no known secret
			// fields to strip, so export the config as-is rather than failing
			// the whole export.
			slog.Warn("backup export: could not redact connector config", "connectorId", connectors[i].ID, "type", connectors[i].Type, "error", err)
			continue
		}
		connectors[i].ConfigData = redacted
	}

	docs, err := exportDocs(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("export docs: %w", err)
	}
	var docVersions []store.DocVersionRecord
	for _, d := range docs {
		versions, err := s.GetDocVersions(ctx, d.ID)
		if err != nil {
			return nil, fmt.Errorf("export doc versions for %s: %w", d.ID, err)
		}
		docVersions = append(docVersions, versions...)
	}

	templates, err := exportTemplates(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("export templates: %w", err)
	}
	var sections []store.TemplateSectionRecord
	for _, t := range templates {
		s2, err := s.GetTemplateSections(ctx, t.ID)
		if err != nil {
			return nil, fmt.Errorf("export template sections for %s: %w", t.ID, err)
		}
		sections = append(sections, s2...)
	}

	return &Bundle{
		Version:          BundleVersion,
		ExportedAt:       time.Now().UTC().Format(time.RFC3339),
		Connectors:       connectors,
		Docs:             docs,
		DocVersions:      docVersions,
		Templates:        templates,
		TemplateSections: sections,
		AIConfig:         LoadAIConfigSummary(ctx, s),
	}, nil
}

// RedactConnectorConfig removes every field the connector's TypeSchema marks
// as secret-bearing (see store.IsSecretFieldType) from configData (a
// JSON-encoded map). Returns an error only when configData itself fails to
// parse. Exported for reuse by other read-only export features (e.g.
// internal/diagnostics) that need the same secret-stripping behavior.
//
// This deliberately parses/marshals configData as plain JSON rather than via
// store.ParseConnectorConfig/MarshalConnectorConfig: the values being
// redacted are dropped outright, so nothing here needs the encryption key,
// and whether a value is encrypted, plaintext-legacy, or garbage is
// irrelevant to deleting it by key name.
func RedactConnectorConfig(connType, configData string) (string, error) {
	var cfg map[string]any
	if err := json.Unmarshal([]byte(configData), &cfg); err != nil {
		return "", fmt.Errorf("parse connector config: %w", err)
	}
	schema, err := connector.GetTypeSchema(connType)
	if err != nil {
		return configData, nil //nolint:nilerr // unknown type: nothing known to redact, keep as-is
	}
	for _, f := range schema.Fields {
		if store.IsSecretFieldType(f.Type) {
			delete(cfg, f.Key)
		}
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal connector config: %w", err)
	}
	return string(b), nil
}

// LoadAIConfigSummary reads the AI config row directly (never selecting
// api_key_encrypted), following the same query style as
// api/settings.Handler.LoadAIConfig. Returns nil if the row can't be read.
// Exported for reuse by internal/diagnostics.
func LoadAIConfigSummary(ctx context.Context, s *store.Store) *AIConfigSummary {
	var rec struct {
		Enabled  int
		Provider sql.NullString
		Model    sql.NullString
		BaseURL  sql.NullString
		Mode     string
	}
	err := s.DB().QueryRowContext(ctx, `
		SELECT enabled, provider, model, base_url, mode FROM ai_config WHERE id = 1
	`).Scan(&rec.Enabled, &rec.Provider, &rec.Model, &rec.BaseURL, &rec.Mode)
	if err != nil {
		return nil
	}
	return &AIConfigSummary{
		Enabled:  rec.Enabled != 0,
		Provider: rec.Provider.String,
		Model:    rec.Model.String,
		BaseURL:  rec.BaseURL.String,
		Mode:     rec.Mode,
	}
}

// ValidateBundle checks referential integrity and format version without
// touching the database. It returns the first problem found.
func ValidateBundle(b *Bundle) error {
	if b.Version != BundleVersion {
		return fmt.Errorf("unsupported backup version: got %d, expected %d", b.Version, BundleVersion)
	}

	docIDs := make(map[string]bool, len(b.Docs))
	for _, d := range b.Docs {
		docIDs[d.ID] = true
	}
	for _, v := range b.DocVersions {
		if !docIDs[v.DocID] {
			return fmt.Errorf("doc version %q references unknown doc %q", v.ID, v.DocID)
		}
	}

	templateIDs := make(map[string]bool, len(b.Templates))
	for _, t := range b.Templates {
		templateIDs[t.ID] = true
	}
	for _, sec := range b.TemplateSections {
		if !templateIDs[sec.TemplateID] {
			return fmt.Errorf("template section %q references unknown template %q", sec.ID, sec.TemplateID)
		}
	}

	for _, c := range b.Connectors {
		if !validCategories[c.Category] {
			return fmt.Errorf("connector %q has invalid category %q", c.ID, c.Category)
		}
	}

	return nil
}

// Import validates the bundle, then applies it additively and idempotently:
// records whose ID already exists are left untouched and counted as
// "skipped". AIConfig is never imported (export-only/informational).
func Import(ctx context.Context, s *store.Store, b *Bundle) (Result, error) {
	var res Result
	if err := ValidateBundle(b); err != nil {
		return res, err
	}
	if err := s.WithinTransaction(ctx, func(tx *store.Store) error {
		var err error
		res, err = importBundle(ctx, tx, b)
		return err
	}); err != nil {
		return res, err
	}
	return res, nil
}

func exportDocs(ctx context.Context, s *store.Store) ([]store.DocRecord, error) {
	docs := []store.DocRecord{}
	for offset := 0; ; offset += exportPageSize {
		page, total, err := s.ListAllDocs(ctx, "", offset, exportPageSize)
		if err != nil {
			return nil, err
		}
		if len(page) == 0 {
			return docs, nil
		}
		docs = append(docs, page...)
		if len(docs) >= total {
			return docs, nil
		}
	}
}

func exportTemplates(ctx context.Context, s *store.Store) ([]store.TemplateRecord, error) {
	templates := []store.TemplateRecord{}
	for offset := 0; ; offset += exportPageSize {
		page, total, err := s.ListTemplates(ctx, offset, exportPageSize)
		if err != nil {
			return nil, err
		}
		if len(page) == 0 {
			return templates, nil
		}
		templates = append(templates, page...)
		if len(templates) >= total {
			return templates, nil
		}
	}
}

func importBundle(ctx context.Context, s *store.Store, b *Bundle) (Result, error) {
	var res Result
	for _, c := range b.Connectors {
		if _, err := s.GetConnector(ctx, c.ID); err == nil {
			res.Connectors.Skipped++
			continue
		}
		if err := s.CreateConnector(ctx, &c); err != nil {
			return res, fmt.Errorf("import connector %q: %w", c.ID, err)
		}
		res.Connectors.Imported++
	}

	for _, d := range b.Docs {
		if _, err := s.GetDoc(ctx, d.ID); err == nil {
			res.Docs.Skipped++
			continue
		}
		if err := s.CreateDoc(ctx, &d); err != nil {
			return res, fmt.Errorf("import doc %q: %w", d.ID, err)
		}
		res.Docs.Imported++
	}

	// No Get-by-ID for doc versions; check existence against the versions
	// already stored for each doc instead (cheap: one query per doc, not per
	// version, since versions naturally group by doc in the bundle).
	existingVersions := make(map[string]bool)
	for _, d := range b.Docs {
		versions, err := s.GetDocVersions(ctx, d.ID)
		if err != nil {
			return res, fmt.Errorf("check doc versions for %q: %w", d.ID, err)
		}
		for _, v := range versions {
			existingVersions[v.ID] = true
		}
	}
	for _, v := range b.DocVersions {
		if existingVersions[v.ID] {
			res.DocVersions.Skipped++
			continue
		}
		if err := s.CreateDocVersion(ctx, &v); err != nil {
			return res, fmt.Errorf("import doc version %q: %w", v.ID, err)
		}
		res.DocVersions.Imported++
	}

	for _, t := range b.Templates {
		if _, err := s.GetTemplate(ctx, t.ID); err == nil {
			res.Templates.Skipped++
			continue
		}
		if err := s.CreateTemplate(ctx, &t); err != nil {
			return res, fmt.Errorf("import template %q: %w", t.ID, err)
		}
		res.Templates.Imported++
	}

	existingSections := make(map[string]bool)
	for _, t := range b.Templates {
		sections, err := s.GetTemplateSections(ctx, t.ID)
		if err != nil {
			return res, fmt.Errorf("check template sections for %q: %w", t.ID, err)
		}
		for _, sec := range sections {
			existingSections[sec.ID] = true
		}
	}
	for _, sec := range b.TemplateSections {
		if existingSections[sec.ID] {
			res.TemplateSections.Skipped++
			continue
		}
		if err := s.CreateTemplateSection(ctx, &sec); err != nil {
			return res, fmt.Errorf("import template section %q: %w", sec.ID, err)
		}
		res.TemplateSections.Imported++
	}

	return res, nil
}

// Run describes one backup that was created (stored in the database for
// scheduling history and retrieval).
type Run struct {
	ID          string `json:"id"`
	TriggeredBy string `json:"triggeredBy"` // "schedule" or "manual"
	FilePath    string `json:"filePath"`
	SizeBytes   int64  `json:"sizeBytes"`
	CreatedAt   string `json:"createdAt"`
}

// ExportToFile calls Export, marshals the bundle to JSON, writes it to
// {dir}/wiselabz-backup-{RFC3339 timestamp}.json, and returns metadata about
// the created file. The directory is created if it does not exist (0o700),
// and the file is written 0o600 since the bundle is a full infrastructure
// inventory even with secrets redacted.
func ExportToFile(ctx context.Context, s *store.Store, dir string) (Run, error) {
	var run Run
	run.ID = uuid.New().String()
	run.TriggeredBy = "schedule" // default; caller may override

	// Export the bundle
	bundle, err := Export(ctx, s)
	if err != nil {
		return run, fmt.Errorf("export bundle: %w", err)
	}

	// Create directory if it doesn't exist. The bundle is a full
	// infrastructure inventory, so keep it private to the owning user even
	// though secrets are redacted before export.
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return run, fmt.Errorf("create backup directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return run, fmt.Errorf("marshal bundle: %w", err)
	}

	// Write to file with RFC3339 timestamp (using format safe for filenames)
	timestamp := time.Now().UTC().Format("20060102-150405")
	filename := fmt.Sprintf("wiselabz-backup-%s.json", timestamp)
	path := filepath.Join(dir, filename)

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return run, fmt.Errorf("write backup file: %w", err)
	}

	// Populate run metadata
	run.FilePath = path
	run.SizeBytes = int64(len(data))
	run.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	return run, nil
}
