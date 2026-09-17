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
	Supabase       SupabaseConfig
	CronSecret     string
}

type DuitkuConfig struct {
	MerchantCode string
	APIKey       string
	BaseURL      string
	ReturnURL    string
}

// SupabaseConfig is only needed when deploying somewhere without a durable
// filesystem (Vercel). Empty URL/ServiceRoleKey means "use local disk
// instead" — see internal/storage.
type SupabaseConfig struct {
	URL            string
	ServiceRoleKey string
	StorageBucket  string
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
		HTTPAddr:       httpAddr(),
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
		Supabase: SupabaseConfig{
			URL:            getenv("SUPABASE_URL", ""),
			ServiceRoleKey: getenv("SUPABASE_SERVICE_ROLE_KEY", ""),
			StorageBucket:  getenv("SUPABASE_STORAGE_BUCKET", "invitation-media"),
		},
		CronSecret: getenv("CRON_SECRET", ""),
	}, nil
}

// httpAddr prefers PORT — the convention hosts like Vercel/Railway/Heroku use
// to tell the app which port they expect it to listen on — falling back to
// HTTP_ADDR (or :8080) for local dev and anywhere PORT isn't set.
func httpAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return getenv("HTTP_ADDR", ":8080")
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
