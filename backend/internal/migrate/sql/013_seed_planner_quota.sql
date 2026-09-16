-- +goose Up
INSERT INTO plan_entitlements (plan_id, key, value) VALUES
    ('00000000-0000-0000-0000-000000000001', 'max_active_planners', '1'),
    ('00000000-0000-0000-0000-000000000002', 'max_active_planners', '5'),
    ('00000000-0000-0000-0000-000000000003', 'max_active_planners', '-1')
ON CONFLICT (plan_id, key) DO NOTHING;

-- +goose Down
DELETE FROM plan_entitlements WHERE key = 'max_active_planners';
