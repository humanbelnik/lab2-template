package main

import (
	"log"
	"net/http"
	"time"

	"gateway-service/internal/client"
	"gateway-service/internal/config"
	deliveryhttp "gateway-service/internal/delivery/http"
	"gateway-service/internal/usecase"
)

func main() {
	cfg := config.Load()

	httpClient := client.NewHTTPClient()
	carsClient := client.NewCarsClient(cfg.CarsServiceURL, httpClient)
	paymentClient := client.NewPaymentClient(cfg.PaymentServiceURL, httpClient)
	rentalClient := client.NewRentalClient(cfg.RentalServiceURL, httpClient)

	carUsecase := usecase.NewCarUsecase(carsClient)
	rentalUsecase := usecase.NewRentalUsecase(carsClient, paymentClient, rentalClient)

	carHandler := deliveryhttp.NewCarHandler(carUsecase)
	rentalHandler := deliveryhttp.NewRentalHandler(rentalUsecase)

	router := deliveryhttp.NewRouter(carHandler, rentalHandler)

	server := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("gateway service listening on port %s", cfg.ServicePort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
