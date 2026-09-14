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
}

const qualityFindingColumns = `id, connector_id, doc_id, check_type, severity, title, description,
	remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at`

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

	err := s.db.QueryRowContext(ctx, `
		INSERT INTO quality_findings (id, connector_id, doc_id, check_type, severity, title, description,
			remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'open', 1, ?, ?, NULL)
		ON CONFLICT(connector_id, check_type) WHERE status = 'open'
		DO UPDATE SET last_seen_at = excluded.last_seen_at,
			detected_count = quality_findings.detected_count + 1,
			doc_id = excluded.doc_id,
			severity = excluded.severity, title = excluded.title,
			description = excluded.description, remediation_link = excluded.remediation_link
		RETURNING id
	`, f.ID, f.ConnectorID, nilToStr(f.DocID), f.CheckType, f.Severity, f.Title,
		f.Description, f.RemediationLink, f.FirstDetectedAt, f.LastSeenAt).Scan(&f.ID)
	if err != nil {
		return fmt.Errorf("upsert quality finding: %w", err)
	}
	return nil
}

// ResolveQualityFinding resolves the currently open finding, if any.
func (s *Store) ResolveQualityFinding(ctx context.Context, connectorID, checkType string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := s.db.ExecContext(ctx, `
		UPDATE quality_findings SET status = 'resolved', resolved_at = ?
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
	var docID, resolvedAt sql.NullString
	err := row.Scan(&f.ID, &f.ConnectorID, &docID, &f.CheckType, &f.Severity, &f.Title,
		&f.Description, &f.RemediationLink, &f.Status, &f.DetectedCount,
		&f.FirstDetectedAt, &f.LastSeenAt, &resolvedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return QualityFindingRecord{}, ErrNotFound
	}
	if err != nil {
		return QualityFindingRecord{}, err
	}
	f.DocID = docID.String
	f.ResolvedAt = resolvedAt.String
	return f, nil
}
