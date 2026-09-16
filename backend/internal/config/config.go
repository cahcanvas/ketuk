package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr       string
	DatabaseURL    string
	JWTSecret      string
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	MigrateOnStart bool
	CORSOrigins    []string
	StorageDir     string
	PublicBaseURL  string
	AuthRateLimit  int
	Duitku         DuitkuConfig
}

type DuitkuConfig struct {
	MerchantCode string
	APIKey       string
	BaseURL      string
	ReturnURL    string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	accessTTL, err := time.ParseDuration(getenv("JWT_ACCESS_TTL", "15m"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_ACCESS_TTL: %w", err)
	}
	refreshTTL, err := time.ParseDuration(getenv("JWT_REFRESH_TTL", "168h"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_REFRESH_TTL: %w", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	secret := getenv("JWT_SECRET", "change-me-in-dev")

	// 0 disables the limiter, which is useful for load tests.
	authRateLimit, err := strconv.Atoi(getenv("AUTH_RATE_LIMIT", "10"))
	if err != nil || authRateLimit < 0 {
		return Config{}, fmt.Errorf("AUTH_RATE_LIMIT must be a non-negative integer")
	}

	return Config{
		HTTPAddr:       getenv("HTTP_ADDR", ":8080"),
		DatabaseURL:    dbURL,
		JWTSecret:      secret,
		AccessTTL:      accessTTL,
		RefreshTTL:     refreshTTL,
		MigrateOnStart: getenv("MIGRATE_ON_START", "true") == "true",
		CORSOrigins:    splitCSV(getenv("CORS_ORIGINS", "http://localhost:5173")),
		StorageDir:     getenv("STORAGE_DIR", "./storage"),
		PublicBaseURL:  strings.TrimRight(getenv("PUBLIC_BASE_URL", "http://localhost:8080"), "/"),
		AuthRateLimit:  authRateLimit,
		Duitku: DuitkuConfig{
			MerchantCode: getenv("DUITKU_MERCHANT_CODE", ""),
			APIKey:       getenv("DUITKU_API_KEY", ""),
			BaseURL:      getenv("DUITKU_BASE_URL", "https://api-sandbox.duitku.com"),
			ReturnURL:    getenv("DUITKU_RETURN_URL", "http://localhost:5173/billing/return"),
		},
	}, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
