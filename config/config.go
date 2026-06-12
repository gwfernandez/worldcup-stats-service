package config

import (
	"fmt"
	"os"
	"strings"
)

// Config contains the runtime configuration loaded from environment variables.
type Config struct {
	DatabaseURL        string
	Port               string
	CORSAllowedOrigins []string
}

// Load reads and validates the application configuration from environment variables.
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL:        dbURL,
		Port:               port,
		CORSAllowedOrigins: parseCORSAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")),
	}, nil
}

func parseCORSAllowedOrigins(value string) []string {
	origins := make([]string, 0)
	for _, origin := range strings.Split(value, ",") {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		origins = append(origins, origin)
	}

	return origins
}
