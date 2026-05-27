package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// NationalTeamRepository defines the contract for national team data access.
type NationalTeamRepository interface {
	List(ctx context.Context, filter domain.NationalTeamFilter) ([]domain.NationalTeam, int64, error)
	GetByCode(ctx context.Context, code string) (*domain.NationalTeam, error)
}
