-- +goose Up
CREATE TABLE planner_sessions (
    id UUID PRIMARY KEY,
    day_id UUID NOT NULL REFERENCES planner_days(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    starts_at TIME NOT NULL,
    ends_at TIME,
    pic_name TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX planner_sessions_day_idx ON planner_sessions (day_id);

-- +goose Down
DROP INDEX IF EXISTS planner_sessions_day_idx;
DROP TABLE IF EXISTS planner_sessions;
