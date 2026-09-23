package sync

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// dueConnectorsLimit bounds how many due connectors are synced per tick.
// ponytail: fine for a self-hosted ops tool's connector count; paginate if that changes.
const dueConnectorsLimit = 50

// RunDueSyncs syncs every connector whose next_run_at has passed. Only the
// sweep-level error (listing due connectors) is returned for job health
// (#384) purposes — a single connector's sync failing is its own concern
// (surfaced via its connector status/quality findings), not the sync job's.
func (e *Engine) RunDueSyncs(ctx context.Context, logger *slog.Logger) error {
	now := time.Now().UTC().Format(time.RFC3339)
	due, err := e.store.ListDueConnectors(ctx, now, dueConnectorsLimit)
	if err != nil {
		return fmt.Errorf("list due connectors: %w", err)
	}
	for _, c := range due {
		if ctx.Err() != nil {
			return nil
		}
		if _, err := e.runSyncFields(ctx, c.ID, uuid.New().String(), nil, true); err != nil && !errors.Is(err, ErrAlreadyRunning) {
			logger.Error("scheduled sync failed", "connector", c.ID, "error", err)
		}
	}
	return nil
}
