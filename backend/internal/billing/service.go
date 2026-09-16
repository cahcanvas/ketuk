package billing

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/pay"
)

const (
	KeyMaxActiveInvitations   = "max_active_invitations"
	KeyMaxGuestsPerInvitation = "max_guests_per_invitation"
	KeyPremiumTemplates       = "premium_templates"
	KeyPlannerDaily           = "planner_daily"
	KeyMaxActivePlanners      = "max_active_planners"
)

type Plan struct {
	ID           uuid.UUID         `json:"id"`
	Slug         string            `json:"slug"`
	Name         string            `json:"name"`
	PriceIDR     int64             `json:"price_idr"`
	Interval     string            `json:"interval"`
	Entitlements map[string]string `json:"entitlements"`
}

type Subscription struct {
	PlanSlug     string            `json:"plan"`
	Status       string            `json:"status"`
	Entitlements map[string]string `json:"entitlements"`
}

type Service struct {
	pool        *pgxpool.Pool
	pay         pay.Gateway
	callbackURL string
	returnURL   string
}

func New(pool *pgxpool.Pool, gw pay.Gateway, callbackURL, returnURL string) *Service {
	if gw == nil {
		gw = pay.Disabled{}
	}
	return &Service{pool: pool, pay: gw, callbackURL: callbackURL, returnURL: returnURL}
}

func (s *Service) ListPlans(ctx context.Context) ([]Plan, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, slug, name, price_idr, billing_interval FROM plans WHERE is_active = true ORDER BY price_idr`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var plans []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.Slug, &p.Name, &p.PriceIDR, &p.Interval); err != nil {
			return nil, err
		}
		ents, err := s.entitlements(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		p.Entitlements = ents
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (s *Service) Subscription(ctx context.Context, userID uuid.UUID) (Subscription, error) {
	var planID uuid.UUID
	var sub Subscription
	err := s.pool.QueryRow(ctx, `
		SELECT p.id, p.slug, s.status
		FROM subscriptions s
		JOIN plans p ON p.id = s.plan_id
		WHERE s.user_id = $1
	`, userID).Scan(&planID, &sub.PlanSlug, &sub.Status)
	if err != nil {
		return Subscription{}, apierr.NotFound("subscription")
	}
	ents, err := s.entitlements(ctx, planID)
	if err != nil {
		return Subscription{}, err
	}
	sub.Entitlements = ents
	return sub, nil
}

func (s *Service) AssertInvitationQuota(ctx context.Context, userID uuid.UUID) error {
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM invitations WHERE owner_id=$1 AND status <> 'archived'`, userID).Scan(&n); err != nil {
		return err
	}
	return s.RequireInt(ctx, userID, KeyMaxActiveInvitations, n, 1)
}

func (s *Service) AssertGuestQuota(ctx context.Context, userID, invitationID uuid.UUID) error {
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM invitation_guests WHERE invitation_id=$1`, invitationID).Scan(&n); err != nil {
		return err
	}
	return s.RequireInt(ctx, userID, KeyMaxGuestsPerInvitation, n, 1)
}

func (s *Service) AssertPremiumTemplate(ctx context.Context, userID uuid.UUID, tier string) error {
	if tier != "premium" {
		return nil
	}
	return s.RequireFlag(ctx, userID, KeyPremiumTemplates)
}

func (s *Service) AssertPlannerQuota(ctx context.Context, userID uuid.UUID) error {
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM planners WHERE owner_id=$1 AND status <> 'archived'`, userID).Scan(&n); err != nil {
		return err
	}
	return s.RequireInt(ctx, userID, KeyMaxActivePlanners, n, 1)
}

func (s *Service) AssertPlannerDaily(ctx context.Context, userID uuid.UUID) error {
	return s.RequireFlag(ctx, userID, KeyPlannerDaily)
}

func (s *Service) RequireInt(ctx context.Context, userID uuid.UUID, key string, used, need int) error {
	ents, err := s.userEntitlements(ctx, userID)
	if err != nil {
		return err
	}
	raw, ok := ents[key]
	if !ok {
		return apierr.Entitlement("missing entitlement " + key)
	}
	limit, err := strconv.Atoi(raw)
	if err != nil {
		return err
	}
	if limit < 0 {
		return nil // unlimited
	}
	if used+need > limit {
		return apierr.Entitlement(key + " quota exceeded")
	}
	return nil
}

func (s *Service) RequireFlag(ctx context.Context, userID uuid.UUID, key string) error {
	ents, err := s.userEntitlements(ctx, userID)
	if err != nil {
		return err
	}
	if ents[key] != "1" && ents[key] != "true" {
		return apierr.Entitlement(key + " not included in current plan")
	}
	return nil
}

func (s *Service) userEntitlements(ctx context.Context, userID uuid.UUID) (map[string]string, error) {
	sub, err := s.Subscription(ctx, userID)
	if err != nil {
		return nil, err
	}
	return sub.Entitlements, nil
}

func (s *Service) entitlements(ctx context.Context, planID uuid.UUID) (map[string]string, error) {
	rows, err := s.pool.Query(ctx, `SELECT key, value FROM plan_entitlements WHERE plan_id = $1`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, rows.Err()
}
