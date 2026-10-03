package usecase

import (
	"context"
	"fmt"
	"testing"

	"payment-service/internal/domain"
)

type fakePaymentRepository struct {
	payments map[string]domain.Payment
	nextID   int
}

func newFakePaymentRepository() *fakePaymentRepository {
	return &fakePaymentRepository{payments: make(map[string]domain.Payment)}
}

func (f *fakePaymentRepository) Create(ctx context.Context, price int) (domain.Payment, error) {
	f.nextID++
	uid := fmt.Sprintf("payment-%d", f.nextID)
	payment := domain.Payment{PaymentUID: uid, Status: domain.PaymentStatusPaid, Price: price}
	f.payments[uid] = payment
	return payment, nil
}

func (f *fakePaymentRepository) GetByUID(ctx context.Context, paymentUID string) (domain.Payment, error) {
	payment, ok := f.payments[paymentUID]
	if !ok {
		return domain.Payment{}, domain.ErrPaymentNotFound
	}
	return payment, nil
}

func (f *fakePaymentRepository) SetStatus(ctx context.Context, paymentUID string, status domain.PaymentStatus) (domain.Payment, error) {
	payment, ok := f.payments[paymentUID]
	if !ok {
		return domain.Payment{}, domain.ErrPaymentNotFound
	}
	payment.Status = status
	f.payments[paymentUID] = payment
	return payment, nil
}

func TestPaymentUsecase_CreateAndCancel(t *testing.T) {
	u := NewPaymentUsecase(newFakePaymentRepository())

	payment, err := u.Create(context.Background(), 10500)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status != domain.PaymentStatusPaid {
		t.Fatalf("expected status PAID, got %s", payment.Status)
	}
	if payment.Price != 10500 {
		t.Fatalf("expected price 10500, got %d", payment.Price)
	}

	canceled, err := u.Cancel(context.Background(), payment.PaymentUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canceled.Status != domain.PaymentStatusCanceled {
		t.Fatalf("expected status CANCELED, got %s", canceled.Status)
	}
}

func TestPaymentUsecase_GetByUID_NotFound(t *testing.T) {
	u := NewPaymentUsecase(newFakePaymentRepository())

	_, err := u.GetByUID(context.Background(), "missing")
	if err != domain.ErrPaymentNotFound {
		t.Fatalf("expected ErrPaymentNotFound, got %v", err)
	}
}
