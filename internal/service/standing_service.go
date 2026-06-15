package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// standingService implements StandingService with business logic.
type standingService struct {
	repo             repository.StandingRepository
	teamNameResolver TeamNameResolver
}

// NewStandingService creates a new StandingService.
func NewStandingService(repo repository.StandingRepository, resolvers ...TeamNameResolver) StandingService {
	var resolver TeamNameResolver
	if len(resolvers) > 0 {
		resolver = resolvers[0]
	}
	return &standingService{repo: repo, teamNameResolver: resolver}
}

// List returns a paginated list of historical standings.
func (s *standingService) List(ctx context.Context, filter domain.StandingFilter) (*domain.StandingListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	filter.ConfederationCode = strings.ToUpper(filter.ConfederationCode)

	data, totalElements, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	if data == nil {
		data = make([]domain.Standing, 0)
	}
	if err := s.hydrateStandings(ctx, data, filter.Language); err != nil {
		return nil, err
	}

	return &domain.StandingListResponse{
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

func (s *standingService) hydrateStandings(ctx context.Context, standings []domain.Standing, language string) error {
	for i := range standings {
		name, err := s.resolveTeamName(ctx, standings[i].Team.Code, language)
		if err != nil {
			return err
		}
		standings[i].Team.Name = name
	}
	return nil
}

func (s *standingService) resolveTeamName(ctx context.Context, code string, language string) (string, error) {
	if s.teamNameResolver == nil {
		return code, nil
	}
	return s.teamNameResolver.Resolve(ctx, code, language)
}
