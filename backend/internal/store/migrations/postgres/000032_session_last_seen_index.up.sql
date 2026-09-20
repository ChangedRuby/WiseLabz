-- 000032_session_last_seen_index.up.sql — index the stale-session retention sweep (#299 follow-up).
CREATE INDEX IF NOT EXISTS idx_sessions_last_seen ON sessions(last_seen_at);
