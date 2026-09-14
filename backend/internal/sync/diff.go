package sync

import (
	"fmt"
	"sort"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// DiffResult describes a detected change between two snapshots.
type DiffResult struct {
	Type     string      `json:"type"`     // "added", "removed", "modified"
	Severity string      `json:"severity"` // "info", "warning", "critical"
	Summary  string      `json:"summary"`
	Detail   string      `json:"detail"`
	Patches  []DiffPatch `json:"patches"`
	// RelatedServiceIDs are the connector/service IDs (from
	// ServiceDependency.Ref) the changed snapshot depends on, so a drift can
	// be clustered with the services it may affect instead of standing alone.
	// ponytail: one hop (this snapshot's own declared dependencies), not a
	// full dependency-graph walk; extend to transitive deps if UI grouping
	// needs it.
	RelatedServiceIDs []string `json:"relatedServiceIds,omitempty"`
}

// DiffPatch is a single text diff hunk.
type DiffPatch struct {
	Section string `json:"section"`
	Old     string `json:"old,omitempty"`
	New     string `json:"new,omitempty"`
	Line    int    `json:"line,omitempty"`
}

// Compare compares two snapshots and returns detected differences.
func Compare(prev, curr *connector.ServiceSnapshot) []DiffResult {
	var results []DiffResult

	// Build section maps
	prevSections := make(map[string]connector.SnapshotSection)
	for _, s := range prev.Sections {
		prevSections[s.Title] = s
	}
	currSections := make(map[string]connector.SnapshotSection)
	for _, s := range curr.Sections {
		currSections[s.Title] = s
	}

	currRelated := relatedServiceIDs(curr)
	prevRelated := relatedServiceIDs(prev)

	currTitles := sortedKeys(currSections)
	prevTitles := sortedKeys(prevSections)

	// Check for added sections
	for _, title := range currTitles {
		cs := currSections[title]
		if _, ok := prevSections[title]; !ok {
			results = append(results, DiffResult{
				Type:     "added",
				Severity: "info",
				Summary:  fmt.Sprintf("New section added: %s", title),
				Detail:   fmt.Sprintf("Section %q was added with %d lines of content.", title, lineCount(cs.Content)),
				Patches: []DiffPatch{{
					Section: title,
					New:     cs.Content,
				}},
				RelatedServiceIDs: currRelated,
			})
		}
	}

	// Check for removed sections
	for _, title := range prevTitles {
		ps := prevSections[title]
		if _, ok := currSections[title]; !ok {
			results = append(results, DiffResult{
				Type:     "removed",
				Severity: "warning",
				Summary:  fmt.Sprintf("Section removed: %s", title),
				Detail:   fmt.Sprintf("Section %q was removed.", title),
				Patches: []DiffPatch{{
					Section: title,
					Old:     ps.Content,
				}},
				RelatedServiceIDs: prevRelated,
			})
		}
	}

	// Check for modified sections
	for _, title := range currTitles {
		cs := currSections[title]
		ps, ok := prevSections[title]
		if !ok {
			continue // already handled as "added"
		}
		if cs.Content != ps.Content {
			results = append(results, DiffResult{
				Type:     "modified",
				Severity: severityForChange(ps.Content, cs.Content),
				Summary:  fmt.Sprintf("Section modified: %s", title),
				Detail:   fmt.Sprintf("Section %q changed (%d → %d lines).", title, lineCount(ps.Content), lineCount(cs.Content)),
				Patches: []DiffPatch{{
					Section: title,
					Old:     ps.Content,
					New:     cs.Content,
				}},
				RelatedServiceIDs: currRelated,
			})
		}
	}

	return results
}

// relatedServiceIDs extracts the resolved connector/service IDs a snapshot
// declares as dependencies (ServiceDependency.Ref), deduplicated.
func relatedServiceIDs(snap *connector.ServiceSnapshot) []string {
	if snap == nil {
		return nil
	}
	seen := make(map[string]bool)
	var ids []string
	for _, d := range snap.Dependencies {
		if d.Ref == "" || seen[d.Ref] {
			continue
		}
		seen[d.Ref] = true
		ids = append(ids, d.Ref)
	}
	return ids
}

// sortedKeys returns a section map's titles in deterministic order, so
// change records for a single sync are always produced in the same order.
func sortedKeys(m map[string]connector.SnapshotSection) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func lineCount(s string) int {
	if s == "" {
		return 0
	}
	return len(strings.Split(s, "\n"))
}

func severityForChange(prev, curr string) string {
	// Simple heuristic: large diffs are more severe
	oldLines := lineCount(prev)
	newLines := lineCount(curr)
	diff := oldLines - newLines
	if diff < 0 {
		diff = -diff
	}
	if diff > 10 {
		return "critical"
	}
	if diff > 3 {
		return "warning"
	}
	return "info"
}
