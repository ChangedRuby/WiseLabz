package doc

import (
	"crypto/sha1" //nolint:gosec // stable short node-ID hash, not a security use
	"fmt"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// renderMermaid emits a Mermaid flowchart (graph LR) with centerEntity as
// the hub node and one edge per link, labeled with the match reason.
func renderMermaid(centerEntity string, links []EntityLink) string {
	var b strings.Builder
	b.WriteString("graph LR\n")

	centerID := "n" + shortHash("center|"+centerEntity)
	fmt.Fprintf(&b, "    %s[%q]\n", centerID, centerEntity)

	for _, l := range links {
		nodeID := entityNodeID(l.ConnectorID, l.Entity)
		label := fmt.Sprintf("%s (%s)", l.Entity.Name, l.Entity.Kind)
		fmt.Fprintf(&b, "    %s[%q]\n", nodeID, label)
		fmt.Fprintf(&b, "    %s -->|%s| %s\n", centerID, l.Reason, nodeID)
	}

	return b.String()
}

// entityNodeID builds a stable, Mermaid-safe node identifier: the entity's
// ExternalID when present (unique within a connector), otherwise a short
// hash of connector+kind+name so unrelated entities never collide.
func entityNodeID(connectorID string, e connector.SnapshotEntity) string {
	if e.ExternalID != "" {
		return "n" + shortHash(connectorID+"|"+e.Kind+"|"+e.ExternalID)
	}
	return "n" + shortHash(connectorID+"|"+e.Kind+"|"+e.Name)
}

func shortHash(s string) string {
	sum := sha1.Sum([]byte(s)) //nolint:gosec
	return fmt.Sprintf("%x", sum[:6])
}

// labEntity is one entity observed by one connector, used to build the
// lab-wide topology diagram (see GenerateLabTopology).
type labEntity struct {
	ConnectorID   string
	ConnectorName string
	Entity        connector.SnapshotEntity
}

// labLink is a matched pair of entities from two different connectors.
type labLink struct {
	A, B   labEntity
	Reason string
}

// renderLabMermaid emits a Mermaid flowchart with one node per entity
// (including entities with no matches, so isolated services still show up)
// and one edge per matched pair.
func renderLabMermaid(entities []labEntity, links []labLink) string {
	var b strings.Builder
	b.WriteString("graph LR\n")

	seen := map[string]bool{}
	node := func(le labEntity) string {
		id := entityNodeID(le.ConnectorID, le.Entity)
		if !seen[id] {
			seen[id] = true
			label := fmt.Sprintf("%s (%s)", le.Entity.Name, le.Entity.Kind)
			fmt.Fprintf(&b, "    %s[%q]\n", id, label)
		}
		return id
	}

	for _, e := range entities {
		node(e)
	}
	for _, l := range links {
		aID, bID := node(l.A), node(l.B)
		fmt.Fprintf(&b, "    %s -->|%s| %s\n", aID, l.Reason, bID)
	}

	return b.String()
}
