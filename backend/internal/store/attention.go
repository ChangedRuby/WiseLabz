package store

import (
	"context"
	"fmt"
)

// AttentionItem represents a merged alert or finding for the attention queue.
type AttentionItem struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"` // "alert" or "finding"
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	ConnectorID string `json:"connectorId"`
	DetectedAt  string `json:"detectedAt"`
	ChangeID    string `json:"changeId,omitempty"`  // alerts only
	FindingType string `json:"checkType,omitempty"` // findings only
	RunbookID   string `json:"runbookId,omitempty"`
}

// attentionUnion is the pending-alert / open-finding union, already narrowed
// to connectors the caller holds a grant on and joined to its runbook. The
// two %s slots take the since and API-key connector filters; both selects
// bind: userID, roles..., [since], [keyIDs...].
const attentionUnion = `
	SELECT a.id AS id, 'alert' AS kind, a.severity AS severity, a.title AS title,
	       a.service_id AS connector_id, a.created_at AS detected_at,
	       COALESCE(a.change_id, '') AS change_id, '' AS check_type,
	       COALESCE(rb.id, '') AS runbook_id
	FROM alerts a
	JOIN user_connector_roles g ON g.connector_id = a.service_id AND g.user_id = ? AND g.role IN (%s)
	LEFT JOIN runbooks rb ON rb.target_type = 'alert_severity' AND rb.target_value = a.severity
	WHERE a.status = 'pending'%s
	UNION ALL
	SELECT f.id, 'finding', f.severity, f.title, f.connector_id, f.last_seen_at,
	       '', f.check_type, COALESCE(rb.id, '')
	FROM quality_findings f
	JOIN user_connector_roles g ON g.connector_id = f.connector_id AND g.user_id = ? AND g.role IN (%s)
	LEFT JOIN runbooks rb ON rb.target_type = 'finding_check_type' AND rb.target_value = f.check_type
	WHERE f.status = 'open'%s`

// MergedAttentionItems merges pending alerts and open quality findings into a
// single severity-then-recency-sorted attention queue, optionally cut off at
// since (RFC3339), keeps only items on connectors userID holds a viewer grant
// on (default deny), and paginates the result. Filtering, runbook lookup,
// ordering and pagination all happen in one UNION ALL query in SQL.
func (s *Store) MergedAttentionItems(ctx context.Context, userID, since string, offset, pageSize int) ([]AttentionItem, int, error) {
	var roles []any
	for role, rank := range connectorRoleRank {
		if rank >= connectorRoleRank["viewer"] {
			roles = append(roles, role)
		}
	}
	rolePH := placeholders(len(roles))
	alertSince, findingSince := "", ""
	if since != "" {
		alertSince = " AND a.created_at >= ?"
		findingSince = " AND f.last_seen_at >= ?"
	}
	alertKeyFilter, keyArgs := apiKeyConnectorFilter(ctx, "a.service_id")
	findingKeyFilter, _ := apiKeyConnectorFilter(ctx, "f.connector_id")
	union := fmt.Sprintf(attentionUnion, rolePH, alertSince+alertKeyFilter, rolePH, findingSince+findingKeyFilter)

	var args []any
	args = append(args, userID)
	args = append(args, roles...)
	if since != "" {
		args = append(args, since)
	}
	args = append(args, keyArgs...)
	args = append(args, userID)
	args = append(args, roles...)
	if since != "" {
		args = append(args, since)
	}
	args = append(args, keyArgs...)

	var total int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ("+union+") u", args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count attention: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM (`+union+`) u
		ORDER BY CASE severity WHEN 'critical' THEN 0 WHEN 'warning' THEN 1 WHEN 'info' THEN 2 ELSE 3 END,
		         detected_at DESC, id
		LIMIT ? OFFSET ?`, append(args, pageSize, offset)...)
	if err != nil {
		return nil, 0, fmt.Errorf("list attention: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	items := make([]AttentionItem, 0)
	for rows.Next() {
		var it AttentionItem
		if err := rows.Scan(&it.ID, &it.Kind, &it.Severity, &it.Title, &it.ConnectorID,
			&it.DetectedAt, &it.ChangeID, &it.FindingType, &it.RunbookID); err != nil {
			return nil, 0, fmt.Errorf("scan attention: %w", err)
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate attention: %w", err)
	}
	return items, total, nil
}
