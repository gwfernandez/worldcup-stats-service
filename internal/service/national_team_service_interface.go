package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// NationalTeamService defines the contract for national team business logic.
type NationalTeamService interface {
	List(ctx context.Context, filter domain.NationalTeamFilter) (*domain.NationalTeamListResponse, error)
	GetByCode(ctx context.Context, code string) (*domain.NationalTeam, error)
}
