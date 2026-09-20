-- Support both branches of DeleteExpiredShareLinks expiry/revocation cleanup (#299).
CREATE INDEX IF NOT EXISTS idx_share_links_expires_at ON share_links(expires_at);
CREATE INDEX IF NOT EXISTS idx_share_links_revoked_at ON share_links(revoked_at);
