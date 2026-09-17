-- 000026_credential_rotation.up.sql — secret rotation reminders (#239 PR1).
-- SQLite doesn't support ALTER TABLE ... ALTER CONSTRAINT, so quality_findings
-- is rebuilt to extend the check_type CHECK (same pattern as 000009).

ALTER TABLE connectors ADD COLUMN secret_rotated_at TEXT;
UPDATE connectors SET secret_rotated_at = created_at WHERE secret_rotated_at IS NULL;
ALTER TABLE connectors ADD COLUMN user_expires_at TEXT;
ALTER TABLE connectors ADD COLUMN rotation_max_age_days INTEGER;

DROP INDEX IF EXISTS idx_quality_findings_open;
DROP INDEX IF EXISTS idx_quality_findings_status;
DROP INDEX IF EXISTS idx_quality_findings_connector;

CREATE TABLE quality_findings_new (
    id                  TEXT PRIMARY KEY,
    connector_id        TEXT NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    doc_id              TEXT REFERENCES docs(id) ON DELETE SET NULL,
    check_type          TEXT NOT NULL CHECK(check_type IN ('stale','empty','failing','ownership_incomplete','credential_rotation')),
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
INSERT INTO quality_findings_new (id, connector_id, doc_id, check_type, severity, title, description,
    remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at, notified_severity)
SELECT id, connector_id, doc_id, check_type, severity, title, description,
    remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at, NULL
FROM quality_findings;
DROP TABLE quality_findings;
ALTER TABLE quality_findings_new RENAME TO quality_findings;

CREATE UNIQUE INDEX idx_quality_findings_open ON quality_findings(connector_id, check_type) WHERE status = 'open';
CREATE INDEX idx_quality_findings_status ON quality_findings(status, last_seen_at DESC);
CREATE INDEX idx_quality_findings_connector ON quality_findings(connector_id);
