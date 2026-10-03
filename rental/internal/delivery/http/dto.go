package http

import "rental-service/internal/domain"

type RentalResponse struct {
	RentalUID  string `json:"rentalUid"`
	Username   string `json:"username"`
	PaymentUID string `json:"paymentUid"`
	CarUID     string `json:"carUid"`
	DateFrom   string `json:"dateFrom"`
	DateTo     string `json:"dateTo"`
	Status     string `json:"status"`
}

func toRentalResponse(rental domain.Rental) RentalResponse {
	return RentalResponse{
		RentalUID:  rental.RentalUID,
		Username:   rental.Username,
		PaymentUID: rental.PaymentUID,
		CarUID:     rental.CarUID,
		DateFrom:   rental.DateFrom.UTC().Format(domain.DateLayout),
		DateTo:     rental.DateTo.UTC().Format(domain.DateLayout),
		Status:     string(rental.Status),
	}
}

func toRentalResponses(rentals []domain.Rental) []RentalResponse {
	result := make([]RentalResponse, 0, len(rentals))
	for _, rental := range rentals {
		result = append(result, toRentalResponse(rental))
	}
	return result
}

type CreateRentalRequest struct {
	Username   string `json:"username"`
	PaymentUID string `json:"paymentUid"`
	CarUID     string `json:"carUid"`
	DateFrom   string `json:"dateFrom"`
	DateTo     string `json:"dateTo"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
