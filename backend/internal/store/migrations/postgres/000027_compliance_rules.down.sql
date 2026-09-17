DROP INDEX idx_quality_findings_open;
DELETE FROM quality_findings WHERE check_type = 'compliance';
ALTER TABLE quality_findings DROP CONSTRAINT quality_findings_check_type_check;
ALTER TABLE quality_findings ADD CONSTRAINT quality_findings_check_type_check
    CHECK(check_type IN ('stale','empty','failing','ownership_incomplete','credential_rotation'));
ALTER TABLE quality_findings DROP COLUMN rule_id;
CREATE UNIQUE INDEX idx_quality_findings_open ON quality_findings(connector_id, check_type) WHERE status = 'open';
DROP TABLE compliance_rules;
