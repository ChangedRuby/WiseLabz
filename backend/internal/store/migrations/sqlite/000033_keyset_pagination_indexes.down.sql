-- 000033_keyset_pagination_indexes.down.sql
DROP INDEX IF EXISTS idx_audit_log_created_id;
DROP INDEX IF EXISTS idx_audit_log_action_created_id;
DROP INDEX IF EXISTS idx_audit_log_target_type_created_id;
DROP INDEX IF EXISTS idx_changes_detected_id;
DROP INDEX IF EXISTS idx_changes_service_detected_id;
DROP INDEX IF EXISTS idx_changes_severity_detected_id;
DROP INDEX IF EXISTS idx_sync_runs_connector_started_id;
