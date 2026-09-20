package sync

import (
	"context"
	"crypto/sha1" //nolint:gosec // non-cryptographic: stable pattern fingerprint, not a security boundary
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/logsafe"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
	"github.com/google/uuid"
)

// repeatDriftWindow is how far back to look for a prior change with the
// same pattern ID when deciding whether a drift is "recurring" (a likely
// misconfiguration loop) rather than a one-off.
// ponytail: fixed 1h window, make configurable if a real deployment needs it.
const repeatDriftWindow = time.Hour

// syncTimeout bounds upstream work for scheduled and manual syncs alike.
const syncTimeout = 5 * time.Minute

// RunSync runs a full sync for a single connector (all fields).
// Flow: Fetch -> Save Snapshot -> Diff -> Create Changes -> Create Alerts
func (e *Engine) RunSync(ctx context.Context, connectorID string, jobID string) (*RunResult, error) {
	return e.RunSyncFields(ctx, connectorID, jobID, nil)
}

// RunSyncFields runs a sync for a single connector, optionally requesting
// only a subset of fields (e.g. []string{"vms","storage"}) so a caller like
// a dashboard quick-check doesn't force a full fetch. A connector type that
// doesn't support selective fetch ignores the hint and returns everything.
func (e *Engine) RunSyncFields(ctx context.Context, connectorID string, jobID string, fields []string) (*RunResult, error) {
	return e.runSyncFields(ctx, connectorID, jobID, fields, false)
}

func (e *Engine) runSyncFields(ctx context.Context, connectorID, jobID string, fields []string, scheduled bool) (*RunResult, error) {
	if _, loaded := e.inFlight.LoadOrStore(connectorID, struct{}{}); loaded {
		if e.hub != nil {
			e.hub.Broadcast(ws.EventSyncComplete, map[string]any{
				"serviceId": connectorID, "jobId": jobID, "error": ErrAlreadyRunning.Error(),
			})
		}
		return nil, ErrAlreadyRunning
	}
	defer e.inFlight.Delete(connectorID)
	ctx, cancel := context.WithTimeout(ctx, syncTimeout)
	defer cancel()
	if scheduled {
		now := time.Now().UTC()
		claimed, err := e.store.ClaimDueConnector(ctx, connectorID, now.Format(time.RFC3339), now.Add(syncTimeout+time.Minute).Format(time.RFC3339))
		if err != nil {
			return nil, err
		}
		if !claimed {
			return nil, nil
		}
	}
	start := time.Now()
	result := &RunResult{ConnectorID: connectorID}

	broadcast := func(phase string, percent int) {
		if e.hub != nil {
			e.hub.Broadcast(ws.EventSyncProgress, map[string]any{
				"serviceId": connectorID,
				"jobId":     jobID,
				"phase":     phase,
				"percent":   percent,
			})
		}
	}

	broadcast("queued", 0)

	// Get connector record
	rec, err := e.store.GetConnector(ctx, connectorID)
	if err != nil {
		slog.Error("sync get connector failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		if e.hub != nil {
			e.hub.Broadcast(ws.EventSyncProgress, map[string]any{
				"serviceId": connectorID,
				"jobId":     jobID,
				"phase":     "error",
				"percent":   0,
				"message":   err.Error(),
			})
		}
		return markError(result, start, fmt.Errorf("get connector: %w", err))
	}

	attempt := rec.RetryCount + 1

	// finish records this run's outcome as a sync_runs history row, broadcasts
	// EventSyncComplete so the UI progress indicator always resolves, and for
	// non-skipped runs persists the resulting retry/backoff + next_run_at
	// schedule state on the connector (see computeNextRun). Called from every
	// exit path below once rec has been loaded.
	finish := func(status string, runErr error) {
		e.finishSync(ctx, connectorID, jobID, rec, result, start, attempt, status, runErr)
	}

	if !rec.Enabled {
		result.Status = "skipped"
		result.Duration = time.Since(start).String()
		finish("skipped", nil)
		return result, nil
	}

	// Parse config and inject top-level fields stored in separate columns
	// so connector factories see url + verify_tls alongside their custom fields.
	cfg, err := e.syncConfig(rec, fields)
	if err != nil {
		slog.Error("sync parse config failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		e.broadcastSyncError(connectorID, jobID, err)
		finish("error", err)
		return markError(result, start, fmt.Errorf("parse config: %w", err))
	}

	// Get connector implementation
	conn, err := connector.Get(rec.Type, cfg)
	if err != nil {
		slog.Error("sync get connector impl failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		e.broadcastSyncError(connectorID, jobID, err)
		finish("error", err)
		return markError(result, start, fmt.Errorf("get connector impl: %w", err))
	}

	refreshedRec, refreshedCfg, refreshedConn, err := e.refreshSyncCredentials(ctx, connectorID, rec, cfg, conn, fields)
	if err != nil {
		finish("error", err)
		return markError(result, start, err)
	}
	rec, cfg, conn = refreshedRec, refreshedCfg, refreshedConn

	sn, err := e.fetchSyncSnapshot(ctx, connectorID, jobID, conn, cfg, broadcast)
	if err != nil {
		finish("error", err)
		return markError(result, start, fmt.Errorf("fetch: %w", err))
	}

	// Post-fetch transform/enrich pipeline: normalize field names, filter
	// PII, merge/derive data — registered per connector category.
	if err := runTransformers(ctx, rec.Category, sn); err != nil {
		slog.Error("sync transform failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
	}

	broadcast("diffing", 60)
	changes, alerts, snapshotID, err := e.persistSyncSnapshot(ctx, connectorID, sn)
	if err != nil {
		finish("error", err)
		return markError(result, start, err)
	}
	result.SnapshotID = snapshotID
	result.ChangesCount = len(changes)
	result.AlertsCount = len(alerts)
	e.broadcastSyncResults(connectorID, changes, alerts)
	e.notifySyncAlerts(ctx, alerts)

	// Update connector status
	now := time.Now().UTC().Format(time.RFC3339)
	_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{
		"status":         "online",
		"status_message": "Sync successful",
		"last_sync_at":   now,
	})

	result.Status = "success"
	result.Duration = time.Since(start).String()
	finish("success", nil)

	return result, nil
}

func (e *Engine) finishSync(ctx context.Context, connectorID, jobID string, rec *store.ConnectorRecord, result *RunResult, start time.Time, attempt int, status string, runErr error) {
	// Persist the outcome even when the fetch exhausted its deadline.
	finishCtx, finishCancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer finishCancel()
	if status != "skipped" && e.hub != nil {
		payload := map[string]any{
			"serviceId":       connectorID,
			"jobId":           jobID,
			"changesDetected": result.ChangesCount,
			"alertsRaised":    result.AlertsCount,
			"durationMs":      time.Since(start).Milliseconds(),
		}
		if runErr != nil {
			payload["error"] = runErr.Error()
		}
		e.hub.Broadcast(ws.EventSyncComplete, payload)
	}
	durationMs := int(time.Since(start).Milliseconds())
	errMsg := ""
	if runErr != nil {
		errMsg = runErr.Error()
	}
	runStatus := store.SyncRunStatusSuccess
	switch status {
	case "error":
		runStatus = store.SyncRunStatusError
	case "skipped":
		runStatus = store.SyncRunStatusSkipped
	}
	if err := e.store.CreateSyncRun(finishCtx, &store.SyncRunRecord{
		ConnectorID:  connectorID,
		StartedAt:    start.UTC().Format(time.RFC3339Nano),
		FinishedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		DurationMs:   &durationMs,
		Status:       runStatus,
		Error:        errMsg,
		Attempt:      attempt,
		ChangesCount: result.ChangesCount,
		AlertsCount:  result.AlertsCount,
	}); err != nil {
		slog.Error("record sync run failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
	}

	// Skipped (disabled connector) runs aren't real attempts: leave the
	// existing retry/schedule state untouched. ListDueConnectors also
	// filters on enabled=1, so a disabled connector's schedule is moot
	// until it's re-enabled anyway.
	if status == "skipped" {
		return
	}

	nextRunAt, newRetryCount := computeNextRun(rec.ScheduleSeconds, rec.RetryCount, status == "success", time.Now())
	updates := map[string]any{
		"last_sync_duration_ms": durationMs,
		"last_sync_error":       errMsg,
		"retry_count":           newRetryCount,
	}
	if nextRunAt != nil {
		updates["next_run_at"] = nextRunAt.UTC().Format(time.RFC3339)
	} else {
		updates["next_run_at"] = nil
	}
	if err := e.store.UpdateConnector(finishCtx, connectorID, updates); err != nil {
		slog.Error("update connector schedule failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
	}
	if e.qualityChecker != nil {
		if err := e.qualityChecker.RunForConnector(ctx, connectorID); err != nil {
			slog.Error("quality check failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		}
	}
	if e.docRegenerator != nil {
		if err := e.docRegenerator.RegenerateForConnector(ctx, connectorID); err != nil {
			slog.Error("doc regeneration failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		}
	}
}

func (e *Engine) syncConfig(rec *store.ConnectorRecord, fields []string) (map[string]any, error) {
	cfg, err := store.ParseConnectorConfig(rec.Type, rec.ConfigData, e.encKey)
	if err != nil {
		return nil, err
	}
	cfg["url"], cfg["verify_tls"] = rec.URL, rec.VerifyTLS
	if len(fields) > 0 {
		cfg["fields"] = fields
	}
	return cfg, nil
}

func (e *Engine) refreshSyncCredentials(ctx context.Context, connectorID string, rec *store.ConnectorRecord, cfg map[string]any, conn connector.Connector, fields []string) (*store.ConnectorRecord, map[string]any, connector.Connector, error) {
	if !rec.IsCredentialExpired(time.Now()) {
		return rec, cfg, conn, nil
	}
	if _, ok := conn.(connector.CredentialRefresher); !ok {
		err := connector.NewAuthError(fmt.Errorf("credentials expired at %s", rec.CredentialExpiresAt))
		_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{"status": "offline", "status_message": err.Error()})
		return nil, nil, nil, err
	}
	if err := e.RefreshCredentials(ctx, connectorID); err != nil {
		authErr := connector.NewAuthError(fmt.Errorf("credential refresh failed: %w", err))
		_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{"status": "offline", "status_message": authErr.Error()})
		return nil, nil, nil, authErr
	}
	var err error
	if rec, err = e.store.GetConnector(ctx, connectorID); err != nil {
		return nil, nil, nil, fmt.Errorf("get connector after refresh: %w", err)
	}
	if cfg, err = e.syncConfig(rec, fields); err != nil {
		return nil, nil, nil, fmt.Errorf("parse config after refresh: %w", err)
	}
	if conn, err = connector.Get(rec.Type, cfg); err != nil {
		return nil, nil, nil, fmt.Errorf("get connector impl after refresh: %w", err)
	}
	return rec, cfg, conn, nil
}

func (e *Engine) fetchSyncSnapshot(ctx context.Context, connectorID, jobID string, conn connector.Connector, cfg map[string]any, broadcast func(string, int)) (*connector.ServiceSnapshot, error) {
	_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{"status": "online", "status_message": "Syncing..."})
	broadcast("fetching", 28)
	sn, err := conn.Fetch(ctx, cfg)
	if err == nil {
		return sn, nil
	}
	_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{"status": "degraded", "status_message": fmt.Sprintf("Fetch failed: %v", err)})
	slog.Error("sync fetch failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
	e.broadcastSyncError(connectorID, jobID, err)
	return nil, err
}

func (e *Engine) persistSyncSnapshot(ctx context.Context, connectorID string, sn *connector.ServiceSnapshot) ([]*store.ChangeRecord, []*store.AlertRecord, string, error) {
	prev, prevErr := e.store.GetLatestSnapshot(ctx, connectorID)
	if prevErr != nil && !errors.Is(prevErr, store.ErrNotFound) {
		return nil, nil, "", fmt.Errorf("get latest snapshot: %w", prevErr)
	}
	data, _ := json.Marshal(sn)
	snapshot := &store.SnapshotRecord{ConnectorID: connectorID, Data: string(data), FetchedAt: sn.FetchedAt.Format(time.RFC3339)}
	maintenance, err := e.store.GetActiveMaintenanceWindow(ctx, connectorID)
	if err != nil {
		slog.Error("get active maintenance window failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
	}
	var changes []*store.ChangeRecord
	var alerts []*store.AlertRecord
	err = e.store.WithinTransaction(ctx, func(tx *store.Store) error {
		if err := tx.CreateSnapshot(ctx, snapshot); err != nil {
			return err
		}
		if prevErr == nil && maintenance == nil {
			var old connector.ServiceSnapshot
			if err := json.Unmarshal([]byte(prev.Data), &old); err != nil {
				slog.Error("sync: previous snapshot unparseable, skipping diff", "connector", logsafe.Sanitize(connectorID), "snapshot", prev.ID, "error", logsafe.Sanitize(err.Error()))
			} else if err := e.createSyncChanges(ctx, tx, connectorID, &old, sn, &changes, &alerts); err != nil {
				return err
			}
		}
		if err := tx.CreateChanges(ctx, changes); err != nil {
			return err
		}
		return tx.CreateAlerts(ctx, alerts)
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("save sync results: %w", err)
	}
	return changes, alerts, snapshot.ID, nil
}

func (e *Engine) createSyncChanges(ctx context.Context, tx *store.Store, connectorID string, old, sn *connector.ServiceSnapshot, changes *[]*store.ChangeRecord, alerts *[]*store.AlertRecord) error {
	counts, err := tx.CountRecentChangePatterns(ctx, connectorID, time.Now().Add(-repeatDriftWindow).UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	for _, d := range Compare(old, sn) {
		diff, _ := json.Marshal(d.Patches)
		related, _ := json.Marshal(d.RelatedServiceIDs)
		pattern := changePatternID(connectorID, d.Type, d.Summary)
		severity, summary := d.Severity, d.Summary
		if n := counts[pattern]; n > 0 {
			summary = fmt.Sprintf("%s (recurring — seen %d time(s) in the last hour)", d.Summary, n)
			switch severity {
			case "info":
				severity = "warning"
			case "warning":
				severity = "critical"
			}
		}
		change := &store.ChangeRecord{ID: uuid.NewString(), ServiceID: connectorID, ChangeType: d.Type, Severity: severity, Summary: summary, Diff: string(diff), AffectedDocIDs: "[]", RelatedServiceIDs: string(related), PatternID: pattern}
		*changes = append(*changes, change)
		counts[pattern]++
		if severity != "info" {
			*alerts = append(*alerts, &store.AlertRecord{ChangeID: change.ID, ServiceID: connectorID, Severity: severity, Title: summary, Description: d.Detail})
		}
	}
	return nil
}

func (e *Engine) broadcastSyncError(connectorID, jobID string, err error) {
	if e.hub != nil {
		e.hub.Broadcast(ws.EventSyncProgress, map[string]any{"serviceId": connectorID, "jobId": jobID, "phase": "error", "percent": 0, "message": err.Error()})
	}
}

func (e *Engine) broadcastSyncResults(connectorID string, changes []*store.ChangeRecord, alerts []*store.AlertRecord) {
	if e.hub == nil {
		return
	}
	for _, change := range changes {
		e.hub.Broadcast(ws.EventChangeDetected, map[string]any{"changeId": change.ID, "serviceId": connectorID, "changeType": change.ChangeType, "severity": change.Severity, "summary": change.Summary, "willTriggerAi": false})
	}
	for _, alert := range alerts {
		e.hub.Broadcast(ws.EventAlertCreated, map[string]any{"alertId": alert.ID, "serviceId": connectorID, "severity": alert.Severity, "title": alert.Title})
	}
}

func (e *Engine) notifySyncAlerts(ctx context.Context, alerts []*store.AlertRecord) {
	if e.notifier == nil || len(alerts) == 0 {
		return
	}
	batch := make([]store.AlertRecord, len(alerts))
	for i, alert := range alerts {
		batch[i] = *alert
	}
	e.notifier.NotifyAlertsCreated(ctx, batch)
}

// changePatternID derives a stable identifier for a drift "shape" — same
// service, same change type, same summary — so recurrences of the same
// drift can be recognized across sync runs.
func changePatternID(connectorID, changeType, summary string) string {
	h := sha1.Sum([]byte(connectorID + "|" + changeType + "|" + summary)) //nolint:gosec // fingerprint, not a security use
	return hex.EncodeToString(h[:])
}

func markError(r *RunResult, start time.Time, err error) (*RunResult, error) {
	r.Status = "error"
	r.Error = err.Error()
	r.Duration = time.Since(start).String()
	return r, err
}
