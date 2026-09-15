-- 000022_dns_connector_category.down.sql
-- See the up migration for why the replacement table is built under a
-- fresh name and renamed into place rather than renaming "connectors" away.

DROP INDEX IF EXISTS idx_connectors_category;

CREATE TABLE connectors_new (
    id              TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    category        TEXT NOT NULL CHECK(category IN ('virtualization','containers_paas','networking')),
    type            TEXT NOT NULL,
    url             TEXT NOT NULL,
    owner           TEXT,
    verify_tls      INTEGER NOT NULL DEFAULT 1 CHECK(verify_tls IN (0,1)),
    config_data     TEXT NOT NULL DEFAULT '{}',
    enabled         INTEGER NOT NULL DEFAULT 1 CHECK(enabled IN (0,1)),
    status          TEXT NOT NULL DEFAULT 'unknown' CHECK(status IN ('online','degraded','offline','unknown')),
    status_message  TEXT NOT NULL DEFAULT '',
    last_sync_at    TEXT,
    schedule_seconds       INTEGER,
    next_run_at            TEXT,
    last_sync_duration_ms  INTEGER,
    last_sync_error        TEXT NOT NULL DEFAULT '',
    retry_count            INTEGER NOT NULL DEFAULT 0,
    credential_expires_at  TEXT,
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL
);

INSERT INTO connectors_new (id, name, category, type, url, owner, verify_tls, config_data, enabled, status, status_message,
    last_sync_at, schedule_seconds, next_run_at, last_sync_duration_ms, last_sync_error, retry_count, credential_expires_at, created_at, updated_at)
SELECT id, name, category, type, url, owner, verify_tls, config_data, enabled, status, status_message,
    last_sync_at, schedule_seconds, next_run_at, last_sync_duration_ms, last_sync_error, retry_count, credential_expires_at, created_at, updated_at
FROM connectors;

DROP TABLE connectors;

ALTER TABLE connectors_new RENAME TO connectors;

CREATE INDEX idx_connectors_category ON connectors(category);
