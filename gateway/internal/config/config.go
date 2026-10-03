package config

import "os"

type Config struct {
	ServicePort       string
	CarsServiceURL    string
	RentalServiceURL  string
	PaymentServiceURL string
}

func Load() Config {
	return Config{
		ServicePort:       getEnv("SERVICE_PORT", "8080"),
		CarsServiceURL:    getEnv("CARS_SERVICE_URL", "http://localhost:8070"),
		RentalServiceURL:  getEnv("RENTAL_SERVICE_URL", "http://localhost:8060"),
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "http://localhost:8050"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
