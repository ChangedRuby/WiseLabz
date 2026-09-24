// Package report builds scheduled/manual instance reports (#280): a
// point-in-time snapshot of docs changed, compliance/quality findings,
// secret rotations, drift, and job health over a time window, rendered as
// Markdown and HTML and delivered through existing notification channels.
//
// This file defines the data shapes only. Nothing here talks to the store —
// callers (the future report.Generator) fetch rows themselves and assemble
// a ReportData, which keeps rendering and window math unit-testable without
// a database.
package report

import "time"

// DefinitionSummary is the minimal snapshot of the report_definitions row
// that produced a report, captured at generation time. It is stored on the
// report itself (not just referenced by ID) so a report stays readable even
// after its definition is edited or deleted (definition_id ON DELETE SET
// NULL, per the plan).
type DefinitionSummary struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// ReportData is the typed payload persisted as JSON on the `reports` row.
// HTML/Markdown are rendered from it (see render.go); the Markdown is also
// stored alongside for fast retrieval.
type ReportData struct { //nolint:revive // name matches the API schema (docs/openapi.yaml ReportData) and the plan's spec
	Definition  DefinitionSummary `json:"definition"`
	PeriodStart time.Time         `json:"periodStart"`
	PeriodEnd   time.Time         `json:"periodEnd"`
	// Truncated is true when the window exceeded the 31-day cap (Q22) and
	// was clamped to end.Add(-31 days); the report header notes this.
	Truncated bool     `json:"truncated"`
	Sections  Sections `json:"sections"`
}

// Sections holds one sub-report per area (plan Q4). Each section carries its
// own Error so a single failing query degrades that section instead of
// failing the whole report (plan Q21): the section is marked unavailable,
// the report is still stored/sent, and the generator returns a joined error.
type Sections struct {
	Docs       DocsSection       `json:"docs"`
	Compliance ComplianceSection `json:"compliance"`
	Rotations  RotationsSection  `json:"rotations"`
	Quality    QualitySection    `json:"quality"`
	Drift      DriftSection      `json:"drift"`
	Jobs       JobsSection       `json:"jobs"`
}

// DocChangeEntry is one doc version created or updated in the report window.
type DocChangeEntry struct {
	DocID         string    `json:"docId"`
	Title         string    `json:"title"`
	ConnectorID   string    `json:"connectorId"`
	ConnectorName string    `json:"connectorName"`
	Version       int       `json:"version"`
	ChangedAt     time.Time `json:"changedAt"`
}

// DocsSection summarizes documentation changed during the window.
type DocsSection struct {
	// Error is non-empty when this section's query failed; Entries/Total
	// are then zero-valued and the section is rendered as unavailable.
	Error string           `json:"error,omitempty"`
	Total int              `json:"total"`
	Items []DocChangeEntry `json:"items"`
}

// FindingSummary is a quality_findings row shape shared by the compliance
// and quality sections (they differ only in which check_type they filter
// to — plan's compliance = check_type=compliance, quality = the rest).
type FindingSummary struct {
	ID            string     `json:"id"`
	ConnectorID   string     `json:"connectorId"`
	ConnectorName string     `json:"connectorName"`
	CheckType     string     `json:"checkType"`
	Severity      string     `json:"severity"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	DetectedAt    time.Time  `json:"detectedAt"`
	ResolvedAt    *time.Time `json:"resolvedAt,omitempty"`
}

// ComplianceSection summarizes compliance findings (check_type=compliance)
// newly detected or resolved during the window, plus a still-open count.
type ComplianceSection struct {
	Error         string           `json:"error,omitempty"`
	OpenCount     int              `json:"openCount"`
	DetectedCount int              `json:"detectedCount"`
	ResolvedCount int              `json:"resolvedCount"`
	Items         []FindingSummary `json:"items"`
}

// QualitySection is ComplianceSection's counterpart for every other
// check_type (stale, empty, failing, ownership_incomplete, config_drift).
type QualitySection struct {
	Error         string           `json:"error,omitempty"`
	OpenCount     int              `json:"openCount"`
	DetectedCount int              `json:"detectedCount"`
	ResolvedCount int              `json:"resolvedCount"`
	Items         []FindingSummary `json:"items"`
}

// RotationEntry is one connector's secret-rotation state relevant to the
// window: either it rotated during the window, or it currently has an open
// credential_rotation finding (due soon or overdue).
type RotationEntry struct {
	ConnectorID   string     `json:"connectorId"`
	ConnectorName string     `json:"connectorName"`
	RotatedAt     *time.Time `json:"rotatedAt,omitempty"`
	Overdue       bool       `json:"overdue"`
	DueSoon       bool       `json:"dueSoon"`
}

// RotationsSection summarizes secret rotations completed in the window plus
// connectors with an open credential_rotation finding right now.
type RotationsSection struct {
	Error   string          `json:"error,omitempty"`
	Rotated []RotationEntry `json:"rotated"`
	Overdue []RotationEntry `json:"overdue"`
	DueSoon []RotationEntry `json:"dueSoon"`
}

// SeverityCounts tallies changes by severity (plan Q19).
type SeverityCounts struct {
	Info     int `json:"info"`
	Warning  int `json:"warning"`
	Critical int `json:"critical"`
}

// ConnectorDrift is one connector's change counts for the window.
type ConnectorDrift struct {
	ConnectorID   string         `json:"connectorId"`
	ConnectorName string         `json:"connectorName"`
	Counts        SeverityCounts `json:"counts"`
}

// ChangeEntry is one drift change, used for the section's top-N list.
type ChangeEntry struct {
	ID            string    `json:"id"`
	ConnectorID   string    `json:"connectorId"`
	ConnectorName string    `json:"connectorName"`
	Severity      string    `json:"severity"`
	Summary       string    `json:"summary"`
	DetectedAt    time.Time `json:"detectedAt"`
}

// DriftSection summarizes detected changes (drift) during the window: per
// -connector severity counts plus the 10 most severe individual changes
// (plan Q19).
type DriftSection struct {
	Error       string           `json:"error,omitempty"`
	Total       int              `json:"total"`
	ByConnector []ConnectorDrift `json:"byConnector"`
	TopChanges  []ChangeEntry    `json:"topChanges"`
}

// JobHealthEntry mirrors one row from PR1's job_health store table
// (backend/internal/store/job_health.go), copied in rather than imported so
// this package has no store dependency.
type JobHealthEntry struct {
	Name          string     `json:"name"`
	Status        string     `json:"status"` // "ok" | "failing"
	LastError     string     `json:"lastError,omitempty"`
	LastRunAt     *time.Time `json:"lastRunAt,omitempty"`
	LastSuccessAt *time.Time `json:"lastSuccessAt,omitempty"`
	LastFailureAt *time.Time `json:"lastFailureAt,omitempty"`
}

// JobsSection is the job-health rollup (plan Q4's sixth section), sourced
// from store.ListJobHealth() — job health itself ignores any
// connector_ids filter (plan Q14).
type JobsSection struct {
	Error        string           `json:"error,omitempty"`
	FailingCount int              `json:"failingCount"`
	Items        []JobHealthEntry `json:"items"`
}
