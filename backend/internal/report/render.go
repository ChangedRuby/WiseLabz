package report

import (
	"bytes"
	"embed"
	htmltemplate "html/template"
	"strconv"
	"text/template"
	"time"
)

//go:embed templates/report.md.tmpl
var mdTemplateFS embed.FS

//go:embed templates/report.html.tmpl
var htmlTemplateFS embed.FS

// funcMap is shared by both templates. Kept identical (rather than two maps)
// so the two renderings never drift in how they format a date or the
// truncation-cap wording.
var funcMap = map[string]any{
	"fmtDate": func(t time.Time) string {
		return t.UTC().Format("2006-01-02 15:04 UTC")
	},
	"maxWindowDays": func() string {
		return strconv.Itoa(int(maxWindow / (24 * time.Hour)))
	},
}

var (
	mdTemplate   = template.Must(template.New("report.md.tmpl").Funcs(funcMap).ParseFS(mdTemplateFS, "templates/report.md.tmpl"))
	htmlTemplate = htmltemplate.Must(htmltemplate.New("report.html.tmpl").Funcs(funcMap).ParseFS(htmlTemplateFS, "templates/report.html.tmpl"))
)

// RenderMarkdown renders data as the Markdown document stored on the report
// row (and shown/downloaded as-is via GET /api/reports/{id}).
func RenderMarkdown(data ReportData) (string, error) {
	var buf bytes.Buffer
	if err := mdTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderHTML renders data as a self-contained HTML page (inline CSS, no
// external assets) for the download-as-HTML endpoint. It is rendered on
// demand from the stored JSON data rather than persisted, so template fixes
// apply retroactively to old reports.
func RenderHTML(data ReportData) (string, error) {
	var buf bytes.Buffer
	if err := htmlTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
