package config

import "os"

type Config struct {
	ServicePort string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
}

func Load() Config {
	return Config{
		ServicePort: getEnv("SERVICE_PORT", "8060"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "program"),
		DBPassword:  getEnv("DB_PASSWORD", "test"),
		DBName:      getEnv("DB_NAME", "rentals"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
