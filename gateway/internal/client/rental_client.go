package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"gateway-service/internal/domain"
)

type RentalClient struct {
	baseURL string
	http    *HTTPClient
}

func NewRentalClient(baseURL string, httpClient *HTTPClient) *RentalClient {
	return &RentalClient{baseURL: baseURL, http: httpClient}
}

type rentalDTO struct {
	RentalUID  string `json:"rentalUid"`
	Username   string `json:"username"`
	PaymentUID string `json:"paymentUid"`
	CarUID     string `json:"carUid"`
	DateFrom   string `json:"dateFrom"`
	DateTo     string `json:"dateTo"`
	Status     string `json:"status"`
}

func (d rentalDTO) toDomain() domain.Rental {
	return domain.Rental{
		RentalUID:  d.RentalUID,
		Username:   d.Username,
		PaymentUID: d.PaymentUID,
		CarUID:     d.CarUID,
		DateFrom:   d.DateFrom,
		DateTo:     d.DateTo,
		Status:     d.Status,
	}
}

var ErrRentalNotFound = fmt.Errorf("rental not found")

func (c *RentalClient) Create(ctx context.Context, rental domain.Rental) (domain.Rental, error) {
	payload, err := json.Marshal(rentalDTO{
		Username:   rental.Username,
		PaymentUID: rental.PaymentUID,
		CarUID:     rental.CarUID,
		DateFrom:   rental.DateFrom,
		DateTo:     rental.DateTo,
	})
	if err != nil {
		return domain.Rental{}, err
	}

	url := fmt.Sprintf("%s/api/v1/rentals", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return domain.Rental{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Rental{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Rental{}, fmt.Errorf("rental service: unexpected status %d", resp.StatusCode)
	}

	var body rentalDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Rental{}, err
	}

	return body.toDomain(), nil
}

func (c *RentalClient) ListByUsername(ctx context.Context, username string) ([]domain.Rental, error) {
	requestURL := fmt.Sprintf("%s/api/v1/rentals?username=%s", c.baseURL, url.QueryEscape(username))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rental service: unexpected status %d", resp.StatusCode)
	}

	var body []rentalDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	rentals := make([]domain.Rental, 0, len(body))
	for _, item := range body {
		rentals = append(rentals, item.toDomain())
	}

	return rentals, nil
}

func (c *RentalClient) GetByUID(ctx context.Context, rentalUID string) (domain.Rental, error) {
	url := fmt.Sprintf("%s/api/v1/rentals/%s", c.baseURL, rentalUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Rental{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Rental{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Rental{}, ErrRentalNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Rental{}, fmt.Errorf("rental service: unexpected status %d", resp.StatusCode)
	}

	var body rentalDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Rental{}, err
	}

	return body.toDomain(), nil
}

func (c *RentalClient) Finish(ctx context.Context, rentalUID string) (domain.Rental, error) {
	return c.postStatusChange(ctx, rentalUID, "finish")
}

func (c *RentalClient) Cancel(ctx context.Context, rentalUID string) (domain.Rental, error) {
	return c.postStatusChange(ctx, rentalUID, "cancel")
}

func (c *RentalClient) postStatusChange(ctx context.Context, rentalUID, action string) (domain.Rental, error) {
	url := fmt.Sprintf("%s/api/v1/rentals/%s/%s", c.baseURL, rentalUID, action)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return domain.Rental{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Rental{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Rental{}, ErrRentalNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Rental{}, fmt.Errorf("rental service: unexpected status %d", resp.StatusCode)
	}

	var body rentalDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Rental{}, err
	}

	return body.toDomain(), nil
}
