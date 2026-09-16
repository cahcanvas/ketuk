-- +goose Up
CREATE TABLE invitation_inviters (
    id UUID PRIMARY KEY,
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX invitation_inviters_invitation_idx ON invitation_inviters (invitation_id);

CREATE TABLE invitation_media (
    id UUID PRIMARY KEY,
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (kind IN ('gallery', 'music')),
    original_name TEXT NOT NULL DEFAULT '',
    content_type TEXT NOT NULL DEFAULT '',
    rel_path TEXT NOT NULL,
    caption TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX invitation_media_invitation_idx ON invitation_media (invitation_id, kind);

CREATE UNIQUE INDEX invitation_media_one_music_idx
    ON invitation_media (invitation_id)
    WHERE kind = 'music';

-- +goose Down
DROP TABLE IF EXISTS invitation_media;
DROP TABLE IF EXISTS invitation_inviters;
