package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// MatchRepository defines the contract for match data access.
type MatchRepository interface {
	ListByYear(ctx context.Context, year int) ([]domain.FixtureMatchRecord, error)
}
