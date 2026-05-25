package service

import (
	"context"
	"fmt"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// championshipService implements ChampionshipService with business logic.
type championshipService struct {
	repo repository.ChampionshipRepository
}

// NewChampionshipService creates a new ChampionshipService.
func NewChampionshipService(repo repository.ChampionshipRepository) ChampionshipService {
	return &championshipService{repo: repo}
}

// List returns a paginated and filtered list of championships.
func (s *championshipService) List(ctx context.Context, filter domain.ChampionshipFilter) (*domain.ChampionshipListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	data, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data = make([]domain.Championship, 0)
	}

	return &domain.ChampionshipListResponse{
		Data: data,
		Pagination: domain.PaginationInfo{
			Page:  filter.Page,
			Size:  filter.Size,
			Total: total,
		},
	}, nil
}

// GetByYear returns a championship by its year, filling stats with default values if they don't exist.
func (s *championshipService) GetByYear(ctx context.Context, year int) (*domain.Championship, error) {
	championship, err := s.repo.GetByYear(ctx, year)
	if err != nil {
		return nil, err
	}
	if championship == nil {
		return nil, fmt.Errorf("%w: championship for year %d not found", domain.ErrNotFound, year)
	}

	// Apply business logic fallback: if stats don't exist, use default values
	if championship.Stats == nil {
		championship.Stats = &domain.ChampionshipStats{
			TotalTeams:      0,
			TotalMatches:    0,
			TotalStadiums:   0,
			TotalPlayers:    0,
			TotalGoals:      0,
			RunnerUpCode:    "",
			ThirdPlaceCode:  "",
			FourthPlaceCode: "",
			TopScorers:      make([]domain.TopScorer, 0),
			TopScorerGoals:  0,
		}
	}

	return championship, nil
}
