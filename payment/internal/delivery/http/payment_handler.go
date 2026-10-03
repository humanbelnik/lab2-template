package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"
)

type PaymentHandler struct {
	usecase *usecase.PaymentUsecase
}

func NewPaymentHandler(usecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{usecase: usecase}
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Price <= 0 {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	payment, err := h.usecase.Create(r.Context(), req.Price)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toPaymentResponse(payment))
}

func (h *PaymentHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	paymentUID := r.PathValue("paymentUid")

	payment, err := h.usecase.GetByUID(r.Context(), paymentUID)
	if errors.Is(err, domain.ErrPaymentNotFound) {
		writeError(w, http.StatusNotFound, "payment not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toPaymentResponse(payment))
}

func (h *PaymentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	paymentUID := r.PathValue("paymentUid")

	payment, err := h.usecase.Cancel(r.Context(), paymentUID)
	if errors.Is(err, domain.ErrPaymentNotFound) {
		writeError(w, http.StatusNotFound, "payment not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toPaymentResponse(payment))
}
