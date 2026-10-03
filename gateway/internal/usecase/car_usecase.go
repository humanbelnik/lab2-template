package usecase

import (
	"context"

	"gateway-service/internal/client"
	"gateway-service/internal/domain"
)

type CarUsecase struct {
	carsClient *client.CarsClient
}

func NewCarUsecase(carsClient *client.CarsClient) *CarUsecase {
	return &CarUsecase{carsClient: carsClient}
}

func (u *CarUsecase) List(ctx context.Context, page, size int, showAll bool) (domain.CarPage, error) {
	return u.carsClient.List(ctx, page, size, showAll)
}
