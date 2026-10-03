package repository

import (
	"context"

	"payment-service/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, price int) (domain.Payment, error)
	GetByUID(ctx context.Context, paymentUID string) (domain.Payment, error)
	SetStatus(ctx context.Context, paymentUID string, status domain.PaymentStatus) (domain.Payment, error)
}
