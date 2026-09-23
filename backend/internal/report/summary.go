package report

import (
	"fmt"
	"strings"
)

// Summary builds the short Markdown counts line notifications.Dispatcher's
// NotifyReport sends to channels (Slack/Discord/webhook), alongside the
// report link. It intentionally does not reuse RenderMarkdown's full
// per-section detail — channel messages should stay skimmable.
func Summary(data ReportData) string {
	var b strings.Builder
	fmt.Fprintf(&b, "**%s** (%s – %s)", data.Definition.Name, data.PeriodStart.UTC().Format("2006-01-02"), data.PeriodEnd.UTC().Format("2006-01-02"))
	if data.Truncated {
		b.WriteString(" _(truncated)_")
	}
	b.WriteString("\n")

	s := data.Sections
	if s.Docs.Error != "" {
		fmt.Fprintf(&b, "- Docs changed: unavailable (%s)\n", s.Docs.Error)
	} else {
		fmt.Fprintf(&b, "- Docs changed: %d\n", s.Docs.Total)
	}

	if s.Compliance.Error != "" {
		fmt.Fprintf(&b, "- Compliance: unavailable (%s)\n", s.Compliance.Error)
	} else {
		fmt.Fprintf(&b, "- Compliance: %d open, %d new, %d resolved\n", s.Compliance.OpenCount, s.Compliance.DetectedCount, s.Compliance.ResolvedCount)
	}

	if s.Rotations.Error != "" {
		fmt.Fprintf(&b, "- Rotations: unavailable (%s)\n", s.Rotations.Error)
	} else {
		fmt.Fprintf(&b, "- Rotations: %d rotated, %d overdue, %d due soon\n", len(s.Rotations.Rotated), len(s.Rotations.Overdue), len(s.Rotations.DueSoon))
	}

	if s.Quality.Error != "" {
		fmt.Fprintf(&b, "- Quality: unavailable (%s)\n", s.Quality.Error)
	} else {
		fmt.Fprintf(&b, "- Quality: %d open, %d new, %d resolved\n", s.Quality.OpenCount, s.Quality.DetectedCount, s.Quality.ResolvedCount)
	}

	if s.Drift.Error != "" {
		fmt.Fprintf(&b, "- Drift: unavailable (%s)\n", s.Drift.Error)
	} else {
		fmt.Fprintf(&b, "- Drift: %d change(s)\n", s.Drift.Total)
	}

	if s.Jobs.Error != "" {
		fmt.Fprintf(&b, "- Job health: unavailable (%s)\n", s.Jobs.Error)
	} else {
		fmt.Fprintf(&b, "- Job health: %d failing\n", s.Jobs.FailingCount)
	}

	return b.String()
}
