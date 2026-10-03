package postgres

import (
	"context"
	"database/sql"

	"cars-service/internal/domain"
)

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{db: db}
}

func (r *CarRepository) List(ctx context.Context, page, size int, showAll bool) (domain.CarPage, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM cars WHERE availability = true OR $1`
	if err := r.db.QueryRowContext(ctx, countQuery, showAll).Scan(&total); err != nil {
		return domain.CarPage{}, err
	}

	offset := (page - 1) * size
	rows, err := r.db.QueryContext(ctx, `
		SELECT car_uid, brand, model, registration_number, power, price, type, availability
		FROM cars
		WHERE availability = true OR $1
		ORDER BY id
		LIMIT $2 OFFSET $3`, showAll, size, offset)
	if err != nil {
		return domain.CarPage{}, err
	}
	defer rows.Close()

	items := make([]domain.Car, 0)
	for rows.Next() {
		car, err := scanCar(rows)
		if err != nil {
			return domain.CarPage{}, err
		}
		items = append(items, car)
	}
	if err := rows.Err(); err != nil {
		return domain.CarPage{}, err
	}

	return domain.CarPage{
		Page:          page,
		PageSize:      len(items),
		TotalElements: total,
		Items:         items,
	}, nil
}

func (r *CarRepository) GetByUID(ctx context.Context, carUID string) (domain.Car, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT car_uid, brand, model, registration_number, power, price, type, availability
		FROM cars
		WHERE car_uid = $1`, carUID)

	car, err := scanCar(row)
	if err == sql.ErrNoRows {
		return domain.Car{}, domain.ErrCarNotFound
	}
	if err != nil {
		return domain.Car{}, err
	}
	return car, nil
}

func (r *CarRepository) SetAvailability(ctx context.Context, carUID string, available bool) (domain.Car, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE cars
		SET availability = $1
		WHERE car_uid = $2
		RETURNING car_uid, brand, model, registration_number, power, price, type, availability`,
		available, carUID)

	car, err := scanCar(row)
	if err == sql.ErrNoRows {
		return domain.Car{}, domain.ErrCarNotFound
	}
	if err != nil {
		return domain.Car{}, err
	}
	return car, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCar(s scanner) (domain.Car, error) {
	var car domain.Car
	var carType string
	if err := s.Scan(
		&car.CarUID,
		&car.Brand,
		&car.Model,
		&car.RegistrationNumber,
		&car.Power,
		&car.Price,
		&carType,
		&car.Available,
	); err != nil {
		return domain.Car{}, err
	}
	car.Type = domain.CarType(carType)
	return car, nil
}
