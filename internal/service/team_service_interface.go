package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// TeamService defines the contract for team business logic.
type TeamService interface {
	List(ctx context.Context, filter domain.TeamFilter) (*domain.TeamListResponse, error)
	GetByCode(ctx context.Context, code string) (*domain.Team, error)
}
