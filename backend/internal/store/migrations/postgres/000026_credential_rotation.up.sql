-- 000026_credential_rotation.up.sql — secret rotation reminders (#239 PR1).

ALTER TABLE connectors ADD COLUMN secret_rotated_at TEXT;
UPDATE connectors SET secret_rotated_at = created_at WHERE secret_rotated_at IS NULL;
ALTER TABLE connectors ADD COLUMN user_expires_at TEXT;
ALTER TABLE connectors ADD COLUMN rotation_max_age_days INTEGER;

ALTER TABLE quality_findings DROP CONSTRAINT quality_findings_check_type_check;
ALTER TABLE quality_findings ADD CONSTRAINT quality_findings_check_type_check
    CHECK(check_type IN ('stale','empty','failing','ownership_incomplete','credential_rotation'));
ALTER TABLE quality_findings ADD COLUMN notified_severity TEXT;
