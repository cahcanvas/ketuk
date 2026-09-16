-- +goose Up
CREATE TABLE invitations (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type_id UUID NOT NULL REFERENCES event_types(id),
    template_id UUID NOT NULL REFERENCES invitation_templates(id),
    template_version INT NOT NULL DEFAULT 1,
    slug TEXT UNIQUE,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    title TEXT NOT NULL,
    content JSONB NOT NULL DEFAULT '{}'::jsonb,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX invitations_owner_idx ON invitations (owner_id);

CREATE TABLE invitation_locations (
    id UUID PRIMARY KEY,
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    label TEXT NOT NULL DEFAULT '',
    venue_name TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    maps_url TEXT NOT NULL DEFAULT '',
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    starts_at TIMESTAMPTZ,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX invitation_locations_invitation_idx ON invitation_locations (invitation_id);

CREATE TABLE invitation_guests (
    id UUID PRIMARY KEY,
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    email TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    guest_group TEXT NOT NULL DEFAULT '',
    rsvp_token TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'invited', 'opened', 'accepted', 'declined')),
    plus_ones INT NOT NULL DEFAULT 0,
    message TEXT NOT NULL DEFAULT '',
    rsvp_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX invitation_guests_invitation_idx ON invitation_guests (invitation_id);

-- +goose Down
DROP TABLE IF EXISTS invitation_guests;
DROP TABLE IF EXISTS invitation_locations;
DROP TABLE IF EXISTS invitations;
