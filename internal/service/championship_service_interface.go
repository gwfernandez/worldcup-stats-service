package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ChampionshipService defines the contract for championship business logic.
type ChampionshipService interface {
	List(ctx context.Context, filter domain.ChampionshipFilter) (*domain.ChampionshipListResponse, error)
	GetByYear(ctx context.Context, year int) (*domain.Championship, error)
}
