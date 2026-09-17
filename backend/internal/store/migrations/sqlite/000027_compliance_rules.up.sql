-- User-defined snapshot compliance rules (#239 PR3).
CREATE TABLE compliance_rules (
    id               TEXT PRIMARY KEY,
    name             TEXT NOT NULL,
    connector_type   TEXT NOT NULL,
    entity_kind      TEXT NOT NULL,
    conditions       TEXT NOT NULL,
    severity         TEXT NOT NULL CHECK(severity IN ('info','warning','critical')),
    title            TEXT NOT NULL,
    remediation_link TEXT NOT NULL DEFAULT '',
    enabled          INTEGER NOT NULL DEFAULT 0,
    created_at       TEXT NOT NULL,
    updated_at       TEXT NOT NULL
);

-- Disabled examples are discoverable without producing findings.
INSERT INTO compliance_rules (id, name, connector_type, entity_kind, conditions, severity, title, enabled, created_at, updated_at) VALUES
    ('seed-pfsense-any-any', 'pfSense allows any-to-any traffic', 'pfsense', 'rule', '[{"attribute":"action","op":"eq","value":"pass"},{"attribute":"source","op":"eq","value":"any"},{"attribute":"destination","op":"eq","value":"any"}]', 'critical', 'Firewall rule allows any-to-any traffic', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('seed-opnsense-unlogged-pass', 'OPNsense pass rule without logging', 'opnsense', 'rule', '[{"attribute":"enabled","op":"eq","value":true},{"attribute":"action","op":"eq","value":"pass"},{"attribute":"log","op":"eq","value":false}]', 'warning', 'Enabled firewall pass rule is not logged', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('seed-proxmox-firewall-disabled', 'Proxmox VM firewall disabled', 'proxmox', 'vm', '[{"attribute":"firewall_enabled","op":"eq","value":false}]', 'warning', 'Guest firewall is disabled', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('seed-docker-privileged', 'Docker privileged container', 'docker', 'container', '[{"attribute":"privileged","op":"eq","value":true}]', 'critical', 'Container runs in privileged mode', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('seed-docker-host-network', 'Docker host network mode', 'docker', 'container', '[{"attribute":"network_mode","op":"eq","value":"host"}]', 'warning', 'Container uses host network mode', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- SQLite rebuild is required to change the check constraint and add rule_id.
DROP INDEX IF EXISTS idx_quality_findings_open;
DROP INDEX IF EXISTS idx_quality_findings_status;
DROP INDEX IF EXISTS idx_quality_findings_connector;

CREATE TABLE quality_findings_new (
    id                  TEXT PRIMARY KEY,
    connector_id        TEXT NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    doc_id              TEXT REFERENCES docs(id) ON DELETE SET NULL,
    -- Preserve resolved finding history when a rule is deleted; callers resolve
    -- open findings first, while ON DELETE SET NULL keeps audit-visible rows.
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
INSERT INTO quality_findings_new (id, connector_id, doc_id, rule_id, check_type, severity, title, description,
    remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at, notified_severity)
SELECT id, connector_id, doc_id, NULL, check_type, severity, title, description,
    remediation_link, status, detected_count, first_detected_at, last_seen_at, resolved_at, notified_severity
FROM quality_findings;
DROP TABLE quality_findings;
ALTER TABLE quality_findings_new RENAME TO quality_findings;

CREATE UNIQUE INDEX idx_quality_findings_open ON quality_findings(connector_id, check_type, COALESCE(rule_id, '')) WHERE status = 'open';
CREATE INDEX idx_quality_findings_status ON quality_findings(status, last_seen_at DESC);
CREATE INDEX idx_quality_findings_connector ON quality_findings(connector_id);
