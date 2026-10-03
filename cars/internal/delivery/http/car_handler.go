package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"cars-service/internal/domain"
	"cars-service/internal/usecase"
)

type CarHandler struct {
	usecase *usecase.CarUsecase
}

func NewCarHandler(usecase *usecase.CarUsecase) *CarHandler {
	return &CarHandler{usecase: usecase}
}

func (h *CarHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := 1
	if v := query.Get("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 1 {
			page = parsed
		}
	}

	size := 10
	if v := query.Get("size"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 1 && parsed <= 100 {
			size = parsed
		}
	}

	showAll := false
	if v := query.Get("showAll"); v != "" {
		showAll, _ = strconv.ParseBool(v)
	}

	result, err := h.usecase.List(r.Context(), page, size, showAll)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toPaginationResponse(result))
}

func (h *CarHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	carUID := r.PathValue("carUid")

	car, err := h.usecase.GetByUID(r.Context(), carUID)
	if errors.Is(err, domain.ErrCarNotFound) {
		writeError(w, http.StatusNotFound, "car not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toCarResponse(car))
}

func (h *CarHandler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	carUID := r.PathValue("carUid")

	var req UpdateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	car, err := h.usecase.SetAvailability(r.Context(), carUID, req.Available)
	if errors.Is(err, domain.ErrCarNotFound) {
		writeError(w, http.StatusNotFound, "car not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toCarResponse(car))
}
