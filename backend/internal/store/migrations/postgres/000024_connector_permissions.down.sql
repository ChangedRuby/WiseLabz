-- 000024_connector_permissions.down.sql

DROP TABLE user_connector_roles;

ALTER TABLE users ALTER COLUMN instance_admin_role SET DEFAULT 'viewer';
ALTER TABLE users DROP CONSTRAINT users_instance_admin_role_check;
UPDATE users SET instance_admin_role = 'operator' WHERE instance_admin_role = 'admin';
UPDATE users SET instance_admin_role = 'viewer' WHERE instance_admin_role = 'user';
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK(instance_admin_role IN ('viewer','operator'));
ALTER TABLE users RENAME COLUMN instance_admin_role TO role;
