// Package quality detects actionable documentation quality gaps.
package quality

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

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
)

// Checker runs documentation quality checks and broadcasts detected findings.
type Checker struct {
	store *store.Store
	hub   *ws.Hub
	now   func() time.Time
}

// NewChecker creates a documentation quality checker.
func NewChecker(s *store.Store, hub *ws.Hub) *Checker {
	return &Checker{store: s, hub: hub, now: time.Now}
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
	c.broadcastChanged(connectorID)
	return errors.Join(errs...)
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
	if finding.ID != candidateID {
		return nil, nil
	}
	return finding, nil
}

// RunStaleSweepOnce performs one pass of stale documentation checks for all connectors.
func RunStaleSweepOnce(ctx context.Context, s *store.Store, hub *ws.Hub, logger *slog.Logger) {
	checker := NewChecker(s, hub)
	connectors, err := checker.store.ListAllConnectors(ctx)
	if err != nil {
		logger.Error("list connectors for stale sweep", "error", err)
		return
	}
	for _, connector := range connectors {
		docs, err := checker.store.ListDocsByService(ctx, connector.ID)
		if err != nil {
			logger.Error("list docs for stale sweep", "connector", connector.ID, "error", err)
			continue
		}
		finding, err := checker.checkStale(ctx, connector.ID, docs)
		if err != nil {
			logger.Error("quality stale check failed", "connector", connector.ID, "error", err)
			continue
		}
		if finding != nil && checker.hub != nil {
			checker.broadcastCreated(connector.ID, finding)
		}
		checker.broadcastChanged(connector.ID)
	}
}
