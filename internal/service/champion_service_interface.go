package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ChampionService defines the contract for champion business logic.
type ChampionService interface {
	List(ctx context.Context, filter domain.ChampionFilter) (*domain.ChampionListResponse, error)
}
