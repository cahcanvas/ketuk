-- +goose Up
ALTER TABLE users ADD COLUMN username TEXT;
UPDATE users SET username = name;
ALTER TABLE users ALTER COLUMN username SET NOT NULL;
CREATE UNIQUE INDEX users_username_lower_idx ON users (lower(username));

DROP INDEX IF EXISTS users_name_lower_idx;
ALTER TABLE users ALTER COLUMN name DROP NOT NULL;
UPDATE users SET name = NULL;

-- +goose Down
DROP INDEX IF EXISTS users_username_lower_idx;
ALTER TABLE users DROP COLUMN IF EXISTS username;
