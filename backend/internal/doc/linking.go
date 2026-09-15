package doc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// EntityLink is a cross-connector relationship between one of the current
// service's entities and an entity observed by another connector, matched
// live at doc-render time (see matchEntities).
type EntityLink struct {
	Entity        connector.SnapshotEntity
	ConnectorID   string
	ConnectorName string
	Reason        string // "external ID", "IP address", or "hostname"
}

// matchEntities loads the latest snapshot of every other connector and
// matches its entities against entities, in precedence order (first hit
// wins per pair): matching ExternalID+Kind, then matching IP, then matching
// Hostname (case-insensitive). Matches across different connectors are all
// kept; only exact (ConnectorID, ExternalID) duplicates are dropped.
func matchEntities(ctx context.Context, s *store.Store, connectorID string, entities []connector.SnapshotEntity) ([]EntityLink, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	connectors, err := s.ListAllConnectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("list connectors: %w", err)
	}

	var links []EntityLink
	seen := map[string]bool{}
	for _, c := range connectors {
		if c.ID == connectorID {
			continue
		}
		sn, err := s.GetLatestSnapshot(ctx, c.ID)
		if err != nil {
			continue // no snapshot yet for this connector; soft-skip
		}
		var snap connector.ServiceSnapshot
		if err := json.Unmarshal([]byte(sn.Data), &snap); err != nil {
			continue
		}
		for _, other := range snap.Entities {
			for _, mine := range entities {
				reason := matchReason(mine, other)
				if reason == "" {
					continue
				}
				key := dedupKey(c.ID, other)
				if seen[key] {
					break
				}
				seen[key] = true
				links = append(links, EntityLink{
					Entity:        other,
					ConnectorID:   c.ID,
					ConnectorName: c.Name,
					Reason:        reason,
				})
				break
			}
		}
	}

	return links, nil
}

// matchReason returns the precedence-ordered reason two entities match, or
// "" if they don't match on any tier.
func matchReason(a, b connector.SnapshotEntity) string {
	if a.ExternalID != "" && b.ExternalID != "" && a.Kind == b.Kind && a.ExternalID == b.ExternalID {
		return "external ID"
	}
	if a.IP != "" && b.IP != "" && a.IP == b.IP {
		return "IP address"
	}
	if a.Hostname != "" && b.Hostname != "" && strings.EqualFold(a.Hostname, b.Hostname) {
		return "hostname"
	}
	return ""
}

func dedupKey(connectorID string, e connector.SnapshotEntity) string {
	if e.ExternalID != "" {
		return connectorID + "|" + e.Kind + "|" + e.ExternalID
	}
	return connectorID + "|" + e.Kind + "|" + e.Name + "|" + e.IP + "|" + e.Hostname
}
