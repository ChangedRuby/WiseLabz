package report

import (
	"strings"
	"testing"
	"time"
)

// sampleData returns a fixed ReportData used by the golden tests below.
// Times are fixed (not time.Now()) so the rendered output is deterministic.
func sampleData() ReportData {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 9, 8, 0, 0, 0, 0, time.UTC)
	changedAt := time.Date(2026, 9, 3, 12, 30, 0, 0, time.UTC)
	rotatedAt := time.Date(2026, 9, 2, 9, 0, 0, 0, time.UTC)

	return ReportData{
		Definition:  DefinitionSummary{Slug: "weekly", Name: "Weekly Ops Report"},
		PeriodStart: start,
		PeriodEnd:   end,
		Truncated:   false,
		Sections: Sections{
			Docs: DocsSection{
				Total: 2,
				Items: []DocChangeEntry{
					{DocID: "d1", Title: "<Runbook> & Notes", ConnectorID: "c1", ConnectorName: "Proxmox", Version: 3, ChangedAt: changedAt},
					{DocID: "d2", Title: "Onboarding Guide", ConnectorID: "c2", ConnectorName: "Docker", Version: 1, ChangedAt: changedAt},
				},
			},
			Compliance: ComplianceSection{
				OpenCount:     1,
				DetectedCount: 1,
				ResolvedCount: 0,
				Items: []FindingSummary{
					{ID: "f1", ConnectorID: "c1", ConnectorName: "Proxmox", CheckType: "compliance", Severity: "warning", Title: "Missing owner tag", Status: "open", DetectedAt: changedAt},
				},
			},
			Rotations: RotationsSection{
				Rotated: []RotationEntry{{ConnectorID: "c1", ConnectorName: "Proxmox", RotatedAt: &rotatedAt}},
				Overdue: []RotationEntry{{ConnectorID: "c3", ConnectorName: "Netbox", Overdue: true}},
				DueSoon: nil,
			},
			Quality: QualitySection{
				OpenCount:     0,
				DetectedCount: 0,
				ResolvedCount: 1,
				Items:         nil,
			},
			Drift: DriftSection{
				Total: 4,
				ByConnector: []ConnectorDrift{
					{ConnectorID: "c1", ConnectorName: "Proxmox", Counts: SeverityCounts{Critical: 1, Warning: 2, Info: 1}},
				},
				TopChanges: []ChangeEntry{
					{ID: "ch1", ConnectorID: "c1", ConnectorName: "Proxmox", Severity: "critical", Summary: "VM <web-01> memory doubled", DetectedAt: changedAt},
				},
			},
			Jobs: JobsSection{
				FailingCount: 1,
				Items: []JobHealthEntry{
					{Name: "sync", Status: "failing", LastError: "timeout"},
					{Name: "backup", Status: "ok"},
				},
			},
		},
	}
}

func TestRenderMarkdown_Golden(t *testing.T) {
	got, err := RenderMarkdown(sampleData())
	if err != nil {
		t.Fatalf("RenderMarkdown: %v", err)
	}

	const want = `# Weekly Ops Report

_2026-09-01 00:00 UTC – 2026-09-08 00:00 UTC_

## Documentation changed
2 doc version(s) changed.
- **<Runbook> & Notes** (Proxmox) — v3, 2026-09-03 12:30 UTC
- **Onboarding Guide** (Docker) — v1, 2026-09-03 12:30 UTC

## Compliance
1 new, 0 resolved, 1 still open.
- [warning] **Missing owner tag** (Proxmox) — open

## Secret rotations
1 rotated, 1 overdue, 0 due soon.
- Rotated: **Proxmox** (2026-09-02 09:00 UTC)
- Overdue: **Netbox**

## Quality
0 new, 1 resolved, 0 still open.
- None.

## Drift
4 change(s) detected.
- **Proxmox** — 1 critical, 2 warning, 1 info

### Most severe changes
- [critical] **Proxmox**: VM <web-01> memory doubled (2026-09-03 12:30 UTC)

## Job health
1 job(s) currently failing.
- **sync** — failing: timeout
- **backup** — ok
`

	if got != want {
		t.Errorf("RenderMarkdown mismatch.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestRenderMarkdown_Truncated(t *testing.T) {
	data := sampleData()
	data.Truncated = true
	got, err := RenderMarkdown(data)
	if err != nil {
		t.Fatalf("RenderMarkdown: %v", err)
	}
	if !strings.Contains(got, "Window truncated to 31 days") {
		t.Errorf("expected truncation notice, got:\n%s", got)
	}
}

func TestRenderMarkdown_SectionUnavailable(t *testing.T) {
	data := sampleData()
	data.Sections.Drift = DriftSection{Error: "query timed out"}
	got, err := RenderMarkdown(data)
	if err != nil {
		t.Fatalf("RenderMarkdown: %v", err)
	}
	if !strings.Contains(got, "_Unavailable: query timed out_") {
		t.Errorf("expected drift section marked unavailable, got:\n%s", got)
	}
	// Other sections still render.
	if !strings.Contains(got, "## Job health") {
		t.Errorf("expected other sections to still render, got:\n%s", got)
	}
}

func TestRenderHTML_EscapesDocTitles(t *testing.T) {
	got, err := RenderHTML(sampleData())
	if err != nil {
		t.Fatalf("RenderHTML: %v", err)
	}

	// html/template must escape the raw "<Runbook> & Notes" title and the
	// "<web-01>" change summary — neither should appear as raw HTML tags.
	if strings.Contains(got, "<Runbook>") {
		t.Errorf("doc title was not escaped, got:\n%s", got)
	}
	if !strings.Contains(got, "&lt;Runbook&gt; &amp; Notes") {
		t.Errorf("expected escaped doc title, got:\n%s", got)
	}
	if strings.Contains(got, "<web-01>") {
		t.Errorf("change summary was not escaped, got:\n%s", got)
	}
	if !strings.Contains(got, "&lt;web-01&gt;") {
		t.Errorf("expected escaped change summary, got:\n%s", got)
	}

	// Sanity: still a well-formed self-contained page.
	if !strings.Contains(got, "<style>") {
		t.Errorf("expected inline <style>, got:\n%s", got)
	}
	if !strings.Contains(got, "<title>Weekly Ops Report</title>") {
		t.Errorf("expected page title, got:\n%s", got)
	}
}

func TestRenderHTML_Truncated(t *testing.T) {
	data := sampleData()
	data.Truncated = true
	got, err := RenderHTML(data)
	if err != nil {
		t.Fatalf("RenderHTML: %v", err)
	}
	if !strings.Contains(got, "Window truncated to 31 days") {
		t.Errorf("expected truncation notice, got:\n%s", got)
	}
}

func TestRenderHTML_SectionUnavailable(t *testing.T) {
	data := sampleData()
	data.Sections.Quality = QualitySection{Error: "db unreachable"}
	got, err := RenderHTML(data)
	if err != nil {
		t.Fatalf("RenderHTML: %v", err)
	}
	if !strings.Contains(got, `Unavailable: db unreachable`) {
		t.Errorf("expected quality section marked unavailable, got:\n%s", got)
	}
}

func TestSummary_Golden(t *testing.T) {
	got := Summary(sampleData())
	const want = `**Weekly Ops Report** (2026-09-01 – 2026-09-08)
- Docs changed: 2
- Compliance: 1 open, 1 new, 0 resolved
- Rotations: 1 rotated, 1 overdue, 0 due soon
- Quality: 0 open, 0 new, 1 resolved
- Drift: 4 change(s)
- Job health: 1 failing
`
	if got != want {
		t.Errorf("Summary mismatch.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestSummary_Truncated(t *testing.T) {
	data := sampleData()
	data.Truncated = true
	got := Summary(data)
	if !strings.Contains(got, "_(truncated)_") {
		t.Errorf("expected truncated marker, got:\n%s", got)
	}
}

func TestSummary_SectionUnavailable(t *testing.T) {
	data := sampleData()
	data.Sections.Jobs = JobsSection{Error: "store error"}
	got := Summary(data)
	if !strings.Contains(got, "- Job health: unavailable (store error)") {
		t.Errorf("expected unavailable job health line, got:\n%s", got)
	}
}
