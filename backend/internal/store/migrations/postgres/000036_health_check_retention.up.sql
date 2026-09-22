-- 000036_health_check_retention.up.sql — adds a configurable retention
-- window for the health_checks time series (#281), following the same
-- *_days pattern as the other retention_settings columns.

ALTER TABLE retention_settings ADD COLUMN health_check_days INTEGER NOT NULL DEFAULT 90;
