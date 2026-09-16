package pay

import "context"

type Invoice struct {
	MerchantOrderID string
	AmountIDR       int64
	ProductDetails  string
	Email           string
	CustomerName    string
	CallbackURL     string
	ReturnURL       string
	AdditionalParam string
}

type Result struct {
	Reference  string
	PaymentURL string
}

type Callback struct {
	MerchantOrderID string
	Amount          string
	ResultCode      string
	Reference       string
	AdditionalParam string
	Success         bool
}

// Status is the gateway's own view of a transaction, used to reconcile orders
// whose callback never arrived.
type Status struct {
	MerchantOrderID string
	Reference       string
	Paid            bool
	Pending         bool
}

type Gateway interface {
	Configured() bool
	CreateInvoice(ctx context.Context, inv Invoice) (Result, error)
	VerifyCallback(merchantCode, amount, merchantOrderID, signature string) bool
	CheckTransaction(ctx context.Context, merchantOrderID string) (Status, error)
}
