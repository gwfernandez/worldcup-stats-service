package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ChampionRepository defines the contract for champion data access.
type ChampionRepository interface {
	List(ctx context.Context, filter domain.ChampionFilter) ([]domain.Champion, int64, error)
}
