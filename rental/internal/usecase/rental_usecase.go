package usecase

import (
	"context"

	"rental-service/internal/domain"
	"rental-service/internal/repository"
)

type RentalUsecase struct {
	repo repository.RentalRepository
}

func NewRentalUsecase(repo repository.RentalRepository) *RentalUsecase {
	return &RentalUsecase{repo: repo}
}

func (u *RentalUsecase) Create(ctx context.Context, rental domain.Rental) (domain.Rental, error) {
	return u.repo.Create(ctx, rental)
}

func (u *RentalUsecase) ListByUsername(ctx context.Context, username string) ([]domain.Rental, error) {
	return u.repo.ListByUsername(ctx, username)
}

func (u *RentalUsecase) GetByUID(ctx context.Context, rentalUID string) (domain.Rental, error) {
	return u.repo.GetByUID(ctx, rentalUID)
}

func (u *RentalUsecase) Finish(ctx context.Context, rentalUID string) (domain.Rental, error) {
	return u.repo.SetStatus(ctx, rentalUID, domain.RentalStatusFinished)
}

func (u *RentalUsecase) Cancel(ctx context.Context, rentalUID string) (domain.Rental, error) {
	return u.repo.SetStatus(ctx, rentalUID, domain.RentalStatusCanceled)
}
