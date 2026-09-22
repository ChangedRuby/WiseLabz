-- Revert 000037: drop golden_snapshots and restrict quality_findings.check_type
-- back to the set that excludes 'config_drift'.
-- Rows with check_type = 'config_drift' are dropped first so the narrower
-- constraint can be re-added without failing on existing data.

DROP TABLE golden_snapshots;

DELETE FROM quality_findings WHERE check_type = 'config_drift';
ALTER TABLE quality_findings DROP CONSTRAINT quality_findings_check_type_check;
ALTER TABLE quality_findings ADD CONSTRAINT quality_findings_check_type_check
    CHECK(check_type IN ('stale','empty','failing','ownership_incomplete','credential_rotation','compliance'));
