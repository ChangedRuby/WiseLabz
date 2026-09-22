package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/auth"
)

// AuditRecord represents a row in the audit_log table: who did what to which
// target, and when. See docs/AUDIT.md for exactly which actions are covered.
type AuditRecord struct {
	ID          string `json:"id"`
	ActorUserID string `json:"actorUserId"`
	ActorRole   string `json:"actorRole"`
	Action      string `json:"action"`
	TargetType  string `json:"targetType"`
	TargetID    string `json:"targetId"`
	Detail      string `json:"detail"`
	CreatedAt   string `json:"createdAt"`
}

// CreateAuditRecord inserts an audit_log row, filling ID/CreatedAt/Detail
// defaults when left zero-valued (same convention as CreateChange/CreateAlert).
func (s *Store) CreateAuditRecord(ctx context.Context, a *AuditRecord) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.CreatedAt == "" {
		a.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if a.Detail == "" {
		a.Detail = "{}"
	}
	if err := s.insertAuditRecords(ctx, []AuditRecord{*a}); err != nil {
		return err
	}
	return nil
}

// CreateAuditRecords inserts audit_log rows in one statement.
func (s *Store) CreateAuditRecords(ctx context.Context, records []AuditRecord) error {
	if len(records) == 0 {
		return nil
	}

	for i := range records {
		a := &records[i]
		if a.ID == "" {
			a.ID = uuid.New().String()
		}
		if a.CreatedAt == "" {
			a.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		}
		if a.Detail == "" {
			a.Detail = "{}"
		}
	}
	return s.insertAuditRecords(ctx, records)
}

func (s *Store) insertAuditRecords(ctx context.Context, records []AuditRecord) error {
	args := make([]any, 0, len(records)*8)
	values := make([]string, 0, len(records))
	for _, a := range records {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args, a.ID, a.ActorUserID, a.ActorRole, a.Action, a.TargetType, a.TargetID, a.Detail, a.CreatedAt)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO audit_log (id, actor_user_id, actor_role, action, target_type, target_id, detail, created_at)
		VALUES `+strings.Join(values, ", "), args...)
	if err != nil {
		return fmt.Errorf("create audit records: %w", err)
	}
	return nil
}

// RecordAuditBatchFromContext records one audit entry for each target using the
// authenticated actor in ctx.
func (s *Store) RecordAuditBatchFromContext(ctx context.Context, action, targetType string, records []AuditRecord) error {
	for i := range records {
		records[i].ActorUserID = auth.UserIDFromContext(ctx)
		records[i].ActorRole = actorRoleLabel(ctx)
		records[i].Action = action
		records[i].TargetType = targetType
	}
	return s.CreateAuditRecords(ctx, records)
}

// RecordAuditFromContext is the one-line call handlers make after an audited
// action succeeds: it marshals detail to JSON (nil -> "{}") and pulls the
// actor from the request context (set by auth.AuthMiddleware).
//
// Called by handlers for audited actions; most call sites record only after
// success, with security-relevant elevation denials as the documented
// exception. Call sites treat a returned error as non-fatal (slog.Error and
// continue); see docs/AUDIT.md.
func (s *Store) RecordAuditFromContext(ctx context.Context, action, targetType, targetID string, detail any) error {
	detailJSON := ""
	if detail != nil {
		data, err := json.Marshal(detail)
		if err != nil {
			return fmt.Errorf("marshal audit detail: %w", err)
		}
		detailJSON = string(data)
	}

	return s.CreateAuditRecord(ctx, &AuditRecord{
		ActorUserID: auth.UserIDFromContext(ctx),
		ActorRole:   actorRoleLabel(ctx),
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Detail:      detailJSON,
	})
}

// actorRoleLabel records the actor's flat instance-admin role for the audit
// trail ("admin"/"user"). Per-connector grant changes are self-describing
// via the audit action/target/detail, so this stays the coarse instance-wide
// label rather than trying to cram a connector-scoped role into one column.
// Empty when there's no authenticated actor in ctx (mirrors ActorUserID).
func actorRoleLabel(ctx context.Context) string {
	if auth.UserIDFromContext(ctx) == "" {
		return ""
	}
	if auth.InstanceAdminFromContext(ctx) {
		return "admin"
	}
	return "user"
}

// ListAuditRecords returns a paginated list of audit records, newest first,
// optionally filtered by action, target type, and/or a created_at range
// (createdAfter/createdBefore are RFC3339 timestamps; empty = unbounded).
func (s *Store) ListAuditRecords(ctx context.Context, action, targetType, createdAfter, createdBefore string, offset, limit int) ([]AuditRecord, int, error) {
	where, args := auditFilterClause(action, targetType, createdAfter, createdBefore)

	return paginatedQuery(ctx, s.db, "audit_log", auditColumns, where, args, "created_at DESC", limit, offset, scanAuditRecord)
}

// ListAuditRecordsKeyset returns one keyset (cursor) page of audit records,
// newest first, using the same filters as ListAuditRecords. Rows strictly
// before cur in (created_at, id) order are returned; a zero cur starts at the
// newest record. total counts every record matching the filters, ignoring cur.
func (s *Store) ListAuditRecordsKeyset(ctx context.Context, action, targetType, createdAfter, createdBefore string, cur Keyset, limit int) ([]AuditRecord, int, error) {
	where, args := auditFilterClause(action, targetType, createdAfter, createdBefore)

	return keysetQuery(ctx, s.db, "audit_log", auditColumns, where, args, "created_at", cur, limit, scanAuditRecord)
}

// ListAllAuditRecords returns every audit record matching the given filters
// (no pagination), newest first. Used by the export endpoint, which needs
// the full matching set rather than one page — kept as a separate method
// instead of overloading ListAuditRecords with a "limit<=0 means unbounded"
// convention, since that would change what a zero/negative limit means for
// its existing paginated caller.
func (s *Store) ListAllAuditRecords(ctx context.Context, action, targetType, createdAfter, createdBefore string) ([]AuditRecord, error) {
	where, args := auditFilterClause(action, targetType, createdAfter, createdBefore)

	query := `SELECT ` + auditColumns + `
		FROM audit_log ` + where + ` ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list all audit records: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	return scanAuditRecordRows(rows)
}

// auditFilterClause builds the shared WHERE clause + args for
// ListAuditRecords and ListAllAuditRecords.
func auditFilterClause(action, targetType, createdAfter, createdBefore string) (string, []any) {
	where := "WHERE 1=1"
	var args []any
	if action != "" {
		where += " AND action = ?"
		args = append(args, action)
	}
	if targetType != "" {
		where += " AND target_type = ?"
		args = append(args, targetType)
	}
	if createdAfter != "" {
		where += " AND created_at >= ?"
		args = append(args, createdAfter)
	}
	if createdBefore != "" {
		where += " AND created_at <= ?"
		args = append(args, createdBefore)
	}
	return where, args
}

// auditColumns is the shared column list for every audit_log SELECT.
const auditColumns = `id, actor_user_id, actor_role, action, target_type, target_id, detail, created_at`

func scanAuditRecord(row rowScanner) (AuditRecord, error) {
	var a AuditRecord
	err := row.Scan(&a.ID, &a.ActorUserID, &a.ActorRole, &a.Action, &a.TargetType, &a.TargetID, &a.Detail, &a.CreatedAt)
	return a, err
}

// scanAuditRecordRows scans all rows of an audit_log query into []AuditRecord,
// returning a non-nil empty slice (never nil) when there are no rows.
func scanAuditRecordRows(rows *sql.Rows) ([]AuditRecord, error) {
	records := make([]AuditRecord, 0)
	for rows.Next() {
		a, err := scanAuditRecord(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		records = append(records, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit records: %w", err)
	}
	return records, nil
}
