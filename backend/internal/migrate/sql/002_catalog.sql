-- +goose Up
CREATE TABLE event_types (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    invitation_schema JSONB NOT NULL DEFAULT '{}'::jsonb,
    planner_categories JSONB NOT NULL DEFAULT '[]'::jsonb
);

CREATE TABLE invitation_templates (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    type_id UUID NOT NULL REFERENCES event_types(id),
    name TEXT NOT NULL,
    tier TEXT NOT NULL CHECK (tier IN ('free', 'premium')),
    version INT NOT NULL DEFAULT 1,
    slots JSONB NOT NULL DEFAULT '[]'::jsonb
);

-- +goose Down
DROP TABLE IF EXISTS invitation_templates;
DROP TABLE IF EXISTS event_types;
