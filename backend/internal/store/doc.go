package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/ws"
)

// Doc lock tuning: TTL is how long an acquired lock stays valid without
// renewal; heartbeat is how often the editor UI is expected to re-acquire
// (renew) it. Locks are advisory-only — see AcquireDocLock.
const (
	DocLockTTL       = 5 * time.Minute
	DocLockHeartbeat = 30 * time.Second
)

// ErrLockHeldByOther is returned by AcquireDocLock when a live lock is held
// by a different user.
var ErrLockHeldByOther = errors.New("doc lock held by another user")

// DocRecord represents a row in the docs table.
type DocRecord struct {
	ID             string `json:"docId"`
	Title          string `json:"title"`
	Kind           string `json:"kind"`
	ServiceID      string `json:"serviceId"`
	Content        string `json:"content"`
	CurrentVersion int    `json:"currentVersion"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// DocLockRecord represents a row in the doc_locks table: an advisory,
// single-writer editing lock on a doc.
type DocLockRecord struct {
	DocID      string `json:"docId"`
	UserID     string `json:"userId"`
	AcquiredAt string `json:"acquiredAt"`
	ExpiresAt  string `json:"expiresAt"`
}

// DocVersionRecord represents a row in the doc_versions table.
type DocVersionRecord struct {
	ID        string `json:"id"`
	DocID     string `json:"docId"`
	Rev       int    `json:"rev"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Trigger   string `json:"trigger"`
	CreatedAt string `json:"createdAt"`
}

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

// --- Doc CRUD ---

// CreateDoc inserts a new documentation record.
func (s *Store) CreateDoc(ctx context.Context, d *DocRecord) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if d.CreatedAt == "" {
		d.CreatedAt = now
	}
	if d.UpdatedAt == "" {
		d.UpdatedAt = now
	}
	if d.CurrentVersion == 0 {
		d.CurrentVersion = 1
	}
	if d.Kind == "" {
		d.Kind = "lab"
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO docs (id, title, kind, service_id, content, current_version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, d.ID, d.Title, d.Kind, nilToStr(d.ServiceID), d.Content, d.CurrentVersion, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create doc: %w", err)
	}
	return nil
}

// ExistingDocIDs returns the subset of ids that already exist as docs, in
// one query — used by backup import to check N records without N round-trips.
func (s *Store) ExistingDocIDs(ctx context.Context, ids []string) (map[string]bool, error) {
	return existingIDs(ctx, s.db, "docs", ids)
}

// GetDoc retrieves a single documentation record by ID.
func (s *Store) GetDoc(ctx context.Context, id string) (*DocRecord, error) {
	d := &DocRecord{}
	var serviceID sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, kind, service_id, content, current_version, created_at, updated_at
		FROM docs WHERE id = ?
	`, id).Scan(&d.ID, &d.Title, &d.Kind, &serviceID, &d.Content,
		&d.CurrentVersion, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get doc: %w", err)
	}
	d.ServiceID = serviceID.String
	return d, nil
}

// UpdateDoc updates the content of a documentation record. If expectedVersion
// is non-nil, the update only applies when current_version matches it
// (optimistic concurrency); a stale version returns ErrVersionConflict instead
// of silently overwriting a newer edit.
func (s *Store) UpdateDoc(ctx context.Context, id, content string, expectedVersion *int) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE docs SET content = ?, updated_at = ?, current_version = current_version + 1 WHERE id = ?`
	args := []any{content, now, id}
	if expectedVersion != nil {
		query += ` AND current_version = ?`
		args = append(args, *expectedVersion)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update doc: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		if expectedVersion != nil {
			if _, getErr := s.GetDoc(ctx, id); getErr == nil {
				return ErrVersionConflict
			}
		}
		return ErrNotFound
	}
	return nil
}

// DeleteDoc removes a documentation record by ID.
func (s *Store) DeleteDoc(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM docs WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete doc: %w", err)
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

// ListDocsByService returns all documentation records for a given service.
func (s *Store) ListDocsByService(ctx context.Context, serviceID string) ([]DocRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, kind, service_id, content, current_version, created_at, updated_at
		FROM docs WHERE service_id = ? ORDER BY updated_at DESC
	`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list docs by service: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var docs []DocRecord
	for rows.Next() {
		var d DocRecord
		var svcID sql.NullString
		if err := rows.Scan(&d.ID, &d.Title, &d.Kind, &svcID, &d.Content,
			&d.CurrentVersion, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		d.ServiceID = svcID.String
		docs = append(docs, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate docs: %w", err)
	}
	if docs == nil {
		docs = []DocRecord{}
	}
	return docs, nil
}

// ListDocsGroupedByService returns every service doc grouped by connector ID.
func (s *Store) ListDocsGroupedByService(ctx context.Context) (map[string][]DocRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, kind, service_id, content, current_version, created_at, updated_at
		FROM docs WHERE service_id IS NOT NULL ORDER BY service_id, updated_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list docs grouped by service: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	docsByService := map[string][]DocRecord{}
	for rows.Next() {
		var d DocRecord
		var svcID sql.NullString
		if err := rows.Scan(&d.ID, &d.Title, &d.Kind, &svcID, &d.Content, &d.CurrentVersion, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		d.ServiceID = svcID.String
		docsByService[d.ServiceID] = append(docsByService[d.ServiceID], d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate docs: %w", err)
	}
	return docsByService, nil
}

// docColumns is the shared column list for the ListAllDocs SELECT.
const docColumns = `id, title, kind, service_id, content, current_version, created_at, updated_at`

func scanDoc(row rowScanner) (DocRecord, error) {
	var d DocRecord
	var svcID sql.NullString
	err := row.Scan(&d.ID, &d.Title, &d.Kind, &svcID, &d.Content, &d.CurrentVersion, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return DocRecord{}, err
	}
	d.ServiceID = svcID.String
	return d, nil
}

// ListAllDocs returns a paginated, optionally search-filtered list of all docs.
func (s *Store) ListAllDocs(ctx context.Context, search string, offset, limit int) ([]DocRecord, int, error) {
	where := "WHERE 1=1"
	var args []any
	if search != "" {
		where += " AND title LIKE ?"
		args = append(args, "%"+search+"%")
	}

	return paginatedQuery(ctx, s.db, "docs", docColumns, where, args, "updated_at DESC", limit, offset, scanDoc)
}

// --- Doc versions ---

// CreateDocVersion inserts a new version record for a document.
func (s *Store) CreateDocVersion(ctx context.Context, v *DocVersionRecord) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	if v.CreatedAt == "" {
		v.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO doc_versions (id, doc_id, rev, content, author, trigger, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, v.ID, v.DocID, v.Rev, v.Content, nilToStr(v.Author), v.Trigger, v.CreatedAt)
	if err != nil {
		return fmt.Errorf("create doc version: %w", err)
	}
	return nil
}

// GetDocVersions returns all version records for a document, newest first.
func (s *Store) GetDocVersions(ctx context.Context, docID string) ([]DocVersionRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, doc_id, rev, content, author, trigger, created_at
		FROM doc_versions WHERE doc_id = ? ORDER BY rev DESC
	`, docID)
	if err != nil {
		return nil, fmt.Errorf("get doc versions: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var versions []DocVersionRecord
	for rows.Next() {
		var v DocVersionRecord
		var author sql.NullString
		if err := rows.Scan(&v.ID, &v.DocID, &v.Rev, &v.Content, &author, &v.Trigger, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		v.Author = author.String
		versions = append(versions, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate doc versions: %w", err)
	}
	if versions == nil {
		versions = []DocVersionRecord{}
	}
	return versions, nil
}

// GetAllDocVersions returns version records for every given doc in one
// query (newest-per-doc first), replacing a GetDocVersions call per doc.
func (s *Store) GetAllDocVersions(ctx context.Context, docIDs []string) ([]DocVersionRecord, error) {
	versions := []DocVersionRecord{}
	if len(docIDs) == 0 {
		return versions, nil
	}
	args := make([]any, len(docIDs))
	for i, id := range docIDs {
		args[i] = id
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, doc_id, rev, content, author, trigger, created_at
		FROM doc_versions WHERE doc_id IN (`+placeholders(len(docIDs))+`) ORDER BY doc_id, rev DESC
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("get all doc versions: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var v DocVersionRecord
		var author sql.NullString
		if err := rows.Scan(&v.ID, &v.DocID, &v.Rev, &v.Content, &author, &v.Trigger, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		v.Author = author.String
		versions = append(versions, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate doc versions: %w", err)
	}
	return versions, nil
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

// CountDocs returns total number of docs.
func (s *Store) CountDocs(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM docs`).Scan(&count)
	return count, err
}

// --- Doc locks (advisory presence/soft-lock, issue #92) ---

// AcquireDocLock takes (or renews) the advisory lock on a doc for userID.
// Succeeds when the doc is unlocked, already expired, or already held by
// userID (renewal); otherwise returns the live lock plus ErrLockHeldByOther.
// This never blocks Save — it's a UI nudge only.
func (s *Store) AcquireDocLock(ctx context.Context, docID, userID string) (*DocLockRecord, error) {
	now := time.Now().UTC()
	acquiredAt := now.Format(time.RFC3339)
	expiresAt := now.Add(DocLockTTL).Format(time.RFC3339)

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO doc_locks (doc_id, user_id, acquired_at, expires_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(doc_id) DO UPDATE SET
			user_id = excluded.user_id,
			acquired_at = excluded.acquired_at,
			expires_at = excluded.expires_at
		WHERE doc_locks.expires_at < ? OR doc_locks.user_id = ?
	`, docID, userID, acquiredAt, expiresAt, acquiredAt, userID)
	if err != nil {
		return nil, fmt.Errorf("acquire doc lock: %w", err)
	}

	lock, err := s.GetDocLock(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("acquire doc lock: %w", err)
	}
	if lock.UserID != userID {
		return lock, ErrLockHeldByOther
	}
	return lock, nil
}

// ReleaseDocLock releases the lock on docID if held by userID; a no-op
// otherwise (not held, held by someone else, or already expired).
func (s *Store) ReleaseDocLock(ctx context.Context, docID, userID string) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM doc_locks WHERE doc_id = ? AND user_id = ?`, docID, userID); err != nil {
		return fmt.Errorf("release doc lock: %w", err)
	}
	return nil
}

// GetDocLock returns the live (non-expired) lock on docID, or ErrNotFound if
// unlocked/expired. Expiry is checked here (lazy), not dependent on a sweep.
func (s *Store) GetDocLock(ctx context.Context, docID string) (*DocLockRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	l := &DocLockRecord{}
	err := s.db.QueryRowContext(ctx, `
		SELECT doc_id, user_id, acquired_at, expires_at
		FROM doc_locks WHERE doc_id = ? AND expires_at > ?
	`, docID, now).Scan(&l.DocID, &l.UserID, &l.AcquiredAt, &l.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get doc lock: %w", err)
	}
	return l, nil
}

// RunDocLockSweep periodically clears expired doc locks and broadcasts their
// expiry so idle viewers' presence banners clear without polling. Expiry
// itself is already enforced lazily by GetDocLock; this sweep is purely row
// hygiene plus the WS notification. Mirrors quality.RunStaleSweep.
func RunDocLockSweep(ctx context.Context, s *Store, hub *ws.Hub, interval time.Duration, logger *slog.Logger) {
	logger.Info("Doc lock sweep started")
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		runDocLockSweep(ctx, s, hub, logger)
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
		}
	}
}

func runDocLockSweep(ctx context.Context, s *Store, hub *ws.Hub, logger *slog.Logger) {
	now := time.Now().UTC().Format(time.RFC3339)
	rows, err := s.db.QueryContext(ctx, `SELECT doc_id, user_id FROM doc_locks WHERE expires_at <= ?`, now)
	if err != nil {
		logger.Error("doc lock sweep: list expired", "error", err)
		return
	}
	var expired []DocLockRecord
	for rows.Next() {
		var l DocLockRecord
		if err := rows.Scan(&l.DocID, &l.UserID); err != nil {
			rows.Close() //nolint:errcheck
			logger.Error("doc lock sweep: scan", "error", err)
			return
		}
		expired = append(expired, l)
	}
	rows.Close() //nolint:errcheck
	if err := rows.Err(); err != nil {
		logger.Error("doc lock sweep: iterate expired", "error", err)
		return
	}

	if len(expired) == 0 {
		return
	}
	if _, err := s.db.ExecContext(ctx, `DELETE FROM doc_locks WHERE expires_at <= ?`, now); err != nil {
		logger.Error("doc lock sweep: delete", "error", err)
		return
	}
	for _, l := range expired {
		if hub != nil {
			hub.Broadcast(ws.EventDocLockExpired, map[string]any{"docId": l.DocID, "userId": l.UserID})
		}
	}
}
