// Package sync provides the sync engine for fetching connector data,
// diffing snapshots, and creating changes/alerts.
package sync

import (
	"context"
	"crypto/sha1" //nolint:gosec // non-cryptographic: stable pattern fingerprint, not a security boundary
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/logsafe"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// repeatDriftWindow is how far back to look for a prior change with the
// same pattern ID when deciding whether a drift is "recurring" (a likely
// misconfiguration loop) rather than a one-off.
// ponytail: fixed 1h window, make configurable if a real deployment needs it.
const repeatDriftWindow = time.Hour

// retrySchedule is the backoff delay before each retry attempt (index 0 =
// delay before the 2nd consecutive attempt, etc), mirroring
// notifications.retrySchedule.
// ponytail: fixed schedule, not exponential-from-config; add jitter/config if a real deployment needs it.
var retrySchedule = []time.Duration{time.Minute, 5 * time.Minute, 15 * time.Minute, 30 * time.Minute, time.Hour}

// maxSyncConcurrency bounds how many connectors RunSyncAll fetches at once,
// matching notifications.maxConcurrentNotifications' fanout pattern.
const maxSyncConcurrency = 4

// syncTimeout bounds upstream work for scheduled and manual syncs alike.
const syncTimeout = 5 * time.Minute

// ErrAlreadyRunning means this connector already has an active sync.
var ErrAlreadyRunning = errors.New("connector sync already running")

// computeNextRun decides a connector's next scheduled run time and updated
// retry count after one sync attempt. Unlike notification delivery retries
// (which give up after a fixed number of attempts), a scheduled connector
// must keep trying forever: once retryCount exceeds the backoff schedule, it
// keeps retrying at the schedule's last (longest) step. The retry delay is
// also never allowed to exceed the connector's own scheduleSeconds cadence,
// so a short-cadence connector doesn't end up waiting an hour to recover.
//
// scheduleSeconds nil means manual-only: nextRunAt is always nil, but
// retryCount is still tracked for display purposes. retryCount is the
// connector's consecutive-failure count *before* this attempt.
func computeNextRun(scheduleSeconds *int, retryCount int, success bool, now time.Time) (nextRunAt *time.Time, newRetryCount int) {
	if success {
		newRetryCount = 0
	} else {
		newRetryCount = retryCount + 1
	}

	if scheduleSeconds == nil {
		return nil, newRetryCount
	}
	cadence := time.Duration(*scheduleSeconds) * time.Second

	delay := cadence
	if !success {
		idx := newRetryCount - 1
		if idx >= len(retrySchedule) {
			idx = len(retrySchedule) - 1
		}
		delay = retrySchedule[idx]
		if cadence < delay {
			delay = cadence
		}
	}

	next := now.Add(delay)
	return &next, newRetryCount
}

// AlertNotifier dispatches notifications for a newly created alert.
type AlertNotifier interface {
	NotifyAlertCreated(ctx context.Context, alertID, title, message string)
}

// QualityChecker evaluates documentation quality after a sync attempt.
type QualityChecker interface {
	RunForConnector(ctx context.Context, connectorID string) error
}

// DocRegenerator re-renders a connector's existing docs from its latest
// snapshot after a sync attempt, recording a new "sync"-triggered doc
// version for any doc whose content changed.
type DocRegenerator interface {
	RegenerateForConnector(ctx context.Context, connectorID string) error
}

// Engine runs sync jobs against connectors.
type Engine struct {
	inFlight       sync.Map // connector ID -> active run; entries are removed on every exit
	store          *store.Store
	hub            *ws.Hub
	notifier       AlertNotifier
	qualityChecker QualityChecker
	docRegenerator DocRegenerator
	// encKey is the base64-encoded AES-256 key (config.Encryption.Key) used
	// to decrypt/re-encrypt secret-bearing connector config fields via
	// store.ParseConnectorConfig/MarshalConnectorConfig.
	encKey string
}

// NewEngine creates a new sync engine.
func NewEngine(s *store.Store, h *ws.Hub, notifier AlertNotifier, qualityChecker QualityChecker, encKey string) *Engine {
	return &Engine{store: s, hub: h, notifier: notifier, qualityChecker: qualityChecker, encKey: encKey}
}

// SetDocRegenerator wires a DocRegenerator into the engine after
// construction. Kept as a setter (rather than a NewEngine parameter) so
// existing call sites don't need to change; a nil regenerator (the default)
// simply skips sync-triggered doc regeneration.
func (e *Engine) SetDocRegenerator(dr DocRegenerator) {
	e.docRegenerator = dr
}

// RunResult holds the outcome of a sync run.
type RunResult struct {
	ConnectorID  string `json:"connectorId"`
	SnapshotID   string `json:"snapshotId"`
	ChangesCount int    `json:"changesCount"`
	AlertsCount  int    `json:"alertsCount"`
	Status       string `json:"status"` // "success", "error"
	Error        string `json:"error,omitempty"`
	Duration     string `json:"duration"`
}

// RunSync runs a full sync for a single connector (all fields).
// Flow: Fetch -> Save Snapshot -> Diff -> Create Changes -> Create Alerts
func (e *Engine) RunSync(ctx context.Context, connectorID string, jobID string) (*RunResult, error) {
	return e.RunSyncFields(ctx, connectorID, jobID, nil)
}

// RefreshCredentials refreshes a connector's credentials via its
// connector.CredentialRefresher implementation and persists the result
// (config_data + credential_expires_at). Shared by the sync-time
// expired-credential path in RunSyncFields and the bulk-reauth API handler.
func (e *Engine) RefreshCredentials(ctx context.Context, connectorID string) error {
	rec, err := e.store.GetConnector(ctx, connectorID)
	if err != nil {
		return fmt.Errorf("get connector: %w", err)
	}
	cfg, err := store.ParseConnectorConfig(rec.Type, rec.ConfigData, e.encKey)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	cfg["url"] = rec.URL
	cfg["verify_tls"] = rec.VerifyTLS

	conn, err := connector.Get(rec.Type, cfg)
	if err != nil {
		return fmt.Errorf("get connector impl: %w", err)
	}
	refresher, ok := conn.(connector.CredentialRefresher)
	if !ok {
		return fmt.Errorf("connector does not support credential refresh")
	}

	newCfg, expiresAt, err := refresher.RefreshCredentials(ctx, cfg)
	if err != nil {
		return fmt.Errorf("credential refresh failed: %w", err)
	}

	// Persist the refreshed credentials on their own — url/verify_tls
	// already live in their own columns and "fields" is a per-request
	// hint, not connector config; strip them regardless of whether the
	// refresher's newConfig (often built by copying its input, which
	// already carries these) included them.
	toStore := make(map[string]any, len(newCfg))
	for k, v := range newCfg {
		toStore[k] = v
	}
	delete(toStore, "url")
	delete(toStore, "verify_tls")
	delete(toStore, "fields")
	configData, err := store.MarshalConnectorConfig(rec.Type, toStore, e.encKey)
	if err != nil {
		return fmt.Errorf("marshal refreshed config: %w", err)
	}
	return e.store.UpdateConnector(ctx, connectorID, map[string]any{
		"config_data":           configData,
		"credential_expires_at": expiresAt.UTC().Format(time.RFC3339),
	})
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

	if !rec.Enabled {
		result.Status = "skipped"
		result.Duration = time.Since(start).String()
		finish("skipped", nil)
		return result, nil
	}

	// Parse config and inject top-level fields stored in separate columns
	// so connector factories see url + verify_tls alongside their custom fields.
	cfg, err := store.ParseConnectorConfig(rec.Type, rec.ConfigData, e.encKey)
	if err != nil {
		slog.Error("sync parse config failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		if e.hub != nil {
			e.hub.Broadcast(ws.EventSyncProgress, map[string]any{
				"serviceId": connectorID,
				"jobId":     jobID,
				"phase":     "error",
				"percent":   0,
				"message":   err.Error(),
			})
		}
		finish("error", err)
		return markError(result, start, fmt.Errorf("parse config: %w", err))
	}
	cfg["url"] = rec.URL
	cfg["verify_tls"] = rec.VerifyTLS
	if len(fields) > 0 {
		cfg["fields"] = fields
	}

	// Get connector implementation
	conn, err := connector.Get(rec.Type, cfg)
	if err != nil {
		slog.Error("sync get connector impl failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		if e.hub != nil {
			e.hub.Broadcast(ws.EventSyncProgress, map[string]any{
				"serviceId": connectorID,
				"jobId":     jobID,
				"phase":     "error",
				"percent":   0,
				"message":   err.Error(),
			})
		}
		finish("error", err)
		return markError(result, start, fmt.Errorf("get connector impl: %w", err))
	}

	// Credential expiry check: refuse (or refresh, if the connector supports
	// it) rather than let an expired credential fail Fetch with a confusing
	// upstream error.
	if rec.IsCredentialExpired(time.Now()) {
		if _, ok := conn.(connector.CredentialRefresher); !ok {
			authErr := connector.NewAuthError(fmt.Errorf("credentials expired at %s", rec.CredentialExpiresAt))
			_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{
				"status":         "offline",
				"status_message": authErr.Error(),
			})
			finish("error", authErr)
			return markError(result, start, authErr)
		}
		if refreshErr := e.RefreshCredentials(ctx, connectorID); refreshErr != nil {
			authErr := connector.NewAuthError(fmt.Errorf("credential refresh failed: %w", refreshErr))
			_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{
				"status":         "offline",
				"status_message": authErr.Error(),
			})
			finish("error", authErr)
			return markError(result, start, authErr)
		}
		// Reload: RefreshCredentials persisted the new config_data/expiry.
		rec, err = e.store.GetConnector(ctx, connectorID)
		if err != nil {
			finish("error", err)
			return markError(result, start, fmt.Errorf("get connector after refresh: %w", err))
		}
		cfg, err = store.ParseConnectorConfig(rec.Type, rec.ConfigData, e.encKey)
		if err != nil {
			finish("error", err)
			return markError(result, start, fmt.Errorf("parse config after refresh: %w", err))
		}
		cfg["url"] = rec.URL
		cfg["verify_tls"] = rec.VerifyTLS
		if len(fields) > 0 {
			cfg["fields"] = fields
		}
		conn, err = connector.Get(rec.Type, cfg)
		if err != nil {
			finish("error", err)
			return markError(result, start, fmt.Errorf("get connector impl after refresh: %w", err))
		}
	}

	// Update status to fetching
	_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{
		"status":         "online",
		"status_message": "Syncing...",
	})

	broadcast("fetching", 28)

	// Fetch data
	sn, err := conn.Fetch(ctx, cfg)
	if err != nil {
		_ = e.store.UpdateConnector(ctx, connectorID, map[string]any{
			"status":         "degraded",
			"status_message": fmt.Sprintf("Fetch failed: %v", err),
		})
		slog.Error("sync fetch failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		if e.hub != nil {
			e.hub.Broadcast(ws.EventSyncProgress, map[string]any{
				"serviceId": connectorID,
				"jobId":     jobID,
				"phase":     "error",
				"percent":   0,
				"message":   err.Error(),
			})
		}
		finish("error", err)
		return markError(result, start, fmt.Errorf("fetch: %w", err))
	}

	// Post-fetch transform/enrich pipeline: normalize field names, filter
	// PII, merge/derive data — registered per connector category.
	if err := runTransformers(ctx, rec.Category, sn); err != nil {
		slog.Error("sync transform failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
	}

	broadcast("diffing", 60)

	// Get previous snapshot for diff
	prevSn, prevErr := e.store.GetLatestSnapshot(ctx, connectorID)
	if prevErr != nil && !errors.Is(prevErr, store.ErrNotFound) {
		slog.Error("sync: previous snapshot unreadable, skipping diff", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(prevErr.Error()))
		wrapped := fmt.Errorf("get latest snapshot: %w", prevErr)
		if e.hub != nil {
			e.hub.Broadcast(ws.EventSyncProgress, map[string]any{
				"serviceId": connectorID,
				"jobId":     jobID,
				"phase":     "error",
				"percent":   0,
				"message":   wrapped.Error(),
			})
		}
		finish("error", wrapped)
		return markError(result, start, wrapped)
	}

	// Save new snapshot
	snData, _ := json.Marshal(sn)
	snRec := &store.SnapshotRecord{
		ConnectorID: connectorID,
		Data:        string(snData),
		FetchedAt:   sn.FetchedAt.Format(time.RFC3339),
	}
	if err := e.store.CreateSnapshot(ctx, snRec); err != nil {
		slog.Error("sync save snapshot failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
		if e.hub != nil {
			e.hub.Broadcast(ws.EventSyncProgress, map[string]any{
				"serviceId": connectorID,
				"jobId":     jobID,
				"phase":     "error",
				"percent":   0,
				"message":   err.Error(),
			})
		}
		finish("error", err)
		return markError(result, start, fmt.Errorf("save snapshot: %w", err))
	}
	result.SnapshotID = snRec.ID

	broadcast("generating", 85)

	// Under an active maintenance window, the snapshot above is still saved
	// (so sync history isn't lost) but drift alerting/change-record creation
	// is suppressed — this is the shared path for both manual "sync now" and
	// scheduled runs, so this is the one place that check needs to live.
	maintenance, err := e.store.GetActiveMaintenanceWindow(ctx, connectorID)
	if err != nil {
		slog.Error("get active maintenance window failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
	}

	// Diff against previous snapshot
	if prevErr == nil && maintenance == nil {
		var prevSnap connector.ServiceSnapshot
		if err := json.Unmarshal([]byte(prevSn.Data), &prevSnap); err != nil {
			slog.Error("sync: previous snapshot unparseable, skipping diff", "connector", logsafe.Sanitize(connectorID), "snapshot", prevSn.ID, "error", logsafe.Sanitize(err.Error()))
		} else {
			diffResults := Compare(&prevSnap, sn)
			for _, d := range diffResults {
				diffJSON, _ := json.Marshal(d.Patches)
				relatedJSON, _ := json.Marshal(d.RelatedServiceIDs)
				patternID := changePatternID(connectorID, d.Type, d.Summary)

				severity := d.Severity
				summary := d.Summary
				if repeatCount, err := e.store.CountRecentChangesByPattern(ctx, connectorID, patternID,
					time.Now().Add(-repeatDriftWindow).UTC().Format(time.RFC3339), ""); err != nil {
					slog.Error("count recent changes by pattern failed", "error", logsafe.Sanitize(err.Error()))
				} else if repeatCount > 0 {
					// Same drift recurring within the window: likely a
					// misconfiguration loop rather than a one-off — flag it
					// as unusual by raising severity and noting the repeat.
					summary = fmt.Sprintf("%s (recurring — seen %d time(s) in the last hour)", d.Summary, repeatCount)
					switch severity {
					case "info":
						severity = "warning"
					case "warning":
						severity = "critical"
					}
				}

				change := &store.ChangeRecord{
					ServiceID:         connectorID,
					ChangeType:        d.Type,
					Severity:          severity,
					Summary:           summary,
					Diff:              string(diffJSON),
					AffectedDocIDs:    "[]",
					RelatedServiceIDs: string(relatedJSON),
					PatternID:         patternID,
				}
				if err := e.store.CreateChange(ctx, change); err != nil {
					slog.Error("failed to create change", "error", logsafe.Sanitize(err.Error()))
					continue
				}
				result.ChangesCount++

				if e.hub != nil {
					e.hub.Broadcast(ws.EventChangeDetected, map[string]any{
						"changeId":      change.ID,
						"serviceId":     connectorID,
						"changeType":    d.Type,
						"severity":      severity,
						"summary":       summary,
						"willTriggerAi": false,
					})
				}

				// Create alert for non-info changes
				if severity != "info" {
					alert := &store.AlertRecord{
						ChangeID:    change.ID,
						ServiceID:   connectorID,
						Severity:    severity,
						Title:       summary,
						Description: d.Detail,
					}
					if err := e.store.CreateAlert(ctx, alert); err != nil {
						slog.Error("failed to create alert", "error", logsafe.Sanitize(err.Error()))
						continue
					}
					result.AlertsCount++

					if e.notifier != nil {
						e.notifier.NotifyAlertCreated(ctx, alert.ID, alert.Title, alert.Description)
					}

					if e.hub != nil {
						e.hub.Broadcast(ws.EventAlertCreated, map[string]any{
							"alertId":   alert.ID,
							"serviceId": connectorID,
							"severity":  severity,
							"title":     summary,
						})
					}
				}
			}
		}
	}

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

// RunSyncAll runs sync for all enabled connectors.
func (e *Engine) RunSyncAll(ctx context.Context, jobID string) ([]RunResult, error) {
	connectors, err := e.store.ListAllConnectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("list connectors: %w", err)
	}

	// Fetch+write each connector concurrently, bounded by a semaphore, mirroring
	// notifications.Dispatcher's fanoutSem pattern. Results are written to
	// per-index slots so ordering stays deterministic (matching connectors order)
	// despite concurrent completion.
	sem := make(chan struct{}, maxSyncConcurrency)
	slots := make([]*RunResult, len(connectors))
	var wg sync.WaitGroup
	for i, c := range connectors {
		if !c.Enabled {
			continue
		}
		sem <- struct{}{}
		wg.Add(1)
		go func(i int, connectorID string) {
			defer wg.Done()
			defer func() { <-sem }()
			result, err := e.RunSync(ctx, connectorID, jobID)
			if err != nil {
				slog.Error("sync failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
				return
			}
			slots[i] = result
		}(i, c.ID)
	}
	wg.Wait()

	var results []RunResult
	for _, r := range slots {
		if r != nil {
			results = append(results, *r)
		}
	}
	return results, nil
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
