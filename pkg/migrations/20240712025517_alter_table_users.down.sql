DROP TABLE IF EXISTS users;
ALTER TABLE users
DROP CONSTRAINT IF EXISTS unique_email;