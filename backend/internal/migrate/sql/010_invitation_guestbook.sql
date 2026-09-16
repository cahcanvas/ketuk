-- +goose Up
CREATE TABLE invitation_guestbook (
    id UUID PRIMARY KEY,
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    guest_id UUID REFERENCES invitation_guests(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX invitation_guestbook_invitation_idx ON invitation_guestbook (invitation_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS invitation_guestbook;
