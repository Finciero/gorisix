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
func NewPaymentService(secret string, receiverID string) *PaymentsService {
	client := httpClient{
		client: &http.Client{},
		secret: secret,
		key:    receiverID,
	}
	return &PaymentsService{&client}
}

// ReceiverID ...
func (ps *PaymentsService) ReceiverID() string { return ps.client.key }

// Payment returns information of the payment with the given id.
func (ps *PaymentsService) Payment(id string) (*PaymentResponse, error) {
	resp, err := ps.client.Get("/merchant/order/"+id, nil)
	if err != nil {
		return nil, err
	}
	var pr PaymentResponse
	if err := unmarshalJSON(resp.Body, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

// Payment returns information of the payment with the given id.
func (ps *PaymentsService) Payments() ([]*PaymentResponse, error) {
	resp, err := ps.client.Get("/merchant/orders/", nil)
	if err != nil {
		return nil, err
	}

	prs := make([]*PaymentResponse, 0)
	if err := unmarshalJSON(resp.Body, &prs); err != nil {
		return nil, err
	}
	return prs, nil
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
	resp, err := ps.client.PostForm("/merchant/orders/", p.Params())
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
	Currency        string `json:"currency"`
	Description     string `json:"description"`
	MerchantOrderId string `json:"merchant_order_id"`
	NotifyURL       string `json:"notify_url"`
	Price           int64  `json:"price"`
	ReturnURL       string `json:"return_url"`
	Timeout         int32  `json:"timeout"`
}

// Params returns a map used to sign the requests
func (p *Payment) Params() url.Values {
	form := url.Values{
		"currency":          {p.Currency},
		"description":       {p.Description},
		"merchant_order_id": {p.MerchantOrderId},
		"notify_url":        {p.NotifyURL},
		"price":             {fmt.Sprintf("%d", p.Price)},
		"return_url":        {p.ReturnURL},
		"timeout":           {fmt.Sprintf("%d", p.Timeout)},
	}

	return form
}

// PaymentResponse represents the information returned by 46d's api after a payment action
type PaymentResponse struct {
	ID              string      `json:"id"`
	Merchant        int32       `json:"merchant"`
	Price           floatString `json:"price"`
	Description     string      `json:"description"`
	MerchantOrderId string      `json:"merchant_order_id"`
	CreationDate    time.Time   `json:"creation_date"`
	ReturnURL       string      `json:"return_url"`
	RedirectURL     string      `json:"redirectURL"`
	Status          string      `json:"status"`
	NotifyURL       string      `json:"notify_url"`
	Timeout         int32       `json:"timeout"`
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
		ID              *string      `json:"id"`
		Merchant        *int32       `json:"merchant"`
		Price           *floatString `json:"price"`
		Description     *string      `json:"description"`
		MerchantOrderId *string      `json:"merchant_order_id"`
		CreationDate    *time.Time   `json:"creation_date"`
		ReturnURL       *string      `json:"return_url"`
		RedirectURL     *string      `json:"redirectURL"`
		Status          *string      `json:"status"`
		NotifyURL       *string      `json:"notify_url"`
		Timeout         *int32       `json:"timeout"`
	}{
		ID:              &p.ID,
		Merchant:        &p.Merchant,
		Price:           &p.Price,
		Description:     &p.Description,
		MerchantOrderId: &p.MerchantOrderId,
		CreationDate:    &p.CreationDate,
		ReturnURL:       &p.ReturnURL,
		RedirectURL:     &p.RedirectURL,
		Status:          &p.Status,
		NotifyURL:       &p.NotifyURL,
		Timeout:         &p.Timeout,
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
