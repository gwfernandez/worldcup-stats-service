package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// nationalTeamService implements NationalTeamService with business logic.
type nationalTeamService struct {
	repo repository.NationalTeamRepository
}

// NewNationalTeamService creates a new service that depends on the repository interface.
func NewNationalTeamService(repo repository.NationalTeamRepository) NationalTeamService {
	return &nationalTeamService{repo: repo}
}

func (s *nationalTeamService) List(ctx context.Context, filter domain.NationalTeamFilter) (*domain.NationalTeamListResponse, error) {
	teams, totalElements, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	return &domain.NationalTeamListResponse{
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

func (s *nationalTeamService) GetByID(ctx context.Context, id int64) (*domain.NationalTeam, error) {
	team, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, fmt.Errorf("%w: national team with id %d not found", domain.ErrNotFound, id)
	}
	team.Code = strings.ToUpper(team.Code)
	team.ConfederationCode = strings.ToUpper(team.ConfederationCode)
	team.FederationCode = strings.ToUpper(team.FederationCode)
	return team, nil
}

func (s *nationalTeamService) GetByCode(ctx context.Context, code string) (*domain.NationalTeam, error) {
	team, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, fmt.Errorf("%w: national team with code %s not found", domain.ErrNotFound, strings.ToUpper(code))
	}
	team.Code = strings.ToUpper(team.Code)
	team.ConfederationCode = strings.ToUpper(team.ConfederationCode)
	team.FederationCode = strings.ToUpper(team.FederationCode)
	return team, nil
}
