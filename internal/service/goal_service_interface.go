package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// GoalService defines the contract for player goals business logic.
type GoalService interface {
	ListByPlayer(ctx context.Context, filter domain.GoalFilter) (*domain.GoalListResponse, error)
}
