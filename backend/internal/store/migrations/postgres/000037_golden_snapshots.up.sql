-- 000037_golden_snapshots.up.sql — pin a snapshot per connector as the
-- known-good baseline, and add 'config_drift' as a valid quality_findings
-- check_type so later snapshots can be flagged when they deviate from it
-- (#275).

CREATE TABLE golden_snapshots (
    connector_id TEXT PRIMARY KEY REFERENCES connectors(id) ON DELETE CASCADE,
    snapshot_id  TEXT NOT NULL REFERENCES service_snapshots(id) ON DELETE CASCADE,
    pinned_by    TEXT NOT NULL,
    pinned_at    TEXT NOT NULL
);

ALTER TABLE quality_findings DROP CONSTRAINT quality_findings_check_type_check;
ALTER TABLE quality_findings ADD CONSTRAINT quality_findings_check_type_check
    CHECK(check_type IN ('stale','empty','failing','ownership_incomplete','credential_rotation','compliance','config_drift'));
