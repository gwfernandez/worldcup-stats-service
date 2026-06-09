package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// FixtureService defines the contract for fixture business logic.
type FixtureService interface {
	GetByYear(ctx context.Context, year int, language string) (*domain.Fixture, error)
}
