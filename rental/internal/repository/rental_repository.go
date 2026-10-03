package repository

import (
	"context"

	"rental-service/internal/domain"
)

type RentalRepository interface {
	Create(ctx context.Context, rental domain.Rental) (domain.Rental, error)
	ListByUsername(ctx context.Context, username string) ([]domain.Rental, error)
	GetByUID(ctx context.Context, rentalUID string) (domain.Rental, error)
	SetStatus(ctx context.Context, rentalUID string, status domain.RentalStatus) (domain.Rental, error)
}
