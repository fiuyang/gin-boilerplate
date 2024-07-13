CREATE TABLE IF NOT EXISTS password_resets (
    id SERIAL PRIMARY KEY,
    email VARCHAR(125) NULL,
    otp BIGINT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_users_email'
          AND connamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public')
    ) THEN
        -- Add the foreign key constraint if it doesn't exist
ALTER TABLE password_resets
    ADD CONSTRAINT fk_users_email
        FOREIGN KEY (email)
            REFERENCES users (email)
            ON UPDATE CASCADE
            ON DELETE CASCADE;
END IF;
END $$;
--
-- ALTER TABLE IF EXISTS password_resets
-- ADD CONSTRAINT IF EXISTS fk_users_email
-- FOREIGN KEY (email)
-- REFERENCES users (email)
-- ON UPDATE CASCADE
-- ON DELETE CASCADE;
