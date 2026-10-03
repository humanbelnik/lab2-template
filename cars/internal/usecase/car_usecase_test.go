package usecase

import (
	"context"
	"testing"

	"cars-service/internal/domain"
)

type fakeCarRepository struct {
	cars map[string]domain.Car
}

func newFakeCarRepository() *fakeCarRepository {
	return &fakeCarRepository{cars: map[string]domain.Car{
		"car-1": {CarUID: "car-1", Brand: "Mercedes Benz", Price: 3500, Available: true},
	}}
}

func (f *fakeCarRepository) List(ctx context.Context, page, size int, showAll bool) (domain.CarPage, error) {
	items := make([]domain.Car, 0)
	for _, car := range f.cars {
		if car.Available || showAll {
			items = append(items, car)
		}
	}
	return domain.CarPage{Page: page, PageSize: len(items), TotalElements: len(items), Items: items}, nil
}

func (f *fakeCarRepository) GetByUID(ctx context.Context, carUID string) (domain.Car, error) {
	car, ok := f.cars[carUID]
	if !ok {
		return domain.Car{}, domain.ErrCarNotFound
	}
	return car, nil
}

func (f *fakeCarRepository) SetAvailability(ctx context.Context, carUID string, available bool) (domain.Car, error) {
	car, ok := f.cars[carUID]
	if !ok {
		return domain.Car{}, domain.ErrCarNotFound
	}
	car.Available = available
	f.cars[carUID] = car
	return car, nil
}

func TestCarUsecase_GetByUID_NotFound(t *testing.T) {
	u := NewCarUsecase(newFakeCarRepository())

	_, err := u.GetByUID(context.Background(), "missing")
	if err != domain.ErrCarNotFound {
		t.Fatalf("expected ErrCarNotFound, got %v", err)
	}
}

func TestCarUsecase_SetAvailability(t *testing.T) {
	u := NewCarUsecase(newFakeCarRepository())

	car, err := u.SetAvailability(context.Background(), "car-1", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if car.Available {
		t.Fatalf("expected car to be unavailable")
	}
}

func TestCarUsecase_List_ShowAll(t *testing.T) {
	u := NewCarUsecase(newFakeCarRepository())

	if _, err := u.SetAvailability(context.Background(), "car-1", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	page, err := u.List(context.Background(), 1, 10, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Items) != 0 {
		t.Fatalf("expected no available cars, got %d", len(page.Items))
	}

	page, err = u.List(context.Background(), 1, 10, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Items) != 1 {
		t.Fatalf("expected 1 car when showAll=true, got %d", len(page.Items))
	}
}
