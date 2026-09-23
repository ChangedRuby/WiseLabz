package report

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// Generator builds and persists reports. Queries are deliberately kept here:
// each section is independent, so a broken historical table makes one section
// partial instead of losing the entire scheduled delivery.
type Generator struct{ Store *store.Store }

func NewGenerator(s *store.Store) *Generator { return &Generator{Store: s} }

func (g *Generator) Generate(ctx context.Context, def store.ReportDefinitionRecord, trigger string) (store.ReportRecord, error) {
	var watermark *time.Time
	if last, err := g.Store.LatestScheduledReport(ctx, def.ID); err == nil {
		if t, e := time.Parse(time.RFC3339, last.PeriodEnd); e == nil {
			watermark = &t
		}
	} else if !errors.Is(err, store.ErrNotFound) {
		return store.ReportRecord{}, err
	}
	start, end, truncated := ComputeWindow(watermark, time.Now().UTC())
	d := ReportData{Definition: DefinitionSummary{Slug: def.Slug, Name: def.Name}, PeriodStart: start, PeriodEnd: end, Truncated: truncated,
		Sections: Sections{Docs: DocsSection{Items: []DocChangeEntry{}}, Compliance: ComplianceSection{Items: []FindingSummary{}}, Rotations: RotationsSection{Rotated: []RotationEntry{}, Overdue: []RotationEntry{}, DueSoon: []RotationEntry{}}, Quality: QualitySection{Items: []FindingSummary{}}, Drift: DriftSection{ByConnector: []ConnectorDrift{}, TopChanges: []ChangeEntry{}}, Jobs: JobsSection{Items: []JobHealthEntry{}}}}
	var sectionNames []string
	_ = json.Unmarshal([]byte(def.Sections), &sectionNames)
	selected := make(map[string]bool, len(sectionNames))
	for _, name := range sectionNames {
		selected[name] = true
	}
	var errs []error
	if selected["docs"] {
		if e := g.docs(ctx, &d, def.ConnectorIDs); e != nil {
			d.Sections.Docs.Error = e.Error()
			errs = append(errs, e)
		}
	}
	if selected["compliance"] {
		if e := g.findings(ctx, &d, def.ConnectorIDs, true); e != nil {
			d.Sections.Compliance.Error = e.Error()
			errs = append(errs, e)
		}
	}
	if selected["quality"] {
		if e := g.findings(ctx, &d, def.ConnectorIDs, false); e != nil {
			d.Sections.Quality.Error = e.Error()
			errs = append(errs, e)
		}
	}
	if selected["rotations"] {
		if e := g.rotations(ctx, &d, def.ConnectorIDs); e != nil {
			d.Sections.Rotations.Error = e.Error()
			errs = append(errs, e)
		}
	}
	if selected["drift"] {
		if e := g.drift(ctx, &d, def.ConnectorIDs); e != nil {
			d.Sections.Drift.Error = e.Error()
			errs = append(errs, e)
		}
	}
	if selected["jobs"] {
		if e := g.jobs(ctx, &d); e != nil {
			d.Sections.Jobs.Error = e.Error()
			errs = append(errs, e)
		}
	}
	data, err := json.Marshal(d)
	if err != nil {
		return store.ReportRecord{}, err
	}
	markdown, err := RenderMarkdown(d)
	if err != nil {
		return store.ReportRecord{}, err
	}
	r := store.ReportRecord{DefinitionID: def.ID, DefinitionName: def.Name, Trigger: trigger, PeriodStart: start.Format(time.RFC3339), PeriodEnd: end.Format(time.RFC3339), Truncated: truncated, Data: string(data), Markdown: markdown, Status: "ok"}
	if len(errs) > 0 {
		r.Status = "partial"
	}
	if err := g.Store.CreateReport(ctx, &r); err != nil {
		return store.ReportRecord{}, err
	}
	return r, errors.Join(errs...)
}

func connectorFilter(raw string) (string, []any) {
	var ids []string
	_ = json.Unmarshal([]byte(raw), &ids)
	if len(ids) == 0 {
		return "", nil
	}
	q := " AND c.id IN ("
	a := make([]any, len(ids))
	for i, id := range ids {
		if i > 0 {
			q += ","
		}
		q += "?"
		a[i] = id
	}
	return q + ")", a
}
func (g *Generator) docs(ctx context.Context, d *ReportData, ids string) error {
	f, a := connectorFilter(ids)
	rows, e := g.Store.DB().QueryContext(ctx, `SELECT d.id,d.title,COALESCE(d.service_id,''),COALESCE(c.name,''),v.rev,v.created_at FROM doc_versions v JOIN docs d ON d.id=v.doc_id LEFT JOIN connectors c ON c.id=d.service_id WHERE v.created_at>=? AND v.created_at<?`+f+` ORDER BY v.created_at DESC`, append([]any{d.PeriodStart.Format(time.RFC3339), d.PeriodEnd.Format(time.RFC3339)}, a...)...)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		var x DocChangeEntry
		var ts string
		if e = rows.Scan(&x.DocID, &x.Title, &x.ConnectorID, &x.ConnectorName, &x.Version, &ts); e != nil {
			return e
		}
		x.ChangedAt, _ = time.Parse(time.RFC3339, ts)
		d.Sections.Docs.Items = append(d.Sections.Docs.Items, x)
	}
	d.Sections.Docs.Total = len(d.Sections.Docs.Items)
	return rows.Err()
}
func (g *Generator) findings(ctx context.Context, d *ReportData, ids string, compliance bool) error {
	f, a := connectorFilter(ids)
	kind := "q.check_type <> 'compliance' AND q.check_type <> 'credential_rotation'"
	if compliance {
		kind = "q.check_type = 'compliance'"
	}
	args := append([]any{d.PeriodStart.Format(time.RFC3339), d.PeriodEnd.Format(time.RFC3339)}, a...)
	rows, e := g.Store.DB().QueryContext(ctx, `SELECT q.id,q.connector_id,c.name,q.check_type,q.severity,q.title,q.status,q.first_detected_at,q.resolved_at FROM quality_findings q JOIN connectors c ON c.id=q.connector_id WHERE (`+kind+`) AND (q.first_detected_at>=? AND q.first_detected_at<? OR q.resolved_at>=? AND q.resolved_at<?)`+f, append(args[:2], append(args[:2], args[2:]...)...)...)
	if e != nil {
		return e
	}
	defer rows.Close()
	var items []FindingSummary
	for rows.Next() {
		var x FindingSummary
		var ds string
		var rs sql.NullString
		if e = rows.Scan(&x.ID, &x.ConnectorID, &x.ConnectorName, &x.CheckType, &x.Severity, &x.Title, &x.Status, &ds, &rs); e != nil {
			return e
		}
		x.DetectedAt, _ = time.Parse(time.RFC3339, ds)
		if rs.Valid {
			t, _ := time.Parse(time.RFC3339, rs.String)
			x.ResolvedAt = &t
		}
		items = append(items, x)
	}
	if e = rows.Err(); e != nil {
		return e
	}
	var open, det, res int
	for _, x := range items {
		if x.DetectedAt.After(d.PeriodStart) || x.DetectedAt.Equal(d.PeriodStart) {
			det++
		}
		if x.ResolvedAt != nil {
			res++
		}
	}
	e = g.Store.DB().QueryRowContext(ctx, `SELECT COUNT(*) FROM quality_findings q JOIN connectors c ON c.id=q.connector_id WHERE (`+kind+`) AND q.status='open'`+f, a...).Scan(&open)
	if e != nil {
		return e
	}
	if compliance {
		d.Sections.Compliance = ComplianceSection{OpenCount: open, DetectedCount: det, ResolvedCount: res, Items: items}
	} else {
		d.Sections.Quality = QualitySection{OpenCount: open, DetectedCount: det, ResolvedCount: res, Items: items}
	}
	return nil
}
func (g *Generator) rotations(ctx context.Context, d *ReportData, ids string) error {
	f, a := connectorFilter(ids)
	rows, e := g.Store.DB().QueryContext(ctx, `SELECT c.id,c.name,c.secret_rotated_at FROM connectors c WHERE c.secret_rotated_at>=? AND c.secret_rotated_at<?`+f, append([]any{d.PeriodStart.Format(time.RFC3339), d.PeriodEnd.Format(time.RFC3339)}, a...)...)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		var x RotationEntry
		var ts string
		if e = rows.Scan(&x.ConnectorID, &x.ConnectorName, &ts); e != nil {
			return e
		}
		t, _ := time.Parse(time.RFC3339, ts)
		x.RotatedAt = &t
		d.Sections.Rotations.Rotated = append(d.Sections.Rotations.Rotated, x)
	}
	if e = rows.Err(); e != nil {
		return e
	}
	rows, e = g.Store.DB().QueryContext(ctx, `SELECT c.id,c.name,q.severity FROM quality_findings q JOIN connectors c ON c.id=q.connector_id WHERE q.check_type='credential_rotation' AND q.status='open'`+f, a...)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		var x RotationEntry
		var severity string
		if e = rows.Scan(&x.ConnectorID, &x.ConnectorName, &severity); e != nil {
			return e
		}
		if severity == "critical" {
			x.Overdue = true
			d.Sections.Rotations.Overdue = append(d.Sections.Rotations.Overdue, x)
		} else {
			x.DueSoon = true
			d.Sections.Rotations.DueSoon = append(d.Sections.Rotations.DueSoon, x)
		}
	}
	return rows.Err()
}
func (g *Generator) drift(ctx context.Context, d *ReportData, ids string) error {
	f, a := connectorFilter(ids)
	if e := g.Store.DB().QueryRowContext(ctx, `SELECT COUNT(*) FROM changes x JOIN connectors c ON c.id=x.service_id WHERE x.detected_at>=? AND x.detected_at<?`+f, append([]any{d.PeriodStart.Format(time.RFC3339), d.PeriodEnd.Format(time.RFC3339)}, a...)...).Scan(&d.Sections.Drift.Total); e != nil {
		return e
	}
	rows, e := g.Store.DB().QueryContext(ctx, `SELECT x.id,x.service_id,c.name,x.severity,x.summary,x.detected_at FROM changes x JOIN connectors c ON c.id=x.service_id WHERE x.detected_at>=? AND x.detected_at<?`+f+` ORDER BY CASE x.severity WHEN 'critical' THEN 0 WHEN 'warning' THEN 1 ELSE 2 END,x.detected_at DESC LIMIT 10`, append([]any{d.PeriodStart.Format(time.RFC3339), d.PeriodEnd.Format(time.RFC3339)}, a...)...)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		var x ChangeEntry
		var ts string
		if e = rows.Scan(&x.ID, &x.ConnectorID, &x.ConnectorName, &x.Severity, &x.Summary, &ts); e != nil {
			return e
		}
		x.DetectedAt, _ = time.Parse(time.RFC3339, ts)
		d.Sections.Drift.TopChanges = append(d.Sections.Drift.TopChanges, x)
	}
	if e = rows.Err(); e != nil {
		return e
	}
	rows, e = g.Store.DB().QueryContext(ctx, `SELECT c.id,c.name,SUM(CASE WHEN x.severity='info' THEN 1 ELSE 0 END),SUM(CASE WHEN x.severity='warning' THEN 1 ELSE 0 END),SUM(CASE WHEN x.severity='critical' THEN 1 ELSE 0 END) FROM changes x JOIN connectors c ON c.id=x.service_id WHERE x.detected_at>=? AND x.detected_at<?`+f+` GROUP BY c.id,c.name ORDER BY c.name`, append([]any{d.PeriodStart.Format(time.RFC3339), d.PeriodEnd.Format(time.RFC3339)}, a...)...)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		var x ConnectorDrift
		if e = rows.Scan(&x.ConnectorID, &x.ConnectorName, &x.Counts.Info, &x.Counts.Warning, &x.Counts.Critical); e != nil {
			return e
		}
		d.Sections.Drift.ByConnector = append(d.Sections.Drift.ByConnector, x)
	}
	return rows.Err()
}
func (g *Generator) jobs(ctx context.Context, d *ReportData) error {
	rows, e := g.Store.DB().QueryContext(ctx, `SELECT name,last_status,last_error,last_run_at,last_success_at,last_failure_at FROM job_health ORDER BY name`)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		var x JobHealthEntry
		var a, b, c sql.NullString
		if e = rows.Scan(&x.Name, &x.Status, &x.LastError, &a, &b, &c); e != nil {
			return e
		}
		for _, p := range []struct {
			s sql.NullString
			d **time.Time
		}{{a, &x.LastRunAt}, {b, &x.LastSuccessAt}, {c, &x.LastFailureAt}} {
			if p.s.Valid {
				t, _ := time.Parse(time.RFC3339, p.s.String)
				*p.d = &t
			}
		}
		if x.Status == "failing" {
			d.Sections.Jobs.FailingCount++
		}
		d.Sections.Jobs.Items = append(d.Sections.Jobs.Items, x)
	}
	return rows.Err()
}
