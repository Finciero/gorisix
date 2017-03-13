package gorisix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// PaymentsService ...
type PaymentsService struct {
	client *httpClient
}

// NewPaymentService ...
func NewPaymentService(secret string, receiverID int) *PaymentsService {
	client := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}
	return &PaymentsService{&client}
}

// ReceiverID ...
func (ps *PaymentsService) ReceiverID() int { return ps.client.recid }

// Payment returns information of the payment with the given id.
func (ps *PaymentsService) Payment(id string) (*PaymentResponse, error) {
	resp, err := ps.client.Get("/payments/"+id, nil)
	if err != nil {
		return nil, err
	}
	var pr PaymentResponse
	if err := unmarshalJSON(resp.Body, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

// PaymentStatus ...
func (ps *PaymentsService) PaymentStatus(notificationToken string) (*PaymentStatusChangeResponse, error) {
	values := url.Values{"notification_token": {notificationToken}}

	resp, err := ps.client.Get("/payments?"+values.Encode(), values)
	if err != nil {
		return nil, err
	}

	var pscr PaymentStatusChangeResponse
	if err := unmarshalJSON(resp.Body, &pscr); err != nil {
		return nil, err
	}

	return &pscr, nil
}

// CreatePayment creates a new payment and returns the URLs to complete the payment.
func (ps *PaymentsService) CreatePayment(p *Payment) (*PaymentResponse, error) {
	resp, err := ps.client.PostForm("/payments", p.Params())
	if err != nil {
		return nil, err
	}

	var pr PaymentResponse
	if err := unmarshalJSON(resp.Body, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

// Payment represents the payment form requires by khipu to make a payment POST
type Payment struct {
	MerchantInvoiceID string  `json:"merchantInvoiceID"`
	Amount            float64 `json:"amount"`
	CurrencyType      string  `json:"currency_type"`
	Description       string  `json:"description"`

	PostbackURL string `json:"postbackURL"`
	ReturnURL   string `json:"returnUrl"`

	TransactionTimeDateStamp time.Time `json:"transactionTimeDateStamp"`
	Timeout                  int32     `json:"timeout"`
}

// Params returns a map used to sign the requests
func (p *Payment) Params() url.Values {
	form := url.Values{
		"merchantInvoiceID": {p.MerchantInvoiceID},
		"amount":            {fmt.Sprintf("%.2f", p.Amount)},
		"currencyType":      {p.CurrencyType},
		"description":       {p.Description},

		"postbackURL": {p.PostbackURL},
		"returnURL":   {p.ReturnURL},

		"transactionsTimeDateStamp": {p.TransactionTimeDateStamp.String()},
		"timeout":                   {fmt.Sprintf("%d", p.Timeout)},
	}

	return form
}

// PaymentResponse represents the information returned by 46d's api after a payment action
type PaymentResponse struct {
	TransactionTimeDateStamp time.Time `json:"transactionTimeDateStamp"`
	TransactionID            string    `json:"transactionID"`
	RedirectURL              string    `json:"redirectURL"`
}

type floatString float64

func (fs *floatString) UnmarshalJSON(b []byte) error {
	str := string(bytes.Trim(b, `"`))

	amount, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*fs = floatString(amount)

	return nil
}

// UnmarshalJSON unmarshal struct
func (p *PaymentResponse) UnmarshalJSON(b []byte) error {
	raw := struct {
		TransactionTimeDateStamp *time.Time `json:"transactionTimeDateStamp"`
		TransactionID            *string    `json:"transactionID"`
		RedirectURL              *string    `json:"redirectURL"`
	}{
		TransactionTimeDateStamp: &p.TransactionTimeDateStamp,
		TransactionID:            &p.TransactionID,
		RedirectURL:              &p.RedirectURL,
	}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	return nil
}

// PaymentStatusChangeRequest represents the information returned by 46d's when a payment has changes status.
type PaymentStatusChangeRequest struct {
	NotificationTimeStamp time.Time `json:"notificationTimeStamp"`
	NotificationID        string    `json:"notificationID"`
}

// PaymentStatusChangeResponse represents a success response defined by 46d.
type PaymentStatusChangeResponse struct {
	Result string `json:"result"`
}
