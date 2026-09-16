-- +goose Up
ALTER TABLE orders ADD COLUMN kind TEXT NOT NULL DEFAULT 'subscription'
    CHECK (kind IN ('subscription', 'gift', 'vendor'));

CREATE INDEX orders_user_kind_idx ON orders (user_id, kind, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS orders_user_kind_idx;
ALTER TABLE orders DROP COLUMN IF EXISTS kind;
