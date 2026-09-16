-- +goose Up
CREATE TABLE resource_links (
    id UUID PRIMARY KEY,
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    planner_id UUID NOT NULL REFERENCES planners(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (invitation_id, planner_id)
);

CREATE INDEX resource_links_planner_idx ON resource_links (planner_id);

-- +goose Down
DROP TABLE IF EXISTS resource_links;
