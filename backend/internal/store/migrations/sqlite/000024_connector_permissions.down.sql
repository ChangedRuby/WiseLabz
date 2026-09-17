-- 000024_connector_permissions.down.sql

DROP TABLE user_connector_roles;

CREATE TABLE users_old (
    id              TEXT PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    display_name    TEXT NOT NULL DEFAULT '',
    email           TEXT NOT NULL DEFAULT '',
    role            TEXT NOT NULL DEFAULT 'viewer' CHECK(role IN ('viewer','operator')),
    auth_source     TEXT NOT NULL DEFAULT 'local' CHECK(auth_source IN ('local','oidc')),
    password_hash   TEXT NOT NULL DEFAULT '',
    disabled        INTEGER NOT NULL DEFAULT 0 CHECK(disabled IN (0,1)),
    can_manage_dashboard_defaults INTEGER NOT NULL DEFAULT 0 CHECK(can_manage_dashboard_defaults IN (0,1)),
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until    TEXT NOT NULL DEFAULT '',
    created_at      TEXT NOT NULL
);

INSERT INTO users_old (id, username, display_name, email, role, auth_source, password_hash,
    disabled, can_manage_dashboard_defaults, failed_login_attempts, locked_until, created_at)
SELECT id, username, display_name, email,
       CASE WHEN instance_admin_role = 'admin' THEN 'operator' ELSE 'viewer' END,
       auth_source, password_hash, disabled, can_manage_dashboard_defaults, failed_login_attempts, locked_until, created_at
FROM users;

DROP TABLE users;
ALTER TABLE users_old RENAME TO users;
