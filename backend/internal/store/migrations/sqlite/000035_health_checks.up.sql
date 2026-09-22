-- 000035_health_checks.up.sql — persists each connector health check
-- (connector/health.go's Validate-only check) as a time-series row, so
-- availability % and MTTR can be computed over a window instead of only
-- exposing the latest status snapshot (#281).

CREATE TABLE health_checks (
    id              TEXT PRIMARY KEY,
    connector_id    TEXT NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    status          TEXT NOT NULL CHECK(status IN ('online','degraded','offline')),
    message         TEXT NOT NULL DEFAULT '',
    latency_ms      INTEGER,
    checked_at      TEXT NOT NULL
);
CREATE INDEX idx_health_checks_connector_checked_at ON health_checks(connector_id, checked_at DESC);
