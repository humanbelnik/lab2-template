package http

import (
	"encoding/json"
	"net/http"

	"gateway-service/internal/domain"
	"gateway-service/internal/usecase"
)

const usernameHeader = "X-User-Name"

type RentalHandler struct {
	usecase *usecase.RentalUsecase
}

func NewRentalHandler(usecase *usecase.RentalUsecase) *RentalHandler {
	return &RentalHandler{usecase: usecase}
}

func (h *RentalHandler) Create(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(usernameHeader)
	if username == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	var req CreateRentalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
		return
	}

	result, err := h.usecase.Create(r.Context(), domain.CreateRentalInput{
		Username: username,
		CarUID:   req.CarUID,
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCreateRentalResponse(result))
}

func (h *RentalHandler) ListByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(usernameHeader)
	if username == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	rentals, err := h.usecase.ListByUsername(r.Context(), username)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toRentalResponses(rentals))
}

func (h *RentalHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(usernameHeader)
	if username == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	rentalUID := r.PathValue("rentalUid")

	rental, err := h.usecase.GetByUID(r.Context(), rentalUID, username)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toRentalResponse(rental))
}

func (h *RentalHandler) Finish(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(usernameHeader)
	if username == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	rentalUID := r.PathValue("rentalUid")

	if err := h.usecase.Finish(r.Context(), rentalUID, username); err != nil {
		writeDomainError(w, err)
		return
	}

	writeNoContent(w)
}

func (h *RentalHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(usernameHeader)
	if username == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	rentalUID := r.PathValue("rentalUid")

	if err := h.usecase.Cancel(r.Context(), rentalUID, username); err != nil {
		writeDomainError(w, err)
		return
	}

	writeNoContent(w)
}
