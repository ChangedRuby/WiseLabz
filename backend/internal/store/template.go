package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TemplateRecord represents a row in the templates table.
type TemplateRecord struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	AppliesTo      string `json:"appliesTo"`
	CurrentVersion int    `json:"currentVersion"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// TemplateVersionRecord represents a versioned snapshot of a template.
type TemplateVersionRecord struct {
	ID          string `json:"id"`
	TemplateID  string `json:"templateId"`
	Rev         int    `json:"rev"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AppliesTo   string `json:"appliesTo"`
	Sections    string `json:"sections"`
	Author      string `json:"author"`
	Trigger     string `json:"trigger"`
	CreatedAt   string `json:"createdAt"`
}

// TemplateVersionSection is a section stored in a template version snapshot.
type TemplateVersionSection struct {
	Title string `json:"title"`
	Order int    `json:"order"`
	Body  string `json:"body"`
}

// TemplateSectionRecord represents a row in the template_sections table.
type TemplateSectionRecord struct {
	ID         string `json:"id"`
	TemplateID string `json:"templateId"`
	Title      string `json:"title"`
	Ord        int    `json:"order"`
	Body       string `json:"body"`
}

// --- Template CRUD ---

// CreateTemplate inserts a new documentation template.
func (s *Store) CreateTemplate(ctx context.Context, t *TemplateRecord) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if t.CreatedAt == "" {
		t.CreatedAt = now
	}
	if t.UpdatedAt == "" {
		t.UpdatedAt = now
	}
	t.CurrentVersion = 1

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO templates (id, name, description, applies_to, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, t.ID, t.Name, t.Description, nilToStr(t.AppliesTo), t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	return nil
}

// ExistingTemplateIDs returns the subset of ids that already exist as
// templates, in one query — used by backup import to check N records
// without N round-trips.
func (s *Store) ExistingTemplateIDs(ctx context.Context, ids []string) (map[string]bool, error) {
	return existingIDs(ctx, s.db, "templates", ids)
}

// GetTemplate retrieves a single template record by ID.
func (s *Store) GetTemplate(ctx context.Context, id string) (*TemplateRecord, error) {
	t := &TemplateRecord{}
	var appliesTo sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, description, applies_to, current_version, created_at, updated_at
		FROM templates WHERE id = ?
	`, id).Scan(&t.ID, &t.Name, &t.Description, &appliesTo, &t.CurrentVersion, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}
	t.AppliesTo = appliesTo.String
	return t, nil
}

// UpdateTemplate updates the name and description of a template.
func (s *Store) UpdateTemplate(ctx context.Context, id string, updates map[string]any) error {
	now := time.Now().UTC().Format(time.RFC3339)
	args := []any{now}
	var parts []string
	for k, v := range updates {
		switch k {
		case "name":
			parts = append(parts, "name = ?")
			args = append(args, v)
		case "description":
			parts = append(parts, "description = ?")
			args = append(args, v)
		case "applies_to":
			parts = append(parts, "applies_to = ?")
			args = append(args, v)
		}
	}
	query := "UPDATE templates SET updated_at = ?, current_version = current_version + 1"
	if len(parts) > 0 {
		query += ", " + strings.Join(parts, ", ")
	}
	query += " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update template: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteTemplate removes a template record by ID.
func (s *Store) DeleteTemplate(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM templates WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete template: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func scanTemplate(row rowScanner) (TemplateRecord, error) {
	var t TemplateRecord
	var appliesTo sql.NullString
	err := row.Scan(&t.ID, &t.Name, &t.Description, &appliesTo, &t.CurrentVersion, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return TemplateRecord{}, err
	}
	t.AppliesTo = appliesTo.String
	return t, nil
}

// ListTemplates returns a paginated list of template records.
func (s *Store) ListTemplates(ctx context.Context, offset, limit int) ([]TemplateRecord, int, error) {
	return paginatedQuery(ctx, s.db, "templates",
		"id, name, description, applies_to, current_version, created_at, updated_at",
		"", nil, "created_at DESC", limit, offset, scanTemplate)
}

// --- Template versions ---

// CreateTemplateVersion inserts a new version record for a template.
func (s *Store) CreateTemplateVersion(ctx context.Context, v *TemplateVersionRecord) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	if v.CreatedAt == "" {
		v.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO template_versions
			(id, template_id, rev, name, description, applies_to, sections, author, trigger, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, v.ID, v.TemplateID, v.Rev, v.Name, v.Description, nilToStr(v.AppliesTo),
		v.Sections, nilToStr(v.Author), v.Trigger, v.CreatedAt)
	if err != nil {
		return fmt.Errorf("create template version: %w", err)
	}
	return nil
}

// GetTemplateVersions returns all version records for a template, newest first.
func (s *Store) GetTemplateVersions(ctx context.Context, templateID string) ([]TemplateVersionRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, template_id, rev, name, description, applies_to, sections, author, trigger, created_at
		FROM template_versions WHERE template_id = ? ORDER BY rev DESC
	`, templateID)
	if err != nil {
		return nil, fmt.Errorf("get template versions: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var versions []TemplateVersionRecord
	for rows.Next() {
		var v TemplateVersionRecord
		var appliesTo, author sql.NullString
		if err := rows.Scan(&v.ID, &v.TemplateID, &v.Rev, &v.Name, &v.Description,
			&appliesTo, &v.Sections, &author, &v.Trigger, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		v.AppliesTo = appliesTo.String
		v.Author = author.String
		versions = append(versions, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate template versions: %w", err)
	}
	if versions == nil {
		versions = []TemplateVersionRecord{}
	}
	return versions, nil
}

// --- Template sections ---

// CreateTemplateSection inserts a new section into a template.
func (s *Store) CreateTemplateSection(ctx context.Context, sec *TemplateSectionRecord) error {
	if sec.ID == "" {
		sec.ID = uuid.New().String()
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO template_sections (id, template_id, title, ord, body)
		VALUES (?, ?, ?, ?, ?)
	`, sec.ID, sec.TemplateID, sec.Title, sec.Ord, sec.Body)
	if err != nil {
		return fmt.Errorf("create template section: %w", err)
	}
	return nil
}

// GetTemplateSections returns all sections for a template, ordered by position.
func (s *Store) GetTemplateSections(ctx context.Context, templateID string) ([]TemplateSectionRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, template_id, title, ord, body
		FROM template_sections WHERE template_id = ? ORDER BY ord
	`, templateID)
	if err != nil {
		return nil, fmt.Errorf("get template sections: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var sections []TemplateSectionRecord
	for rows.Next() {
		var sec TemplateSectionRecord
		if err := rows.Scan(&sec.ID, &sec.TemplateID, &sec.Title, &sec.Ord, &sec.Body); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		sections = append(sections, sec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate template sections: %w", err)
	}
	if sections == nil {
		sections = []TemplateSectionRecord{}
	}
	return sections, nil
}

// GetAllTemplateSections returns sections for every given template in one
// query (ordered per template), replacing a GetTemplateSections call per
// template.
func (s *Store) GetAllTemplateSections(ctx context.Context, templateIDs []string) ([]TemplateSectionRecord, error) {
	sections := []TemplateSectionRecord{}
	if len(templateIDs) == 0 {
		return sections, nil
	}
	args := make([]any, len(templateIDs))
	for i, id := range templateIDs {
		args[i] = id
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, template_id, title, ord, body
		FROM template_sections WHERE template_id IN (`+placeholders(len(templateIDs))+`) ORDER BY template_id, ord
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("get all template sections: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var sec TemplateSectionRecord
		if err := rows.Scan(&sec.ID, &sec.TemplateID, &sec.Title, &sec.Ord, &sec.Body); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		sections = append(sections, sec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate template sections: %w", err)
	}
	return sections, nil
}

// DeleteTemplateSections removes all sections belonging to a template.
func (s *Store) DeleteTemplateSections(ctx context.Context, templateID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM template_sections WHERE template_id = ?`, templateID)
	return err
}
