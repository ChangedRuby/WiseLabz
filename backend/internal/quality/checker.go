// Package quality detects actionable documentation quality gaps.
package quality

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/compliance"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
	"github.com/google/uuid"
)

const (
	// StaleThreshold is fixed for v1; ponytail: make it configurable if labs need different freshness policies.
	StaleThreshold = 30 * 24 * time.Hour
	// EmptyContentMinChars is the minimum trimmed document length considered useful.
	EmptyContentMinChars = 40
	// ConsecutiveFailuresThreshold is the failed-sync streak that opens a finding.
	ConsecutiveFailuresThreshold = 3
	complianceEntityNameLimit    = 10
)

// severityRank orders finding severities so the notification hook can tell
// an escalation (rank increases) from a de-escalation or repeat (it doesn't).
var severityRank = map[string]int{"info": 0, "warning": 1, "critical": 2}

// FindingNotifier is implemented by notifications.Dispatcher. A narrow
// interface keeps quality decoupled from the notifications package.
type FindingNotifier interface {
	NotifyFindingCreated(ctx context.Context, findingID, title, message string)
}

// RotationConfig is the credential_rotation check's default policy —
// backend/internal/config's rotation.max_age_days / rotation.warn_days,
// overridable per connector via ConnectorRecord.RotationMaxAgeDays.
type RotationConfig struct {
	MaxAgeDays int
	WarnDays   int
}

// Checker runs documentation quality checks and broadcasts detected findings.
type Checker struct {
	store    *store.Store
	hub      *ws.Hub
	notifier FindingNotifier
	rotation RotationConfig
	now      func() time.Time
}

// NewChecker creates a documentation quality checker. notifier may be nil
// (no finding notifications dispatched, e.g. in tests that don't need them).
func NewChecker(s *store.Store, hub *ws.Hub, notifier FindingNotifier, rotation RotationConfig) *Checker {
	return &Checker{store: s, hub: hub, notifier: notifier, rotation: rotation, now: time.Now}
}

// RunForConnector runs every quality check for one connector.
func (c *Checker) RunForConnector(ctx context.Context, connectorID string) error {
	docs, err := c.store.ListDocsByService(ctx, connectorID)
	if err != nil {
		return fmt.Errorf("list docs: %w", err)
	}

	checks := []struct {
		name string
		run  func(context.Context, string, []store.DocRecord) (*store.QualityFindingRecord, error)
	}{
		{name: "stale", run: c.checkStale},
		{name: "empty", run: c.checkEmpty},
		{name: "failing", run: func(ctx context.Context, connectorID string, _ []store.DocRecord) (*store.QualityFindingRecord, error) {
			return c.checkFailing(ctx, connectorID)
		}},
		{name: "ownership", run: func(ctx context.Context, connectorID string, _ []store.DocRecord) (*store.QualityFindingRecord, error) {
			return c.checkOwnership(ctx, connectorID)
		}},
		{name: "credential_rotation", run: func(ctx context.Context, connectorID string, _ []store.DocRecord) (*store.QualityFindingRecord, error) {
			return c.checkCredentialRotation(ctx, connectorID)
		}},
	}

	var errs []error
	for _, check := range checks {
		finding, err := check.run(ctx, connectorID, docs)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s check: %w", check.name, err))
			continue
		}
		if finding != nil && c.hub != nil {
			c.broadcastCreated(connectorID, finding)
		}
	}
	complianceFindings, err := c.checkCompliance(ctx, connectorID)
	if err != nil {
		errs = append(errs, fmt.Errorf("compliance check: %w", err))
	}
	for _, finding := range complianceFindings {
		if finding != nil && c.hub != nil {
			c.broadcastCreated(connectorID, finding)
		}
	}
	c.broadcastChanged(connectorID)
	return errors.Join(errs...)
}

// EvaluateRule evaluates an enabled rule against the latest snapshots for its
// matching connectors. It is used after a rule is created, changed, or enabled.
func (c *Checker) EvaluateRule(ctx context.Context, ruleID string) error {
	record, err := c.store.GetComplianceRule(ctx, ruleID)
	if err != nil {
		return err
	}
	if !record.Enabled {
		return c.ResolveRule(ctx, ruleID)
	}
	rule, err := complianceRule(record)
	if err != nil {
		return err
	}
	connectors, err := c.store.ListAllConnectors(ctx)
	if err != nil {
		return err
	}
	var errs []error
	for _, conn := range connectors {
		if conn.Type != rule.ConnectorType {
			continue
		}
		snapshot, err := c.loadComplianceSnapshot(ctx, conn.ID)
		if err != nil {
			errs = append(errs, fmt.Errorf("connector %s: %w", conn.ID, err))
			continue
		}
		finding, err := c.evaluateComplianceRule(ctx, conn.ID, rule, snapshot)
		if err != nil {
			errs = append(errs, fmt.Errorf("connector %s: %w", conn.ID, err))
			continue
		}
		if finding != nil && c.hub != nil {
			c.broadcastCreated(conn.ID, finding)
		}
		c.broadcastChanged(conn.ID)
	}
	return errors.Join(errs...)
}

// ResolveRule resolves all open findings for a disabled or deleted rule.
func (c *Checker) ResolveRule(ctx context.Context, ruleID string) error {
	return c.store.ResolveQualityFindingsForRule(ctx, ruleID)
}

func (c *Checker) checkCompliance(ctx context.Context, connectorID string) ([]*store.QualityFindingRecord, error) {
	conn, err := c.store.GetConnector(ctx, connectorID)
	if err != nil {
		return nil, err
	}
	records, err := c.store.ListComplianceRules(ctx)
	if err != nil {
		return nil, err
	}
	findings := make([]*store.QualityFindingRecord, 0)
	var errs []error
	var snapshot *compliance.Snapshot
	var snapshotErr error
	snapshotLoaded := false
	for i := range records {
		if !records[i].Enabled || records[i].ConnectorType != conn.Type {
			continue
		}
		rule, err := complianceRule(&records[i])
		if err != nil {
			errs = append(errs, fmt.Errorf("rule %s: %w", records[i].ID, err))
			continue
		}
		if !snapshotLoaded {
			snapshot, snapshotErr = c.loadComplianceSnapshot(ctx, connectorID)
			snapshotLoaded = true
		}
		if snapshotErr != nil {
			errs = append(errs, fmt.Errorf("rule %s: %w", records[i].ID, snapshotErr))
			continue
		}
		finding, err := c.evaluateComplianceRule(ctx, connectorID, rule, snapshot)
		if err != nil {
			errs = append(errs, fmt.Errorf("rule %s: %w", records[i].ID, err))
			continue
		}
		if finding != nil {
			findings = append(findings, finding)
		}
	}
	return findings, errors.Join(errs...)
}

func complianceRule(record *store.ComplianceRuleRecord) (compliance.Rule, error) {
	var conditions []compliance.Condition
	if err := json.Unmarshal([]byte(record.Conditions), &conditions); err != nil {
		return compliance.Rule{}, fmt.Errorf("decode conditions: %w", err)
	}
	return compliance.Rule{
		ID: record.ID, Name: record.Name, ConnectorType: record.ConnectorType,
		EntityKind: record.EntityKind, Conditions: conditions, Severity: record.Severity,
		Title: record.Title, RemediationLink: record.RemediationLink, Enabled: record.Enabled,
	}, nil
}

func (c *Checker) loadComplianceSnapshot(ctx context.Context, connectorID string) (*compliance.Snapshot, error) {
	record, err := c.store.GetLatestSnapshot(ctx, connectorID)
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var snapshot connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(record.Data), &snapshot); err != nil {
		slog.Warn("skipping malformed snapshot for compliance rule", "error", err)
		return nil, nil
	}
	entities := make([]compliance.Entity, len(snapshot.Entities))
	for i, entity := range snapshot.Entities {
		entities[i] = compliance.Entity{Kind: entity.Kind, Name: entity.Name, Attributes: entity.Attributes}
	}
	return &compliance.Snapshot{Entities: entities}, nil
}

func (c *Checker) evaluateComplianceRule(ctx context.Context, connectorID string, rule compliance.Rule, snapshot *compliance.Snapshot) (*store.QualityFindingRecord, error) {
	// Missing or malformed snapshots must not resolve existing findings.
	if snapshot == nil {
		return nil, nil
	}
	matches := compliance.Evaluate(rule, *snapshot)
	if len(matches) == 0 {
		return nil, c.store.ResolveQualityFindingForRule(ctx, connectorID, rule.ID)
	}
	finding := &store.QualityFindingRecord{
		ConnectorID: connectorID, RuleID: rule.ID, CheckType: "compliance",
		Severity: rule.Severity, Title: rule.Title, Description: complianceDescription(matches),
		RemediationLink: rule.RemediationLink,
	}
	return c.upsert(ctx, finding)
}

func complianceDescription(matches []compliance.Entity) string {
	limit := len(matches)
	if limit > complianceEntityNameLimit {
		limit = complianceEntityNameLimit
	}
	names := make([]string, limit)
	for i := range names {
		names[i] = matches[i].Name
	}
	description := "Violating entities: " + strings.Join(names, ", ")
	if remaining := len(matches) - limit; remaining > 0 {
		description += fmt.Sprintf(", and %d more", remaining)
	}
	return description + "."
}

func (c *Checker) broadcastCreated(connectorID string, finding *store.QualityFindingRecord) {
	c.hub.Broadcast(ws.EventQualityFindingCreated, map[string]any{
		"findingId": finding.ID, "connectorId": connectorID,
		"checkType": finding.CheckType, "severity": finding.Severity,
	})
}

func (c *Checker) broadcastChanged(connectorID string) {
	if c.hub != nil {
		c.hub.Broadcast(ws.EventQualityFindingsChanged, map[string]any{"connectorId": connectorID})
	}
}

func (c *Checker) checkStale(ctx context.Context, connectorID string, docs []store.DocRecord) (*store.QualityFindingRecord, error) {
	var selected *store.DocRecord
	var selectedUpdatedAt time.Time
	now := c.now().UTC()
	for i := range docs {
		updatedAt, err := time.Parse(time.RFC3339, docs[i].UpdatedAt)
		if err != nil {
			slog.Warn("skipping doc with malformed updated_at", "doc", docs[i].ID, "error", err)
			continue
		}
		if now.Sub(updatedAt) <= StaleThreshold {
			continue
		}
		if selected == nil || updatedAt.Before(selectedUpdatedAt) ||
			(updatedAt.Equal(selectedUpdatedAt) && docs[i].ID < selected.ID) {
			selected = &docs[i]
			selectedUpdatedAt = updatedAt
		}
	}

	if selected == nil {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "stale")
	}
	finding := &store.QualityFindingRecord{
		ConnectorID:     connectorID,
		DocID:           selected.ID,
		CheckType:       "stale",
		Severity:        "warning",
		Title:           "Documentation is stale",
		Description:     fmt.Sprintf("%q has not been updated in more than 30 days.", selected.Title),
		RemediationLink: "/docs/" + selected.ID,
	}
	return c.upsert(ctx, finding)
}

func (c *Checker) checkEmpty(ctx context.Context, connectorID string, docs []store.DocRecord) (*store.QualityFindingRecord, error) {
	var selected *store.DocRecord
	selectedLength := 0
	for i := range docs {
		contentLength := len(strings.TrimSpace(docs[i].Content))
		if contentLength >= EmptyContentMinChars {
			continue
		}
		if selected == nil || contentLength < selectedLength ||
			(contentLength == selectedLength && docs[i].ID < selected.ID) {
			selected = &docs[i]
			selectedLength = contentLength
		}
	}

	if selected == nil {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "empty")
	}
	finding := &store.QualityFindingRecord{
		ConnectorID:     connectorID,
		DocID:           selected.ID,
		CheckType:       "empty",
		Severity:        "warning",
		Title:           "Documentation is empty",
		Description:     fmt.Sprintf("%q has fewer than %d characters of content.", selected.Title, EmptyContentMinChars),
		RemediationLink: "/docs/" + selected.ID + "/edit",
	}
	return c.upsert(ctx, finding)
}

func (c *Checker) checkFailing(ctx context.Context, connectorID string) (*store.QualityFindingRecord, error) {
	runs, err := c.store.ListSyncRunsByConnector(ctx, connectorID, ConsecutiveFailuresThreshold)
	if err != nil {
		return nil, err
	}
	if len(runs) < ConsecutiveFailuresThreshold {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "failing")
	}
	for _, run := range runs {
		if run.Status != store.SyncRunStatusError {
			return nil, c.store.ResolveQualityFinding(ctx, connectorID, "failing")
		}
	}

	finding := &store.QualityFindingRecord{
		ConnectorID:     connectorID,
		CheckType:       "failing",
		Severity:        "critical",
		Title:           "Connector sync is failing",
		Description:     fmt.Sprintf("The last %d sync attempts failed.", ConsecutiveFailuresThreshold),
		RemediationLink: "/services/" + connectorID,
	}
	return c.upsert(ctx, finding)
}

func (c *Checker) checkOwnership(ctx context.Context, connectorID string) (*store.QualityFindingRecord, error) {
	connector, err := c.store.GetConnector(ctx, connectorID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(connector.Owner) != "" {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "ownership_incomplete")
	}

	finding := &store.QualityFindingRecord{
		ConnectorID:     connectorID,
		CheckType:       "ownership_incomplete",
		Severity:        "info",
		Title:           "Documentation owner is missing",
		Description:     "Assign an owner so this connector's documentation has a clear maintainer.",
		RemediationLink: "/connectors/" + connectorID + "/edit",
	}
	return c.upsert(ctx, finding)
}

func (c *Checker) upsert(ctx context.Context, finding *store.QualityFindingRecord) (*store.QualityFindingRecord, error) {
	candidateID := uuid.New().String()
	finding.ID = candidateID
	if err := c.store.UpsertQualityFinding(ctx, finding); err != nil {
		return nil, err
	}
	isNew := finding.ID == candidateID
	c.maybeNotify(ctx, finding)
	if !isNew {
		return nil, nil
	}
	return finding, nil
}

// maybeNotify dispatches a finding notification the first time a finding is
// opened, and again whenever its severity escalates past what was last
// notified — never on a repeat detection at the same (or lower) severity.
// ResolveQualityFinding clears NotifiedSeverity, so a resolve-then-reopen is
// treated as new. Generic across every check type: it only looks at
// Severity/NotifiedSeverity, never CheckType.
func (c *Checker) maybeNotify(ctx context.Context, finding *store.QualityFindingRecord) {
	if c.notifier == nil {
		return
	}
	if finding.NotifiedSeverity != "" && severityRank[finding.Severity] <= severityRank[finding.NotifiedSeverity] {
		return
	}
	c.notifier.NotifyFindingCreated(ctx, finding.ID, finding.Title, finding.Description)
	if err := c.store.SetQualityFindingNotifiedSeverity(ctx, finding.ID, finding.Severity); err != nil {
		slog.Error("failed to record finding notification", "finding", finding.ID, "error", err)
	}
}

// checkCredentialRotation flags a connector whose secret hasn't been rotated
// within its rotation window. Connectors whose implementation refreshes its
// own credentials (CredentialRefresher) are never flagged — nothing for a
// human to rotate. due = min(UserExpiresAt, SecretRotatedAt + maxAgeDays),
// using the connector's RotationMaxAgeDays override when set, else the
// configured global default.
func (c *Checker) checkCredentialRotation(ctx context.Context, connectorID string) (*store.QualityFindingRecord, error) {
	conn, err := c.store.GetConnector(ctx, connectorID)
	if err != nil {
		return nil, err
	}
	if connector.IsCredentialRefresherType(conn.Type) {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "credential_rotation")
	}

	rotatedAt, err := time.Parse(time.RFC3339, conn.SecretRotatedAt)
	if err != nil {
		// No usable rotation timestamp: nothing to flag rather than a false positive.
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "credential_rotation")
	}

	maxAgeDays := c.rotation.MaxAgeDays
	if conn.RotationMaxAgeDays != nil {
		maxAgeDays = *conn.RotationMaxAgeDays
	}
	due := rotatedAt.AddDate(0, 0, maxAgeDays)
	if conn.UserExpiresAt != "" {
		if userDue, err := time.Parse(time.RFC3339, conn.UserExpiresAt); err == nil && userDue.Before(due) {
			due = userDue
		}
	}

	now := c.now().UTC()
	var severity string
	switch {
	case !now.Before(due):
		severity = "critical"
	case !now.Before(due.AddDate(0, 0, -c.rotation.WarnDays)):
		severity = "warning"
	default:
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "credential_rotation")
	}

	ageDays := int(now.Sub(rotatedAt).Hours() / 24)
	finding := &store.QualityFindingRecord{
		ConnectorID:     connectorID,
		CheckType:       "credential_rotation",
		Severity:        severity,
		Title:           "Credential rotation due",
		Description:     fmt.Sprintf("Secret last rotated %d days ago; due %s.", ageDays, due.Format(time.RFC3339)),
		RemediationLink: "/connectors/" + connectorID + "/edit",
	}
	return c.upsert(ctx, finding)
}

// RunStaleSweepOnce performs one pass of every quality check for all connectors.
// Its legacy name is kept because the quality cron already calls it.
func RunStaleSweepOnce(ctx context.Context, s *store.Store, hub *ws.Hub, notifier FindingNotifier, logger *slog.Logger) {
	checker := NewChecker(s, hub, notifier, RotationConfig{})
	connectors, err := checker.store.ListAllConnectors(ctx)
	if err != nil {
		logger.Error("list connectors for stale sweep", "error", err)
		return
	}
	for _, connector := range connectors {
		if err := checker.RunForConnector(ctx, connector.ID); err != nil {
			logger.Error("quality sweep failed", "connector", connector.ID, "error", err)
			continue
		}
	}
}
