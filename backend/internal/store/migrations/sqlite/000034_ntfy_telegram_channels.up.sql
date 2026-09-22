-- Add ntfy and telegram as valid notification_deliveries channels.
-- SQLite doesn't support ALTER TABLE ... ALTER CONSTRAINT, so we rebuild the table.

DROP INDEX IF EXISTS idx_deliveries_notification;
DROP INDEX IF EXISTS idx_deliveries_retry;

CREATE TABLE notification_deliveries_new (
    id              TEXT PRIMARY KEY,
    notification_id TEXT NOT NULL REFERENCES in_app_notifications(id) ON DELETE CASCADE,
    channel         TEXT NOT NULL CHECK(channel IN ('in_app','smtp','webhook','discord','slack','ntfy','telegram')),
    status          TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','sent','failed')),
    attempts        INTEGER NOT NULL DEFAULT 0,
    last_error      TEXT NOT NULL DEFAULT '',
    next_attempt_at TEXT,
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL
);
INSERT INTO notification_deliveries_new SELECT * FROM notification_deliveries;
DROP TABLE notification_deliveries;
ALTER TABLE notification_deliveries_new RENAME TO notification_deliveries;

CREATE INDEX idx_deliveries_notification ON notification_deliveries(notification_id);
CREATE INDEX idx_deliveries_retry ON notification_deliveries(status, next_attempt_at);
