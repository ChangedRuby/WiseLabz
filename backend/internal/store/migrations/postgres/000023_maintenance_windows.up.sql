-- 000023_maintenance_windows.up.sql — per-connector maintenance windows that
-- suppress scheduled sync and drift alerting/change-record creation while
-- active (#236).

CREATE TABLE maintenance_windows (
    id TEXT PRIMARY KEY,
    connector_id TEXT NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    starts_at TEXT NOT NULL,
    ends_at TEXT NOT NULL,
    created_by TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE INDEX idx_maintenance_windows_connector_active ON maintenance_windows(connector_id, ends_at);
