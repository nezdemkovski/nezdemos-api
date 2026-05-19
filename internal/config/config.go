package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Settings struct {
	Addr        string
	DatabaseURL string
	APIKey      string
}

func LoadDotenv() error {
	return godotenv.Load()
}

func Load() Settings {
	return Settings{
		Addr:        env("NEZDEMOS_ADDR", ":8080"),
		DatabaseURL: databaseURL(),
		APIKey:      os.Getenv("NEZDEMOS_API_KEY"),
	}
}

func databaseURL() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}
	host := env("POSTGRES_HOST", "localhost")
	port := env("POSTGRES_PORT", "5432")
	name := env("POSTGRES_DB", "health")
	user := env("POSTGRES_USER", "health")
	password := os.Getenv("POSTGRES_PASSWORD")
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   name,
	}
	return u.String()
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
