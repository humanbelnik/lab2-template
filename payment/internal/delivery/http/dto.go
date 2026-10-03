package http

import "payment-service/internal/domain"

type PaymentResponse struct {
	PaymentUID string `json:"paymentUid"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}

func toPaymentResponse(payment domain.Payment) PaymentResponse {
	return PaymentResponse{
		PaymentUID: payment.PaymentUID,
		Status:     string(payment.Status),
		Price:      payment.Price,
	}
}

type CreatePaymentRequest struct {
	Price int `json:"price"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
