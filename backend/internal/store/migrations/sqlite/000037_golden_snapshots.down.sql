-- Revert 000037: drop golden_snapshots and restrict quality_findings.check_type
-- back to the set that excludes 'config_drift'.
-- Note: any open config_drift findings are dropped since they'd violate the
-- narrower constraint. This is expected for a down migration.

DROP TABLE golden_snapshots;

DROP INDEX IF EXISTS idx_quality_findings_open;
DROP INDEX IF EXISTS idx_quality_findings_status;
DROP INDEX IF EXISTS idx_quality_findings_connector;

CREATE TABLE quality_findings_new (
    id                  TEXT PRIMARY KEY,
    connector_id        TEXT NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    doc_id              TEXT REFERENCES docs(id) ON DELETE SET NULL,
    rule_id             TEXT REFERENCES compliance_rules(id) ON DELETE SET NULL,
    check_type          TEXT NOT NULL CHECK(check_type IN ('stale','empty','failing','ownership_incomplete','credential_rotation','compliance')),
    severity            TEXT NOT NULL CHECK(severity IN ('info','warning','critical')),
    title               TEXT NOT NULL,
    description         TEXT NOT NULL DEFAULT '',
    remediation_link    TEXT NOT NULL DEFAULT '',
    status              TEXT NOT NULL DEFAULT 'open' CHECK(status IN ('open','resolved')),
    detected_count      INTEGER NOT NULL DEFAULT 1,
    first_detected_at   TEXT NOT NULL,
    last_seen_at        TEXT NOT NULL,
    resolved_at         TEXT,
    notified_severity   TEXT
);
INSERT INTO quality_findings_new SELECT * FROM quality_findings WHERE check_type != 'config_drift';
DROP TABLE quality_findings;
ALTER TABLE quality_findings_new RENAME TO quality_findings;

CREATE UNIQUE INDEX idx_quality_findings_open ON quality_findings(connector_id, check_type, COALESCE(rule_id, '')) WHERE status = 'open';
CREATE INDEX idx_quality_findings_status ON quality_findings(status, last_seen_at DESC);
CREATE INDEX idx_quality_findings_connector ON quality_findings(connector_id);
