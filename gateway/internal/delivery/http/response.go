package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"gateway-service/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeDomainError(w http.ResponseWriter, err error) {
	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		writeJSON(w, http.StatusBadRequest, toValidationErrorResponse(validationErr))
		return
	}

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Message: notFoundErr.Message})
		return
	}

	writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Message: "upstream service is unavailable"})
}
