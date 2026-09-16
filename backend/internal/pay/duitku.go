package pay

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ketuk.id/api/internal/apierr"
)

type Duitku struct {
	MerchantCode string
	APIKey       string
	BaseURL      string
	HTTP         *http.Client
}

func NewDuitku(merchantCode, apiKey, baseURL string) *Duitku {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api-sandbox.duitku.com"
	}
	return &Duitku{
		MerchantCode: strings.TrimSpace(merchantCode),
		APIKey:       strings.TrimSpace(apiKey),
		BaseURL:      strings.TrimRight(baseURL, "/"),
		HTTP:         &http.Client{Timeout: 20 * time.Second},
	}
}

func (d *Duitku) Configured() bool {
	return d != nil && d.MerchantCode != "" && d.APIKey != ""
}

func (d *Duitku) CreateInvoice(ctx context.Context, inv Invoice) (Result, error) {
	if !d.Configured() {
		return Result{}, apierr.Unavailable("payment_unavailable", "Duitku is not configured")
	}
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sig := hmacSHA256(d.MerchantCode+ts, d.APIKey)
	body := map[string]any{
		"paymentAmount":   inv.AmountIDR,
		"merchantOrderId": inv.MerchantOrderID,
		"productDetails":  inv.ProductDetails,
		"email":           inv.Email,
		"customerVaName":  inv.CustomerName,
		"callbackUrl":     inv.CallbackURL,
		"returnUrl":       inv.ReturnURL,
		"additionalParam": inv.AdditionalParam,
		"expiryPeriod":    60,
		"itemDetails": []map[string]any{{
			"name":     inv.ProductDetails,
			"price":    inv.AmountIDR,
			"quantity": 1,
		}},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return Result{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.BaseURL+"/api/merchant/createInvoice", bytes.NewReader(raw))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-duitku-timestamp", ts)
	req.Header.Set("x-duitku-signature", sig)
	req.Header.Set("x-duitku-merchantcode", d.MerchantCode)

	res, err := d.HTTP.Do(req)
	if err != nil {
		return Result{}, apierr.Unavailable("payment_unavailable", "Duitku request failed")
	}
	defer res.Body.Close()
	payload, _ := io.ReadAll(res.Body)
	var parsed struct {
		StatusCode    string `json:"statusCode"`
		StatusMessage string `json:"statusMessage"`
		Reference     string `json:"reference"`
		PaymentURL    string `json:"paymentUrl"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return Result{}, fmt.Errorf("duitku: invalid response")
	}
	if res.StatusCode >= 300 || (parsed.StatusCode != "" && parsed.StatusCode != "00") {
		msg := parsed.StatusMessage
		if msg == "" {
			msg = "Duitku rejected invoice"
		}
		return Result{}, apierr.Unavailable("payment_unavailable", msg)
	}
	return Result{Reference: parsed.Reference, PaymentURL: parsed.PaymentURL}, nil
}

func (d *Duitku) VerifyCallback(merchantCode, amount, merchantOrderID, signature string) bool {
	if !d.Configured() {
		return false
	}
	if merchantCode != d.MerchantCode {
		return false
	}
	wantHMAC := hmacSHA256(merchantCode+amount+merchantOrderID, d.APIKey)
	wantMD5 := md5Hex(merchantCode + amount + merchantOrderID + d.APIKey)
	return hmac.Equal([]byte(strings.ToLower(signature)), []byte(wantHMAC)) ||
		hmac.Equal([]byte(strings.ToLower(signature)), []byte(wantMD5))
}

func (d *Duitku) CheckTransaction(ctx context.Context, merchantOrderID string) (Status, error) {
	if !d.Configured() {
		return Status{}, apierr.Unavailable("payment_unavailable", "Duitku is not configured")
	}
	body := map[string]any{
		"merchantCode":    d.MerchantCode,
		"merchantOrderId": merchantOrderID,
		"signature":       md5Hex(d.MerchantCode + merchantOrderID + d.APIKey),
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return Status{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.BaseURL+"/api/merchant/transactionStatus", bytes.NewReader(raw))
	if err != nil {
		return Status{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := d.HTTP.Do(req)
	if err != nil {
		return Status{}, apierr.Unavailable("payment_unavailable", "Duitku request failed")
	}
	defer res.Body.Close()
	payload, _ := io.ReadAll(res.Body)
	var parsed struct {
		MerchantOrderID string `json:"merchantOrderId"`
		Reference       string `json:"reference"`
		StatusCode      string `json:"statusCode"`
		StatusMessage   string `json:"statusMessage"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return Status{}, fmt.Errorf("duitku: invalid response")
	}
	if res.StatusCode >= 300 {
		msg := parsed.StatusMessage
		if msg == "" {
			msg = "Duitku rejected the status request"
		}
		return Status{}, apierr.Unavailable("payment_unavailable", msg)
	}
	return Status{
		MerchantOrderID: merchantOrderID,
		Reference:       parsed.Reference,
		Paid:            parsed.StatusCode == "00",
		Pending:         parsed.StatusCode == "01",
	}, nil
}

func hmacSHA256(message, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

type Disabled struct{}

func (Disabled) Configured() bool { return false }

func (Disabled) CreateInvoice(context.Context, Invoice) (Result, error) {
	return Result{}, apierr.Unavailable("payment_unavailable", "Duitku is not configured")
}

func (Disabled) VerifyCallback(string, string, string, string) bool { return false }

func (Disabled) CheckTransaction(context.Context, string) (Status, error) {
	return Status{}, apierr.Unavailable("payment_unavailable", "Duitku is not configured")
}
