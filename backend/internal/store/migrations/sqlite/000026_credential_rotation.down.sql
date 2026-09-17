-- 000026_credential_rotation.down.sql
-- Note: rows with check_type = 'credential_rotation' are dropped by the
-- WHERE clause below. This is expected for a down migration.

DROP INDEX IF EXISTS idx_quality_findings_open;
DROP INDEX IF EXISTS idx_quality_findings_status;
DROP INDEX IF EXISTS idx_quality_findings_connector;

CREATE TABLE quality_findings_new (
    id                  TEXT PRIMARY KEY,
    connector_id        TEXT NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    doc_id              TEXT REFERENCES docs(id) ON DELETE SET NULL,
    check_type          TEXT NOT NULL CHECK(check_type IN ('stale','empty','failing','ownership_incomplete')),
    severity            TEXT NOT NULL CHECK(severity IN ('info','warning','critical')),
    title               TEXT NOT NULL,
    description         TEXT NOT NULL DEFAULT '',
    remediation_link    TEXT NOT NULL DEFAULT '',
    status              TEXT NOT NULL DEFAULT 'open' CHECK(status IN ('open','resolved')),
    detected_count      INTEGER NOT NULL DEFAULT 1,
    first_detected_at   TEXT NOT NULL,
    last_seen_at        TEXT NOT NULL,
    resolved_at         TEXT
);
INSERT INTO quality_findings_new (id, connector_id, doc_id, check_type, severity, title, description,
    remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at)
SELECT id, connector_id, doc_id, check_type, severity, title, description,
    remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at
FROM quality_findings WHERE check_type != 'credential_rotation';
DROP TABLE quality_findings;
ALTER TABLE quality_findings_new RENAME TO quality_findings;

CREATE UNIQUE INDEX idx_quality_findings_open ON quality_findings(connector_id, check_type) WHERE status = 'open';
CREATE INDEX idx_quality_findings_status ON quality_findings(status, last_seen_at DESC);
CREATE INDEX idx_quality_findings_connector ON quality_findings(connector_id);

ALTER TABLE connectors DROP COLUMN rotation_max_age_days;
ALTER TABLE connectors DROP COLUMN user_expires_at;
ALTER TABLE connectors DROP COLUMN secret_rotated_at;
