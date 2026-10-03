package domain

import "errors"

type PaymentStatus string

const (
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusCanceled PaymentStatus = "CANCELED"
)

type Payment struct {
	ID         int64
	PaymentUID string
	Status     PaymentStatus
	Price      int
}

var ErrPaymentNotFound = errors.New("payment not found")
