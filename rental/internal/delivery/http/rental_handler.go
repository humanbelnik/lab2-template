package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"rental-service/internal/domain"
	"rental-service/internal/usecase"
)

type RentalHandler struct {
	usecase *usecase.RentalUsecase
}

func NewRentalHandler(usecase *usecase.RentalUsecase) *RentalHandler {
	return &RentalHandler{usecase: usecase}
}

func (h *RentalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRentalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dateFrom, err := time.Parse(domain.DateLayout, req.DateFrom)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid dateFrom")
		return
	}
	dateTo, err := time.Parse(domain.DateLayout, req.DateTo)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid dateTo")
		return
	}

	rental, err := h.usecase.Create(r.Context(), domain.Rental{
		Username:   req.Username,
		PaymentUID: req.PaymentUID,
		CarUID:     req.CarUID,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toRentalResponse(rental))
}

func (h *RentalHandler) ListByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}

	rentals, err := h.usecase.ListByUsername(r.Context(), username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toRentalResponses(rentals))
}

func (h *RentalHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	rentalUID := r.PathValue("rentalUid")

	rental, err := h.usecase.GetByUID(r.Context(), rentalUID)
	if errors.Is(err, domain.ErrRentalNotFound) {
		writeError(w, http.StatusNotFound, "rental not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toRentalResponse(rental))
}

func (h *RentalHandler) Finish(w http.ResponseWriter, r *http.Request) {
	rentalUID := r.PathValue("rentalUid")

	rental, err := h.usecase.Finish(r.Context(), rentalUID)
	if errors.Is(err, domain.ErrRentalNotFound) {
		writeError(w, http.StatusNotFound, "rental not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toRentalResponse(rental))
}

func (h *RentalHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	rentalUID := r.PathValue("rentalUid")

	rental, err := h.usecase.Cancel(r.Context(), rentalUID)
	if errors.Is(err, domain.ErrRentalNotFound) {
		writeError(w, http.StatusNotFound, "rental not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toRentalResponse(rental))
}
