package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"rental-service/internal/config"
	deliveryhttp "rental-service/internal/delivery/http"
	"rental-service/internal/repository/postgres"
	"rental-service/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := connectWithRetry(cfg, 10, 2*time.Second)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := postgres.Migrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	rentalRepo := postgres.NewRentalRepository(db)
	rentalUsecase := usecase.NewRentalUsecase(rentalRepo)
	rentalHandler := deliveryhttp.NewRentalHandler(rentalUsecase)

	router := deliveryhttp.NewRouter(rentalHandler)

	server := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("rental service listening on port %s", cfg.ServicePort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func connectWithRetry(cfg config.Config, attempts int, delay time.Duration) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *sql.DB
	var err error
	for i := range attempts {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				return db, nil
			} else {
				err = pingErr
			}
		}
		log.Printf("waiting for database (attempt %d/%d): %v", i+1, attempts, err)
		time.Sleep(delay)
	}
	return nil, err
}
