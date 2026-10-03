package postgres

import (
	"context"
	"database/sql"

	"rental-service/internal/domain"
	"rental-service/internal/util"
)

type RentalRepository struct {
	db *sql.DB
}

func NewRentalRepository(db *sql.DB) *RentalRepository {
	return &RentalRepository{db: db}
}

func (r *RentalRepository) Create(ctx context.Context, rental domain.Rental) (domain.Rental, error) {
	rental.RentalUID = util.NewUUID()
	rental.Status = domain.RentalStatusInProgress

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO rental (rental_uid, username, payment_uid, car_uid, date_from, date_to, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		rental.RentalUID, rental.Username, rental.PaymentUID, rental.CarUID,
		rental.DateFrom, rental.DateTo, rental.Status)
	if err != nil {
		return domain.Rental{}, err
	}

	return rental, nil
}

func (r *RentalRepository) ListByUsername(ctx context.Context, username string) ([]domain.Rental, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT rental_uid, username, payment_uid, car_uid, date_from, date_to, status
		FROM rental
		WHERE username = $1
		ORDER BY id`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rentals := make([]domain.Rental, 0)
	for rows.Next() {
		rental, err := scanRental(rows)
		if err != nil {
			return nil, err
		}
		rentals = append(rentals, rental)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rentals, nil
}

func (r *RentalRepository) GetByUID(ctx context.Context, rentalUID string) (domain.Rental, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT rental_uid, username, payment_uid, car_uid, date_from, date_to, status
		FROM rental
		WHERE rental_uid = $1`, rentalUID)

	rental, err := scanRental(row)
	if err == sql.ErrNoRows {
		return domain.Rental{}, domain.ErrRentalNotFound
	}
	if err != nil {
		return domain.Rental{}, err
	}
	return rental, nil
}

func (r *RentalRepository) SetStatus(ctx context.Context, rentalUID string, status domain.RentalStatus) (domain.Rental, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE rental
		SET status = $1
		WHERE rental_uid = $2
		RETURNING rental_uid, username, payment_uid, car_uid, date_from, date_to, status`,
		status, rentalUID)

	rental, err := scanRental(row)
	if err == sql.ErrNoRows {
		return domain.Rental{}, domain.ErrRentalNotFound
	}
	if err != nil {
		return domain.Rental{}, err
	}
	return rental, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRental(s scanner) (domain.Rental, error) {
	var rental domain.Rental
	var status string
	if err := s.Scan(
		&rental.RentalUID,
		&rental.Username,
		&rental.PaymentUID,
		&rental.CarUID,
		&rental.DateFrom,
		&rental.DateTo,
		&status,
	); err != nil {
		return domain.Rental{}, err
	}
	rental.Status = domain.RentalStatus(status)
	return rental, nil
}
