package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ChampionshipRepository defines the contract for championship data access.
type ChampionshipRepository interface {
	List(ctx context.Context, filter domain.ChampionshipFilter) ([]domain.Championship, int64, error)
	GetByYear(ctx context.Context, year int) (*domain.Championship, error)
	ListTeamsByYear(ctx context.Context, filter domain.ChampionshipTeamFilter) ([]domain.ChampionshipTeam, int64, error)
	ListScorersByYear(ctx context.Context, filter domain.ChampionshipScorerFilter) ([]domain.ChampionshipScorer, int64, error)
}
