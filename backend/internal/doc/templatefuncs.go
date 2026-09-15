package doc

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// TemplateFuncs returns a FuncMap of built-in template functions for doc rendering.
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"dateFormat":      dateFormat,
		"truncate":        truncate,
		"toJSON":          toJSON,
		"filterByTitle":   filterByTitle,
		"join":            join,
		"relatedEntities": relatedEntities,
	}
}

// dateFormat formats a time value or string using a Go layout string.
// If the input is a string, it's parsed as RFC3339 first.
func dateFormat(layout string, t any) (string, error) {
	switch v := t.(type) {
	case time.Time:
		return v.Format(layout), nil
	case string:
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return "", err
		}
		return parsed.Format(layout), nil
	default:
		return "", nil
	}
}

// truncate truncates a string to maxChars, appending "..." if truncated.
func truncate(maxChars int, s string) string {
	if len(s) <= maxChars {
		return s
	}
	return s[:maxChars] + "..."
}

// toJSON marshals a value to a JSON string for embedding structured data.
func toJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// filterByTitle returns a subset of sections matching the given title.
func filterByTitle(title string, sections []connector.SnapshotSection) []connector.SnapshotSection {
	var result []connector.SnapshotSection
	for _, sec := range sections {
		if sec.Title == title {
			result = append(result, sec)
		}
	}
	return result
}

// join joins a slice of strings with a separator (stdlib wrapper for template piping).
func join(sep string, strs []string) string {
	return strings.Join(strs, sep)
}

// relatedEntities renders matched cross-connector entity links as a markdown
// bullet list.
func relatedEntities(links []EntityLink) string {
	if len(links) == 0 {
		return "_No related entities found._\n"
	}
	var b strings.Builder
	for _, l := range links {
		fmt.Fprintf(&b, "- **%s** (%s) via %s — matched on %s\n", l.Entity.Name, l.Entity.Kind, l.ConnectorName, l.Reason)
	}
	return b.String()
}
