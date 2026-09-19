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

// ChangeRecord represents a row in the changes table.
type ChangeRecord struct {
	ID             string `json:"id"`
	ServiceID      string `json:"serviceId"`
	ChangeType     string `json:"changeType"`
	Severity       string `json:"severity"`
	Summary        string `json:"summary"`
	Diff           string `json:"diff"`
	Status         string `json:"status"`
	DetectedAt     string `json:"detectedAt"`
	AffectedDocIDs string `json:"affectedDocIds"`
	// RelatedServiceIDs is a JSON array of connector IDs this change's
	// service depends on (from ServiceSnapshot.Dependencies), letting a
	// dependent-service drift (e.g. a Docker config change affecting a
	// Proxmox VM) be recognized as one logical change instead of N
	// unrelated ones. "" defaults to "[]".
	RelatedServiceIDs string `json:"relatedServiceIds"`
	// PatternID is a stable identifier derived from (service, change type,
	// summary) shared by recurrences of the same drift, so the UI can
	// eventually surface "seen this before". "" means not computed.
	PatternID string `json:"patternId"`
	// Narration is the cached plain-English "explain this change" AI
	// narration (issue #238, piece 2/3). "" means it hasn't been generated yet.
	Narration string `json:"narration"`
}

// changeColumns is the shared column list for every change SELECT.
const changeColumns = `id, service_id, change_type, severity, summary, diff, status, detected_at, affected_doc_ids, related_service_ids, pattern_id, narration`

// scanChange scans one changeColumns row. Used by ListChanges's paginatedQuery.
func scanChange(row rowScanner) (ChangeRecord, error) {
	var c ChangeRecord
	err := row.Scan(&c.ID, &c.ServiceID, &c.ChangeType, &c.Severity, &c.Summary,
		&c.Diff, &c.Status, &c.DetectedAt, &c.AffectedDocIDs, &c.RelatedServiceIDs, &c.PatternID, &c.Narration)
	return c, err
}

// alertColumns is the shared column list for every alert SELECT.
const alertColumns = `id, change_id, service_id, severity, title, description, status, snoozed_until, created_at`

// scanAlert scans one alertColumns row. Used by ListAlerts's paginatedQuery.
func scanAlert(row rowScanner) (AlertRecord, error) {
	var a AlertRecord
	var changeID, snoozedUntil sql.NullString
	err := row.Scan(&a.ID, &changeID, &a.ServiceID, &a.Severity, &a.Title,
		&a.Description, &a.Status, &snoozedUntil, &a.CreatedAt)
	if err != nil {
		return AlertRecord{}, err
	}
	a.ChangeID = changeID.String
	a.SnoozedUntil = snoozedUntil.String
	return a, nil
}

// AlertRecord represents a row in the alerts table.
type AlertRecord struct {
	ID           string `json:"id"`
	ChangeID     string `json:"changeId"`
	ServiceID    string `json:"serviceId"`
	Severity     string `json:"severity"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	SnoozedUntil string `json:"snoozedUntil"`
	CreatedAt    string `json:"createdAt"`
}

// --- Change CRUD ---

// CreateChange inserts a new infrastructure change record.
func (s *Store) CreateChange(ctx context.Context, c *ChangeRecord) error {
	return s.CreateChanges(ctx, []*ChangeRecord{c})
}

// CreateChanges inserts bounded batches; use WithinTransaction for atomic multi-batch writes.
func (s *Store) CreateChanges(ctx context.Context, records []*ChangeRecord) error {
	// Stay below SQLite's portable 999-parameter limit.
	const batchSize = 999 / 11
	for start := 0; start < len(records); start += batchSize {
		end := min(start+batchSize, len(records))
		values := make([]string, 0, end-start)
		args := make([]any, 0, (end-start)*11)
		for _, c := range records[start:end] {
			if c.ID == "" {
				c.ID = uuid.New().String()
			}
			if c.DetectedAt == "" {
				c.DetectedAt = time.Now().UTC().Format(time.RFC3339)
			}
			if c.Status == "" {
				c.Status = "new"
			}
			if c.Diff == "" {
				c.Diff = "{}"
			}
			if c.AffectedDocIDs == "" {
				c.AffectedDocIDs = "[]"
			}
			if c.RelatedServiceIDs == "" {
				c.RelatedServiceIDs = "[]"
			}

			values = append(values, "("+placeholders(11)+")")
			args = append(args, c.ID, c.ServiceID, c.ChangeType, c.Severity, c.Summary, c.Diff, c.Status, c.DetectedAt, c.AffectedDocIDs, c.RelatedServiceIDs, c.PatternID)
		}
		if _, err := s.db.ExecContext(ctx, "INSERT INTO changes (id, service_id, change_type, severity, summary, diff, status, detected_at, affected_doc_ids, related_service_ids, pattern_id) VALUES "+strings.Join(values, ","), args...); err != nil {
			return fmt.Errorf("create changes: %w", err)
		}
	}
	return nil
}

// GetChange retrieves a single change record by ID.
func (s *Store) GetChange(ctx context.Context, id string) (*ChangeRecord, error) {
	c := &ChangeRecord{}
	err := s.db.QueryRowContext(ctx, `
		SELECT `+changeColumns+`
		FROM changes WHERE id = ?
	`, id).Scan(&c.ID, &c.ServiceID, &c.ChangeType, &c.Severity, &c.Summary,
		&c.Diff, &c.Status, &c.DetectedAt, &c.AffectedDocIDs, &c.RelatedServiceIDs, &c.PatternID, &c.Narration)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get change: %w", err)
	}
	return c, nil
}

// UpdateChangeNarration persists the generated "explain this change"
// narration so later views reuse it instead of re-calling the AI provider.
func (s *Store) UpdateChangeNarration(ctx context.Context, id, narration string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE changes SET narration = ? WHERE id = ?`, narration, id)
	if err != nil {
		return fmt.Errorf("update change narration: %w", err)
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

// CountRecentChangesByPattern returns how many changes with the given
// patternID (for a service) were detected at or after `since`, excluding
// excludeID (the change just created for this drift). Used to flag a repeat
// drift — the same pattern recurring within a short window, e.g. a
// misconfiguration loop.
func (s *Store) CountRecentChangesByPattern(ctx context.Context, serviceID, patternID, since, excludeID string) (int, error) {
	if patternID == "" {
		return 0, nil
	}
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM changes
		WHERE service_id = ? AND pattern_id = ? AND detected_at >= ? AND id != ?
	`, serviceID, patternID, since, excludeID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count recent changes by pattern: %w", err)
	}
	return count, nil
}

// UpdateChangeStatus updates the status of a change record.
func (s *Store) UpdateChangeStatus(ctx context.Context, id, status string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE changes SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update change status: %w", err)
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

// ListChangesByID returns the requested changes keyed by ID.
func (s *Store) ListChangesByID(ctx context.Context, ids []string) (map[string]ChangeRecord, error) {
	changes := make(map[string]ChangeRecord, len(ids))
	if len(ids) == 0 {
		return changes, nil
	}

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := s.db.QueryContext(ctx, `SELECT `+changeColumns+` FROM changes WHERE id IN (`+placeholders(len(ids))+`)`, args...)
	if err != nil {
		return nil, fmt.Errorf("list changes by id: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var c ChangeRecord
		if err := rows.Scan(&c.ID, &c.ServiceID, &c.ChangeType, &c.Severity, &c.Summary, &c.Diff, &c.Status, &c.DetectedAt, &c.AffectedDocIDs, &c.RelatedServiceIDs, &c.PatternID, &c.Narration); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		changes[c.ID] = c
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate changes: %w", err)
	}
	return changes, nil
}

// UpdateChangeStatuses updates all requested change IDs in one statement.
func (s *Store) UpdateChangeStatuses(ctx context.Context, ids []string, status string) error {
	if len(ids) == 0 {
		return nil
	}
	args := make([]any, 1, len(ids)+1)
	args[0] = status
	for _, id := range ids {
		args = append(args, id)
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE changes SET status = ? WHERE id IN (`+placeholders(len(ids))+`)`, args...); err != nil {
		return fmt.Errorf("update change statuses: %w", err)
	}
	return nil
}

// ListChanges returns a paginated list of change records, optionally filtered by service and severity.
func (s *Store) ListChanges(ctx context.Context, serviceID, severity string, offset, limit int) ([]ChangeRecord, int, error) {
	where := "WHERE 1=1"
	var args []any
	if serviceID != "" {
		where += " AND service_id = ?"
		args = append(args, serviceID)
	}
	if severity != "" {
		where += " AND severity = ?"
		args = append(args, severity)
	}

	return paginatedQuery(ctx, s.db, "changes", changeColumns, where, args, "detected_at DESC", limit, offset, scanChange)
}

// CountChanges returns the total number of change records.
func (s *Store) CountChanges(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM changes`).Scan(&count)
	return count, err
}

// CountChangesNew returns count of unacknowledged changes.
func (s *Store) CountChangesNew(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM changes WHERE status = 'new'`).Scan(&count)
	return count, err
}

// --- Alert CRUD ---

// CreateAlert inserts a new alert record.
func (s *Store) CreateAlert(ctx context.Context, a *AlertRecord) error {
	return s.CreateAlerts(ctx, []*AlertRecord{a})
}

// CreateAlerts inserts bounded batches; use WithinTransaction for atomic multi-batch writes.
func (s *Store) CreateAlerts(ctx context.Context, records []*AlertRecord) error {
	// Stay below SQLite's portable 999-parameter limit.
	const batchSize = 999 / 9
	for start := 0; start < len(records); start += batchSize {
		end := min(start+batchSize, len(records))
		values := make([]string, 0, end-start)
		args := make([]any, 0, (end-start)*9)
		for _, a := range records[start:end] {
			if a.ID == "" {
				a.ID = uuid.New().String()
			}
			if a.CreatedAt == "" {
				a.CreatedAt = time.Now().UTC().Format(time.RFC3339)
			}
			if a.Status == "" {
				a.Status = "pending"
			}

			values = append(values, "("+placeholders(9)+")")
			args = append(args, a.ID, nilToStr(a.ChangeID), a.ServiceID, a.Severity, a.Title, a.Description, a.Status, nilToStr(a.SnoozedUntil), a.CreatedAt)
		}
		if _, err := s.db.ExecContext(ctx, "INSERT INTO alerts (id, change_id, service_id, severity, title, description, status, snoozed_until, created_at) VALUES "+strings.Join(values, ","), args...); err != nil {
			return fmt.Errorf("create alerts: %w", err)
		}
	}
	return nil
}

// GetAlert retrieves a single alert record by ID.
func (s *Store) GetAlert(ctx context.Context, id string) (*AlertRecord, error) {
	a := &AlertRecord{}
	var changeID, snoozedUntil sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, change_id, service_id, severity, title, description, status, snoozed_until, created_at
		FROM alerts WHERE id = ?
	`, id).Scan(&a.ID, &changeID, &a.ServiceID, &a.Severity, &a.Title,
		&a.Description, &a.Status, &snoozedUntil, &a.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get alert: %w", err)
	}
	a.ChangeID = changeID.String
	a.SnoozedUntil = snoozedUntil.String
	return a, nil
}

// UpdateAlertStatus updates the status and snooze deadline of an alert.
func (s *Store) UpdateAlertStatus(ctx context.Context, id, status, snoozedUntil string) error {
	args := []any{status}
	query := "UPDATE alerts SET status = ?"
	if snoozedUntil != "" {
		query += ", snoozed_until = ?"
		args = append(args, snoozedUntil)
	}
	query += " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update alert: %w", err)
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

// UpdateAlertStatuses updates all requested alert IDs in one statement and
// returns the IDs that existed before the update.
func (s *Store) UpdateAlertStatuses(ctx context.Context, ids []string, status, snoozedUntil string) (map[string]bool, error) {
	found := make(map[string]bool, len(ids))
	if len(ids) == 0 {
		return found, nil
	}
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id FROM alerts WHERE id IN (`+placeholders(len(ids))+`)`, args...)
	if err != nil {
		return nil, fmt.Errorf("list alerts by id: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		found[id] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate alerts: %w", err)
	}
	if len(found) == 0 {
		return found, nil
	}

	args = []any{status, snoozedUntil}
	for id := range found {
		args = append(args, id)
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE alerts SET status = ?, snoozed_until = ? WHERE id IN (`+placeholders(len(found))+`)`, args...); err != nil {
		return nil, fmt.Errorf("update alert statuses: %w", err)
	}
	return found, nil
}

// ListAlerts returns a paginated list of alerts, optionally filtered by service, severity, status, and since.
func (s *Store) ListAlerts(ctx context.Context, serviceID, severity, status, since string, offset, limit int) ([]AlertRecord, int, error) {
	where := "WHERE 1=1"
	var args []any
	if serviceID != "" {
		where += " AND service_id = ?"
		args = append(args, serviceID)
	}
	if severity != "" {
		where += " AND severity = ?"
		args = append(args, severity)
	}
	if status != "" {
		where += " AND status = ?"
		args = append(args, status)
	}
	if since != "" {
		where += " AND created_at >= ?"
		args = append(args, since)
	}

	return paginatedQuery(ctx, s.db, "alerts", alertColumns, where, args, "created_at DESC", limit, offset, scanAlert)
}

// UnsnoozeExpiredAlerts moves every snoozed alert whose snoozed_until has
// passed back to pending in one statement and returns how many changed.
func (s *Store) UnsnoozeExpiredAlerts(ctx context.Context) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := s.db.ExecContext(ctx, `
		UPDATE alerts SET status = 'pending'
		WHERE status = 'snoozed' AND snoozed_until IS NOT NULL AND snoozed_until <= ?
	`, now)
	if err != nil {
		return 0, fmt.Errorf("unsnooze expired alerts: %w", err)
	}
	return rowsAffected(res), nil
}

// GetExpiredSnoozedAlerts returns alerts where snoozed_until has passed.
func (s *Store) GetExpiredSnoozedAlerts(ctx context.Context) ([]AlertRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, change_id, service_id, severity, title, description, status, snoozed_until, created_at
		FROM alerts WHERE status = 'snoozed' AND snoozed_until IS NOT NULL AND snoozed_until <= ?
	`, now)
	if err != nil {
		return nil, fmt.Errorf("get expired snoozed: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var alerts []AlertRecord
	for rows.Next() {
		var a AlertRecord
		var changeID, snoozedUntil sql.NullString
		if err := rows.Scan(&a.ID, &changeID, &a.ServiceID, &a.Severity, &a.Title, &a.Description,
			&a.Status, &snoozedUntil, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		a.ChangeID = changeID.String
		a.SnoozedUntil = snoozedUntil.String
		alerts = append(alerts, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate expired snoozed alerts: %w", err)
	}
	if alerts == nil {
		alerts = []AlertRecord{}
	}
	return alerts, nil
}

// CountAlerts returns the total number of alert records.
func (s *Store) CountAlerts(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts`).Scan(&count)
	return count, err
}

// CountAlertsPending returns count of pending alerts.
func (s *Store) CountAlertsPending(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE status = 'pending'`).Scan(&count)
	return count, err
}

// GetLatestChanges returns the most recent N changes detected at or after
// since (RFC3339, matching the stored detected_at format). An empty since
// applies no lower bound.
func (s *Store) GetLatestChanges(ctx context.Context, n int, since string) ([]ChangeRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+changeColumns+`
		FROM changes WHERE detected_at >= ? ORDER BY detected_at DESC LIMIT ?
	`, since, n)
	if err != nil {
		return nil, fmt.Errorf("get latest changes: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var changes []ChangeRecord
	for rows.Next() {
		var c ChangeRecord
		if err := rows.Scan(&c.ID, &c.ServiceID, &c.ChangeType, &c.Severity, &c.Summary, &c.Diff, &c.Status, &c.DetectedAt, &c.AffectedDocIDs, &c.RelatedServiceIDs, &c.PatternID, &c.Narration); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		changes = append(changes, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate latest changes: %w", err)
	}
	if changes == nil {
		changes = []ChangeRecord{}
	}
	return changes, nil
}

// GetLastSyncTimestamp returns the most recent sync timestamp across all connectors.
func (s *Store) GetLastSyncTimestamp(ctx context.Context) (string, error) {
	var ts sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT MAX(last_sync_at) FROM connectors`).Scan(&ts)
	if err != nil {
		return "", err
	}
	return ts.String, nil
}

// CountConnectorsByStatus counts connectors grouped by status.
func (s *Store) CountConnectorsByStatus(ctx context.Context) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT status, COUNT(*) FROM connectors GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	counts := map[string]int{"online": 0, "degraded": 0, "offline": 0, "unknown": 0}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("scan status count: %w", err)
		}
		counts[status] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connector status counts: %w", err)
	}
	return counts, nil
}

// CountRecentChangePatterns counts every recent pattern for one service in one query.
func (s *Store) CountRecentChangePatterns(ctx context.Context, serviceID, since string) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT pattern_id, COUNT(*) FROM changes WHERE service_id = ? AND detected_at >= ? AND pattern_id != '' GROUP BY pattern_id`, serviceID, since)
	if err != nil {
		return nil, fmt.Errorf("count recent change patterns: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	counts := make(map[string]int)
	for rows.Next() {
		var pattern string
		var count int
		if err := rows.Scan(&pattern, &count); err != nil {
			return nil, err
		}
		counts[pattern] = count
	}
	return counts, rows.Err()
}
