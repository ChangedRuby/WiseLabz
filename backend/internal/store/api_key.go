package store

import (
	"context"
	"database/sql"
	"encoding/json"
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
	// Scope is auth.APIKeyScopeFull or auth.APIKeyScopeRead (#278).
	Scope string `json:"scope"`
	// ConnectorIDs, when non-empty, restricts the key to these connectors.
	ConnectorIDs []string `json:"connectorIds"`
}

const apiKeyColumns = `id, user_id, name, token_hash, role, created_at, expires_at, last_used_at, revoked_at, scope, connector_ids`

func scanAPIKey(row rowScanner, key *APIKey) error {
	var connectorIDs string
	if err := row.Scan(&key.ID, &key.UserID, &key.Name, &key.TokenHash, &key.Role,
		&key.CreatedAt, &key.ExpiresAt, &key.LastUsedAt, &key.RevokedAt, &key.Scope, &connectorIDs); err != nil {
		return err
	}
	ids, err := decodeConnectorIDs(connectorIDs)
	if err != nil {
		return err
	}
	key.ConnectorIDs = ids
	return nil
}

func decodeConnectorIDs(raw string) ([]string, error) {
	ids := make([]string, 0)
	if raw == "" {
		return ids, nil
	}
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil, fmt.Errorf("decode api key connector ids: %w", err)
	}
	return ids, nil
}

// CreateAPIKey inserts a new API key. An empty Scope means full access and
// nil ConnectorIDs means no connector restriction.
func (s *Store) CreateAPIKey(ctx context.Context, key *APIKey) error {
	if key.ID == "" {
		key.ID = uuid.New().String()
	}
	if key.CreatedAt == "" {
		key.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if key.Scope == "" {
		key.Scope = auth.APIKeyScopeFull
	}
	if key.ConnectorIDs == nil {
		key.ConnectorIDs = []string{}
	}
	connectorIDs, err := json.Marshal(key.ConnectorIDs)
	if err != nil {
		return fmt.Errorf("encode api key connector ids: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO api_keys (`+apiKeyColumns+`)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, key.ID, key.UserID, key.Name, key.TokenHash, key.Role, key.CreatedAt,
		key.ExpiresAt, key.LastUsedAt, key.RevokedAt, key.Scope, string(connectorIDs))
	if err != nil {
		return fmt.Errorf("create api key: %w", err)
	}
	return nil
}

// GetAPIKeyByHash retrieves an API key by its stored token hash.
func (s *Store) GetAPIKeyByHash(ctx context.Context, tokenHash string) (*APIKey, error) {
	key := &APIKey{}
	err := scanAPIKey(s.db.QueryRowContext(ctx, `
		SELECT `+apiKeyColumns+`
		FROM api_keys WHERE token_hash = ?
	`, tokenHash), key)
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
	err := scanAPIKey(s.db.QueryRowContext(ctx, `
		SELECT `+apiKeyColumns+`
		FROM api_keys WHERE id = ?
	`, id), key)
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
// instance-admin role always reflects the user's current role, and a
// disabled user's keys are rejected even if the key itself is still active.
func (s *Store) LookupAPIKey(ctx context.Context, tokenHash string) (*auth.APIKeyClaims, error) {
	claims := &auth.APIKeyClaims{}
	var role, scope, connectorIDs string
	var disabled int
	err := s.db.QueryRowContext(ctx, `
		SELECT api_keys.id, api_keys.user_id, users.instance_admin_role, api_keys.expires_at, api_keys.revoked_at, users.disabled, api_keys.last_used_at,
			api_keys.scope, api_keys.connector_ids
		FROM api_keys
		JOIN users ON users.id = api_keys.user_id
		WHERE api_keys.token_hash = ?
	`, tokenHash).Scan(&claims.KeyID, &claims.UserID, &role, &claims.ExpiresAt, &claims.RevokedAt, &disabled, &claims.LastUsedAt,
		&scope, &connectorIDs)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("lookup api key: %w", err)
	}
	if disabled != 0 {
		return nil, ErrNotFound
	}
	claims.InstanceAdmin = role == "admin"
	ids, err := decodeConnectorIDs(connectorIDs)
	if err != nil {
		return nil, err
	}
	// Anything but an explicit "full" scope is treated as read-only, so an
	// unexpected value fails closed.
	claims.Restriction = auth.APIKeyRestriction{
		ReadOnly:     scope != auth.APIKeyScopeFull,
		ConnectorIDs: ids,
	}
	return claims, nil
}

// ListAPIKeysForUser returns all API keys owned by a user, newest first.
func (s *Store) ListAPIKeysForUser(ctx context.Context, userID string) ([]APIKey, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+apiKeyColumns+`
		FROM api_keys WHERE user_id = ? ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	keys := make([]APIKey, 0)
	for rows.Next() {
		var key APIKey
		if err := scanAPIKey(rows, &key); err != nil {
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

// TouchAPIKeyLastUsed records successful authentication at most once a minute.
// The predicate also handles concurrent requests that read the same old timestamp.
func (s *Store) TouchAPIKeyLastUsed(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := s.db.ExecContext(ctx, `
		UPDATE api_keys SET last_used_at = ? WHERE id = ? AND (last_used_at = '' OR last_used_at <= ?)
	`, now.Format(time.RFC3339), id, now.Add(-time.Minute).Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("touch api key last used: %w", err)
	}
	return nil
}
