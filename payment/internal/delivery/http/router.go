package http

import "net/http"

func NewRouter(paymentHandler *PaymentHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /manage/health", HealthCheck)

	mux.HandleFunc("POST /api/v1/payments", paymentHandler.Create)
	mux.HandleFunc("GET /api/v1/payments/{paymentUid}", paymentHandler.GetByUID)
	mux.HandleFunc("POST /api/v1/payments/{paymentUid}/cancel", paymentHandler.Cancel)

	return mux
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
