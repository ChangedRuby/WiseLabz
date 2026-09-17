-- 000025_share_links.up.sql — read-only doc share links (#240 PR2).
-- Same shape as api_keys: a hashed opaque token, never the raw secret.
-- doc_tree_root identifies the shared subtree using the same node IDs the
-- doc tree API already returns ("root" for the whole lab, a connector ID,
-- or a doc ID) — resolved at creation and at view time by the docs handler,
-- not by the store.

CREATE TABLE share_links (
    id                TEXT PRIMARY KEY,
    token_hash        TEXT NOT NULL,
    doc_tree_root     TEXT NOT NULL,
    created_by        TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at        TEXT NOT NULL,
    expires_at        TEXT NOT NULL DEFAULT '',
    revoked_at        TEXT NOT NULL DEFAULT '',
    last_accessed_at  TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_share_links_token_hash ON share_links(token_hash);
CREATE INDEX idx_share_links_created_by ON share_links(created_by);
