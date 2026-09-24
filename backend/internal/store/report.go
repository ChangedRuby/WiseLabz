package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReportDefinitionRecord is the persisted schedule and delivery configuration.
type ReportDefinitionRecord struct {
	ID, Slug, Name, CronExpr, Timezone, Sections, ConnectorIDs, Channels, CreatedBy, CreatedAt, UpdatedAt string
	Enabled                                                                                               bool
}

// ReportRecord is one generated, immutable report snapshot.
type ReportRecord struct {
	ID, DefinitionID, DefinitionName, Trigger, PeriodStart, PeriodEnd, Data, Markdown, Status, CreatedAt string
	Truncated                                                                                            bool
}

// ListReportDefinitions returns configured schedules ordered by display name.
func (s *Store) ListReportDefinitions(ctx context.Context) ([]ReportDefinitionRecord, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, slug, name, enabled, cron_expr, timezone, sections, connector_ids, channels, created_by, created_at, updated_at FROM report_definitions ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list report definitions: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	var out []ReportDefinitionRecord
	for rows.Next() {
		var r ReportDefinitionRecord
		var enabled int
		if err := rows.Scan(&r.ID, &r.Slug, &r.Name, &enabled, &r.CronExpr, &r.Timezone, &r.Sections, &r.ConnectorIDs, &r.Channels, &r.CreatedBy, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan report definition: %w", err)
		}
		r.Enabled = enabled != 0
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate report definitions: %w", err)
	}
	return out, nil
}

// GetReportDefinition retrieves one schedule by ID.
func (s *Store) GetReportDefinition(ctx context.Context, id string) (ReportDefinitionRecord, error) {
	var r ReportDefinitionRecord
	var enabled int
	err := s.db.QueryRowContext(ctx, `SELECT id, slug, name, enabled, cron_expr, timezone, sections, connector_ids, channels, created_by, created_at, updated_at FROM report_definitions WHERE id = ?`, id).Scan(&r.ID, &r.Slug, &r.Name, &enabled, &r.CronExpr, &r.Timezone, &r.Sections, &r.ConnectorIDs, &r.Channels, &r.CreatedBy, &r.CreatedAt, &r.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return r, ErrNotFound
	}
	if err != nil {
		return r, fmt.Errorf("get report definition: %w", err)
	}
	r.Enabled = enabled != 0
	return r, nil
}

// CreateReportDefinition persists a schedule and fills its ID and timestamps.
func (s *Store) CreateReportDefinition(ctx context.Context, r *ReportDefinitionRecord) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}
	if r.UpdatedAt == "" {
		r.UpdatedAt = now
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO report_definitions (id,slug,name,enabled,cron_expr,timezone,sections,connector_ids,channels,created_by,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, r.ID, r.Slug, r.Name, boolToInt(r.Enabled), r.CronExpr, r.Timezone, r.Sections, r.ConnectorIDs, r.Channels, r.CreatedBy, r.CreatedAt, r.UpdatedAt)
	if isUniqueViolation(err) {
		return ErrConflict
	}
	if err != nil {
		return fmt.Errorf("create report definition: %w", err)
	}
	return nil
}

// UpdateReportDefinition updates a schedule's mutable fields.
func (s *Store) UpdateReportDefinition(ctx context.Context, r ReportDefinitionRecord) error {
	r.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	res, err := s.db.ExecContext(ctx, `UPDATE report_definitions SET name=?,enabled=?,cron_expr=?,timezone=?,sections=?,connector_ids=?,channels=?,updated_at=? WHERE id=?`, r.Name, boolToInt(r.Enabled), r.CronExpr, r.Timezone, r.Sections, r.ConnectorIDs, r.Channels, r.UpdatedAt, r.ID)
	if err != nil {
		return fmt.Errorf("update report definition: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteReportDefinition removes a schedule by ID.
func (s *Store) DeleteReportDefinition(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM report_definitions WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("delete report definition: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateReport stores an immutable generated report snapshot.
func (s *Store) CreateReport(ctx context.Context, r *ReportRecord) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.CreatedAt == "" {
		r.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO reports (id,definition_id,definition_name,trigger,period_start,period_end,truncated,data,markdown,status,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`, r.ID, nilToStr(r.DefinitionID), r.DefinitionName, r.Trigger, r.PeriodStart, r.PeriodEnd, boolToInt(r.Truncated), r.Data, r.Markdown, r.Status, r.CreatedAt)
	if err != nil {
		return fmt.Errorf("create report: %w", err)
	}
	return nil
}

// LatestScheduledReport finds the newest scheduled snapshot for a definition.
func (s *Store) LatestScheduledReport(ctx context.Context, definitionID string) (ReportRecord, error) {
	return s.getReport(ctx, `SELECT id,definition_id,definition_name,trigger,period_start,period_end,truncated,data,markdown,status,created_at FROM reports WHERE definition_id=? AND trigger='scheduled' ORDER BY period_end DESC LIMIT 1`, definitionID)
}

// GetReport retrieves a report snapshot by ID.
func (s *Store) GetReport(ctx context.Context, id string) (ReportRecord, error) {
	return s.getReport(ctx, `SELECT id,definition_id,definition_name,trigger,period_start,period_end,truncated,data,markdown,status,created_at FROM reports WHERE id=?`, id)
}

// ListReports returns newest reports first, optionally restricted to a definition.
func (s *Store) ListReports(ctx context.Context, definitionID string, limit, offset int) ([]ReportRecord, int, error) {
	where, args := "", []any{}
	if definitionID != "" {
		where, args = " WHERE definition_id = ?", append(args, definitionID)
	}
	var total int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reports"+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reports: %w", err)
	}
	args = append(args, limit, offset)
	rows, err := s.db.QueryContext(ctx, `SELECT id, definition_id, definition_name, trigger, period_start, period_end, truncated, data, markdown, status, created_at FROM reports`+where+` ORDER BY period_end DESC, id DESC LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	out := []ReportRecord{}
	for rows.Next() {
		var r ReportRecord
		var definitionID sql.NullString
		var truncated int
		if err := rows.Scan(&r.ID, &definitionID, &r.DefinitionName, &r.Trigger, &r.PeriodStart, &r.PeriodEnd, &truncated, &r.Data, &r.Markdown, &r.Status, &r.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan report: %w", err)
		}
		r.DefinitionID = definitionID.String
		r.Truncated = truncated != 0
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate reports: %w", err)
	}
	return out, total, nil
}
func (s *Store) getReport(ctx context.Context, q string, args ...any) (ReportRecord, error) {
	var r ReportRecord
	var id sql.NullString
	var trunc int
	err := s.db.QueryRowContext(ctx, q, args...).Scan(&r.ID, &id, &r.DefinitionName, &r.Trigger, &r.PeriodStart, &r.PeriodEnd, &trunc, &r.Data, &r.Markdown, &r.Status, &r.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return r, ErrNotFound
	}
	if err != nil {
		return r, fmt.Errorf("get report: %w", err)
	}
	r.DefinitionID = id.String
	r.Truncated = trunc != 0
	return r, nil
}
