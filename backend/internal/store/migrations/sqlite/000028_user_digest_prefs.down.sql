-- 000028_user_digest_prefs.down.sql
-- SQLite's ALTER TABLE DROP COLUMN (3.35+) is used rather than the
-- rebuild-and-rename dance from earlier migrations, since these are plain
-- column drops with no constraint/rename semantics to preserve.

ALTER TABLE users DROP COLUMN digest_timezone;
ALTER TABLE users DROP COLUMN digest_last_sent_at;
ALTER TABLE users DROP COLUMN digest_cadence;
