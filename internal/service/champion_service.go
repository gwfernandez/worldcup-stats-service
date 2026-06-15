package service

import (
	"context"
	"fmt"
	"math"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// championService implements ChampionService with business logic.
type championService struct {
	repo             repository.ChampionRepository
	teamNameResolver TeamNameResolver
}

// NewChampionService creates a new ChampionService.
func NewChampionService(repo repository.ChampionRepository, resolvers ...TeamNameResolver) ChampionService {
	var resolver TeamNameResolver
	if len(resolvers) > 0 {
		resolver = resolvers[0]
	}
	return &championService{repo: repo, teamNameResolver: resolver}
}

// List returns a paginated list of champions.
func (s *championService) List(ctx context.Context, filter domain.ChampionFilter) (*domain.ChampionListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	data, totalElements, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	if data == nil {
		data = make([]domain.Champion, 0)
	}
	if err := s.hydrateChampions(ctx, data, filter.Language); err != nil {
		return nil, err
	}

	return &domain.ChampionListResponse{
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

func (s *championService) hydrateChampions(ctx context.Context, champions []domain.Champion, language string) error {
	for i := range champions {
		name, err := s.resolveTeamName(ctx, champions[i].Team.Code, language)
		if err != nil {
			return err
		}
		champions[i].Team.Name = name
	}
	return nil
}

func (s *championService) resolveTeamName(ctx context.Context, code string, language string) (string, error) {
	if s.teamNameResolver == nil {
		return code, nil
	}
	return s.teamNameResolver.Resolve(ctx, code, language)
}
