package store

import (
	"context"
	"testing"
)

func TestComplianceRuleCRUD(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	rule := &ComplianceRuleRecord{
		Name:          "Privileged containers",
		ConnectorType: "docker",
		EntityKind:    "container",
		Conditions:    `[{"attribute":"privileged","op":"eq","value":true}]`,
		Severity:      "critical",
		Title:         "Privileged container",
		Enabled:       true,
	}
	if err := s.CreateComplianceRule(ctx, rule); err != nil {
		t.Fatalf("CreateComplianceRule() error: %v", err)
	}
	got, err := s.GetComplianceRule(ctx, rule.ID)
	if err != nil {
		t.Fatalf("GetComplianceRule() error: %v", err)
	}
	if got.Conditions != rule.Conditions || !got.Enabled || got.CreatedAt == "" || got.UpdatedAt == "" {
		t.Fatalf("stored rule = %+v", got)
	}
	rule.Name = "Host network containers"
	rule.Enabled = false
	if err := s.UpdateComplianceRule(ctx, rule); err != nil {
		t.Fatalf("UpdateComplianceRule() error: %v", err)
	}
	rules, err := s.ListComplianceRules(ctx)
	if err != nil || len(rules) != 6 { // five disabled migration examples plus this rule
		t.Fatalf("ListComplianceRules() = %d, %v; want 6, nil", len(rules), err)
	}
	if err := s.DeleteComplianceRule(ctx, rule.ID); err != nil {
		t.Fatalf("DeleteComplianceRule() error: %v", err)
	}
}

func TestComplianceFindingRuleDedupAndResolve(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	c := seedQualityConnector(t, s, "rule findings")
	rule := &ComplianceRuleRecord{Name: "one", ConnectorType: "proxmox", EntityKind: "vm", Conditions: "[]", Severity: "warning", Title: "one"}
	if err := s.CreateComplianceRule(ctx, rule); err != nil {
		t.Fatal(err)
	}
	first := &QualityFindingRecord{ConnectorID: c.ID, RuleID: rule.ID, CheckType: "compliance", Severity: "warning", Title: "one"}
	second := &QualityFindingRecord{ConnectorID: c.ID, RuleID: rule.ID, CheckType: "compliance", Severity: "critical", Title: "two"}
	if err := s.UpsertQualityFinding(ctx, first); err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertQualityFinding(ctx, second); err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID {
		t.Fatalf("dedup IDs = %q, %q", first.ID, second.ID)
	}
	if err := s.ResolveQualityFindingForRule(ctx, c.ID, rule.ID); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetQualityFinding(ctx, first.ID)
	if err != nil || got.Status != "resolved" || got.RuleID != rule.ID {
		t.Fatalf("resolved finding = %+v, %v", got, err)
	}
}
