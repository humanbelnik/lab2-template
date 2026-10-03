package http

import "net/http"

func NewRouter(rentalHandler *RentalHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /manage/health", HealthCheck)

	mux.HandleFunc("POST /api/v1/rentals", rentalHandler.Create)
	mux.HandleFunc("GET /api/v1/rentals", rentalHandler.ListByUsername)
	mux.HandleFunc("GET /api/v1/rentals/{rentalUid}", rentalHandler.GetByUID)
	mux.HandleFunc("POST /api/v1/rentals/{rentalUid}/finish", rentalHandler.Finish)
	mux.HandleFunc("POST /api/v1/rentals/{rentalUid}/cancel", rentalHandler.Cancel)

	return mux
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
