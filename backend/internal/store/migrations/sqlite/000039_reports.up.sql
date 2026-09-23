CREATE TABLE report_definitions (
    id TEXT PRIMARY KEY, slug TEXT NOT NULL UNIQUE, name TEXT NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 1, cron_expr TEXT NOT NULL, timezone TEXT NOT NULL,
    sections TEXT NOT NULL, connector_ids TEXT NOT NULL DEFAULT '[]', channels TEXT NOT NULL DEFAULT '[]',
    created_by TEXT NOT NULL, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
);
CREATE TABLE reports (
    id TEXT PRIMARY KEY, definition_id TEXT REFERENCES report_definitions(id) ON DELETE SET NULL,
    definition_name TEXT NOT NULL, trigger TEXT NOT NULL CHECK(trigger IN ('scheduled','manual')),
    period_start TEXT NOT NULL, period_end TEXT NOT NULL, truncated INTEGER NOT NULL DEFAULT 0,
    data TEXT NOT NULL, markdown TEXT NOT NULL, status TEXT NOT NULL CHECK(status IN ('ok','partial')), created_at TEXT NOT NULL
);
CREATE INDEX reports_definition_trigger_period_idx ON reports(definition_id, trigger, period_end DESC);
CREATE INDEX reports_created_at_idx ON reports(created_at);
ALTER TABLE retention_settings ADD COLUMN report_days INTEGER NOT NULL DEFAULT 90;
