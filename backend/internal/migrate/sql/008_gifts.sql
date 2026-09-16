-- +goose Up
CREATE TABLE gift_methods (
    id UUID PRIMARY KEY,
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (kind IN ('bank', 'ewallet', 'duitku')),
    label TEXT NOT NULL,
    account_name TEXT NOT NULL DEFAULT '',
    account_number TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX gift_methods_invitation_idx ON gift_methods (invitation_id);

CREATE TABLE gift_payments (
    id UUID PRIMARY KEY,
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    method_id UUID NOT NULL REFERENCES gift_methods(id),
    guest_id UUID REFERENCES invitation_guests(id) ON DELETE SET NULL,
    sender_name TEXT NOT NULL DEFAULT '',
    message TEXT NOT NULL DEFAULT '',
    amount_idr BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'failed', 'expired')),
    merchant_order_id TEXT NOT NULL UNIQUE,
    duitku_reference TEXT NOT NULL DEFAULT '',
    payment_url TEXT NOT NULL DEFAULT '',
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX gift_payments_invitation_idx ON gift_payments (invitation_id);

-- +goose Down
DROP TABLE IF EXISTS gift_payments;
DROP TABLE IF EXISTS gift_methods;
