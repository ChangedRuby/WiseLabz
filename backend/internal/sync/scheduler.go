package sync

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// dueConnectorsLimit bounds how many due connectors are synced per tick.
// ponytail: fine for a self-hosted ops tool's connector count; paginate if that changes.
const dueConnectorsLimit = 50

// RunDueSyncs syncs every connector whose next_run_at has passed.
func (e *Engine) RunDueSyncs(ctx context.Context, logger *slog.Logger) {
	now := time.Now().UTC().Format(time.RFC3339)
	due, err := e.store.ListDueConnectors(ctx, now, dueConnectorsLimit)
	if err != nil {
		logger.Error("list due connectors", "error", err)
		return
	}
	for _, c := range due {
		if ctx.Err() != nil {
			return
		}
		if _, err := e.runSyncFields(ctx, c.ID, uuid.New().String(), nil, true); err != nil && !errors.Is(err, ErrAlreadyRunning) {
			logger.Error("scheduled sync failed", "connector", c.ID, "error", err)
		}
	}
}
