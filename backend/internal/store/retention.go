package store

import (
	"context"
	"fmt"
)

// retentionBatchSize bounds how many rows a single retention DELETE removes,
// so one statement never holds the (single, on SQLite) connection or a
// Postgres transaction for long.
const retentionBatchSize = 500

// batchDelete deletes rows of table matching where (which may reference the
// table by the alias "t") in retentionBatchSize chunks, looping until a batch
// removes nothing and checking ctx between batches. It returns the total
// number of rows deleted, including on ctx cancellation or error.
func (s *Store) batchDelete(ctx context.Context, table, where string, args ...any) (int64, error) {
	query := fmt.Sprintf(
		`DELETE FROM %[1]s WHERE id IN (SELECT t.id FROM %[1]s t WHERE %[2]s LIMIT %[3]d)`,
		table, where, retentionBatchSize)
	var total int64
	for {
		if err := ctx.Err(); err != nil {
			return total, err
		}
		res, err := s.db.ExecContext(ctx, query, args...)
		if err != nil {
			return total, err
		}
		n := rowsAffected(res)
		total += n
		if n < retentionBatchSize {
			return total, nil
		}
	}
}

// DeleteOldSnapshots removes service_snapshots rows fetched before cutoff,
// never deleting a connector's most recent snapshot (GetLatestSnapshot's
// "current data" for that service).
func (s *Store) DeleteOldSnapshots(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "service_snapshots", `
		t.fetched_at < ?
		AND EXISTS (
			SELECT 1 FROM service_snapshots n
			WHERE n.connector_id = t.connector_id AND n.fetched_at > t.fetched_at
		)`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old snapshots: %w", err)
	}
	return n, nil
}

// DeleteOldDocVersions removes doc_versions rows created before cutoff,
// never deleting the revision matching the doc's current_version.
func (s *Store) DeleteOldDocVersions(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "doc_versions", `
		t.created_at < ?
		AND t.rev != (SELECT current_version FROM docs WHERE docs.id = t.doc_id)`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old doc versions: %w", err)
	}
	return n, nil
}

// DeleteOldAlerts removes resolved/dismissed alerts created before cutoff.
// Pending and snoozed alerts are active workflow state and are never purged
// regardless of age.
func (s *Store) DeleteOldAlerts(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "alerts",
		`t.created_at < ? AND t.status IN ('resolved', 'dismissed')`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old alerts: %w", err)
	}
	return n, nil
}

// DeleteOldSyncRuns removes sync_runs rows started before cutoff. A
// connector's last-sync status lives on the connectors row itself, so no
// "keep latest" guard is needed here.
func (s *Store) DeleteOldSyncRuns(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "sync_runs", `t.started_at < ?`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old sync runs: %w", err)
	}
	return n, nil
}

// DeleteOldAuditRecords removes audit_log rows created before cutoff. Like
// sync runs, there's no "keep latest" guard — the audit trail has no
// current-state row that needs preserving.
func (s *Store) DeleteOldAuditRecords(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "audit_log", `t.created_at < ?`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old audit records: %w", err)
	}
	return n, nil
}

// DeleteStaleSessions removes sessions not seen since cutoff.
func (s *Store) DeleteStaleSessions(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "sessions", `t.last_seen_at < ?`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete stale sessions: %w", err)
	}
	return n, nil
}

// DeleteOldNotifications removes read in-app notifications created before
// cutoff (their deliveries cascade). Unread notifications are kept.
func (s *Store) DeleteOldNotifications(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "in_app_notifications", `t.created_at < ? AND t.read = 1`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old notifications: %w", err)
	}
	return n, nil
}

// DeleteOldDeliveries removes finished (sent or failed with no retry
// scheduled) notification deliveries created before cutoff.
func (s *Store) DeleteOldDeliveries(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "notification_deliveries", `
		t.created_at < ?
		AND (t.status = 'sent' OR (t.status = 'failed' AND t.next_attempt_at IS NULL))`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old deliveries: %w", err)
	}
	return n, nil
}

// DeleteOldChanges removes acknowledged/dismissed changes detected before
// cutoff. Linked alerts survive with change_id set to NULL.
func (s *Store) DeleteOldChanges(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "changes",
		`t.detected_at < ? AND t.status IN ('acknowledged', 'dismissed')`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old changes: %w", err)
	}
	return n, nil
}

// DeleteExpiredShareLinks removes share links that expired or were revoked
// before cutoff. An empty expires_at/revoked_at means "never".
func (s *Store) DeleteExpiredShareLinks(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "share_links", `
		(t.expires_at != '' AND t.expires_at < ?)
		OR (t.revoked_at != '' AND t.revoked_at < ?)`, cutoff, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete expired share links: %w", err)
	}
	return n, nil
}

// DeleteOldChatConversations removes chat conversations created before
// cutoff (their messages cascade).
func (s *Store) DeleteOldChatConversations(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "chat_conversations", `t.created_at < ?`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old chat conversations: %w", err)
	}
	return n, nil
}

// rowsAffected returns res.RowsAffected(), or 0 if the driver doesn't
// support it — retention counts are informational, not worth failing over.
func rowsAffected(res interface{ RowsAffected() (int64, error) }) int64 {
	n, err := res.RowsAffected()
	if err != nil {
		return 0
	}
	return n
}
