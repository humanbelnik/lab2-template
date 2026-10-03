package http

import "gateway-service/internal/domain"

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
		Type:               car.Type,
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

type CarInfoResponse struct {
	CarUID             string `json:"carUid"`
	Brand              string `json:"brand"`
	Model              string `json:"model"`
	RegistrationNumber string `json:"registrationNumber"`
}

func toCarInfoResponse(car domain.Car) CarInfoResponse {
	return CarInfoResponse{
		CarUID:             car.CarUID,
		Brand:              car.Brand,
		Model:              car.Model,
		RegistrationNumber: car.RegistrationNumber,
	}
}

type PaymentInfoResponse struct {
	PaymentUID string `json:"paymentUid"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}

func toPaymentInfoResponse(payment domain.Payment) PaymentInfoResponse {
	return PaymentInfoResponse{
		PaymentUID: payment.PaymentUID,
		Status:     payment.Status,
		Price:      payment.Price,
	}
}

type RentalResponse struct {
	RentalUID string              `json:"rentalUid"`
	Status    string              `json:"status"`
	DateFrom  string              `json:"dateFrom"`
	DateTo    string              `json:"dateTo"`
	Car       CarInfoResponse     `json:"car"`
	Payment   PaymentInfoResponse `json:"payment"`
}

func toRentalResponse(details domain.RentalDetails) RentalResponse {
	return RentalResponse{
		RentalUID: details.RentalUID,
		Status:    details.Status,
		DateFrom:  details.DateFrom,
		DateTo:    details.DateTo,
		Car:       toCarInfoResponse(details.Car),
		Payment:   toPaymentInfoResponse(details.Payment),
	}
}

func toRentalResponses(details []domain.RentalDetails) []RentalResponse {
	result := make([]RentalResponse, 0, len(details))
	for _, d := range details {
		result = append(result, toRentalResponse(d))
	}
	return result
}

type CreateRentalRequest struct {
	CarUID   string `json:"carUid"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type CreateRentalResponse struct {
	RentalUID string              `json:"rentalUid"`
	Status    string              `json:"status"`
	CarUID    string              `json:"carUid"`
	DateFrom  string              `json:"dateFrom"`
	DateTo    string              `json:"dateTo"`
	Payment   PaymentInfoResponse `json:"payment"`
}

func toCreateRentalResponse(result domain.CreateRentalResult) CreateRentalResponse {
	return CreateRentalResponse{
		RentalUID: result.RentalUID,
		Status:    result.Status,
		CarUID:    result.CarUID,
		DateFrom:  result.DateFrom,
		DateTo:    result.DateTo,
		Payment:   toPaymentInfoResponse(result.Payment),
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ErrorDescription struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Message string             `json:"message"`
	Errors  []ErrorDescription `json:"errors"`
}

func toValidationErrorResponse(err *domain.ValidationError) ValidationErrorResponse {
	errors := make([]ErrorDescription, 0, len(err.Errors))
	for _, fieldError := range err.Errors {
		errors = append(errors, ErrorDescription{Field: fieldError.Field, Error: fieldError.Error})
	}
	return ValidationErrorResponse{
		Message: err.Message,
		Errors:  errors,
	}
}
