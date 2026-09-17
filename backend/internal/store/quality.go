package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// QualityFindingRecord represents an actionable documentation quality gap.
type QualityFindingRecord struct {
	ID              string `json:"id"`
	ConnectorID     string `json:"connectorId"`
	DocID           string `json:"docId,omitempty"`
	CheckType       string `json:"checkType"`
	Severity        string `json:"severity"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	RemediationLink string `json:"remediationLink"`
	Status          string `json:"status"`
	DetectedCount   int    `json:"detectedCount"`
	FirstDetectedAt string `json:"firstDetectedAt"`
	LastSeenAt      string `json:"lastSeenAt"`
	ResolvedAt      string `json:"resolvedAt,omitempty"`
	// NotifiedSeverity is the severity ("info"/"warning"/"critical") this
	// finding was last notified at, or "" if it has never been notified (or
	// was reset by a resolve). The notification hook compares this against
	// the current Severity to decide whether an escalation deserves a new
	// notification, then calls SetQualityFindingNotifiedSeverity.
	NotifiedSeverity string `json:"-"`
}

const qualityFindingColumns = `id, connector_id, doc_id, check_type, severity, title, description,
	remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at, notified_severity`

// UpsertQualityFinding atomically inserts a new open finding or records another
// detection of the existing open finding for the same connector and check.
func (s *Store) UpsertQualityFinding(ctx context.Context, f *QualityFindingRecord) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if f.FirstDetectedAt == "" {
		f.FirstDetectedAt = now
	}
	if f.LastSeenAt == "" {
		f.LastSeenAt = now
	}

	var notifiedSeverity sql.NullString
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO quality_findings (id, connector_id, doc_id, check_type, severity, title, description,
			remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at, notified_severity)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'open', 1, ?, ?, NULL, NULL)
		ON CONFLICT(connector_id, check_type) WHERE status = 'open'
		DO UPDATE SET last_seen_at = excluded.last_seen_at,
			detected_count = quality_findings.detected_count + 1,
			doc_id = excluded.doc_id,
			severity = excluded.severity, title = excluded.title,
			description = excluded.description, remediation_link = excluded.remediation_link
		RETURNING id, notified_severity
	`, f.ID, f.ConnectorID, nilToStr(f.DocID), f.CheckType, f.Severity, f.Title,
		f.Description, f.RemediationLink, f.FirstDetectedAt, f.LastSeenAt).Scan(&f.ID, &notifiedSeverity)
	if err != nil {
		return fmt.Errorf("upsert quality finding: %w", err)
	}
	f.NotifiedSeverity = notifiedSeverity.String
	return nil
}

// SetQualityFindingNotifiedSeverity records the severity a finding was just
// notified at, so a later re-detection at the same severity is deduplicated
// while an escalation (or a resolve-then-reopen, which resets this to "")
// still notifies.
func (s *Store) SetQualityFindingNotifiedSeverity(ctx context.Context, id, severity string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE quality_findings SET notified_severity = ? WHERE id = ?`, severity, id)
	if err != nil {
		return fmt.Errorf("set quality finding notified severity: %w", err)
	}
	return nil
}

// ResolveQualityFinding resolves the currently open finding, if any, and
// clears notified_severity so a later re-open is treated as a new finding
// for notification purposes.
func (s *Store) ResolveQualityFinding(ctx context.Context, connectorID, checkType string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := s.db.ExecContext(ctx, `
		UPDATE quality_findings SET status = 'resolved', resolved_at = ?, notified_severity = NULL
		WHERE connector_id = ? AND check_type = ? AND status = 'open'
	`, now, connectorID, checkType)
	if err != nil {
		return fmt.Errorf("resolve quality finding: %w", err)
	}
	return nil
}

// GetQualityFinding retrieves a quality finding by ID.
func (s *Store) GetQualityFinding(ctx context.Context, id string) (*QualityFindingRecord, error) {
	f, err := scanQualityFinding(s.db.QueryRowContext(ctx,
		`SELECT `+qualityFindingColumns+` FROM quality_findings WHERE id = ?`, id))
	if err != nil {
		return nil, fmt.Errorf("get quality finding: %w", err)
	}
	return &f, nil
}

// ListQualityFindings returns findings newest-seen first with optional filters.
// since, when non-empty, is an RFC3339 cutoff applied to last_seen_at.
func (s *Store) ListQualityFindings(ctx context.Context, connectorID, checkType, status, since string, offset, limit int) ([]QualityFindingRecord, int, error) {
	where := "WHERE 1=1"
	var args []any
	for _, filter := range []struct {
		column string
		value  string
	}{
		{"connector_id", connectorID},
		{"check_type", checkType},
		{"status", status},
	} {
		if filter.value != "" {
			where += " AND " + filter.column + " = ?"
			args = append(args, filter.value)
		}
	}
	if since != "" {
		where += " AND last_seen_at >= ?"
		args = append(args, since)
	}

	return paginatedQuery(ctx, s.db, "quality_findings", qualityFindingColumns, where, args, "last_seen_at DESC", limit, offset, scanQualityFinding)
}

// UpdateQualityFindingStatus changes a finding's status and resolution time.
func (s *Store) UpdateQualityFindingStatus(ctx context.Context, id, status string) error {
	resolvedAt := any(nil)
	if status == "resolved" {
		resolvedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	result, err := s.db.ExecContext(ctx,
		`UPDATE quality_findings SET status = ?, resolved_at = ? WHERE id = ?`, status, resolvedAt, id)
	if err != nil {
		return fmt.Errorf("update quality finding status: %w", err)
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

// CountQualityFindingsOpen returns the number of unresolved findings.
func (s *Store) CountQualityFindingsOpen(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quality_findings WHERE status = 'open'`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count open quality findings: %w", err)
	}
	return count, nil
}

func scanQualityFinding(row rowScanner) (QualityFindingRecord, error) {
	var f QualityFindingRecord
	var docID, resolvedAt, notifiedSeverity sql.NullString
	err := row.Scan(&f.ID, &f.ConnectorID, &docID, &f.CheckType, &f.Severity, &f.Title,
		&f.Description, &f.RemediationLink, &f.Status, &f.DetectedCount,
		&f.FirstDetectedAt, &f.LastSeenAt, &resolvedAt, &notifiedSeverity)
	if errors.Is(err, sql.ErrNoRows) {
		return QualityFindingRecord{}, ErrNotFound
	}
	if err != nil {
		return QualityFindingRecord{}, err
	}
	f.DocID = docID.String
	f.ResolvedAt = resolvedAt.String
	f.NotifiedSeverity = notifiedSeverity.String
	return f, nil
}
