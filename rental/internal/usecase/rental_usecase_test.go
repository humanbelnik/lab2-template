package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"rental-service/internal/domain"
)

type fakeRentalRepository struct {
	rentals map[string]domain.Rental
	nextID  int
}

func newFakeRentalRepository() *fakeRentalRepository {
	return &fakeRentalRepository{rentals: make(map[string]domain.Rental)}
}

func (f *fakeRentalRepository) Create(ctx context.Context, rental domain.Rental) (domain.Rental, error) {
	f.nextID++
	rental.RentalUID = fmt.Sprintf("rental-%d", f.nextID)
	rental.Status = domain.RentalStatusInProgress
	f.rentals[rental.RentalUID] = rental
	return rental, nil
}

func (f *fakeRentalRepository) ListByUsername(ctx context.Context, username string) ([]domain.Rental, error) {
	result := make([]domain.Rental, 0)
	for _, rental := range f.rentals {
		if rental.Username == username {
			result = append(result, rental)
		}
	}
	return result, nil
}

func (f *fakeRentalRepository) GetByUID(ctx context.Context, rentalUID string) (domain.Rental, error) {
	rental, ok := f.rentals[rentalUID]
	if !ok {
		return domain.Rental{}, domain.ErrRentalNotFound
	}
	return rental, nil
}

func (f *fakeRentalRepository) SetStatus(ctx context.Context, rentalUID string, status domain.RentalStatus) (domain.Rental, error) {
	rental, ok := f.rentals[rentalUID]
	if !ok {
		return domain.Rental{}, domain.ErrRentalNotFound
	}
	rental.Status = status
	f.rentals[rentalUID] = rental
	return rental, nil
}

func TestRentalUsecase_CreateAndFinish(t *testing.T) {
	u := NewRentalUsecase(newFakeRentalRepository())

	created, err := u.Create(context.Background(), domain.Rental{
		Username:   "Alex",
		PaymentUID: "payment-1",
		CarUID:     "car-1",
		DateFrom:   time.Date(2021, 10, 8, 0, 0, 0, 0, time.UTC),
		DateTo:     time.Date(2021, 10, 11, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Status != domain.RentalStatusInProgress {
		t.Fatalf("expected status IN_PROGRESS, got %s", created.Status)
	}

	finished, err := u.Finish(context.Background(), created.RentalUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if finished.Status != domain.RentalStatusFinished {
		t.Fatalf("expected status FINISHED, got %s", finished.Status)
	}
}

func TestRentalUsecase_ListByUsername(t *testing.T) {
	u := NewRentalUsecase(newFakeRentalRepository())

	if _, err := u.Create(context.Background(), domain.Rental{Username: "Alex", CarUID: "car-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := u.Create(context.Background(), domain.Rental{Username: "Max", CarUID: "car-2"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rentals, err := u.ListByUsername(context.Background(), "Alex")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rentals) != 1 {
		t.Fatalf("expected 1 rental for Alex, got %d", len(rentals))
	}
}

func TestRentalUsecase_GetByUID_NotFound(t *testing.T) {
	u := NewRentalUsecase(newFakeRentalRepository())

	_, err := u.GetByUID(context.Background(), "missing")
	if err != domain.ErrRentalNotFound {
		t.Fatalf("expected ErrRentalNotFound, got %v", err)
	}
}
