package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gateway-service/internal/domain"
)

type PaymentClient struct {
	baseURL string
	http    *HTTPClient
}

func NewPaymentClient(baseURL string, httpClient *HTTPClient) *PaymentClient {
	return &PaymentClient{baseURL: baseURL, http: httpClient}
}

type paymentDTO struct {
	PaymentUID string `json:"paymentUid"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}

func (d paymentDTO) toDomain() domain.Payment {
	return domain.Payment{
		PaymentUID: d.PaymentUID,
		Status:     d.Status,
		Price:      d.Price,
	}
}

var ErrPaymentNotFound = fmt.Errorf("payment not found")

func (c *PaymentClient) Create(ctx context.Context, price int) (domain.Payment, error) {
	payload, err := json.Marshal(map[string]int{"price": price})
	if err != nil {
		return domain.Payment{}, err
	}

	url := fmt.Sprintf("%s/api/v1/payments", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return domain.Payment{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Payment{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Payment{}, fmt.Errorf("payment service: unexpected status %d", resp.StatusCode)
	}

	var body paymentDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Payment{}, err
	}

	return body.toDomain(), nil
}

func (c *PaymentClient) GetByUID(ctx context.Context, paymentUID string) (domain.Payment, error) {
	url := fmt.Sprintf("%s/api/v1/payments/%s", c.baseURL, paymentUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Payment{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Payment{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Payment{}, ErrPaymentNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Payment{}, fmt.Errorf("payment service: unexpected status %d", resp.StatusCode)
	}

	var body paymentDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Payment{}, err
	}

	return body.toDomain(), nil
}

func (c *PaymentClient) Cancel(ctx context.Context, paymentUID string) (domain.Payment, error) {
	url := fmt.Sprintf("%s/api/v1/payments/%s/cancel", c.baseURL, paymentUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return domain.Payment{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Payment{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Payment{}, ErrPaymentNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Payment{}, fmt.Errorf("payment service: unexpected status %d", resp.StatusCode)
	}

	var body paymentDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Payment{}, err
	}

	return body.toDomain(), nil
}
