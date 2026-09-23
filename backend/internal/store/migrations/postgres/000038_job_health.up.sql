-- 000038_job_health.up.sql — persists scheduled-job health across restarts
-- (#384). Previously each job tracked its own ok/failing state in memory
-- (e.g. docexport.Exporter.failing), so a restart forgot a failing job had
-- already notified and would notify again on the next failure. The
-- scheduler now owns one row per job name here.

CREATE TABLE job_health (
    name            TEXT PRIMARY KEY,
    cron_expr       TEXT NOT NULL,
    last_status     TEXT NOT NULL CHECK(last_status IN ('ok','failing')),
    last_error      TEXT NOT NULL DEFAULT '',
    last_run_at     TEXT,
    last_success_at TEXT,
    last_failure_at TEXT,
    updated_at      TEXT NOT NULL
);
