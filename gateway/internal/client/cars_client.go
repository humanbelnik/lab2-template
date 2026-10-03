package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gateway-service/internal/domain"
)

type CarsClient struct {
	baseURL string
	http    *HTTPClient
}

func NewCarsClient(baseURL string, httpClient *HTTPClient) *CarsClient {
	return &CarsClient{baseURL: baseURL, http: httpClient}
}

type carDTO struct {
	CarUID             string `json:"carUid"`
	Brand              string `json:"brand"`
	Model              string `json:"model"`
	RegistrationNumber string `json:"registrationNumber"`
	Power              *int   `json:"power"`
	Type               string `json:"type"`
	Price              int    `json:"price"`
	Available          bool   `json:"available"`
}

func (d carDTO) toDomain() domain.Car {
	return domain.Car{
		CarUID:             d.CarUID,
		Brand:              d.Brand,
		Model:              d.Model,
		RegistrationNumber: d.RegistrationNumber,
		Power:              d.Power,
		Type:               d.Type,
		Price:              d.Price,
		Available:          d.Available,
	}
}

type carPageDTO struct {
	Page          int      `json:"page"`
	PageSize      int      `json:"pageSize"`
	TotalElements int      `json:"totalElements"`
	Items         []carDTO `json:"items"`
}

func (c *CarsClient) List(ctx context.Context, page, size int, showAll bool) (domain.CarPage, error) {
	url := fmt.Sprintf("%s/api/v1/cars?page=%d&size=%d&showAll=%t", c.baseURL, page, size, showAll)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.CarPage{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.CarPage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.CarPage{}, fmt.Errorf("cars service: unexpected status %d", resp.StatusCode)
	}

	var body carPageDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.CarPage{}, err
	}

	items := make([]domain.Car, 0, len(body.Items))
	for _, item := range body.Items {
		items = append(items, item.toDomain())
	}

	return domain.CarPage{
		Page:          body.Page,
		PageSize:      body.PageSize,
		TotalElements: body.TotalElements,
		Items:         items,
	}, nil
}

var ErrCarNotFound = fmt.Errorf("car not found")

func (c *CarsClient) GetByUID(ctx context.Context, carUID string) (domain.Car, error) {
	url := fmt.Sprintf("%s/api/v1/cars/%s", c.baseURL, carUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Car{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Car{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Car{}, ErrCarNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Car{}, fmt.Errorf("cars service: unexpected status %d", resp.StatusCode)
	}

	var body carDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Car{}, err
	}

	return body.toDomain(), nil
}

func (c *CarsClient) SetAvailability(ctx context.Context, carUID string, available bool) (domain.Car, error) {
	payload, err := json.Marshal(map[string]bool{"available": available})
	if err != nil {
		return domain.Car{}, err
	}

	url := fmt.Sprintf("%s/api/v1/cars/%s", c.baseURL, carUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(payload))
	if err != nil {
		return domain.Car{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Car{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Car{}, ErrCarNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Car{}, fmt.Errorf("cars service: unexpected status %d", resp.StatusCode)
	}

	var body carDTO
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Car{}, err
	}

	return body.toDomain(), nil
}
