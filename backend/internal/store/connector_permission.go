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

// ConnectorGrant represents a row in the user_connector_roles table: a
// user's role (viewer/operator) scoped to a single connector. Absence of a
// grant means no access at all — see UserHasConnectorRole.
type ConnectorGrant struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	ConnectorID string `json:"connectorId"`
	Role        string `json:"role"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// connectorRoleRank orders connector roles so a higher role satisfies a
// lower minimum, mirroring auth.roleSatisfies' viewer <= operator ordering.
var connectorRoleRank = map[string]int{"viewer": 1, "operator": 2}

// GetUserConnectorRole returns the user's role on a connector, or "" if no
// grant exists (default deny).
func (s *Store) GetUserConnectorRole(ctx context.Context, userID, connectorID string) (string, error) {
	var role string
	err := s.db.QueryRowContext(ctx,
		`SELECT role FROM user_connector_roles WHERE user_id = ? AND connector_id = ?`,
		userID, connectorID,
	).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get user connector role: %w", err)
	}
	// An API key can narrow its owner's grant (#278), never widen it.
	return auth.ClampConnectorRole(ctx, connectorID, role), nil
}

// UserHasConnectorRole reports whether userID has at least minRole on
// connectorID. Implements the per-connector counterpart to roleSatisfies.
func (s *Store) UserHasConnectorRole(ctx context.Context, userID, connectorID, minRole string) (bool, error) {
	role, err := s.GetUserConnectorRole(ctx, userID, connectorID)
	if err != nil {
		return false, err
	}
	if role == "" {
		return false, nil
	}
	return connectorRoleRank[role] >= connectorRoleRank[minRole], nil
}

// ListConnectorGrants returns every user's grant on a connector, for the
// connector's Permissions tab.
func (s *Store) ListConnectorGrants(ctx context.Context, connectorID string) ([]ConnectorGrant, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, connector_id, role, created_at, updated_at
		FROM user_connector_roles WHERE connector_id = ? ORDER BY created_at ASC
	`, connectorID)
	if err != nil {
		return nil, fmt.Errorf("list connector grants: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	return scanConnectorGrants(rows)
}

// ListUserConnectorGrants returns every connector grant held by a user.
func (s *Store) ListUserConnectorGrants(ctx context.Context, userID string) ([]ConnectorGrant, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, connector_id, role, created_at, updated_at
		FROM user_connector_roles WHERE user_id = ? ORDER BY created_at ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list user connector grants: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	return scanConnectorGrants(rows)
}

func scanConnectorGrants(rows *sql.Rows) ([]ConnectorGrant, error) {
	grants := make([]ConnectorGrant, 0)
	for rows.Next() {
		var g ConnectorGrant
		if err := rows.Scan(&g.ID, &g.UserID, &g.ConnectorID, &g.Role, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan connector grant: %w", err)
		}
		grants = append(grants, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connector grants: %w", err)
	}
	return grants, nil
}

// apiKeyConnectorFilter returns an " AND <column> IN (...)" SQL fragment and
// its args narrowing a query to the connectors the request's API key is
// restricted to, or "" when the key has no connector restriction. Queries
// that join user_connector_roles directly (rather than going through
// GetUserConnectorRole) append it so key restrictions hold there too.
func apiKeyConnectorFilter(ctx context.Context, column string) (string, []any) {
	ids := auth.APIKeyRestrictionFromContext(ctx).ConnectorIDs
	if len(ids) == 0 {
		return "", nil
	}
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	return " AND " + column + " IN (" + placeholders(len(ids)) + ")", args
}

// FilterConnectorIDsByGrant narrows ids down to the ones userID has at
// least minRole on. Used by list endpoints to enforce default-deny instead
// of returning every connector's data to any authenticated user.
func (s *Store) FilterConnectorIDsByGrant(ctx context.Context, userID string, ids []string, minRole string) ([]string, error) {
	if len(ids) == 0 {
		return ids, nil
	}
	grants, err := s.ListUserConnectorGrants(ctx, userID)
	if err != nil {
		return nil, err
	}
	allowed := make(map[string]bool, len(grants))
	for _, g := range grants {
		role := auth.ClampConnectorRole(ctx, g.ConnectorID, g.Role)
		if role != "" && connectorRoleRank[role] >= connectorRoleRank[minRole] {
			allowed[g.ConnectorID] = true
		}
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if allowed[id] {
			out = append(out, id)
		}
	}
	return out, nil
}

// UpsertConnectorGrant creates or updates a user's role on a connector.
func (s *Store) UpsertConnectorGrant(ctx context.Context, userID, connectorID, role string) (ConnectorGrant, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	g := ConnectorGrant{
		ID:          uuid.New().String(),
		UserID:      userID,
		ConnectorID: connectorID,
		Role:        role,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO user_connector_roles (id, user_id, connector_id, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (user_id, connector_id) DO UPDATE SET role = excluded.role, updated_at = excluded.updated_at
	`, g.ID, g.UserID, g.ConnectorID, g.Role, g.CreatedAt, g.UpdatedAt)
	if err != nil {
		return ConnectorGrant{}, fmt.Errorf("upsert connector grant: %w", err)
	}
	// A conflict path doesn't return the pre-existing row's real ID/CreatedAt;
	// read back so callers (and the audit log) see the persisted values.
	return s.getConnectorGrant(ctx, userID, connectorID)
}

func (s *Store) getConnectorGrant(ctx context.Context, userID, connectorID string) (ConnectorGrant, error) {
	var g ConnectorGrant
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, connector_id, role, created_at, updated_at
		FROM user_connector_roles WHERE user_id = ? AND connector_id = ?
	`, userID, connectorID).Scan(&g.ID, &g.UserID, &g.ConnectorID, &g.Role, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return ConnectorGrant{}, fmt.Errorf("get connector grant: %w", err)
	}
	return g, nil
}

// DeleteConnectorGrant revokes a user's access to a connector.
func (s *Store) DeleteConnectorGrant(ctx context.Context, userID, connectorID string) error {
	result, err := s.db.ExecContext(ctx,
		`DELETE FROM user_connector_roles WHERE user_id = ? AND connector_id = ?`, userID, connectorID)
	if err != nil {
		return fmt.Errorf("delete connector grant: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete connector grant rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
