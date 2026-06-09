package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// TeamRepository defines the contract for team data access.
type TeamRepository interface {
	List(ctx context.Context, filter domain.TeamFilter) ([]domain.Team, int64, error)
	GetByCode(ctx context.Context, code, language string) (*domain.Team, error)
}
