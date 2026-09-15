-- "Ask your lab" chat: per-section doc embeddings for retrieval, plus
-- persisted chat conversations/messages. Separate embed_* config lets the
-- embedding backend (e.g. local Ollama) differ from the chat provider.

ALTER TABLE ai_config ADD COLUMN embed_provider TEXT;
ALTER TABLE ai_config ADD COLUMN embed_model TEXT;
ALTER TABLE ai_config ADD COLUMN embed_api_key_encrypted TEXT NOT NULL DEFAULT '';
ALTER TABLE ai_config ADD COLUMN embed_base_url TEXT;

CREATE TABLE doc_section_embeddings (
    doc_id      TEXT NOT NULL REFERENCES docs(id) ON DELETE CASCADE,
    section_key TEXT NOT NULL,
    content     TEXT NOT NULL,
    vector      BYTEA NOT NULL,
    model       TEXT NOT NULL,
    updated_at  TEXT NOT NULL,
    PRIMARY KEY (doc_id, section_key)
);

CREATE TABLE chat_conversations (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scope_type TEXT NOT NULL CHECK(scope_type IN ('doc','lab')),
    scope_id   TEXT,
    created_at TEXT NOT NULL
);

CREATE TABLE chat_messages (
    id              TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL REFERENCES chat_conversations(id) ON DELETE CASCADE,
    role            TEXT NOT NULL CHECK(role IN ('user','assistant')),
    content         TEXT NOT NULL,
    provider        TEXT NOT NULL DEFAULT '',
    created_at      TEXT NOT NULL
);

CREATE INDEX idx_chat_conversations_user ON chat_conversations(user_id, created_at);
CREATE INDEX idx_chat_messages_conversation ON chat_messages(conversation_id, created_at);
