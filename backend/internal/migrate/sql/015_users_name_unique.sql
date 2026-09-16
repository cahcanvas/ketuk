-- +goose Up
CREATE UNIQUE INDEX users_name_lower_idx ON users (lower(name));

-- +goose Down
DROP INDEX IF EXISTS users_name_lower_idx;
