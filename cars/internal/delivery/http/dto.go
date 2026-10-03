package http

import "cars-service/internal/domain"

type CarResponse struct {
	CarUID             string `json:"carUid"`
	Brand              string `json:"brand"`
	Model              string `json:"model"`
	RegistrationNumber string `json:"registrationNumber"`
	Power              *int   `json:"power"`
	Type               string `json:"type"`
	Price              int    `json:"price"`
	Available          bool   `json:"available"`
}

func toCarResponse(car domain.Car) CarResponse {
	return CarResponse{
		CarUID:             car.CarUID,
		Brand:              car.Brand,
		Model:              car.Model,
		RegistrationNumber: car.RegistrationNumber,
		Power:              car.Power,
		Type:               string(car.Type),
		Price:              car.Price,
		Available:          car.Available,
	}
}

type PaginationResponse struct {
	Page          int           `json:"page"`
	PageSize      int           `json:"pageSize"`
	TotalElements int           `json:"totalElements"`
	Items         []CarResponse `json:"items"`
}

func toPaginationResponse(page domain.CarPage) PaginationResponse {
	items := make([]CarResponse, 0, len(page.Items))
	for _, car := range page.Items {
		items = append(items, toCarResponse(car))
	}
	return PaginationResponse{
		Page:          page.Page,
		PageSize:      page.PageSize,
		TotalElements: page.TotalElements,
		Items:         items,
	}
}

type UpdateAvailabilityRequest struct {
	Available bool `json:"available"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
