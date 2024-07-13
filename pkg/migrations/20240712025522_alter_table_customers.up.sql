CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    username VARCHAR(125) NULL,
    email VARCHAR(125) NULL,
    phone VARCHAR(125) NULL,
    address VARCHAR(125) NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NULL
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'unique_email'
          AND connamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public')
    ) THEN
ALTER TABLE customers
    ADD CONSTRAINT unique_email UNIQUE (email);
END IF;
END $$;
