// Package reports serves scheduled-report definitions and generated snapshots.
package reports

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/report"
	"github.com/WiseLabz/wiselabz/internal/store"
)

type Handler struct {
	Store   *store.Store
	Manager *report.Manager
}

func NewHandler(s *store.Store, m *report.Manager) *Handler { return &Handler{Store: s, Manager: m} }

type input struct {
	Slug         string   `json:"slug"`
	Name         string   `json:"name"`
	Enabled      bool     `json:"enabled"`
	CronExpr     string   `json:"cronExpr"`
	Timezone     string   `json:"timezone"`
	Sections     []string `json:"sections"`
	ConnectorIDs []string `json:"connectorIds"`
	Channels     []string `json:"channels"`
}

var slugRE = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}$`)

func valid(in input) string {
	if !slugRE.MatchString(in.Slug) {
		return "slug must contain lowercase letters, digits, and hyphens"
	}
	if strings.TrimSpace(in.Name) == "" {
		return "name is required"
	}
	if len(in.Sections) == 0 {
		return "at least one section is required"
	}
	if _, e := time.LoadLocation(in.Timezone); e != nil {
		return "timezone must be an IANA timezone"
	}
	for _, s := range in.Sections {
		if !map[string]bool{"docs": true, "compliance": true, "rotations": true, "quality": true, "drift": true, "jobs": true}[s] {
			return "invalid section"
		}
	}
	for _, c := range in.Channels {
		if c != "slack" && c != "discord" && c != "webhook" {
			return "invalid channel"
		}
	}
	return ""
}
func record(in input, old store.ReportDefinitionRecord) store.ReportDefinitionRecord {
	sec, _ := json.Marshal(in.Sections)
	ids, _ := json.Marshal(in.ConnectorIDs)
	ch, _ := json.Marshal(in.Channels)
	return store.ReportDefinitionRecord{ID: old.ID, Slug: in.Slug, Name: in.Name, Enabled: in.Enabled, CronExpr: in.CronExpr, Timezone: in.Timezone, Sections: string(sec), ConnectorIDs: string(ids), Channels: string(ch), CreatedBy: old.CreatedBy, CreatedAt: old.CreatedAt, UpdatedAt: old.UpdatedAt}
}
func definition(r store.ReportDefinitionRecord) map[string]any {
	var s, ids, ch []string
	_ = json.Unmarshal([]byte(r.Sections), &s)
	_ = json.Unmarshal([]byte(r.ConnectorIDs), &ids)
	_ = json.Unmarshal([]byte(r.Channels), &ch)
	return map[string]any{"id": r.ID, "slug": r.Slug, "name": r.Name, "enabled": r.Enabled, "cronExpr": r.CronExpr, "timezone": r.Timezone, "sections": s, "connectorIds": ids, "channels": ch, "createdBy": r.CreatedBy, "createdAt": r.CreatedAt, "updatedAt": r.UpdatedAt}
}
func reportJSON(r store.ReportRecord, full bool) map[string]any {
	x := map[string]any{"id": r.ID, "definitionId": nil, "definitionName": r.DefinitionName, "trigger": r.Trigger, "periodStart": r.PeriodStart, "periodEnd": r.PeriodEnd, "truncated": r.Truncated, "status": r.Status, "createdAt": r.CreatedAt}
	if r.DefinitionID != "" {
		x["definitionId"] = r.DefinitionID
	}
	if full {
		var d any
		_ = json.Unmarshal([]byte(r.Data), &d)
		x["data"] = d
		x["markdown"] = r.Markdown
	}
	return x
}
func (h *Handler) ListDefinitions(w http.ResponseWriter, r *http.Request) {
	xs, e := h.Store.ListReportDefinitions(r.Context())
	if e != nil {
		httputil.Errorf(w, e)
		return
	}
	out := make([]map[string]any, len(xs))
	for i, x := range xs {
		out[i] = definition(x)
	}
	httputil.JSON(w, http.StatusOK, out)
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	in, ok := httputil.DecodeJSON[input](w, r)
	if !ok {
		return
	}
	if m := valid(in); m != "" {
		httputil.Error(w, 400, "invalid_request", m)
		return
	}
	x := record(in, store.ReportDefinitionRecord{CreatedBy: auth.UserIDFromContext(r.Context())})
	if e := h.Store.CreateReportDefinition(r.Context(), &x); e != nil {
		if errors.Is(e, store.ErrConflict) {
			httputil.Error(w, 409, "slug_taken", "Report slug is already in use")
		} else {
			httputil.Errorf(w, e)
		}
		return
	}
	if e := h.Manager.Register(x); e != nil {
		_ = h.Store.DeleteReportDefinition(r.Context(), x.ID)
		httputil.Error(w, 400, "invalid_request", "Invalid cron expression")
		return
	}
	_ = h.Store.RecordAuditFromContext(r.Context(), "report.definition.create", "report_definition", x.ID, nil)
	httputil.JSON(w, 201, definition(x))
}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	old, e := h.Store.GetReportDefinition(r.Context(), r.PathValue("id"))
	if errors.Is(e, store.ErrNotFound) {
		httputil.Error(w, 404, "not_found", "Report definition not found")
		return
	}
	if e != nil {
		httputil.Errorf(w, e)
		return
	}
	in, ok := httputil.DecodeJSON[input](w, r)
	if !ok {
		return
	}
	if in.Slug != old.Slug {
		httputil.Error(w, 400, "slug_immutable", "Report slug is immutable")
		return
	}
	if m := valid(in); m != "" {
		httputil.Error(w, 400, "invalid_request", m)
		return
	}
	x := record(in, old)
	if e = h.Manager.Register(x); e != nil {
		httputil.Error(w, 400, "invalid_request", "Invalid cron expression")
		return
	}
	if e = h.Store.UpdateReportDefinition(r.Context(), x); e != nil {
		httputil.Errorf(w, e)
		return
	}
	x.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	_ = h.Store.RecordAuditFromContext(r.Context(), "report.definition.update", "report_definition", x.ID, nil)
	httputil.JSON(w, 200, definition(x))
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	x, e := h.Store.GetReportDefinition(r.Context(), r.PathValue("id"))
	if errors.Is(e, store.ErrNotFound) {
		httputil.Error(w, 404, "not_found", "Report definition not found")
		return
	}
	if e != nil {
		httputil.Errorf(w, e)
		return
	}
	h.Manager.Unregister(x.Slug)
	if e = h.Store.DeleteReportDefinition(r.Context(), x.ID); e != nil {
		httputil.Errorf(w, e)
		return
	}
	_ = h.Store.DeleteJobHealth(r.Context(), report.JobName(x.Slug))
	_ = h.Store.RecordAuditFromContext(r.Context(), "report.definition.delete", "report_definition", x.ID, nil)
	httputil.NoContent(w)
}
func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	x, e := h.Store.GetReportDefinition(r.Context(), r.PathValue("id"))
	if errors.Is(e, store.ErrNotFound) {
		httputil.Error(w, 404, "not_found", "Report definition not found")
		return
	}
	if e != nil {
		httputil.Errorf(w, e)
		return
	}
	out, e := h.Manager.Run(r.Context(), x)
	report.LogPartial(e)
	_ = h.Store.RecordAuditFromContext(r.Context(), "report.definition.run", "report_definition", x.ID, nil)
	httputil.JSON(w, 201, reportJSON(out, true))
}
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	_, n, o := httputil.Paginate(r)
	xs, _, e := h.Store.ListReports(r.Context(), r.URL.Query().Get("definitionId"), n, o)
	if e != nil {
		httputil.Errorf(w, e)
		return
	}
	out := make([]map[string]any, len(xs))
	for i, x := range xs {
		out[i] = reportJSON(x, false)
	}
	httputil.JSON(w, 200, map[string]any{"items": out})
}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	x, e := h.Store.GetReport(r.Context(), r.PathValue("id"))
	if errors.Is(e, store.ErrNotFound) {
		httputil.Error(w, 404, "not_found", "Report not found")
		return
	}
	if e != nil {
		httputil.Errorf(w, e)
		return
	}
	httputil.JSON(w, 200, reportJSON(x, true))
}
func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	x, e := h.Store.GetReport(r.Context(), r.PathValue("id"))
	if errors.Is(e, store.ErrNotFound) {
		httputil.Error(w, 404, "not_found", "Report not found")
		return
	}
	if e != nil {
		httputil.Errorf(w, e)
		return
	}
	format := r.URL.Query().Get("format")
	if format != "md" && format != "html" {
		httputil.Error(w, 400, "invalid_request", "format must be md or html")
		return
	}
	body := x.Markdown
	ct := "text/markdown; charset=utf-8"
	if format == "html" {
		var d report.ReportData
		if e = json.Unmarshal([]byte(x.Data), &d); e != nil {
			httputil.Errorf(w, e)
			return
		}
		body, e = report.RenderHTML(d)
		if e != nil {
			httputil.Errorf(w, e)
			return
		}
		ct = "text/html; charset=utf-8"
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Disposition", `attachment; filename="report-`+x.ID+`.`+format+`"`)
	_, _ = w.Write([]byte(body))
}
