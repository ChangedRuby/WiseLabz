package sync

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RunDueSyncs syncs every connector whose next_run_at has passed. Only the
// sweep-level error (listing due connectors) is returned for job health
// (#384) purposes — a single connector's sync failing is its own concern
// (surfaced via its connector status/quality findings), not the sync job's.
func (e *Engine) RunDueSyncs(ctx context.Context, logger *slog.Logger) error {
	now := time.Now().UTC().Format(time.RFC3339)
	due, err := e.store.ListDueConnectors(ctx, now, e.dueBatchSize)
	if err != nil {
		return fmt.Errorf("list due connectors: %w", err)
	}
	sem := make(chan struct{}, e.maxConcurrency)
	var wg sync.WaitGroup
loop:
	for _, c := range due {
		if ctx.Err() != nil {
			break
		}
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			break loop
		}
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			defer func() { <-sem }()
			if _, err := e.runSyncFields(ctx, id, uuid.New().String(), nil, true); err != nil && !errors.Is(err, ErrAlreadyRunning) {
				logger.Error("scheduled sync failed", "connector", id, "error", err)
			}
		}(c.ID)
	}
	wg.Wait()
	return nil
}
