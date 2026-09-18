package store

import (
	"context"
	"fmt"
	"sort"
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

// attentionSeverityRank returns a numeric rank for sorting: critical=0, warning=1, info=2.
func attentionSeverityRank(severity string) int {
	switch severity {
	case "critical":
		return 0
	case "warning":
		return 1
	case "info":
		return 2
	default:
		return 3
	}
}

// MergedAttentionItems merges pending alerts and open quality findings into a
// single severity-then-recency-sorted attention queue, optionally cut off at
// since (RFC3339), keeps only items on connectors userID holds a viewer grant
// on (default deny), and paginates the result. It runs one runbook lookup per
// distinct (kind, severity/checkType) pair rather than one per item.
func (s *Store) MergedAttentionItems(ctx context.Context, userID, since string, offset, pageSize int) ([]AttentionItem, int, error) {
	alerts, _, err := s.ListAlerts(ctx, "", "", "pending", since, 0, 1000) // ponytail: unbounded fetch for merge, paginate at the store level if this becomes a bottleneck
	if err != nil {
		return nil, 0, fmt.Errorf("list alerts for attention: %w", err)
	}

	findings, _, err := s.ListQualityFindings(ctx, "", "", "open", since, 0, 1000) // ponytail: unbounded fetch for merge, paginate at the store level if this becomes a bottleneck
	if err != nil {
		return nil, 0, fmt.Errorf("list findings for attention: %w", err)
	}

	items := make([]AttentionItem, 0, len(alerts)+len(findings))
	for _, a := range alerts {
		items = append(items, AttentionItem{
			ID:          a.ID,
			Kind:        "alert",
			Severity:    a.Severity,
			Title:       a.Title,
			ConnectorID: a.ServiceID,
			DetectedAt:  a.CreatedAt,
			ChangeID:    a.ChangeID,
		})
	}
	for _, f := range findings {
		items = append(items, AttentionItem{
			ID:          f.ID,
			Kind:        "finding",
			Severity:    f.Severity,
			Title:       f.Title,
			ConnectorID: f.ConnectorID,
			DetectedAt:  f.LastSeenAt,
			FindingType: f.CheckType,
		})
	}

	seen := make(map[string]bool, len(items))
	connectorIDs := make([]string, 0, len(items))
	for _, item := range items {
		if !seen[item.ConnectorID] {
			seen[item.ConnectorID] = true
			connectorIDs = append(connectorIDs, item.ConnectorID)
		}
	}
	allowedIDs, err := s.FilterConnectorIDsByGrant(ctx, userID, connectorIDs, "viewer")
	if err != nil {
		return nil, 0, fmt.Errorf("filter attention by grant: %w", err)
	}
	allowed := make(map[string]bool, len(allowedIDs))
	for _, id := range allowedIDs {
		allowed[id] = true
	}
	visible := items[:0]
	for _, item := range items {
		if allowed[item.ConnectorID] {
			visible = append(visible, item)
		}
	}
	items = visible

	sort.Slice(items, func(i, j int) bool {
		if attentionSeverityRank(items[i].Severity) != attentionSeverityRank(items[j].Severity) {
			return attentionSeverityRank(items[i].Severity) < attentionSeverityRank(items[j].Severity)
		}
		return items[i].DetectedAt > items[j].DetectedAt
	})

	runbooks := make(map[string]*RunbookRecord)
	for i, item := range items {
		targetType, targetValue := "alert_severity", item.Severity
		if item.Kind == "finding" {
			targetType, targetValue = "finding_check_type", item.FindingType
		}
		key := targetType + "|" + targetValue
		rb, cached := runbooks[key]
		if !cached {
			rb, _ = s.GetRunbookByTarget(ctx, targetType, targetValue)
			runbooks[key] = rb
		}
		if rb != nil {
			items[i].RunbookID = rb.ID
		}
	}

	total := len(items)
	start := offset
	if start > total {
		start = total
	}
	end := offset + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total, nil
}
