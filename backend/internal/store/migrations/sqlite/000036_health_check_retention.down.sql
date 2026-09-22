-- 000036_health_check_retention.down.sql
ALTER TABLE retention_settings DROP COLUMN health_check_days;
