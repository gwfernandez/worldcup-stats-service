package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// championshipService implements ChampionshipService with business logic.
type championshipService struct {
	repo             repository.ChampionshipRepository
	teamNameResolver TeamNameResolver
}

// NewChampionshipService creates a new ChampionshipService.
func NewChampionshipService(repo repository.ChampionshipRepository, resolvers ...TeamNameResolver) ChampionshipService {
	resolver := TeamNameResolver(NewCachedTeamNameResolver(repo))
	if len(resolvers) > 0 && resolvers[0] != nil {
		resolver = resolvers[0]
	}
	return &championshipService{repo: repo, teamNameResolver: resolver}
}

// List returns a paginated and filtered list of championships.
func (s *championshipService) List(ctx context.Context, filter domain.ChampionshipFilter) (*domain.ChampionshipListResponse, error) {
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
		data = make([]domain.Championship, 0)
	}
	if err := s.hydrateChampionshipHosts(ctx, data, filter.Language); err != nil {
		return nil, err
	}

	return &domain.ChampionshipListResponse{
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

// ListTeamsByYear returns a paginated and filtered list of teams for a championship year.
func (s *championshipService) ListTeamsByYear(ctx context.Context, filter domain.ChampionshipTeamFilter) (*domain.ChampionshipTeamListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	data, totalElements, err := s.repo.ListTeamsByYear(ctx, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	if data == nil {
		data = make([]domain.ChampionshipTeam, 0)
	}
	if err := s.hydrateChampionshipTeams(ctx, data, filter.Language); err != nil {
		return nil, err
	}

	return &domain.ChampionshipTeamListResponse{
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

// ListStadiumsByYear returns a paginated and filtered list of stadiums for a championship year.
func (s *championshipService) ListStadiumsByYear(ctx context.Context, filter domain.ChampionshipStadiumFilter) (*domain.ChampionshipStadiumListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	data, totalElements, err := s.repo.ListStadiumsByYear(ctx, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	if data == nil {
		data = make([]domain.ChampionshipStadium, 0)
	}

	return &domain.ChampionshipStadiumListResponse{
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

// ListScorersByYear returns a paginated and filtered list of scorers for a championship year.
func (s *championshipService) ListScorersByYear(ctx context.Context, filter domain.ChampionshipScorerFilter) (*domain.ChampionshipScorerListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	data, totalElements, err := s.repo.ListScorersByYear(ctx, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	if data == nil {
		data = make([]domain.ChampionshipScorer, 0)
	}
	if err := s.hydrateChampionshipScorers(ctx, data, filter.Language); err != nil {
		return nil, err
	}

	return &domain.ChampionshipScorerListResponse{
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

// ListSquadByYearAndTeam returns a paginated list of squad players for a team in a championship year.
func (s *championshipService) ListSquadByYearAndTeam(ctx context.Context, filter domain.ChampionshipSquadFilter) (*domain.ChampionshipSquadListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	data, totalElements, err := s.repo.ListSquadByYearAndTeam(ctx, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	if data == nil {
		data = make([]domain.ChampionshipSquadPlayer, 0)
	}

	return &domain.ChampionshipSquadListResponse{
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

// ListStandingsByYear returns a paginated list of standings for a championship year.
func (s *championshipService) ListStandingsByYear(ctx context.Context, filter domain.ChampionshipStandingFilter) (*domain.ChampionshipStandingListResponse, error) {
	if filter.Page < 1 {
		return nil, fmt.Errorf("%w: page must be greater than or equal to 1", domain.ErrInvalidInput)
	}
	if filter.Size < 1 || filter.Size > 100 {
		return nil, fmt.Errorf("%w: size must be between 1 and 100", domain.ErrInvalidInput)
	}

	data, totalElements, err := s.repo.ListStandingsByYear(ctx, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalElements) / float64(filter.Size)))
	if totalElements == 0 {
		totalPages = 0
	}

	if data == nil {
		data = make([]domain.ChampionshipStanding, 0)
	}
	if err := s.hydrateChampionshipStandings(ctx, data, filter.Language); err != nil {
		return nil, err
	}

	return &domain.ChampionshipStandingListResponse{
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

// GetByYear returns a championship by its year, filling stats with default values if they don't exist.
func (s *championshipService) GetByYear(ctx context.Context, year int, language string) (*domain.Championship, error) {
	championship, err := s.repo.GetByYear(ctx, year)
	if err != nil {
		return nil, err
	}
	if championship == nil {
		return nil, fmt.Errorf("%w: championship for year %d not found", domain.ErrNotFound, year)
	}

	// Apply business logic fallback: if stats don't exist, use default values
	if championship.Stats == nil {
		championship.Stats = &domain.ChampionshipsStats{
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
	if err := s.hydrateChampionship(ctx, championship, language); err != nil {
		return nil, err
	}

	return championship, nil
}

func (s *championshipService) hydrateChampionshipHosts(ctx context.Context, championships []domain.Championship, language string) error {
	for i := range championships {
		if err := s.hydrateChampionship(ctx, &championships[i], language); err != nil {
			return err
		}
	}
	return nil
}

func (s *championshipService) hydrateChampionship(ctx context.Context, championship *domain.Championship, language string) error {
	hosts, err := s.resolveHosts(ctx, championship.HostCodes, language)
	if err != nil {
		return err
	}
	championship.Hosts = hosts
	champion, err := s.resolveChampion(ctx, championship.ChampionCode, language)
	if err != nil {
		return err
	}
	championship.Champion = champion
	return s.hydrateChampionshipStats(ctx, championship.Stats, language)
}

func (s *championshipService) hydrateChampionshipTeams(ctx context.Context, teams []domain.ChampionshipTeam, language string) error {
	for i := range teams {
		name, err := s.resolveTeamName(ctx, teams[i].Team.Code, language)
		if err != nil {
			return err
		}
		teams[i].Team.Name = name
	}
	return nil
}

func (s *championshipService) hydrateChampionshipScorers(ctx context.Context, scorers []domain.ChampionshipScorer, language string) error {
	for i := range scorers {
		name, err := s.resolveTeamName(ctx, scorers[i].Team.Code, language)
		if err != nil {
			return err
		}
		scorers[i].Team.Name = name
	}
	return nil
}

func (s *championshipService) hydrateChampionshipStandings(ctx context.Context, standings []domain.ChampionshipStanding, language string) error {
	for i := range standings {
		name, err := s.resolveTeamName(ctx, standings[i].Team.Code, language)
		if err != nil {
			return err
		}
		standings[i].Team.Name = name
	}
	return nil
}

func (s *championshipService) resolveHosts(ctx context.Context, hostCodes []string, language string) ([]domain.SimpleTeam, error) {
	hosts := make([]domain.SimpleTeam, 0, len(hostCodes))
	for _, hostCode := range hostCodes {
		code := strings.ToUpper(strings.TrimSpace(hostCode))
		if code == "" {
			continue
		}
		name, err := s.resolveTeamName(ctx, code, language)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, domain.SimpleTeam{
			Code: code,
			Name: name,
		})
	}
	return hosts, nil
}

func (s *championshipService) resolveChampion(ctx context.Context, championCode *string, language string) (*domain.SimpleTeam, error) {
	if championCode == nil {
		return nil, nil
	}
	code := strings.ToUpper(strings.TrimSpace(*championCode))
	if code == "" {
		return nil, nil
	}
	name, err := s.resolveTeamName(ctx, code, language)
	if err != nil {
		return nil, err
	}
	return &domain.SimpleTeam{
		Code: code,
		Name: name,
	}, nil
}

func (s *championshipService) hydrateChampionshipStats(ctx context.Context, stats *domain.ChampionshipsStats, language string) error {
	if stats == nil {
		return nil
	}
	var err error
	stats.RunnerUp, err = s.resolvePodiumTeam(ctx, stats.RunnerUpCode, language)
	if err != nil {
		return err
	}
	stats.ThirdPlace, err = s.resolvePodiumTeam(ctx, stats.ThirdPlaceCode, language)
	if err != nil {
		return err
	}
	stats.FourthPlace, err = s.resolvePodiumTeam(ctx, stats.FourthPlaceCode, language)
	return err
}

func (s *championshipService) resolvePodiumTeam(ctx context.Context, teamCode string, language string) (*domain.SimpleTeam, error) {
	code := strings.ToUpper(strings.TrimSpace(teamCode))
	if code == "" {
		return nil, nil
	}
	name, err := s.resolveTeamName(ctx, code, language)
	if err != nil {
		return nil, err
	}
	return &domain.SimpleTeam{
		Code: code,
		Name: name,
	}, nil
}

func (s *championshipService) resolveTeamName(ctx context.Context, code string, language string) (string, error) {
	return s.teamNameResolver.Resolve(ctx, code, language)
}
