-- +goose Up
CREATE TABLE planners (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type_id UUID REFERENCES event_types(id),
    title TEXT NOT NULL,
    budget_mode TEXT NOT NULL DEFAULT 'event' CHECK (budget_mode IN ('event', 'daily')),
    currency TEXT NOT NULL DEFAULT 'IDR',
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'archived')),
    starts_at DATE,
    ends_at DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX planners_owner_idx ON planners (owner_id);

CREATE TABLE planner_days (
    id UUID PRIMARY KEY,
    planner_id UUID NOT NULL REFERENCES planners(id) ON DELETE CASCADE,
    day_date DATE NOT NULL,
    label TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    UNIQUE (planner_id, day_date)
);

CREATE TABLE budget_categories (
    id UUID PRIMARY KEY,
    planner_id UUID NOT NULL REFERENCES planners(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE budget_items (
    id UUID PRIMARY KEY,
    category_id UUID NOT NULL REFERENCES budget_categories(id) ON DELETE CASCADE,
    day_id UUID REFERENCES planner_days(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0,
    paid_at TIMESTAMPTZ,
    note TEXT NOT NULL DEFAULT ''
);

CREATE INDEX budget_items_category_idx ON budget_items (category_id);

CREATE TABLE checklist_items (
    id UUID PRIMARY KEY,
    planner_id UUID NOT NULL REFERENCES planners(id) ON DELETE CASCADE,
    day_id UUID REFERENCES planner_days(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    done BOOLEAN NOT NULL DEFAULT false,
    due_at TIMESTAMPTZ,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX checklist_items_planner_idx ON checklist_items (planner_id);

-- +goose Down
DROP TABLE IF EXISTS checklist_items;
DROP TABLE IF EXISTS budget_items;
DROP TABLE IF EXISTS budget_categories;
DROP TABLE IF EXISTS planner_days;
DROP TABLE IF EXISTS planners;
