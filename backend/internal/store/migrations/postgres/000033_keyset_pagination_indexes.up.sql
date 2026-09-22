-- 000033_keyset_pagination_indexes.up.sql — (sort_key, id) indexes for the
-- opt-in keyset (cursor) pagination on audit records, changes and sync runs
-- (#267). Without the trailing id column the existing single-column indexes
-- still satisfy the range seek, but the planner has to sort each page for the
-- id tie-break — exactly the per-page sort keyset pagination exists to avoid.
CREATE INDEX IF NOT EXISTS idx_audit_log_created_id ON audit_log(created_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_audit_log_action_created_id ON audit_log(action, created_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_audit_log_target_type_created_id ON audit_log(target_type, created_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_changes_detected_id ON changes(detected_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_changes_service_detected_id ON changes(service_id, detected_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_changes_severity_detected_id ON changes(severity, detected_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_sync_runs_connector_started_id ON sync_runs(connector_id, started_at DESC, id DESC);
