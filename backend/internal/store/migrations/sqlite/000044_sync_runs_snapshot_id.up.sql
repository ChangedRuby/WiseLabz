ALTER TABLE sync_runs ADD COLUMN snapshot_id TEXT REFERENCES service_snapshots(id) ON DELETE SET NULL;
