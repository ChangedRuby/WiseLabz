-- Provider fallback routing: an ordered list of extra providers to try when
-- the primary (ai_config.provider) errors with 429/5xx/timeout. Priority 1 is
-- implicitly the primary provider already stored on ai_config; rows here
-- start at priority 2.

CREATE TABLE ai_config_providers (
    id                TEXT PRIMARY KEY,
    config_id         INTEGER NOT NULL DEFAULT 1 REFERENCES ai_config(id) ON DELETE CASCADE,
    priority          INTEGER NOT NULL,
    provider          TEXT NOT NULL,
    model             TEXT NOT NULL DEFAULT '',
    api_key_encrypted TEXT NOT NULL DEFAULT '',
    base_url          TEXT NOT NULL DEFAULT '',
    UNIQUE (config_id, priority)
);
