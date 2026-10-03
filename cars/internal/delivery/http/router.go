package http

import "net/http"

func NewRouter(carHandler *CarHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /manage/health", HealthCheck)

	mux.HandleFunc("GET /api/v1/cars", carHandler.List)
	mux.HandleFunc("GET /api/v1/cars/{carUid}", carHandler.GetByUID)
	mux.HandleFunc("PATCH /api/v1/cars/{carUid}", carHandler.UpdateAvailability)

	return mux
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
