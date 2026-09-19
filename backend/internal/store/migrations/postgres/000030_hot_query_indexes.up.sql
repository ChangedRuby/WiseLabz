-- 000030_hot_query_indexes.up.sql — indexes for hot queries and retention scans (#299).
CREATE INDEX IF NOT EXISTS idx_sync_runs_status_started ON sync_runs(status, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_sync_runs_started ON sync_runs(started_at);
CREATE INDEX IF NOT EXISTS idx_alerts_status_created ON alerts(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_alerts_status_snoozed ON alerts(status, snoozed_until);
CREATE INDEX IF NOT EXISTS idx_changes_status ON changes(status);
CREATE INDEX IF NOT EXISTS idx_doc_versions_created ON doc_versions(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_log_action_created ON audit_log(action, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_log_target_type_created ON audit_log(target_type, created_at);
CREATE INDEX IF NOT EXISTS idx_deliveries_created ON notification_deliveries(created_at);
CREATE INDEX IF NOT EXISTS idx_inapp_user_read_created ON in_app_notifications(user_id, read, created_at DESC);
