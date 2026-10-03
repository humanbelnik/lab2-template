package http

import (
	"net/http"
	"strconv"

	"gateway-service/internal/usecase"
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
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPaginationResponse(result))
}
