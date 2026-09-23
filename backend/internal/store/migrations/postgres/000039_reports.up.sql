CREATE TABLE report_definitions (
    id TEXT PRIMARY KEY, slug TEXT NOT NULL UNIQUE, name TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE, cron_expr TEXT NOT NULL, timezone TEXT NOT NULL,
    sections JSONB NOT NULL, connector_ids JSONB NOT NULL DEFAULT '[]', channels JSONB NOT NULL DEFAULT '[]',
    created_by TEXT NOT NULL, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
);
CREATE TABLE reports (
    id TEXT PRIMARY KEY, definition_id TEXT REFERENCES report_definitions(id) ON DELETE SET NULL,
    definition_name TEXT NOT NULL, trigger TEXT NOT NULL CHECK(trigger IN ('scheduled','manual')),
    period_start TEXT NOT NULL, period_end TEXT NOT NULL, truncated BOOLEAN NOT NULL DEFAULT FALSE,
    data JSONB NOT NULL, markdown TEXT NOT NULL, status TEXT NOT NULL CHECK(status IN ('ok','partial')), created_at TEXT NOT NULL
);
CREATE INDEX reports_definition_trigger_period_idx ON reports(definition_id, trigger, period_end DESC);
CREATE INDEX reports_created_at_idx ON reports(created_at);
ALTER TABLE retention_settings ADD COLUMN report_days INTEGER NOT NULL DEFAULT 90;
