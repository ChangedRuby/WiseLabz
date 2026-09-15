DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_conversations;
DROP TABLE IF EXISTS doc_section_embeddings;

ALTER TABLE ai_config DROP COLUMN embed_provider;
ALTER TABLE ai_config DROP COLUMN embed_model;
ALTER TABLE ai_config DROP COLUMN embed_api_key_encrypted;
ALTER TABLE ai_config DROP COLUMN embed_base_url;
