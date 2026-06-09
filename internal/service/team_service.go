package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// TeamService implements TeamService with business logic.
type teamService struct {
	repo repository.TeamRepository
}

// NewTeamService creates a new service that depends on the repository interface.
func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{repo: repo}
}

func (s *teamService) List(ctx context.Context, filter domain.TeamFilter) (*domain.TeamListResponse, error) {
	teams, totalElements, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	return &domain.TeamListResponse{
		Data: teams,
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

func (s *teamService) GetByCode(ctx context.Context, code, language string) (*domain.Team, error) {
	team, err := s.repo.GetByCode(ctx, code, language)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, fmt.Errorf("%w: team with code %s not found", domain.ErrNotFound, strings.ToUpper(code))
	}
	team.Code = strings.ToUpper(team.Code)
	team.ConfederationCode = strings.ToUpper(team.ConfederationCode)
	team.FederationCode = strings.ToUpper(team.FederationCode)
	return team, nil
}
