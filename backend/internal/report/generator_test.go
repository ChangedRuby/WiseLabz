package report_test

import (
	"context"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/report"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestGeneratorPersistsPartialReportWhenASectionQueryFails(t *testing.T) {
	s := apitest.NewStore(t)
	if _, err := s.DB().ExecContext(context.Background(), `DROP TABLE doc_versions`); err != nil {
		t.Fatal(err)
	}
	r, err := report.NewGenerator(s).Generate(context.Background(), store.ReportDefinitionRecord{
		ID: "definition", Slug: "weekly", Name: "Weekly", Sections: `["docs"]`,
	}, "manual")
	if err == nil || r.Status != "partial" || r.ID == "" {
		t.Fatalf("Generate() = %+v, %v; want persisted partial report", r, err)
	}
	got, err := s.GetReport(context.Background(), r.ID)
	if err != nil || got.Status != "partial" {
		t.Fatalf("GetReport() = %+v, %v; want partial", got, err)
	}
}
