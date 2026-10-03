package postgres

import (
	"context"
	"database/sql"

	"payment-service/internal/domain"
	"payment-service/internal/util"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, price int) (domain.Payment, error) {
	paymentUID := util.NewUUID()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO payment (payment_uid, status, price)
		VALUES ($1, $2, $3)`, paymentUID, domain.PaymentStatusPaid, price)
	if err != nil {
		return domain.Payment{}, err
	}

	return domain.Payment{
		PaymentUID: paymentUID,
		Status:     domain.PaymentStatusPaid,
		Price:      price,
	}, nil
}

func (r *PaymentRepository) GetByUID(ctx context.Context, paymentUID string) (domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT payment_uid, status, price
		FROM payment
		WHERE payment_uid = $1`, paymentUID)

	payment, err := scanPayment(row)
	if err == sql.ErrNoRows {
		return domain.Payment{}, domain.ErrPaymentNotFound
	}
	if err != nil {
		return domain.Payment{}, err
	}
	return payment, nil
}

func (r *PaymentRepository) SetStatus(ctx context.Context, paymentUID string, status domain.PaymentStatus) (domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE payment
		SET status = $1
		WHERE payment_uid = $2
		RETURNING payment_uid, status, price`, status, paymentUID)

	payment, err := scanPayment(row)
	if err == sql.ErrNoRows {
		return domain.Payment{}, domain.ErrPaymentNotFound
	}
	if err != nil {
		return domain.Payment{}, err
	}
	return payment, nil
}

func scanPayment(row *sql.Row) (domain.Payment, error) {
	var payment domain.Payment
	var status string
	if err := row.Scan(&payment.PaymentUID, &status, &payment.Price); err != nil {
		return domain.Payment{}, err
	}
	payment.Status = domain.PaymentStatus(status)
	return payment, nil
}
