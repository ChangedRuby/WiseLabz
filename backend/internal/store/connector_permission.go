package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/google/uuid"
)

// ConnectorGrant represents a row in the user_connector_roles table: a
// user's role (viewer/operator) scoped to a single connector. Absence of a
// grant means no access at all — see UserHasConnectorRole.
//
// A (user, connector) pair may hold up to two rows, one per Source: a
// 'manual' grant an admin made by hand through the permissions API, and an
// 'oidc' grant synced from the user's IdP groups at login (#279 part 3).
// GetUserConnectorRole/FilterConnectorIDsByGrant treat the pair's effective
// role as the highest of the two.
type ConnectorGrant struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	ConnectorID string `json:"connectorId"`
	Role        string `json:"role"`
	Source      string `json:"source"` // "manual" or "oidc"
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// connectorRoleRank orders connector roles so a higher role satisfies a
// lower minimum, mirroring auth.roleSatisfies' viewer <= operator ordering.
var connectorRoleRank = map[string]int{"viewer": 1, "operator": 2}

// highestConnectorRole returns the highest-ranked role in roles, or "" if
// roles is empty or every entry is unranked.
func highestConnectorRole(roles []string) string {
	best := ""
	for _, role := range roles {
		if best == "" || connectorRoleRank[role] > connectorRoleRank[best] {
			best = role
		}
	}
	return best
}

// GetUserConnectorRole returns the user's highest role on a connector across
// every source (manual + oidc), or "" if no grant exists (default deny).
func (s *Store) GetUserConnectorRole(ctx context.Context, userID, connectorID string) (string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT role FROM user_connector_roles WHERE user_id = ? AND connector_id = ?`,
		userID, connectorID,
	)
	if err != nil {
		return "", fmt.Errorf("get user connector role: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return "", fmt.Errorf("scan user connector role: %w", err)
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate user connector role: %w", err)
	}
	role := highestConnectorRole(roles)
	if role == "" {
		return "", nil
	}
	// An API key can narrow its owner's grant (#278), never widen it.
	return auth.ClampConnectorRole(ctx, connectorID, role), nil
}

// ConnectorReaderIDs returns the IDs of the enabled users holding any grant
// (manual or oidc, viewer or operator) on connectorID. Instance admins get no
// implicit access, matching the REST list filter.
func (s *Store) ConnectorReaderIDs(ctx context.Context, connectorID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT DISTINCT r.user_id
		FROM user_connector_roles r JOIN users u ON u.id = r.user_id
		WHERE r.connector_id = ? AND u.disabled = 0`, connectorID)
	if err != nil {
		return nil, fmt.Errorf("list connector readers: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan connector reader: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connector readers: %w", err)
	}
	return ids, nil
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

// ListConnectorGrants returns every user's grant on a connector (one row per
// source), for the connector's Permissions tab.
func (s *Store) ListConnectorGrants(ctx context.Context, connectorID string) ([]ConnectorGrant, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, connector_id, role, source, created_at, updated_at
		FROM user_connector_roles WHERE connector_id = ? ORDER BY created_at ASC
	`, connectorID)
	if err != nil {
		return nil, fmt.Errorf("list connector grants: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	return scanConnectorGrants(rows)
}

// ListUserConnectorGrants returns every connector grant held by a user (one
// row per source).
func (s *Store) ListUserConnectorGrants(ctx context.Context, userID string) ([]ConnectorGrant, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, connector_id, role, source, created_at, updated_at
		FROM user_connector_roles WHERE user_id = ? ORDER BY created_at ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list user connector grants: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	return scanConnectorGrants(rows)
}

// listOIDCConnectorGrants returns userID's current 'oidc'-sourced grants,
// keyed by connector ID. Used by SyncOIDCConnectorGrants to diff against the
// desired state.
func listOIDCConnectorGrants(ctx context.Context, tx DBTX, userID string) (map[string]ConnectorGrant, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, user_id, connector_id, role, source, created_at, updated_at
		FROM user_connector_roles WHERE user_id = ? AND source = 'oidc'
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list oidc connector grants: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	grants, err := scanConnectorGrants(rows)
	if err != nil {
		return nil, err
	}
	byConnector := make(map[string]ConnectorGrant, len(grants))
	for _, g := range grants {
		byConnector[g.ConnectorID] = g
	}
	return byConnector, nil
}

func scanConnectorGrants(rows *sql.Rows) ([]ConnectorGrant, error) {
	grants := make([]ConnectorGrant, 0)
	for rows.Next() {
		var g ConnectorGrant
		if err := rows.Scan(&g.ID, &g.UserID, &g.ConnectorID, &g.Role, &g.Source, &g.CreatedAt, &g.UpdatedAt); err != nil {
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

// connectorGrantSourceManual and connectorGrantSourceOIDC are the two values
// user_connector_roles.source can hold (#279 part 3).
const (
	connectorGrantSourceManual = "manual"
	connectorGrantSourceOIDC   = "oidc"
)

// UpsertConnectorGrant creates or updates a user's manually-granted role on a
// connector — the admin permissions API. It never touches an 'oidc' row for
// the same pair; those are only ever written by SyncOIDCConnectorGrants.
func (s *Store) UpsertConnectorGrant(ctx context.Context, userID, connectorID, role string) (ConnectorGrant, error) {
	return upsertConnectorGrant(ctx, s.db, userID, connectorID, role, connectorGrantSourceManual)
}

func upsertConnectorGrant(ctx context.Context, db DBTX, userID, connectorID, role, source string) (ConnectorGrant, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	g := ConnectorGrant{
		ID:          uuid.New().String(),
		UserID:      userID,
		ConnectorID: connectorID,
		Role:        role,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := db.ExecContext(ctx, `
		INSERT INTO user_connector_roles (id, user_id, connector_id, role, source, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (user_id, connector_id, source) DO UPDATE SET role = excluded.role, updated_at = excluded.updated_at
	`, g.ID, g.UserID, g.ConnectorID, g.Role, g.Source, g.CreatedAt, g.UpdatedAt)
	if err != nil {
		return ConnectorGrant{}, fmt.Errorf("upsert connector grant: %w", err)
	}
	// A conflict path doesn't return the pre-existing row's real ID/CreatedAt;
	// read back so callers (and the audit log) see the persisted values.
	return getConnectorGrant(ctx, db, userID, connectorID, source)
}

func getConnectorGrant(ctx context.Context, db DBTX, userID, connectorID, source string) (ConnectorGrant, error) {
	var g ConnectorGrant
	err := db.QueryRowContext(ctx, `
		SELECT id, user_id, connector_id, role, source, created_at, updated_at
		FROM user_connector_roles WHERE user_id = ? AND connector_id = ? AND source = ?
	`, userID, connectorID, source).Scan(&g.ID, &g.UserID, &g.ConnectorID, &g.Role, &g.Source, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return ConnectorGrant{}, fmt.Errorf("get connector grant: %w", err)
	}
	return g, nil
}

// DeleteConnectorGrant revokes a user's manually-granted access to a
// connector — the admin permissions API. An 'oidc' grant on the same pair,
// if any, is untouched and still governs access until the next login sync.
func (s *Store) DeleteConnectorGrant(ctx context.Context, userID, connectorID string) error {
	result, err := s.db.ExecContext(ctx,
		`DELETE FROM user_connector_roles WHERE user_id = ? AND connector_id = ? AND source = ?`,
		userID, connectorID, connectorGrantSourceManual)
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

// ConnectorGrantDiff is what changed in a call to SyncOIDCConnectorGrants,
// for the audit log.
type ConnectorGrantDiff struct {
	Added   []ConnectorGrant
	Removed []ConnectorGrant
}

// Empty reports whether the sync changed nothing.
func (d ConnectorGrantDiff) Empty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0
}

// SyncOIDCConnectorGrants replaces userID's 'oidc'-sourced connector grants
// with desired (connector ID -> role), in one transaction. Rows not in
// desired are deleted, rows in desired are upserted (skipped if already
// correct), and 'manual' rows are never touched. Called on every OIDC login
// for AuthSource=="oidc" users (#279 part 3); an empty desired map revokes
// every oidc grant the user previously held.
func (s *Store) SyncOIDCConnectorGrants(ctx context.Context, userID string, desired map[string]string) (ConnectorGrantDiff, error) {
	var diff ConnectorGrantDiff
	err := s.WithinTransaction(ctx, func(tx *Store) error {
		existing, err := listOIDCConnectorGrants(ctx, tx.db, userID)
		if err != nil {
			return err
		}
		for connectorID, g := range existing {
			if _, wanted := desired[connectorID]; wanted {
				continue
			}
			if _, err := tx.db.ExecContext(ctx, `DELETE FROM user_connector_roles WHERE id = ?`, g.ID); err != nil {
				return fmt.Errorf("delete stale oidc connector grant: %w", err)
			}
			diff.Removed = append(diff.Removed, g)
		}
		for connectorID, role := range desired {
			if g, ok := existing[connectorID]; ok && g.Role == role {
				continue // already correct, no-op
			}
			g, err := upsertConnectorGrant(ctx, tx.db, userID, connectorID, role, connectorGrantSourceOIDC)
			if err != nil {
				return fmt.Errorf("upsert oidc connector grant: %w", err)
			}
			diff.Added = append(diff.Added, g)
		}
		return nil
	})
	if err != nil {
		return ConnectorGrantDiff{}, err
	}
	return diff, nil
}

// ListConnectorIDs returns every connector's ID, for expanding the "*"
// wildcard in an OIDC provider's group_connector_roles mapping.
func (s *Store) ListConnectorIDs(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id FROM connectors`)
	if err != nil {
		return nil, fmt.Errorf("list connector ids: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan connector id: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connector ids: %w", err)
	}
	return ids, nil
}
