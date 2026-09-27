package sync

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// EntityChange describes one changed entity field. Added and removed entities
// use the field "entity" and carry the complete entity as their value.
type EntityChange struct {
	Kind   string `json:"kind"`
	Key    string `json:"key"`
	Name   string `json:"name"`
	Field  string `json:"field"`
	Change string `json:"change"`
	Old    any    `json:"old,omitempty"`
	New    any    `json:"new,omitempty"`
}

// DependencyChange records an added or removed snapshot dependency.
type DependencyChange struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Ref    string `json:"ref"`
	Change string `json:"change"`
}

// SnapshotDiffSummary counts changed sections, distinct entities, and dependencies.
type SnapshotDiffSummary struct {
	SectionsAdded       int `json:"sectionsAdded"`
	SectionsRemoved     int `json:"sectionsRemoved"`
	SectionsModified    int `json:"sectionsModified"`
	EntitiesAdded       int `json:"entitiesAdded"`
	EntitiesRemoved     int `json:"entitiesRemoved"`
	EntitiesModified    int `json:"entitiesModified"`
	DependenciesAdded   int `json:"dependenciesAdded"`
	DependenciesRemoved int `json:"dependenciesRemoved"`
}

// SnapshotDiff is the complete comparison returned by the snapshot diff API.
type SnapshotDiff struct {
	Provenance   SnapshotDiffProvenance `json:"provenance"`
	Summary      SnapshotDiffSummary    `json:"summary"`
	Sections     []DiffResult           `json:"sections"`
	Entities     []EntityChange         `json:"entities"`
	Dependencies []DependencyChange     `json:"dependencies"`
}

func entityKey(e connector.SnapshotEntity) string {
	if e.ExternalID != "" {
		return "externalId:" + e.ExternalID
	}
	return "name:" + e.Name
}

func entityMap(entities []connector.SnapshotEntity) map[string]connector.SnapshotEntity {
	result := make(map[string]connector.SnapshotEntity, len(entities))
	for _, e := range entities {
		result[e.Kind+"\x00"+entityKey(e)] = e
	}
	return result
}

// CompareEntities returns deterministic field-level changes between entity sets.
func CompareEntities(prev, curr []connector.SnapshotEntity) []EntityChange {
	before, after := entityMap(prev), entityMap(curr)
	keys := make([]string, 0, len(before)+len(after))
	seen := make(map[string]bool)
	for k := range before {
		keys = append(keys, k)
		seen[k] = true
	}
	for k := range after {
		if !seen[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	changes := make([]EntityChange, 0)
	for _, k := range keys {
		old, hasOld := before[k]
		newer, hasNew := after[k]
		e := newer
		if !hasNew {
			e = old
		}
		base := EntityChange{Kind: e.Kind, Key: entityKey(e), Name: e.Name}
		if !hasOld {
			base.Field = "entity"
			base.Change = "added"
			base.New = newer
			changes = append(changes, base)
			continue
		}
		if !hasNew {
			base.Field = "entity"
			base.Change = "removed"
			base.Old = old
			changes = append(changes, base)
			continue
		}
		addField := func(field string, a, b any, oldExists, newExists bool) {
			if oldExists && newExists && reflect.DeepEqual(a, b) {
				return
			}
			c := base
			c.Field = field
			c.Change = "modified"
			if !oldExists {
				c.Change = "added"
			} else if !newExists {
				c.Change = "removed"
			}
			if oldExists {
				c.Old = a
			}
			if newExists {
				c.New = b
			}
			changes = append(changes, c)
		}
		addField("name", old.Name, newer.Name, true, true)
		addField("ip", old.IP, newer.IP, true, true)
		addField("hostname", old.Hostname, newer.Hostname, true, true)
		attributeKeys := make([]string, 0, len(old.Attributes)+len(newer.Attributes))
		attributeSeen := make(map[string]bool)
		for attr := range old.Attributes {
			attributeKeys = append(attributeKeys, attr)
			attributeSeen[attr] = true
		}
		for attr := range newer.Attributes {
			if !attributeSeen[attr] {
				attributeKeys = append(attributeKeys, attr)
			}
		}
		sort.Strings(attributeKeys)
		for _, attr := range attributeKeys {
			a, aok := old.Attributes[attr]
			b, bok := newer.Attributes[attr]
			addField("attributes."+attr, a, b, aok, bok)
		}
	}
	sort.Slice(changes, func(i, j int) bool {
		a, b := changes[i], changes[j]
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		if a.Key != b.Key {
			return a.Key < b.Key
		}
		return a.Field < b.Field
	})
	return changes
}

// CompareDependencies returns deterministic added and removed dependencies.
func CompareDependencies(prev, curr []connector.ServiceDependency) []DependencyChange {
	key := func(d connector.ServiceDependency) string { return d.Kind + "\x00" + d.Name + "\x00" + d.Ref }
	before := make(map[string]connector.ServiceDependency, len(prev))
	after := make(map[string]connector.ServiceDependency, len(curr))
	for _, d := range prev {
		before[key(d)] = d
	}
	for _, d := range curr {
		after[key(d)] = d
	}
	changes := make([]DependencyChange, 0)
	for k, d := range before {
		if _, ok := after[k]; !ok {
			changes = append(changes, DependencyChange{d.Kind, d.Name, d.Ref, "removed"})
		}
	}
	for k, d := range after {
		if _, ok := before[k]; !ok {
			changes = append(changes, DependencyChange{d.Kind, d.Name, d.Ref, "added"})
		}
	}
	sort.Slice(changes, func(i, j int) bool {
		a, b := changes[i], changes[j]
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		if a.Ref != b.Ref {
			return a.Ref < b.Ref
		}
		return a.Change < b.Change
	})
	return changes
}

// BuildSnapshotDiff compares sections, entities, and dependencies in two snapshots.
func BuildSnapshotDiff(prev, curr *connector.ServiceSnapshot) SnapshotDiff {
	d := SnapshotDiff{Sections: Compare(prev, curr), Entities: CompareEntities(prev.Entities, curr.Entities), Dependencies: CompareDependencies(prev.Dependencies, curr.Dependencies)}
	if d.Sections == nil {
		d.Sections = []DiffResult{}
	}
	if d.Entities == nil {
		d.Entities = []EntityChange{}
	}
	if d.Dependencies == nil {
		d.Dependencies = []DependencyChange{}
	}
	for _, section := range d.Sections {
		switch section.Type {
		case "added":
			d.Summary.SectionsAdded++
		case "removed":
			d.Summary.SectionsRemoved++
		case "modified":
			d.Summary.SectionsModified++
		}
	}
	entityCounts := make(map[string]map[string]bool)
	for _, e := range d.Entities {
		change := e.Change
		if e.Field != "entity" {
			change = "modified"
		}
		if entityCounts[change] == nil {
			entityCounts[change] = make(map[string]bool)
		}
		entityCounts[change][e.Kind+"\x00"+e.Key] = true
	}
	d.Summary.EntitiesAdded = len(entityCounts["added"])
	d.Summary.EntitiesRemoved = len(entityCounts["removed"])
	d.Summary.EntitiesModified = len(entityCounts["modified"])
	for _, dependency := range d.Dependencies {
		if dependency.Change == "added" {
			d.Summary.DependenciesAdded++
		} else {
			d.Summary.DependenciesRemoved++
		}
	}
	return d
}

// UnifiedHunks returns line-oriented unified hunks with three context lines.
func UnifiedHunks(old, newer string) string {
	if old == newer {
		return ""
	}
	dmp := diffmatchpatch.New()
	a, b, lines := dmp.DiffLinesToChars(old, newer)
	diffs := dmp.DiffCharsToLines(dmp.DiffMain(a, b, false), lines)
	type line struct {
		prefix byte
		text   string
	}
	var all []line
	for _, diff := range diffs {
		prefix := byte(' ')
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			prefix = '+'
		case diffmatchpatch.DiffDelete:
			prefix = '-'
		}
		parts := strings.SplitAfter(diff.Text, "\n")
		for _, part := range parts {
			if part != "" {
				all = append(all, line{prefix, part})
			}
		}
	}
	var out strings.Builder
	for i := 0; i < len(all); {
		for i < len(all) && all[i].prefix == ' ' {
			i++
		}
		if i == len(all) {
			break
		}
		start := i - 3
		if start < 0 {
			start = 0
		}
		end := i + 4
		if end > len(all) {
			end = len(all)
		}
		for j := i + 1; j < len(all); j++ {
			if all[j].prefix != ' ' && j <= end+3 {
				end = j + 4
				if end > len(all) {
					end = len(all)
				}
			}
		}
		oldStart, newStart := 1, 1
		for _, l := range all[:start] {
			if l.prefix != '+' {
				oldStart++
			}
			if l.prefix != '-' {
				newStart++
			}
		}
		oldCount, newCount := 0, 0
		for _, l := range all[start:end] {
			if l.prefix != '+' {
				oldCount++
			}
			if l.prefix != '-' {
				newCount++
			}
		}
		fmt.Fprintf(&out, "@@ -%d,%d +%d,%d @@\n", oldStart, oldCount, newStart, newCount)
		for _, l := range all[start:end] {
			out.WriteByte(l.prefix)
			out.WriteString(l.text)
			if !strings.HasSuffix(l.text, "\n") {
				out.WriteString("\n\\ No newline at end of file\n")
			}
		}
		i = end
	}
	return out.String()
}
