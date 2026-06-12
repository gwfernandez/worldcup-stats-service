package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jendrix/worldcup-stats-service/config"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("PORT", "9090")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173, https://example.com")

	cfg, err := config.Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, []string{"http://localhost:5173", "https://example.com"}, cfg.CORSAllowedOrigins)
}

func TestLoad_SuccessDefaultPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	// No PORT set

	cfg, err := config.Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "8080", cfg.Port)
	assert.Empty(t, cfg.CORSAllowedOrigins)
}

func TestLoad_SuccessSanitizesCORSAllowedOrigins(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("CORS_ALLOWED_ORIGINS", " http://localhost:5173, ,https://app.example.com,, ")

	cfg, err := config.Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, []string{"http://localhost:5173", "https://app.example.com"}, cfg.CORSAllowedOrigins)
}

func TestLoad_ErrorMissingDatabaseURL(t *testing.T) {
	// No DATABASE_URL set

	cfg, err := config.Load()

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Equal(t, "DATABASE_URL environment variable is required", err.Error())
}
