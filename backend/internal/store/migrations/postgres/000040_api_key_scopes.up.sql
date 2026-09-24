-- 000040_api_key_scopes.up.sql — API-key scopes and per-connector
-- restriction (#278). scope 'read' limits a key to safe HTTP methods and
-- caps its connector role at viewer; connector_ids (a JSON array) limits
-- the connectors a key can reach. The defaults keep existing keys
-- unrestricted.

ALTER TABLE api_keys ADD COLUMN scope TEXT NOT NULL DEFAULT 'full' CHECK(scope IN ('full','read'));
ALTER TABLE api_keys ADD COLUMN connector_ids TEXT NOT NULL DEFAULT '[]';
