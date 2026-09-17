package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ShareLink represents a row in the share_links table: a hashed, revocable
// token granting read-only access to a doc-tree subtree without an
// account. TokenHash is never returned by an API handler; it is only used
// at the authentication boundary, same convention as APIKey.
type ShareLink struct {
	ID             string `json:"id"`
	TokenHash      string `json:"-"`
	DocTreeRoot    string `json:"docTreeRoot"`
	CreatedBy      string `json:"createdBy"`
	CreatedAt      string `json:"createdAt"`
	ExpiresAt      string `json:"expiresAt"`
	RevokedAt      string `json:"revokedAt"`
	LastAccessedAt string `json:"lastAccessedAt"`
}

// CreateShareLink inserts a new share link.
func (s *Store) CreateShareLink(ctx context.Context, link *ShareLink) error {
	if link.ID == "" {
		link.ID = uuid.New().String()
	}
	if link.CreatedAt == "" {
		link.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO share_links (id, token_hash, doc_tree_root, created_by, created_at, expires_at, revoked_at, last_accessed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, link.ID, link.TokenHash, link.DocTreeRoot, link.CreatedBy, link.CreatedAt,
		link.ExpiresAt, link.RevokedAt, link.LastAccessedAt)
	if err != nil {
		return fmt.Errorf("create share link: %w", err)
	}
	return nil
}

// GetShareLinkByHash retrieves a share link by its stored token hash.
func (s *Store) GetShareLinkByHash(ctx context.Context, tokenHash string) (*ShareLink, error) {
	link := &ShareLink{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, token_hash, doc_tree_root, created_by, created_at, expires_at, revoked_at, last_accessed_at
		FROM share_links WHERE token_hash = ?
	`, tokenHash).Scan(&link.ID, &link.TokenHash, &link.DocTreeRoot, &link.CreatedBy, &link.CreatedAt,
		&link.ExpiresAt, &link.RevokedAt, &link.LastAccessedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get share link by hash: %w", err)
	}
	return link, nil
}

// GetShareLinkByID retrieves a share link by ID.
func (s *Store) GetShareLinkByID(ctx context.Context, id string) (*ShareLink, error) {
	link := &ShareLink{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, token_hash, doc_tree_root, created_by, created_at, expires_at, revoked_at, last_accessed_at
		FROM share_links WHERE id = ?
	`, id).Scan(&link.ID, &link.TokenHash, &link.DocTreeRoot, &link.CreatedBy, &link.CreatedAt,
		&link.ExpiresAt, &link.RevokedAt, &link.LastAccessedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get share link by id: %w", err)
	}
	return link, nil
}

// ListShareLinksCreatedBy returns every share link a user has created,
// newest first.
func (s *Store) ListShareLinksCreatedBy(ctx context.Context, userID string) ([]ShareLink, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, token_hash, doc_tree_root, created_by, created_at, expires_at, revoked_at, last_accessed_at
		FROM share_links WHERE created_by = ? ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list share links: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	links := make([]ShareLink, 0)
	for rows.Next() {
		var link ShareLink
		if err := rows.Scan(&link.ID, &link.TokenHash, &link.DocTreeRoot, &link.CreatedBy, &link.CreatedAt,
			&link.ExpiresAt, &link.RevokedAt, &link.LastAccessedAt); err != nil {
			return nil, fmt.Errorf("scan share link: %w", err)
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate share links: %w", err)
	}
	return links, nil
}

// RevokeShareLink marks a share link as revoked.
func (s *Store) RevokeShareLink(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE share_links SET revoked_at = ? WHERE id = ?
	`, time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("revoke share link: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("revoke share link rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// TouchShareLinkLastAccessed records the most recent successful view.
func (s *Store) TouchShareLinkLastAccessed(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE share_links SET last_accessed_at = ? WHERE id = ?
	`, time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("touch share link last accessed: %w", err)
	}
	return nil
}
