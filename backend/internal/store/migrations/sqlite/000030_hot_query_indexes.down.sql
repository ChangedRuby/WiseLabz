-- 000030_hot_query_indexes.down.sql
DROP INDEX IF EXISTS idx_sync_runs_status_started;
DROP INDEX IF EXISTS idx_sync_runs_started;
DROP INDEX IF EXISTS idx_alerts_status_created;
DROP INDEX IF EXISTS idx_alerts_status_snoozed;
DROP INDEX IF EXISTS idx_changes_status;
DROP INDEX IF EXISTS idx_doc_versions_created;
DROP INDEX IF EXISTS idx_audit_log_action_created;
DROP INDEX IF EXISTS idx_audit_log_target_type_created;
DROP INDEX IF EXISTS idx_deliveries_created;
DROP INDEX IF EXISTS idx_inapp_user_read_created;
