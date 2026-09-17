-- 000024_connector_permissions.up.sql — per-connector permissions (#240 PR1).
--
-- The flat instance-wide users.role (viewer/operator) is replaced by
-- per-connector grants in user_connector_roles. A minimal instance-admin
-- flag survives on users (renamed from role) for actions that aren't
-- connector-scoped: user management, API keys, and granting/revoking
-- connector permissions.
--
-- Existing users keep their current access: every user gets their old role
-- (operator->admin, viewer->user) backfilled as a grant on every existing
-- connector, so nobody loses access on deploy.

ALTER TABLE users RENAME COLUMN role TO instance_admin_role;
UPDATE users SET instance_admin_role = 'admin' WHERE instance_admin_role = 'operator';
UPDATE users SET instance_admin_role = 'user' WHERE instance_admin_role = 'viewer';
ALTER TABLE users DROP CONSTRAINT users_role_check;
ALTER TABLE users ADD CONSTRAINT users_instance_admin_role_check CHECK(instance_admin_role IN ('admin','user'));
ALTER TABLE users ALTER COLUMN instance_admin_role SET DEFAULT 'user';

CREATE TABLE user_connector_roles (
    id              TEXT PRIMARY KEY,
    user_id         TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    connector_id    TEXT NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    role            TEXT NOT NULL CHECK(role IN ('viewer','operator')),
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL,
    UNIQUE(user_id, connector_id)
);
CREATE INDEX idx_user_connector_roles_user ON user_connector_roles(user_id);
CREATE INDEX idx_user_connector_roles_connector ON user_connector_roles(connector_id);

INSERT INTO user_connector_roles (id, user_id, connector_id, role, created_at, updated_at)
SELECT gen_random_uuid()::text, u.id, c.id,
       CASE WHEN u.instance_admin_role = 'admin' THEN 'operator' ELSE 'viewer' END,
       now()::text, now()::text
FROM users u CROSS JOIN connectors c;
