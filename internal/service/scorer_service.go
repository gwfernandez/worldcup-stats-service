package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// scorerService implements ScorerService with business logic.
type scorerService struct {
	repo             repository.ScorerRepository
	teamNameResolver TeamNameResolver
}

// NewScorerService creates a new ScorerService.
func NewScorerService(repo repository.ScorerRepository, resolvers ...TeamNameResolver) ScorerService {
	var resolver TeamNameResolver
	if len(resolvers) > 0 {
		resolver = resolvers[0]
	}
	return &scorerService{repo: repo, teamNameResolver: resolver}
}

// List returns a paginated list of historical scorers.
func (s *scorerService) List(ctx context.Context, filter domain.ScorerFilter) (*domain.ScorerListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	filter.TeamCode = strings.ToUpper(filter.TeamCode)
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
		data = make([]domain.Scorer, 0)
	}
	if err := s.hydrateScorers(ctx, data, filter.Language); err != nil {
		return nil, err
	}

	return &domain.ScorerListResponse{
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

// GetByID returns a scorer with personal data and all valid goals.
func (s *scorerService) GetByID(ctx context.Context, playerID int64, language string) (*domain.ScorerDetail, error) {
	if playerID < 1 {
		return nil, fmt.Errorf("%w: playerId must be greater than or equal to 1", domain.ErrInvalidInput)
	}

	scorer, err := s.repo.GetByID(ctx, playerID)
	if err != nil {
		return nil, err
	}
	if scorer == nil {
		return nil, fmt.Errorf("%w: scorer with playerId %d not found", domain.ErrNotFound, playerID)
	}

	if scorer.Championships == nil {
		scorer.Championships = make([]int32, 0)
	}
	if scorer.Teams == nil {
		scorer.Teams = make([]domain.SimpleTeam, 0)
	}
	if scorer.Goals == nil {
		scorer.Goals = make([]domain.Goal, 0)
	}

	if err := s.hydrateScorerDetail(ctx, scorer, language); err != nil {
		return nil, err
	}

	return scorer, nil
}

func (s *scorerService) hydrateScorers(ctx context.Context, scorers []domain.Scorer, language string) error {
	for i := range scorers {
		name, err := s.resolveTeamName(ctx, scorers[i].Team.Code, language)
		if err != nil {
			return err
		}
		scorers[i].Team.Name = name
	}
	return nil
}

func (s *scorerService) hydrateScorerDetail(ctx context.Context, scorer *domain.ScorerDetail, language string) error {
	for i := range scorer.Teams {
		name, err := s.resolveTeamName(ctx, scorer.Teams[i].Code, language)
		if err != nil {
			return err
		}
		scorer.Teams[i].Name = name
	}

	for i := range scorer.Goals {
		if scorer.Goals[i].Hosts == nil {
			scorer.Goals[i].Hosts = make([]domain.SimpleTeam, 0)
		}
		for j := range scorer.Goals[i].Hosts {
			name, err := s.resolveTeamName(ctx, scorer.Goals[i].Hosts[j].Code, language)
			if err != nil {
				return err
			}
			scorer.Goals[i].Hosts[j].Name = name
		}

		name, err := s.resolveTeamName(ctx, scorer.Goals[i].OpponentTeam.Code, language)
		if err != nil {
			return err
		}
		scorer.Goals[i].OpponentTeam.Name = name
	}

	return nil
}

func (s *scorerService) resolveTeamName(ctx context.Context, code string, language string) (string, error) {
	if s.teamNameResolver == nil {
		return code, nil
	}
	return s.teamNameResolver.Resolve(ctx, code, language)
}
