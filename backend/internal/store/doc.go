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
// Content is deliberately not loaded (Content is empty): callers render titles.
func (s *Store) ListDocsGroupedByService(ctx context.Context) (map[string][]DocRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+docSummaryColumns+`
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
		if err := rows.Scan(&d.ID, &d.Title, &d.Kind, &svcID, &d.CurrentVersion, &d.CreatedAt, &d.UpdatedAt); err != nil {
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

// docColumns is the full column list, including the (potentially large) content.
const docColumns = `id, title, kind, service_id, content, current_version, created_at, updated_at`

// docSummaryColumns omits content for list/tree views that never render it.
const docSummaryColumns = `id, title, kind, service_id, current_version, created_at, updated_at`

func scanDocSummary(row rowScanner) (DocRecord, error) {
	var d DocRecord
	var svcID sql.NullString
	if err := row.Scan(&d.ID, &d.Title, &d.Kind, &svcID, &d.CurrentVersion, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return DocRecord{}, err
	}
	d.ServiceID = svcID.String
	return d, nil
}

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

// ListAllDocs returns a paginated, optionally search-filtered list of all docs
// without their content (Content is empty). Use ListAllDocsWithContent when the
// body is needed.
func (s *Store) ListAllDocs(ctx context.Context, search string, offset, limit int) ([]DocRecord, int, error) {
	where, args := docSearchWhere(search)
	return paginatedQuery(ctx, s.db, "docs", docSummaryColumns, where, args, "updated_at DESC", limit, offset, scanDocSummary)
}

// ListAllDocsWithContent is like ListAllDocs but also loads each doc's content
// (e.g. for backup export).
func (s *Store) ListAllDocsWithContent(ctx context.Context, search string, offset, limit int) ([]DocRecord, int, error) {
	where, args := docSearchWhere(search)
	return paginatedQuery(ctx, s.db, "docs", docColumns, where, args, "updated_at DESC", limit, offset, scanDoc)
}

// likeEscaper escapes LIKE wildcards and the escape character itself.
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func escapeLike(s string) string { return likeEscaper.Replace(s) }

func docSearchWhere(search string) (string, []any) {
	where := "WHERE 1=1"
	var args []any
	if search != "" {
		where += ` AND LOWER(title) LIKE LOWER(?) ESCAPE '\'`
		args = append(args, "%"+escapeLike(search)+"%")
	}
	return where, args
}

// CountDocs returns total number of docs.
func (s *Store) CountDocs(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM docs`).Scan(&count)
	return count, err
}
