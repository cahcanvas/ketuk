package billing

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/pay"
)

type Order struct {
	ID              uuid.UUID  `json:"id"`
	Kind            string     `json:"kind"`
	PlanSlug        string     `json:"plan"`
	AmountIDR       int64      `json:"amount_idr"`
	Status          string     `json:"status"`
	MerchantOrderID string     `json:"merchant_order_id"`
	DuitkuReference string     `json:"duitku_reference"`
	PaymentURL      string     `json:"payment_url"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (s *Service) Checkout(ctx context.Context, userID uuid.UUID, planSlug, email, name string) (Order, error) {
	planSlug = strings.TrimSpace(planSlug)
	p, err := s.planBySlug(ctx, planSlug)
	if err != nil {
		return Order{}, err
	}
	if p.PriceIDR <= 0 {
		return Order{}, apierr.BadRequest("invalid_input", "cannot checkout a free plan")
	}
	sub, err := s.Subscription(ctx, userID)
	if err != nil {
		return Order{}, err
	}
	if sub.PlanSlug == p.Slug {
		return Order{}, apierr.Conflict("already_on_plan", "already subscribed to this plan")
	}

	id := uuid.New()
	merchantOrderID := "saas-" + id.String()
	_, err = s.pool.Exec(ctx, `
		INSERT INTO orders (id, user_id, plan_id, amount_idr, status, merchant_order_id, kind)
		VALUES ($1,$2,$3,$4,'pending',$5,'subscription')`,
		id, userID, p.ID, p.PriceIDR, merchantOrderID)
	if err != nil {
		return Order{}, err
	}

	res, err := s.pay.CreateInvoice(ctx, pay.Invoice{
		MerchantOrderID: merchantOrderID,
		AmountIDR:       p.PriceIDR,
		ProductDetails:  "Ketuk " + p.Name,
		Email:           email,
		CustomerName:    name,
		CallbackURL:     s.callbackURL,
		ReturnURL:       s.returnURL,
		AdditionalParam: "saas",
	})
	if err != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE orders SET status='failed', updated_at=now() WHERE id=$1`, id)
		return Order{}, err
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE orders SET duitku_reference=$2, payment_url=$3, updated_at=now() WHERE id=$1`,
		id, res.Reference, res.PaymentURL)
	if err != nil {
		return Order{}, err
	}
	return Order{
		ID: id, Kind: "subscription", PlanSlug: p.Slug, AmountIDR: p.PriceIDR, Status: "pending",
		MerchantOrderID: merchantOrderID, DuitkuReference: res.Reference, PaymentURL: res.PaymentURL,
		CreatedAt: time.Now(),
	}, nil
}

// ListOrders returns the host's own transaction history. An empty kind means
// every kind; guest gift money lives in gift_payments and is never listed here.
func (s *Service) ListOrders(ctx context.Context, userID uuid.UUID, kind string) ([]Order, error) {
	kind = strings.TrimSpace(kind)
	if kind != "" && kind != "subscription" && kind != "gift" && kind != "vendor" {
		return nil, apierr.BadRequest("invalid_input", "kind must be subscription, gift, or vendor")
	}
	rows, err := s.pool.Query(ctx, `
		SELECT o.id, o.kind, p.slug, o.amount_idr, o.status, o.merchant_order_id, o.duitku_reference, o.payment_url, o.paid_at, o.created_at
		FROM orders o JOIN plans p ON p.id = o.plan_id
		WHERE o.user_id=$1 AND ($2='' OR o.kind=$2)
		ORDER BY o.created_at DESC`, userID, kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Order{}
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Kind, &o.PlanSlug, &o.AmountIDR, &o.Status, &o.MerchantOrderID, &o.DuitkuReference, &o.PaymentURL, &o.PaidAt, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// SyncOrder asks Duitku what really happened to a pending order. It is the way
// out when a callback never arrives, and applies the same effect as one.
func (s *Service) SyncOrder(ctx context.Context, userID, orderID uuid.UUID) (Order, error) {
	var merchantOrderID, status string
	err := s.pool.QueryRow(ctx, `
		SELECT merchant_order_id, status FROM orders WHERE id=$1 AND user_id=$2`,
		orderID, userID).Scan(&merchantOrderID, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return Order{}, apierr.NotFound("order")
	}
	if err != nil {
		return Order{}, err
	}
	if status == "pending" {
		st, err := s.pay.CheckTransaction(ctx, merchantOrderID)
		if err != nil {
			return Order{}, err
		}
		if !st.Pending {
			err = s.ConfirmCallback(ctx, pay.Callback{
				MerchantOrderID: merchantOrderID,
				Reference:       st.Reference,
				Success:         st.Paid,
			})
			if err != nil {
				return Order{}, err
			}
		}
	}
	return s.order(ctx, userID, orderID)
}

// ExpirePending closes out orders the host never paid, so that 'pending' does
// not accumulate forever when Duitku's invoice window has passed.
func (s *Service) ExpirePending(ctx context.Context, olderThan time.Duration) (int64, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE orders SET status='expired', updated_at=now()
		WHERE status='pending' AND created_at < now() - make_interval(mins => $1)`,
		int(olderThan.Minutes()))
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (s *Service) order(ctx context.Context, userID, orderID uuid.UUID) (Order, error) {
	var o Order
	err := s.pool.QueryRow(ctx, `
		SELECT o.id, o.kind, p.slug, o.amount_idr, o.status, o.merchant_order_id, o.duitku_reference, o.payment_url, o.paid_at, o.created_at
		FROM orders o JOIN plans p ON p.id = o.plan_id
		WHERE o.id=$1 AND o.user_id=$2`, orderID, userID).
		Scan(&o.ID, &o.Kind, &o.PlanSlug, &o.AmountIDR, &o.Status, &o.MerchantOrderID, &o.DuitkuReference, &o.PaymentURL, &o.PaidAt, &o.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Order{}, apierr.NotFound("order")
	}
	return o, err
}

func (s *Service) ConfirmCallback(ctx context.Context, cb pay.Callback) error {
	if !strings.HasPrefix(cb.MerchantOrderID, "saas-") {
		return nil
	}
	if !cb.Success {
		_, err := s.pool.Exec(ctx, `
			UPDATE orders SET status='failed', updated_at=now()
			WHERE merchant_order_id=$1 AND status='pending'`, cb.MerchantOrderID)
		return err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var orderID, userID, planID uuid.UUID
	var status string
	err = tx.QueryRow(ctx, `
		SELECT id, user_id, plan_id, status FROM orders WHERE merchant_order_id=$1 FOR UPDATE`,
		cb.MerchantOrderID).Scan(&orderID, &userID, &planID, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		// An unknown merchant order id never becomes known, so a retry cannot succeed.
		return nil
	}
	if err != nil {
		return err
	}
	if status == "paid" {
		return nil
	}
	_, err = tx.Exec(ctx, `
		UPDATE orders SET status='paid', duitku_reference=COALESCE(NULLIF($2,''), duitku_reference), paid_at=now(), updated_at=now()
		WHERE id=$1`, orderID, cb.Reference)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		UPDATE subscriptions SET plan_id=$2, status='active', current_period_end=now() + interval '1 month', updated_at=now()
		WHERE user_id=$1`, userID, planID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) planBySlug(ctx context.Context, slug string) (Plan, error) {
	var p Plan
	err := s.pool.QueryRow(ctx, `
		SELECT id, slug, name, price_idr, billing_interval FROM plans WHERE slug=$1 AND is_active=true`, slug).
		Scan(&p.ID, &p.Slug, &p.Name, &p.PriceIDR, &p.Interval)
	if err != nil {
		return Plan{}, apierr.NotFound("plan")
	}
	ents, err := s.entitlements(ctx, p.ID)
	if err != nil {
		return Plan{}, err
	}
	p.Entitlements = ents
	return p, nil
}
