-- +goose Up
INSERT INTO plans (id, slug, name, price_idr, billing_interval) VALUES
    ('00000000-0000-0000-0000-000000000001', 'free', 'Free', 0, 'month'),
    ('00000000-0000-0000-0000-000000000002', 'starter', 'Starter', 79000, 'month'),
    ('00000000-0000-0000-0000-000000000003', 'pro', 'Pro', 199000, 'month');

INSERT INTO plan_entitlements (plan_id, key, value) VALUES
    ('00000000-0000-0000-0000-000000000001', 'max_active_invitations', '1'),
    ('00000000-0000-0000-0000-000000000001', 'max_guests_per_invitation', '50'),
    ('00000000-0000-0000-0000-000000000001', 'premium_templates', 'false'),
    ('00000000-0000-0000-0000-000000000001', 'planner_daily', 'true'),
    ('00000000-0000-0000-0000-000000000002', 'max_active_invitations', '5'),
    ('00000000-0000-0000-0000-000000000002', 'max_guests_per_invitation', '300'),
    ('00000000-0000-0000-0000-000000000002', 'premium_templates', 'true'),
    ('00000000-0000-0000-0000-000000000002', 'planner_daily', 'true'),
    ('00000000-0000-0000-0000-000000000003', 'max_active_invitations', '-1'),
    ('00000000-0000-0000-0000-000000000003', 'max_guests_per_invitation', '2000'),
    ('00000000-0000-0000-0000-000000000003', 'premium_templates', 'true'),
    ('00000000-0000-0000-0000-000000000003', 'planner_daily', 'true');

INSERT INTO event_types (id, slug, name, invitation_schema, planner_categories) VALUES
    ('10000000-0000-0000-0000-000000000001', 'wedding', 'Pernikahan', '{"slots":["couple","rundown","dress_code"]}', '["Venue","Catering","MUA","Dokumentasi","Souvenir"]'),
    ('10000000-0000-0000-0000-000000000002', 'engagement', 'Lamaran', '{"slots":["couple","story"]}', '["Dekor","Dokumentasi","Seserahan"]'),
    ('10000000-0000-0000-0000-000000000003', 'birthday', 'Ulang tahun', '{"slots":["honoree","age"]}', '["Venue","Cake","MC","Souvenir"]'),
    ('10000000-0000-0000-0000-000000000004', 'aqiqah', 'Aqiqah', '{"slots":["child","parents"]}', '["Katering","Dokumentasi"]'),
    ('10000000-0000-0000-0000-000000000005', 'graduation', 'Wisuda', '{"slots":["graduate"]}', '["Venue","Dokumentasi","Gift"]'),
    ('10000000-0000-0000-0000-000000000006', 'corporate', 'Corporate', '{"slots":["company","agenda"]}', '["Venue","F&B","AV","Talent"]'),
    ('10000000-0000-0000-0000-000000000007', 'custom', 'Lainnya', '{"slots":["title","notes"]}', '[]');

INSERT INTO invitation_templates (id, slug, type_id, name, tier, version, slots) VALUES
    ('20000000-0000-0000-0000-000000000001', 'wedding-classic', '10000000-0000-0000-0000-000000000001', 'Wedding Classic', 'free', 1, '["couple","date","venue"]'),
    ('20000000-0000-0000-0000-000000000002', 'wedding-floral', '10000000-0000-0000-0000-000000000001', 'Wedding Floral', 'premium', 1, '["couple","date","venue","gallery"]'),
    ('20000000-0000-0000-0000-000000000003', 'birthday-fun', '10000000-0000-0000-0000-000000000003', 'Birthday Fun', 'free', 1, '["honoree","age","date"]'),
    ('20000000-0000-0000-0000-000000000004', 'corporate-clean', '10000000-0000-0000-0000-000000000006', 'Corporate Clean', 'free', 1, '["company","agenda"]'),
    ('20000000-0000-0000-0000-000000000005', 'custom-minimal', '10000000-0000-0000-0000-000000000007', 'Custom Minimal', 'free', 1, '["title","notes"]');

-- +goose Down
DELETE FROM invitation_templates;
DELETE FROM event_types;
DELETE FROM plan_entitlements;
DELETE FROM plans;
