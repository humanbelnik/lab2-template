package http

import "net/http"

func NewRouter(carHandler *CarHandler, rentalHandler *RentalHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /manage/health", HealthCheck)

	mux.HandleFunc("GET /api/v1/cars", carHandler.List)

	mux.HandleFunc("GET /api/v1/rental", rentalHandler.ListByUsername)
	mux.HandleFunc("POST /api/v1/rental", rentalHandler.Create)
	mux.HandleFunc("GET /api/v1/rental/{rentalUid}", rentalHandler.GetByUID)
	mux.HandleFunc("DELETE /api/v1/rental/{rentalUid}", rentalHandler.Cancel)
	mux.HandleFunc("POST /api/v1/rental/{rentalUid}/finish", rentalHandler.Finish)

	return mux
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
