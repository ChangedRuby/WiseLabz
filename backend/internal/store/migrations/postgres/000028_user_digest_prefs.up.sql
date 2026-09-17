-- 000028_user_digest_prefs.up.sql — per-user notification digest
-- preferences (#237 Phase 4).

ALTER TABLE users ADD COLUMN digest_cadence TEXT NOT NULL DEFAULT 'off' CHECK (digest_cadence IN ('off','daily','weekly'));
ALTER TABLE users ADD COLUMN digest_last_sent_at TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN digest_timezone TEXT NOT NULL DEFAULT 'UTC';
