package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

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
