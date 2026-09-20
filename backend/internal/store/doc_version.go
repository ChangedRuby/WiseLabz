package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

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
