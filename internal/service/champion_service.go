package service

import (
	"context"
	"fmt"
	"math"
	"strings"

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

// ListFinalsWonByTeam returns a paginated list of World Cup finals won by a team.
func (s *championService) ListFinalsWonByTeam(ctx context.Context, filter domain.ChampionFinalFilter) (*domain.ChampionFinalListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	filter.TeamCode = strings.ToUpper(strings.TrimSpace(filter.TeamCode))
	finals, totalElements, err := s.repo.ListFinalsWonByTeam(ctx, filter)
	if err != nil {
		return nil, err
	}

	if finals == nil {
		finals = make([]domain.ChampionFinal, 0)
	}
	if err := s.hydrateChampionFinals(ctx, finals, filter.Language); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	return &domain.ChampionFinalListResponse{
		Data: finals,
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

func (s *championService) hydrateChampionFinals(ctx context.Context, finals []domain.ChampionFinal, language string) error {
	for i := range finals {
		for j := range finals[i].HostCodes {
			hostName, err := s.resolveTeamName(ctx, finals[i].HostCodes[j].Code, language)
			if err != nil {
				return err
			}
			finals[i].HostCodes[j].Name = hostName
		}

		homeTeamName, err := s.resolveTeamName(ctx, finals[i].HomeTeam.Code, language)
		if err != nil {
			return err
		}
		finals[i].HomeTeam.Name = homeTeamName

		awayTeamName, err := s.resolveTeamName(ctx, finals[i].AwayTeam.Code, language)
		if err != nil {
			return err
		}
		finals[i].AwayTeam.Name = awayTeamName
	}
	return nil
}

func (s *championService) resolveTeamName(ctx context.Context, code string, language string) (string, error) {
	if s.teamNameResolver == nil {
		return code, nil
	}
	return s.teamNameResolver.Resolve(ctx, code, language)
}
