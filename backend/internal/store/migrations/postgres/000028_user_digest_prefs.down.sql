-- 000028_user_digest_prefs.down.sql

ALTER TABLE users DROP COLUMN digest_timezone;
ALTER TABLE users DROP COLUMN digest_last_sent_at;
ALTER TABLE users DROP COLUMN digest_cadence;
