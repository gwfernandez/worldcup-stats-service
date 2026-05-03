package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jendrix/worldcup-stats-service/config"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("PORT", "9090")

	cfg, err := config.Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "9090", cfg.Port)
}

func TestLoad_SuccessDefaultPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	// No PORT set

	cfg, err := config.Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "8080", cfg.Port)
}

func TestLoad_ErrorMissingDatabaseURL(t *testing.T) {
	// No DATABASE_URL set

	cfg, err := config.Load()

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Equal(t, "DATABASE_URL environment variable is required", err.Error())
}
