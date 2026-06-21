package service

import (
	"context"
	"fmt"
	"math"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// goalService implements GoalService with business logic.
type goalService struct {
	repo             repository.GoalRepository
	teamNameResolver TeamNameResolver
}

// NewGoalService creates a new GoalService.
func NewGoalService(repo repository.GoalRepository, resolver TeamNameResolver) GoalService {
	return &goalService{repo: repo, teamNameResolver: resolver}
}

// ListByPlayer returns a paginated list of goals scored by a player.
func (s *goalService) ListByPlayer(ctx context.Context, filter domain.GoalFilter) (*domain.GoalListResponse, error) {
	if filter.PlayerID < 1 {
		return nil, fmt.Errorf("%w: playerId must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Year < 0 {
		return nil, fmt.Errorf("%w: year must be greater than or equal to 1 when provided", domain.ErrInvalidInput)
	}
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	data, totalElements, err := s.repo.ListByPlayer(ctx, filter)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data = make([]domain.Goal, 0)
	}
	if err := s.hydrateGoals(ctx, data, filter.Language); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	return &domain.GoalListResponse{
		Data: data,
		Pagination: domain.PaginationInfo{
			Page:          filter.Page,
			Size:          filter.Size,
			TotalElements: totalElements,
			TotalPages:    totalPages,
			HasNext:       filter.Page < totalPages,
			HasPrevious:   filter.Page > 1,
		},
	}, nil
}

func (s *goalService) hydrateGoals(ctx context.Context, goals []domain.Goal, language string) error {
	for i := range goals {
		if goals[i].Hosts == nil {
			goals[i].Hosts = make([]domain.SimpleTeam, 0)
		}
		for j := range goals[i].Hosts {
			name, err := s.teamNameResolver.Resolve(ctx, goals[i].Hosts[j].Code, language)
			if err != nil {
				return err
			}
			goals[i].Hosts[j].Name = name
		}

		name, err := s.teamNameResolver.Resolve(ctx, goals[i].OpponentTeam.Code, language)
		if err != nil {
			return err
		}
		goals[i].OpponentTeam.Name = name
	}
	return nil
}
