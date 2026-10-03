package usecase

import (
	"errors"
	"testing"

	"gateway-service/internal/domain"
)

func TestValidateRentalDates_Valid(t *testing.T) {
	_, _, err := validateRentalDates(domain.CreateRentalInput{
		CarUID:   "car-1",
		DateFrom: "2021-10-08",
		DateTo:   "2021-10-11",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRentalDates_MissingCarUID(t *testing.T) {
	_, _, err := validateRentalDates(domain.CreateRentalInput{
		DateFrom: "2021-10-08",
		DateTo:   "2021-10-11",
	})
	assertValidationError(t, err)
}

func TestValidateRentalDates_InvalidFormat(t *testing.T) {
	_, _, err := validateRentalDates(domain.CreateRentalInput{
		CarUID:   "car-1",
		DateFrom: "08-10-2021",
		DateTo:   "2021-10-11",
	})
	assertValidationError(t, err)
}

func TestValidateRentalDates_DateToBeforeDateFrom(t *testing.T) {
	_, _, err := validateRentalDates(domain.CreateRentalInput{
		CarUID:   "car-1",
		DateFrom: "2021-10-11",
		DateTo:   "2021-10-08",
	})
	assertValidationError(t, err)
}

func assertValidationError(t *testing.T, err error) {
	t.Helper()
	var validationErr *domain.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if len(validationErr.Errors) == 0 {
		t.Fatalf("expected field errors, got none")
	}
}
