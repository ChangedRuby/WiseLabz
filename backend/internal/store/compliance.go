package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ComplianceRuleRecord is a persisted user-defined snapshot compliance rule.
// Conditions is the JSON condition array; rule evaluation deliberately lives
// outside the store package.
type ComplianceRuleRecord struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ConnectorType   string `json:"connectorType"`
	EntityKind      string `json:"entityKind"`
	Conditions      string `json:"conditions"`
	Severity        string `json:"severity"`
	Title           string `json:"title"`
	RemediationLink string `json:"remediationLink"`
	Enabled         bool   `json:"enabled"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

const complianceRuleColumns = `id, name, connector_type, entity_kind, conditions, severity, title, remediation_link, enabled, created_at, updated_at`

// CreateComplianceRule persists a new compliance rule.
func (s *Store) CreateComplianceRule(ctx context.Context, r *ComplianceRuleRecord) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}
	if r.UpdatedAt == "" {
		r.UpdatedAt = now
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO compliance_rules (id, name, connector_type, entity_kind, conditions, severity, title, remediation_link, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.ID, r.Name, r.ConnectorType, r.EntityKind, r.Conditions, r.Severity, r.Title, r.RemediationLink, r.Enabled, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create compliance rule: %w", err)
	}
	return nil
}

// GetComplianceRule returns one rule by ID.
func (s *Store) GetComplianceRule(ctx context.Context, id string) (*ComplianceRuleRecord, error) {
	r, err := scanComplianceRule(s.db.QueryRowContext(ctx, `SELECT `+complianceRuleColumns+` FROM compliance_rules WHERE id = ?`, id))
	if err != nil {
		return nil, fmt.Errorf("get compliance rule: %w", err)
	}
	return r, nil
}

// ListComplianceRules returns every rule, newest first.
func (s *Store) ListComplianceRules(ctx context.Context) ([]ComplianceRuleRecord, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+complianceRuleColumns+` FROM compliance_rules ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list compliance rules: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	rules := make([]ComplianceRuleRecord, 0)
	for rows.Next() {
		r, err := scanComplianceRule(rows)
		if err != nil {
			return nil, fmt.Errorf("scan compliance rule: %w", err)
		}
		rules = append(rules, *r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate compliance rules: %w", err)
	}
	return rules, nil
}

// UpdateComplianceRule replaces the mutable fields of a rule.
func (s *Store) UpdateComplianceRule(ctx context.Context, r *ComplianceRuleRecord) error {
	r.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	result, err := s.db.ExecContext(ctx, `
		UPDATE compliance_rules
		SET name = ?, connector_type = ?, entity_kind = ?, conditions = ?, severity = ?, title = ?, remediation_link = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`, r.Name, r.ConnectorType, r.EntityKind, r.Conditions, r.Severity, r.Title, r.RemediationLink, r.Enabled, r.UpdatedAt, r.ID)
	if err != nil {
		return fmt.Errorf("update compliance rule: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteComplianceRule removes a rule by ID.
func (s *Store) DeleteComplianceRule(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM compliance_rules WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete compliance rule: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func scanComplianceRule(row rowScanner) (*ComplianceRuleRecord, error) {
	var r ComplianceRuleRecord
	if err := row.Scan(&r.ID, &r.Name, &r.ConnectorType, &r.EntityKind, &r.Conditions, &r.Severity, &r.Title, &r.RemediationLink, &r.Enabled, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &r, nil
}
