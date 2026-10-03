package usecase

import (
	"context"

	"cars-service/internal/domain"
	"cars-service/internal/repository"
)

type CarUsecase struct {
	repo repository.CarRepository
}

func NewCarUsecase(repo repository.CarRepository) *CarUsecase {
	return &CarUsecase{repo: repo}
}

func (u *CarUsecase) List(ctx context.Context, page, size int, showAll bool) (domain.CarPage, error) {
	return u.repo.List(ctx, page, size, showAll)
}

func (u *CarUsecase) GetByUID(ctx context.Context, carUID string) (domain.Car, error) {
	return u.repo.GetByUID(ctx, carUID)
}

func (u *CarUsecase) SetAvailability(ctx context.Context, carUID string, available bool) (domain.Car, error) {
	return u.repo.SetAvailability(ctx, carUID, available)
}
