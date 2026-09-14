package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/google/uuid"
)

// APIKey represents a row in the api_keys table. TokenHash is never returned
// by an API handler; it is only used at the authentication boundary.
type APIKey struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	Name       string `json:"name"`
	TokenHash  string `json:"-"`
	Role       string `json:"role"`
	CreatedAt  string `json:"createdAt"`
	ExpiresAt  string `json:"expiresAt"`
	LastUsedAt string `json:"lastUsedAt"`
	RevokedAt  string `json:"revokedAt"`
}

// CreateAPIKey inserts a new API key.
func (s *Store) CreateAPIKey(ctx context.Context, key *APIKey) error {
	if key.ID == "" {
		key.ID = uuid.New().String()
	}
	if key.CreatedAt == "" {
		key.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO api_keys (id, user_id, name, token_hash, role, created_at, expires_at, last_used_at, revoked_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, key.ID, key.UserID, key.Name, key.TokenHash, key.Role, key.CreatedAt,
		key.ExpiresAt, key.LastUsedAt, key.RevokedAt)
	if err != nil {
		return fmt.Errorf("create api key: %w", err)
	}
	return nil
}

// GetAPIKeyByHash retrieves an API key by its stored token hash.
func (s *Store) GetAPIKeyByHash(ctx context.Context, tokenHash string) (*APIKey, error) {
	key := &APIKey{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, token_hash, role, created_at, expires_at, last_used_at, revoked_at
		FROM api_keys WHERE token_hash = ?
	`, tokenHash).Scan(&key.ID, &key.UserID, &key.Name, &key.TokenHash, &key.Role,
		&key.CreatedAt, &key.ExpiresAt, &key.LastUsedAt, &key.RevokedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get api key by hash: %w", err)
	}
	return key, nil
}

// GetAPIKeyByID retrieves an API key by ID.
func (s *Store) GetAPIKeyByID(ctx context.Context, id string) (*APIKey, error) {
	key := &APIKey{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, token_hash, role, created_at, expires_at, last_used_at, revoked_at
		FROM api_keys WHERE id = ?
	`, id).Scan(&key.ID, &key.UserID, &key.Name, &key.TokenHash, &key.Role,
		&key.CreatedAt, &key.ExpiresAt, &key.LastUsedAt, &key.RevokedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get api key by id: %w", err)
	}
	return key, nil
}

// LookupAPIKey adapts the stored key to the auth middleware without coupling
// the auth package to the store package. It joins users so a key's effective
// role always reflects the user's current role, and a disabled user's keys
// are rejected even if the key itself is still active.
func (s *Store) LookupAPIKey(ctx context.Context, tokenHash string) (*auth.APIKeyClaims, error) {
	claims := &auth.APIKeyClaims{}
	var disabled int
	err := s.db.QueryRowContext(ctx, `
		SELECT api_keys.id, api_keys.user_id, users.role, api_keys.expires_at, api_keys.revoked_at, users.disabled
		FROM api_keys
		JOIN users ON users.id = api_keys.user_id
		WHERE api_keys.token_hash = ?
	`, tokenHash).Scan(&claims.KeyID, &claims.UserID, &claims.Role, &claims.ExpiresAt, &claims.RevokedAt, &disabled)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("lookup api key: %w", err)
	}
	if disabled != 0 {
		return nil, ErrNotFound
	}
	return claims, nil
}

// ListAPIKeysForUser returns all API keys owned by a user, newest first.
func (s *Store) ListAPIKeysForUser(ctx context.Context, userID string) ([]APIKey, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, name, token_hash, role, created_at, expires_at, last_used_at, revoked_at
		FROM api_keys WHERE user_id = ? ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	keys := make([]APIKey, 0)
	for rows.Next() {
		var key APIKey
		if err := rows.Scan(&key.ID, &key.UserID, &key.Name, &key.TokenHash, &key.Role,
			&key.CreatedAt, &key.ExpiresAt, &key.LastUsedAt, &key.RevokedAt); err != nil {
			return nil, fmt.Errorf("scan api key: %w", err)
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate api keys: %w", err)
	}
	return keys, nil
}

// RevokeAPIKey marks an API key as revoked.
func (s *Store) RevokeAPIKey(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE api_keys SET revoked_at = ? WHERE id = ?
	`, time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("revoke api key rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// RevokeAllAPIKeysForUser revokes every active API key owned by a user, e.g.
// when the account is disabled.
func (s *Store) RevokeAllAPIKeysForUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE api_keys SET revoked_at = ? WHERE user_id = ? AND revoked_at = ''
	`, time.Now().UTC().Format(time.RFC3339), userID)
	if err != nil {
		return fmt.Errorf("revoke all api keys for user: %w", err)
	}
	return nil
}

// TouchAPIKeyLastUsed records the most recent successful authentication.
func (s *Store) TouchAPIKeyLastUsed(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE api_keys SET last_used_at = ? WHERE id = ?
	`, time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("touch api key last used: %w", err)
	}
	return nil
}
