package usecase

import (
	"context"

	"payment-service/internal/domain"
	"payment-service/internal/repository"
)

type PaymentUsecase struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: repo}
}

func (u *PaymentUsecase) Create(ctx context.Context, price int) (domain.Payment, error) {
	return u.repo.Create(ctx, price)
}

func (u *PaymentUsecase) GetByUID(ctx context.Context, paymentUID string) (domain.Payment, error) {
	return u.repo.GetByUID(ctx, paymentUID)
}

func (u *PaymentUsecase) Cancel(ctx context.Context, paymentUID string) (domain.Payment, error) {
	return u.repo.SetStatus(ctx, paymentUID, domain.PaymentStatusCanceled)
}
