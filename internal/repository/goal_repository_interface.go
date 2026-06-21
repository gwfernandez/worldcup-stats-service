package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// GoalRepository defines the contract for player goals data access.
type GoalRepository interface {
	ListByPlayer(ctx context.Context, filter domain.GoalFilter) ([]domain.Goal, int64, error)
}
