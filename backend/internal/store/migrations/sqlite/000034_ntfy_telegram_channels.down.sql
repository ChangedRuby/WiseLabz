-- Revert to supporting only in_app, smtp, webhook, discord, slack.
-- Note: rows with channel = 'ntfy' or 'telegram' will violate the new constraint.
-- This is expected for a down migration.

DROP INDEX IF EXISTS idx_deliveries_notification;
DROP INDEX IF EXISTS idx_deliveries_retry;

CREATE TABLE notification_deliveries_new (
    id              TEXT PRIMARY KEY,
    notification_id TEXT NOT NULL REFERENCES in_app_notifications(id) ON DELETE CASCADE,
    channel         TEXT NOT NULL CHECK(channel IN ('in_app','smtp','webhook','discord','slack')),
    status          TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','sent','failed')),
    attempts        INTEGER NOT NULL DEFAULT 0,
    last_error      TEXT NOT NULL DEFAULT '',
    next_attempt_at TEXT,
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL
);
INSERT INTO notification_deliveries_new SELECT * FROM notification_deliveries WHERE channel IN ('in_app','smtp','webhook','discord','slack');
DROP TABLE notification_deliveries;
ALTER TABLE notification_deliveries_new RENAME TO notification_deliveries;

CREATE INDEX idx_deliveries_notification ON notification_deliveries(notification_id);
CREATE INDEX idx_deliveries_retry ON notification_deliveries(status, next_attempt_at);
