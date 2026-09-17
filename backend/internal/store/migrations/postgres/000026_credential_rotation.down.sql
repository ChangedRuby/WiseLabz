-- 000026_credential_rotation.down.sql
-- Note: rows with check_type = 'credential_rotation' will violate the
-- reverted constraint. This is expected for a down migration.

ALTER TABLE quality_findings DROP COLUMN notified_severity;
ALTER TABLE quality_findings DROP CONSTRAINT quality_findings_check_type_check;
ALTER TABLE quality_findings ADD CONSTRAINT quality_findings_check_type_check
    CHECK(check_type IN ('stale','empty','failing','ownership_incomplete'));

ALTER TABLE connectors DROP COLUMN rotation_max_age_days;
ALTER TABLE connectors DROP COLUMN user_expires_at;
ALTER TABLE connectors DROP COLUMN secret_rotated_at;
