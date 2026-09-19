package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// User represents a row in the users table.
type User struct {
	ID                         string `json:"id"`
	Username                   string `json:"username"`
	DisplayName                string `json:"displayName"`
	Email                      string `json:"email"`
	InstanceAdminRole          string `json:"instanceAdminRole"`
	AuthSource                 string `json:"authSource"`
	PasswordHash               string `json:"-"`
	Disabled                   bool   `json:"disabled"`
	CanManageDashboardDefaults bool   `json:"canManageDashboardDefaults"`
	CreatedAt                  string `json:"createdAt"`
	FailedLoginAttempts        int    `json:"-"`
	LockedUntil                string `json:"-"`
	DigestCadence              string `json:"digestCadence"`
	DigestLastSentAt           string `json:"-"`
	DigestTimezone             string `json:"digestTimezone"`
}

// Session represents a row in the sessions table.
type Session struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	TokenHash      string `json:"-"`
	AuthProviderID string `json:"authProviderId"`
	UserAgent      string `json:"userAgent"`
	IP             string `json:"ip"`
	CreatedAt      string `json:"createdAt"`
	LastSeenAt     string `json:"lastSeenAt"`
}

// --- User operations ---

// CreateUser inserts a new user and returns it.
func (s *Store) CreateUser(ctx context.Context, user *User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	if user.CreatedAt == "" {
		user.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if user.InstanceAdminRole == "" {
		user.InstanceAdminRole = "user"
	}
	if user.AuthSource == "" {
		user.AuthSource = "local"
	}
	if user.DigestCadence == "" {
		user.DigestCadence = "off"
	}
	if user.DigestTimezone == "" {
		user.DigestTimezone = "UTC"
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (id, username, display_name, email, instance_admin_role, auth_source, password_hash, disabled, can_manage_dashboard_defaults, created_at, digest_cadence, digest_last_sent_at, digest_timezone)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, user.ID, user.Username, user.DisplayName, user.Email, user.InstanceAdminRole, user.AuthSource,
		user.PasswordHash, boolToInt(user.Disabled), boolToInt(user.CanManageDashboardDefaults), user.CreatedAt, user.DigestCadence, user.DigestLastSentAt, user.DigestTimezone)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrConflict
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID.
func (s *Store) GetUserByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	var disabled, canManageDashboardDefaults int
	err := s.db.QueryRowContext(ctx, `
		SELECT id, username, display_name, email, instance_admin_role, auth_source, password_hash, disabled, can_manage_dashboard_defaults, created_at, digest_cadence, digest_last_sent_at, digest_timezone
		FROM users WHERE id = ?
	`, id).Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.InstanceAdminRole,
		&u.AuthSource, &u.PasswordHash, &disabled, &canManageDashboardDefaults, &u.CreatedAt, &u.DigestCadence, &u.DigestLastSentAt, &u.DigestTimezone)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	u.Disabled = disabled != 0
	u.CanManageDashboardDefaults = canManageDashboardDefaults != 0
	return u, nil
}

// GetUserByUsername retrieves a user by username.
func (s *Store) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	u := &User{}
	var disabled, canManageDashboardDefaults int
	err := s.db.QueryRowContext(ctx, `
		SELECT id, username, display_name, email, instance_admin_role, auth_source, password_hash, disabled, can_manage_dashboard_defaults, created_at, failed_login_attempts, locked_until
		FROM users WHERE username = ?
	`, username).Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.InstanceAdminRole,
		&u.AuthSource, &u.PasswordHash, &disabled, &canManageDashboardDefaults, &u.CreatedAt, &u.FailedLoginAttempts, &u.LockedUntil)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	u.Disabled = disabled != 0
	u.CanManageDashboardDefaults = canManageDashboardDefaults != 0
	return u, nil
}

// RegisterFailedLogin increments the account's failed-login counter and, once
// it reaches maxAttempts, locks the account until now+lockDuration and resets
// the counter. Returns true when this call caused the account to lock.
func (s *Store) RegisterFailedLogin(ctx context.Context, userID string, maxAttempts int, lockDuration time.Duration) (bool, error) {
	var attempts int
	err := s.db.QueryRowContext(ctx,
		`UPDATE users SET failed_login_attempts = failed_login_attempts + 1 WHERE id = ? RETURNING failed_login_attempts`,
		userID,
	).Scan(&attempts)
	if err != nil {
		return false, fmt.Errorf("register failed login: %w", err)
	}
	if attempts < maxAttempts {
		return false, nil
	}
	lockedUntil := time.Now().UTC().Add(lockDuration).Format(time.RFC3339)
	if _, err := s.db.ExecContext(ctx,
		`UPDATE users SET failed_login_attempts = 0, locked_until = ? WHERE id = ?`,
		lockedUntil, userID,
	); err != nil {
		return false, fmt.Errorf("lock user: %w", err)
	}
	return true, nil
}

// ClearFailedLogins resets the failed-login counter and any lockout after a
// successful authentication.
func (s *Store) ClearFailedLogins(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET failed_login_attempts = 0, locked_until = '' WHERE id = ?`, userID)
	if err != nil {
		return fmt.Errorf("clear failed logins: %w", err)
	}
	return nil
}

// GetUserRoleStatus returns the current instance-admin role and disabled flag
// for a user. Implements auth.UserStatusChecker so AuthMiddleware can detect
// a role change or account disable that happened after an access token was
// issued. Successful reads are cached for 30 seconds; user mutations invalidate
// the cache immediately. Errors are never cached.
func (s *Store) GetUserRoleStatus(ctx context.Context, userID string) (string, bool, error) {
	now := time.Now()
	s.userStatusMu.Lock()
	cached, ok := s.userStatuses[userID]
	generation := s.userStatusGeneration
	s.userStatusMu.Unlock()
	if ok && now.Before(cached.expires) {
		return cached.role, cached.disabled, nil
	}
	var role string
	var disabled int
	err := s.db.QueryRowContext(ctx, `SELECT instance_admin_role, disabled FROM users WHERE id = ?`, userID).Scan(&role, &disabled)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, ErrNotFound
	}
	if err != nil {
		return "", false, fmt.Errorf("get user role status: %w", err)
	}
	s.userStatusMu.Lock()
	if generation == s.userStatusGeneration {
		// Bound memory even if users stop making requests before their entries expire.
		if len(s.userStatuses) >= 4096 {
			s.userStatuses = nil
		}
		if s.userStatuses == nil {
			s.userStatuses = make(map[string]userStatus)
		}
		s.userStatuses[userID] = userStatus{role: role, disabled: disabled != 0, expires: now.Add(30 * time.Second)}
	}
	s.userStatusMu.Unlock()
	return role, disabled != 0, nil
}

// UpdateUser updates fields on an existing user.
func (s *Store) UpdateUser(ctx context.Context, id string, updates map[string]any) error {
	// Build dynamic UPDATE for simplicity
	query := `UPDATE users SET `
	args := make([]any, 0)
	setClauses := make([]string, 0)

	for k, v := range updates {
		switch k {
		case "display_name":
			setClauses = append(setClauses, "display_name = ?")
			args = append(args, v)
		case "email":
			setClauses = append(setClauses, "email = ?")
			args = append(args, v)
		case "instance_admin_role":
			setClauses = append(setClauses, "instance_admin_role = ?")
			args = append(args, v)
		case "password_hash":
			setClauses = append(setClauses, "password_hash = ?")
			args = append(args, v)
		case "disabled":
			setClauses = append(setClauses, "disabled = ?")
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf("update user: field %q must be a bool", k)
			}
			args = append(args, boolToInt(b))
		case "username":
			setClauses = append(setClauses, "username = ?")
			args = append(args, v)
		case "can_manage_dashboard_defaults":
			setClauses = append(setClauses, "can_manage_dashboard_defaults = ?")
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf("update user: field %q must be a bool", k)
			}
			args = append(args, boolToInt(b))
		case "digest_cadence":
			setClauses = append(setClauses, "digest_cadence = ?")
			args = append(args, v)
		case "digest_last_sent_at":
			setClauses = append(setClauses, "digest_last_sent_at = ?")
			args = append(args, v)
		case "digest_timezone":
			setClauses = append(setClauses, "digest_timezone = ?")
			args = append(args, v)
		}
	}

	if len(setClauses) == 0 {
		return nil
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrConflict
		}
		return fmt.Errorf("update user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	s.invalidateUserStatuses()
	return nil
}

// DeleteUser deletes a user by ID.
func (s *Store) DeleteUser(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	s.invalidateUserStatuses()
	return nil
}

// userColumns is the shared column list for the ListUsers SELECT.
const userColumns = `id, username, display_name, email, instance_admin_role, auth_source, password_hash, disabled, can_manage_dashboard_defaults, created_at, digest_cadence, digest_last_sent_at, digest_timezone`

func scanUser(row rowScanner) (User, error) {
	var u User
	var disabled, canManageDashboardDefaults int
	err := row.Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.InstanceAdminRole,
		&u.AuthSource, &u.PasswordHash, &disabled, &canManageDashboardDefaults, &u.CreatedAt, &u.DigestCadence, &u.DigestLastSentAt, &u.DigestTimezone)
	if err != nil {
		return User{}, err
	}
	u.Disabled = disabled != 0
	u.CanManageDashboardDefaults = canManageDashboardDefaults != 0
	return u, nil
}

// ListUsers returns a paginated list of users.
func (s *Store) ListUsers(ctx context.Context, offset, limit int) ([]User, int, error) {
	return paginatedQuery(ctx, s.db, "users", userColumns, "", nil, "created_at DESC", limit, offset, scanUser)
}

// CountUsers returns the total number of users.
func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

// UserHasPermission checks a named boolean permission column on the users table.
// Implements auth.PermissionChecker. Only "can_manage_dashboard_defaults" is
// supported today; unknown permission names return false.
func (s *Store) UserHasPermission(ctx context.Context, userID, permission string) (bool, error) {
	if permission != "can_manage_dashboard_defaults" {
		return false, nil
	}
	var flag int
	err := s.db.QueryRowContext(ctx,
		`SELECT can_manage_dashboard_defaults FROM users WHERE id = ?`, userID,
	).Scan(&flag)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check user permission: %w", err)
	}
	return flag != 0, nil
}

// --- Session operations ---

// CreateSession inserts a new session.
func (s *Store) CreateSession(ctx context.Context, session *Session) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if session.CreatedAt == "" {
		session.CreatedAt = now
	}
	if session.LastSeenAt == "" {
		session.LastSeenAt = now
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sessions (id, user_id, token_hash, auth_provider_id, user_agent, ip, created_at, last_seen_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, session.ID, session.UserID, session.TokenHash, session.AuthProviderID, session.UserAgent, session.IP,
		session.CreatedAt, session.LastSeenAt)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

// GetSession retrieves a session by ID.
func (s *Store) GetSession(ctx context.Context, id string) (*Session, error) {
	sess := &Session{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, auth_provider_id, user_agent, ip, created_at, last_seen_at
		FROM sessions WHERE id = ?
	`, id).Scan(&sess.ID, &sess.UserID, &sess.TokenHash, &sess.AuthProviderID, &sess.UserAgent, &sess.IP,
		&sess.CreatedAt, &sess.LastSeenAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	return sess, nil
}

// DeleteSession deletes a session by ID.
func (s *Store) DeleteSession(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
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

// DeleteUserSessions deletes all sessions for a user (logout all).
func (s *Store) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE user_id = ?`, userID)
	return err
}

// RevokeSessionsByAuthSource deletes all sessions created by an OIDC provider.
func (s *Store) RevokeSessionsByAuthSource(ctx context.Context, providerID string) (int64, error) {
	if providerID == "" {
		return 0, nil
	}
	result, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE auth_provider_id = ?`, providerID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// RotateSessionToken atomically replaces the current refresh-token hash.
func (s *Store) RotateSessionToken(ctx context.Context, userID, oldHash, newHash string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE sessions SET token_hash = ?, last_seen_at = ? WHERE user_id = ? AND token_hash = ?`, newHash, time.Now().UTC().Format(time.RFC3339), userID, oldHash)
	if err != nil {
		return fmt.Errorf("rotate session token: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rotate session token rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// HasSessionTokenHash reports whether an active session owns a refresh token.
func (s *Store) HasSessionTokenHash(ctx context.Context, userID, tokenHash string) (bool, error) {
	var found int
	err := s.db.QueryRowContext(ctx, `SELECT 1 FROM sessions WHERE user_id = ? AND token_hash = ?`, userID, tokenHash).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("find session token: %w", err)
	}
	return true, nil
}

// ListUserSessions returns all active sessions for a user.
func (s *Store) ListUserSessions(ctx context.Context, userID string) ([]Session, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, token_hash, auth_provider_id, user_agent, ip, created_at, last_seen_at
		FROM sessions WHERE user_id = ? ORDER BY last_seen_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.TokenHash, &s.AuthProviderID, &s.UserAgent, &s.IP,
			&s.CreatedAt, &s.LastSeenAt); err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sessions: %w", err)
	}
	if sessions == nil {
		sessions = []Session{}
	}
	return sessions, nil
}

// UpdateSessionLastSeen bumps the last_seen_at timestamp.
func (s *Store) UpdateSessionLastSeen(ctx context.Context, id string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.ExecContext(ctx, `UPDATE sessions SET last_seen_at = ? WHERE id = ?`, now, id)
	return err
}

// --- OIDC operations ---

// GetUserByOIDCIdentity finds a user by the issuer and subject verified by OIDC.
func (s *Store) GetUserByOIDCIdentity(ctx context.Context, issuer, subject string) (*User, error) {
	u := &User{}
	var disabled, canManageDashboardDefaults int
	err := s.db.QueryRowContext(ctx, `
		SELECT id, username, display_name, email, instance_admin_role, auth_source, password_hash, disabled, can_manage_dashboard_defaults, created_at
		FROM users JOIN oidc_identities ON oidc_identities.user_id = users.id
		WHERE oidc_identities.issuer = ? AND oidc_identities.subject = ?
	`, issuer, subject).Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.InstanceAdminRole,
		&u.AuthSource, &u.PasswordHash, &disabled, &canManageDashboardDefaults, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by oidc identity: %w", err)
	}
	u.Disabled = disabled != 0
	u.CanManageDashboardDefaults = canManageDashboardDefaults != 0
	return u, nil
}

// CreateOIDCUser creates an OIDC user and its verified identity atomically.
// Returns true if a new user was created, false if it already existed (conflict).
func (s *Store) CreateOIDCUser(ctx context.Context, user *User, issuer, subject string) (bool, error) {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	if user.CreatedAt == "" {
		user.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if user.InstanceAdminRole == "" {
		user.InstanceAdminRole = "user"
	}
	user.AuthSource = "oidc"
	err := s.WithinTransaction(ctx, func(tx *Store) error {
		_, err := tx.db.ExecContext(ctx, `INSERT INTO users (id, username, display_name, email, instance_admin_role, auth_source, password_hash, disabled, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, user.ID, user.Username, user.DisplayName, user.Email, user.InstanceAdminRole, user.AuthSource, user.PasswordHash, boolToInt(user.Disabled), user.CreatedAt)
		if err == nil {
			_, err = tx.db.ExecContext(ctx, `INSERT INTO oidc_identities (issuer, subject, user_id) VALUES (?, ?, ?)`, issuer, subject, user.ID)
		}
		if isUniqueViolation(err) {
			return ErrConflict
		}
		if err != nil {
			return fmt.Errorf("create oidc user: %w", err)
		}
		return nil
	})
	return err == nil, err
}

// --- Helpers ---

// HashToken hashes a token string for storage.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		// PostgreSQL: 23505 = unique_violation
		return pgErr.Code == "23505"
	}

	// SQLite reports UNIQUE constraint violations this way
	return contains(err.Error(), "UNIQUE constraint failed")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
