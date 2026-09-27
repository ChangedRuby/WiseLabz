CREATE TABLE runbook_steps (
    id           TEXT PRIMARY KEY,
    runbook_id   TEXT NOT NULL REFERENCES runbooks(id) ON DELETE CASCADE,
    position     INTEGER NOT NULL,
    title        TEXT NOT NULL,
    connector_id TEXT NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    verb         TEXT NOT NULL CHECK(verb IN ('restart','start','stop')),
    entity_ref   TEXT NOT NULL DEFAULT '',
    created_at   TEXT NOT NULL,
    updated_at   TEXT NOT NULL
);
CREATE INDEX idx_runbook_steps_runbook ON runbook_steps(runbook_id, position);
CREATE INDEX idx_runbook_steps_connector ON runbook_steps(connector_id);
