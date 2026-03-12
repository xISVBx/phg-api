package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName            string
	Env                string
	Port               string
	CORSAllowedOrigins []string
	JWTSecret          string
	JWTExpireSeconds   int64
	DBDSN              string
	FilesBasePath      string
	FilesMaxSizeMB     int64
	FilesAllowedMIME   []string
	CustomerKeyMode    string
	SeedEnabled        bool
	SeedAdminUser      string
	SeedAdminPass      string
	SeedSellerUser     string
	SeedSellerPass     string
}

func Load() *Config {
	_ = godotenv.Load()
	expire, _ := strconv.ParseInt(getenv("JWT_EXPIRE_SECONDS", "3600"), 10, 64)
	maxSize, _ := strconv.ParseInt(getenv("FILES_MAX_SIZE_MB", "25"), 10, 64)
	allowed := strings.Split(getenv("FILES_ALLOWED_MIME", "image/jpeg,image/png,application/pdf"), ",")
	corsOrigins := splitCSV(getenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173"))
	return &Config{
		AppName:            getenv("APP_NAME", "photo-gallery-api"),
		Env:                getenv("APP_ENV", "dev"),
		Port:               getenv("APP_PORT", "8080"),
		CORSAllowedOrigins: corsOrigins,
		JWTSecret:          getenv("JWT_SECRET", "change_me"),
		JWTExpireSeconds:   expire,
		DBDSN:              getenv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=photogallery port=5432 sslmode=disable"),
		FilesBasePath:      getenv("FILES_BASE_PATH", "./storage"),
		FilesMaxSizeMB:     maxSize,
		FilesAllowedMIME:   allowed,
		CustomerKeyMode:    getenv("CUSTOMER_KEY_MODE", "both"),
		SeedEnabled:        getenvBool("SEED_ENABLED", false),
		SeedAdminUser:      getenv("SEED_ADMIN_USERNAME", "admin"),
		SeedAdminPass:      getenv("SEED_ADMIN_PASSWORD", "Admin123*"),
		SeedSellerUser:     getenv("SEED_SELLER_USERNAME", "vendedor"),
		SeedSellerPass:     getenv("SEED_SELLER_PASSWORD", "Vendedor123*"),
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func getenvBool(k string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(k)))
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
