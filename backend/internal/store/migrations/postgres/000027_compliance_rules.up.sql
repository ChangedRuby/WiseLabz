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
    enabled          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TEXT NOT NULL,
    updated_at       TEXT NOT NULL
);

INSERT INTO compliance_rules (id, name, connector_type, entity_kind, conditions, severity, title, enabled, created_at, updated_at) VALUES
    ('seed-pfsense-any-any', 'pfSense allows any-to-any traffic', 'pfsense', 'rule', '[{"attribute":"action","op":"eq","value":"pass"},{"attribute":"source","op":"eq","value":"any"},{"attribute":"destination","op":"eq","value":"any"}]', 'critical', 'Firewall rule allows any-to-any traffic', FALSE, CURRENT_TIMESTAMP::TEXT, CURRENT_TIMESTAMP::TEXT),
    ('seed-opnsense-unlogged-pass', 'OPNsense pass rule without logging', 'opnsense', 'rule', '[{"attribute":"enabled","op":"eq","value":true},{"attribute":"action","op":"eq","value":"pass"},{"attribute":"log","op":"eq","value":false}]', 'warning', 'Enabled firewall pass rule is not logged', FALSE, CURRENT_TIMESTAMP::TEXT, CURRENT_TIMESTAMP::TEXT),
    ('seed-proxmox-firewall-disabled', 'Proxmox VM firewall disabled', 'proxmox', 'vm', '[{"attribute":"firewall_enabled","op":"eq","value":false}]', 'warning', 'Guest firewall is disabled', FALSE, CURRENT_TIMESTAMP::TEXT, CURRENT_TIMESTAMP::TEXT),
    ('seed-docker-privileged', 'Docker privileged container', 'docker', 'container', '[{"attribute":"privileged","op":"eq","value":true}]', 'critical', 'Container runs in privileged mode', FALSE, CURRENT_TIMESTAMP::TEXT, CURRENT_TIMESTAMP::TEXT),
    ('seed-docker-host-network', 'Docker host network mode', 'docker', 'container', '[{"attribute":"network_mode","op":"eq","value":"host"}]', 'warning', 'Container uses host network mode', FALSE, CURRENT_TIMESTAMP::TEXT, CURRENT_TIMESTAMP::TEXT);

ALTER TABLE quality_findings ADD COLUMN rule_id TEXT REFERENCES compliance_rules(id) ON DELETE SET NULL;
ALTER TABLE quality_findings DROP CONSTRAINT quality_findings_check_type_check;
ALTER TABLE quality_findings ADD CONSTRAINT quality_findings_check_type_check
    CHECK(check_type IN ('stale','empty','failing','ownership_incomplete','credential_rotation','compliance'));
DROP INDEX idx_quality_findings_open;
-- A deleted rule keeps its resolved finding history via ON DELETE SET NULL;
-- application code resolves open findings before deletion.
CREATE UNIQUE INDEX idx_quality_findings_open ON quality_findings(connector_id, check_type, COALESCE(rule_id, '')) WHERE status = 'open';
