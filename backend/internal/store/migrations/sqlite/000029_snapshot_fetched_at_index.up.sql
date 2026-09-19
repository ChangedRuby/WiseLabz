-- 000029_snapshot_fetched_at_index.up.sql — supports the retention job's
-- fetched_at range scan over service_snapshots (#290).
CREATE INDEX IF NOT EXISTS idx_snapshots_fetched_at ON service_snapshots(fetched_at);
