DROP TABLE IF EXISTS password_resets;

ALTER TABLE password_resets
DROP CONSTRAINT IF EXISTS fk_users_email;
