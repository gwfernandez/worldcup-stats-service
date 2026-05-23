package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// NationalTeamService defines the contract for national team business logic.
type NationalTeamService interface {
	List(ctx context.Context, filter domain.NationalTeamFilter) (*domain.PaginatedNationalTeams, error)
	GetByID(ctx context.Context, id int64) (*domain.NationalTeam, error)
	GetByCode(ctx context.Context, code string) (*domain.NationalTeam, error)
}
