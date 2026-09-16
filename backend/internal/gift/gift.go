package gift

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ketuk.id/api/internal/apierr"
	"ketuk.id/api/internal/invitation"
	"ketuk.id/api/internal/pay"
)

type Service struct {
	pool        *pgxpool.Pool
	invitations *invitation.Service
	pay         pay.Gateway
	callbackURL string
	returnURL   string
}

func New(pool *pgxpool.Pool, invitations *invitation.Service, gw pay.Gateway, callbackURL, returnURL string) *Service {
	if gw == nil {
		gw = pay.Disabled{}
	}
	return &Service{pool: pool, invitations: invitations, pay: gw, callbackURL: callbackURL, returnURL: returnURL}
}

type Method struct {
	ID            uuid.UUID `json:"id"`
	Kind          string    `json:"kind"`
	Label         string    `json:"label"`
	AccountName   string    `json:"account_name"`
	AccountNumber string    `json:"account_number"`
	IsActive      bool      `json:"is_active"`
	SortOrder     int       `json:"sort_order"`
}

type Payment struct {
	ID              uuid.UUID  `json:"id"`
	MethodID        uuid.UUID  `json:"method_id"`
	SenderName      string     `json:"sender_name"`
	Message         string     `json:"message"`
	AmountIDR       int64      `json:"amount_idr"`
	Status          string     `json:"status"`
	MerchantOrderID string     `json:"merchant_order_id"`
	DuitkuReference string     `json:"duitku_reference,omitempty"`
	PaymentURL      string     `json:"payment_url,omitempty"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (s *Service) AddMethod(ctx context.Context, owner, invitationID uuid.UUID, m Method) (Method, error) {
	if _, err := s.invitations.Get(ctx, owner, invitationID); err != nil {
		return Method{}, err
	}
	m.Kind = strings.ToLower(strings.TrimSpace(m.Kind))
	if m.Kind != "bank" && m.Kind != "ewallet" && m.Kind != "duitku" {
		return Method{}, apierr.BadRequest("invalid_input", "kind must be bank, ewallet, or duitku")
	}
	m.Label = strings.TrimSpace(m.Label)
	if m.Label == "" {
		return Method{}, apierr.BadRequest("invalid_input", "label is required")
	}
	m.ID = uuid.New()
	m.IsActive = true
	_, err := s.pool.Exec(ctx, `
		INSERT INTO gift_methods (id, invitation_id, kind, label, account_name, account_number, is_active, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,true,$7)`,
		m.ID, invitationID, m.Kind, m.Label, m.AccountName, m.AccountNumber, m.SortOrder)
	return m, err
}

func (s *Service) ListMethods(ctx context.Context, owner, invitationID uuid.UUID) ([]Method, error) {
	if _, err := s.invitations.Get(ctx, owner, invitationID); err != nil {
		return nil, err
	}
	return s.listMethods(ctx, invitationID, false)
}

func (s *Service) ListPublicMethods(ctx context.Context, invitationID uuid.UUID) ([]Method, error) {
	return s.listMethods(ctx, invitationID, true)
}

func (s *Service) DeleteMethod(ctx context.Context, owner, invitationID, methodID uuid.UUID) error {
	if _, err := s.invitations.Get(ctx, owner, invitationID); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM gift_methods WHERE id=$1 AND invitation_id=$2`, methodID, invitationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apierr.NotFound("gift_method")
	}
	return nil
}

func (s *Service) ListPayments(ctx context.Context, owner, invitationID uuid.UUID) ([]Payment, error) {
	if _, err := s.invitations.Get(ctx, owner, invitationID); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, method_id, sender_name, message, amount_idr, status, merchant_order_id, duitku_reference, payment_url, paid_at, created_at
		FROM gift_payments WHERE invitation_id=$1 ORDER BY created_at DESC`, invitationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Payment{}
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.MethodID, &p.SenderName, &p.Message, &p.AmountIDR, &p.Status, &p.MerchantOrderID, &p.DuitkuReference, &p.PaymentURL, &p.PaidAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Service) CreatePublic(ctx context.Context, invitationID, guestID uuid.UUID, guestName string, methodID uuid.UUID, amount int64, sender, message string) (Payment, error) {
	if amount <= 0 {
		return Payment{}, apierr.BadRequest("invalid_input", "amount_idr must be positive")
	}
	var m Method
	var active bool
	err := s.pool.QueryRow(ctx, `
		SELECT id, kind, label, account_name, account_number, is_active, sort_order
		FROM gift_methods WHERE id=$1 AND invitation_id=$2`, methodID, invitationID).
		Scan(&m.ID, &m.Kind, &m.Label, &m.AccountName, &m.AccountNumber, &active, &m.SortOrder)
	if err != nil {
		return Payment{}, apierr.NotFound("gift_method")
	}
	if !active {
		return Payment{}, apierr.BadRequest("invalid_input", "gift method is inactive")
	}
	sender = strings.TrimSpace(sender)
	if sender == "" {
		sender = guestName
	}
	id := uuid.New()
	merchantOrderID := "gift-" + id.String()
	status := "pending"
	p := Payment{
		ID: id, MethodID: methodID, SenderName: sender, Message: strings.TrimSpace(message),
		AmountIDR: amount, Status: status, MerchantOrderID: merchantOrderID, CreatedAt: time.Now(),
	}
	if m.Kind != "duitku" {
		status = "paid"
		now := time.Now()
		p.Status = status
		p.PaidAt = &now
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO gift_payments (id, invitation_id, method_id, guest_id, sender_name, message, amount_idr, status, merchant_order_id, paid_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		p.ID, invitationID, methodID, guestID, p.SenderName, p.Message, p.AmountIDR, p.Status, p.MerchantOrderID, p.PaidAt)
	if err != nil {
		return Payment{}, err
	}
	if m.Kind != "duitku" {
		return p, nil
	}
	res, err := s.pay.CreateInvoice(ctx, pay.Invoice{
		MerchantOrderID: merchantOrderID,
		AmountIDR:       amount,
		ProductDetails:  "Amplop hadiah " + m.Label,
		Email:           "guest@ketuk.id",
		CustomerName:    sender,
		CallbackURL:     s.callbackURL,
		ReturnURL:       s.returnURL,
		AdditionalParam: "gift",
	})
	if err != nil {
		_, _ = s.pool.Exec(ctx, `UPDATE gift_payments SET status='failed', updated_at=now() WHERE id=$1`, p.ID)
		return Payment{}, err
	}
	p.DuitkuReference = res.Reference
	p.PaymentURL = res.PaymentURL
	_, err = s.pool.Exec(ctx, `
		UPDATE gift_payments SET duitku_reference=$2, payment_url=$3, updated_at=now() WHERE id=$1`,
		p.ID, res.Reference, res.PaymentURL)
	return p, err
}

func (s *Service) ConfirmCallback(ctx context.Context, cb pay.Callback) error {
	if !strings.HasPrefix(cb.MerchantOrderID, "gift-") {
		return nil
	}
	status := "failed"
	if cb.Success {
		status = "paid"
	}
	// No rows updated means the payment was already settled or never existed;
	// neither case can be fixed by letting Duitku retry the callback.
	_, err := s.pool.Exec(ctx, `
		UPDATE gift_payments
		SET status=$2,
		    duitku_reference=COALESCE(NULLIF($3,''), duitku_reference),
		    paid_at=CASE WHEN $2='paid' THEN now() ELSE paid_at END,
		    updated_at=now()
		WHERE merchant_order_id=$1 AND status='pending'`,
		cb.MerchantOrderID, status, cb.Reference)
	return err
}

// ExpirePending closes out Duitku envelopes the guest never paid. Bank and
// e-wallet envelopes are stored as paid on creation, so they are never touched.
func (s *Service) ExpirePending(ctx context.Context, olderThan time.Duration) (int64, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE gift_payments SET status='expired', updated_at=now()
		WHERE status='pending' AND created_at < now() - make_interval(mins => $1)`,
		int(olderThan.Minutes()))
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (s *Service) listMethods(ctx context.Context, invitationID uuid.UUID, activeOnly bool) ([]Method, error) {
	q := `SELECT id, kind, label, account_name, account_number, is_active, sort_order
		FROM gift_methods WHERE invitation_id=$1`
	if activeOnly {
		q += ` AND is_active=true`
	}
	q += ` ORDER BY sort_order, label`
	rows, err := s.pool.Query(ctx, q, invitationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Method{}
	for rows.Next() {
		var m Method
		if err := rows.Scan(&m.ID, &m.Kind, &m.Label, &m.AccountName, &m.AccountNumber, &m.IsActive, &m.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
