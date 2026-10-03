package usecase

import (
	"context"
	"errors"
	"time"

	"gateway-service/internal/client"
	"gateway-service/internal/domain"
)

const dateLayout = "2006-01-02"

type RentalUsecase struct {
	carsClient    *client.CarsClient
	paymentClient *client.PaymentClient
	rentalClient  *client.RentalClient
}

func NewRentalUsecase(carsClient *client.CarsClient, paymentClient *client.PaymentClient, rentalClient *client.RentalClient) *RentalUsecase {
	return &RentalUsecase{
		carsClient:    carsClient,
		paymentClient: paymentClient,
		rentalClient:  rentalClient,
	}
}

func (u *RentalUsecase) Create(ctx context.Context, input domain.CreateRentalInput) (domain.CreateRentalResult, error) {
	dateFrom, dateTo, err := validateRentalDates(input)
	if err != nil {
		return domain.CreateRentalResult{}, err
	}

	car, err := u.carsClient.GetByUID(ctx, input.CarUID)
	if errors.Is(err, client.ErrCarNotFound) {
		return domain.CreateRentalResult{}, &domain.ValidationError{
			Message: "car not found",
			Errors:  []domain.FieldError{{Field: "carUid", Error: "car with given carUid does not exist"}},
		}
	}
	if err != nil {
		return domain.CreateRentalResult{}, err
	}
	if !car.Available {
		return domain.CreateRentalResult{}, &domain.ValidationError{
			Message: "car is not available",
			Errors:  []domain.FieldError{{Field: "carUid", Error: "car is not available for booking"}},
		}
	}

	if _, err := u.carsClient.SetAvailability(ctx, input.CarUID, false); err != nil {
		return domain.CreateRentalResult{}, err
	}

	days := int(dateTo.Sub(dateFrom).Hours() / 24)
	price := days * car.Price

	payment, err := u.paymentClient.Create(ctx, price)
	if err != nil {
		_, _ = u.carsClient.SetAvailability(ctx, input.CarUID, true)
		return domain.CreateRentalResult{}, err
	}

	rental, err := u.rentalClient.Create(ctx, domain.Rental{
		Username:   input.Username,
		PaymentUID: payment.PaymentUID,
		CarUID:     input.CarUID,
		DateFrom:   input.DateFrom,
		DateTo:     input.DateTo,
	})
	if err != nil {
		_, _ = u.paymentClient.Cancel(ctx, payment.PaymentUID)
		_, _ = u.carsClient.SetAvailability(ctx, input.CarUID, true)
		return domain.CreateRentalResult{}, err
	}

	return domain.CreateRentalResult{
		RentalUID: rental.RentalUID,
		Status:    rental.Status,
		CarUID:    rental.CarUID,
		DateFrom:  rental.DateFrom,
		DateTo:    rental.DateTo,
		Payment:   payment,
	}, nil
}

func (u *RentalUsecase) ListByUsername(ctx context.Context, username string) ([]domain.RentalDetails, error) {
	rentals, err := u.rentalClient.ListByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	details := make([]domain.RentalDetails, 0, len(rentals))
	for _, rental := range rentals {
		detail, err := u.buildRentalDetails(ctx, rental)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (u *RentalUsecase) GetByUID(ctx context.Context, rentalUID, username string) (domain.RentalDetails, error) {
	rental, err := u.getOwnedRental(ctx, rentalUID, username)
	if err != nil {
		return domain.RentalDetails{}, err
	}

	return u.buildRentalDetails(ctx, rental)
}

func (u *RentalUsecase) Finish(ctx context.Context, rentalUID, username string) error {
	rental, err := u.getOwnedRental(ctx, rentalUID, username)
	if err != nil {
		return err
	}

	if _, err := u.rentalClient.Finish(ctx, rentalUID); err != nil {
		return err
	}

	if _, err := u.carsClient.SetAvailability(ctx, rental.CarUID, true); err != nil {
		return err
	}

	return nil
}

func (u *RentalUsecase) Cancel(ctx context.Context, rentalUID, username string) error {
	rental, err := u.getOwnedRental(ctx, rentalUID, username)
	if err != nil {
		return err
	}

	if _, err := u.rentalClient.Cancel(ctx, rentalUID); err != nil {
		return err
	}

	if _, err := u.carsClient.SetAvailability(ctx, rental.CarUID, true); err != nil {
		return err
	}

	if _, err := u.paymentClient.Cancel(ctx, rental.PaymentUID); err != nil {
		return err
	}

	return nil
}

func (u *RentalUsecase) getOwnedRental(ctx context.Context, rentalUID, username string) (domain.Rental, error) {
	rental, err := u.rentalClient.GetByUID(ctx, rentalUID)
	if errors.Is(err, client.ErrRentalNotFound) {
		return domain.Rental{}, &domain.NotFoundError{Message: "rental not found"}
	}
	if err != nil {
		return domain.Rental{}, err
	}
	if rental.Username != username {
		return domain.Rental{}, &domain.NotFoundError{Message: "rental not found"}
	}
	return rental, nil
}

func (u *RentalUsecase) buildRentalDetails(ctx context.Context, rental domain.Rental) (domain.RentalDetails, error) {
	car, err := u.carsClient.GetByUID(ctx, rental.CarUID)
	if err != nil {
		return domain.RentalDetails{}, err
	}

	payment, err := u.paymentClient.GetByUID(ctx, rental.PaymentUID)
	if err != nil {
		return domain.RentalDetails{}, err
	}

	return domain.RentalDetails{
		RentalUID: rental.RentalUID,
		Status:    rental.Status,
		DateFrom:  rental.DateFrom,
		DateTo:    rental.DateTo,
		Car:       car,
		Payment:   payment,
	}, nil
}

func validateRentalDates(input domain.CreateRentalInput) (time.Time, time.Time, error) {
	var fieldErrors []domain.FieldError

	if input.CarUID == "" {
		fieldErrors = append(fieldErrors, domain.FieldError{Field: "carUid", Error: "must not be empty"})
	}

	dateFrom, errFrom := time.Parse(dateLayout, input.DateFrom)
	if errFrom != nil {
		fieldErrors = append(fieldErrors, domain.FieldError{Field: "dateFrom", Error: "must be a valid date in format YYYY-MM-DD"})
	}

	dateTo, errTo := time.Parse(dateLayout, input.DateTo)
	if errTo != nil {
		fieldErrors = append(fieldErrors, domain.FieldError{Field: "dateTo", Error: "must be a valid date in format YYYY-MM-DD"})
	}

	if errFrom == nil && errTo == nil && !dateFrom.Before(dateTo) {
		fieldErrors = append(fieldErrors, domain.FieldError{Field: "dateTo", Error: "must be after dateFrom"})
	}

	if len(fieldErrors) > 0 {
		return time.Time{}, time.Time{}, &domain.ValidationError{
			Message: "invalid request",
			Errors:  fieldErrors,
		}
	}

	return dateFrom, dateTo, nil
}
