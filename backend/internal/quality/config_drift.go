package quality

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
)

// configDriftEntryLimit caps how many drifted sections are named in a
// finding's description, mirroring complianceEntityNameLimit.
const configDriftEntryLimit = 10

// checkConfigDrift compares the connector's latest snapshot against its
// pinned golden snapshot (not just the previous snapshot, which
// sync/run.go's change detection already covers) and raises a finding when
// they deviate. No golden snapshot pinned, no latest snapshot, or the
// latest snapshot IS the golden snapshot all resolve any existing finding
// rather than raise a false positive.
func (c *Checker) checkConfigDrift(ctx context.Context, connectorID string) (*store.QualityFindingRecord, error) {
	golden, err := c.store.GetGoldenSnapshot(ctx, connectorID)
	if errors.Is(err, store.ErrNotFound) {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "config_drift")
	}
	if err != nil {
		return nil, err
	}

	latest, err := c.store.GetLatestSnapshot(ctx, connectorID)
	if errors.Is(err, store.ErrNotFound) {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "config_drift")
	}
	if err != nil {
		return nil, err
	}
	if latest.ID == golden.SnapshotID {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "config_drift")
	}

	goldenRecord, err := c.store.GetSnapshotByID(ctx, golden.SnapshotID)
	if errors.Is(err, store.ErrNotFound) {
		// The pinned snapshot itself was removed (e.g. connector data purge).
		// Nothing usable to compare against.
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "config_drift")
	}
	if err != nil {
		return nil, err
	}

	var goldenSnap, latestSnap connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(goldenRecord.Data), &goldenSnap); err != nil {
		slog.Warn("skipping malformed golden snapshot for drift check", "connector", connectorID, "error", err)
		return nil, nil
	}
	if err := json.Unmarshal([]byte(latest.Data), &latestSnap); err != nil {
		slog.Warn("skipping malformed latest snapshot for drift check", "connector", connectorID, "error", err)
		return nil, nil
	}

	results := sync.Compare(&goldenSnap, &latestSnap)
	if len(results) == 0 {
		return nil, c.store.ResolveQualityFinding(ctx, connectorID, "config_drift")
	}

	finding := &store.QualityFindingRecord{
		ConnectorID:     connectorID,
		CheckType:       "config_drift",
		Severity:        highestDriftSeverity(results),
		Title:           "Configuration has drifted from golden baseline",
		Description:     driftDescription(results),
		RemediationLink: "/connectors/" + connectorID + "/data",
	}
	return c.upsert(ctx, finding)
}

// highestDriftSeverity returns the most severe of the diff results, using
// the same info < warning < critical ranking as maybeNotify.
func highestDriftSeverity(results []sync.DiffResult) string {
	severity := "info"
	for _, r := range results {
		if severityRank[r.Severity] > severityRank[severity] {
			severity = r.Severity
		}
	}
	return severity
}

func driftDescription(results []sync.DiffResult) string {
	limit := len(results)
	if limit > configDriftEntryLimit {
		limit = configDriftEntryLimit
	}
	summaries := make([]string, limit)
	for i := range summaries {
		summaries[i] = results[i].Summary
	}
	description := fmt.Sprintf("%d section(s) deviate from the golden baseline: %s", len(results), strings.Join(summaries, "; "))
	if remaining := len(results) - limit; remaining > 0 {
		description += fmt.Sprintf(", and %d more", remaining)
	}
	return description + "."
}
