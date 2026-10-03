package repository

import (
	"context"

	"cars-service/internal/domain"
)

type CarRepository interface {
	List(ctx context.Context, page, size int, showAll bool) (domain.CarPage, error)
	GetByUID(ctx context.Context, carUID string) (domain.Car, error)
	SetAvailability(ctx context.Context, carUID string, available bool) (domain.Car, error)
}
