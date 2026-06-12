package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

const defaultHostLanguage = "es"

// championshipService implements ChampionshipService with business logic.
type championshipService struct {
	repo               repository.ChampionshipRepository
	hostCacheMu        sync.RWMutex
	hostCacheLoaded    bool
	hostNameByLanguage map[string]map[string]string
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
	if err := s.ensureHostCache(ctx); err != nil {
		return nil, err
	}
	championship.Hosts = s.resolveHosts(championship.HostCodes, language)
	championship.Champion = s.resolveChampion(championship.ChampionCode, language)
	s.hydrateChampionshipStats(championship.Stats, language)

	return championship, nil
}

func (s *championshipService) hydrateChampionshipHosts(ctx context.Context, championships []domain.Championship, language string) error {
	if err := s.ensureHostCache(ctx); err != nil {
		return err
	}
	for i := range championships {
		championships[i].Hosts = s.resolveHosts(championships[i].HostCodes, language)
		championships[i].Champion = s.resolveChampion(championships[i].ChampionCode, language)
	}
	return nil
}

func (s *championshipService) ensureHostCache(ctx context.Context) error {
	s.hostCacheMu.RLock()
	if s.hostCacheLoaded {
		s.hostCacheMu.RUnlock()
		return nil
	}
	s.hostCacheMu.RUnlock()

	s.hostCacheMu.Lock()
	defer s.hostCacheMu.Unlock()
	if s.hostCacheLoaded {
		return nil
	}

	translations, err := s.repo.ListTeamTranslations(ctx)
	if err != nil {
		return err
	}
	hostNameByLanguage := make(map[string]map[string]string)
	for _, translation := range translations {
		language := normalizeHostLanguage(translation.Language)
		code := strings.ToUpper(strings.TrimSpace(translation.TeamCode))
		if language == "" || code == "" {
			continue
		}
		if _, ok := hostNameByLanguage[language]; !ok {
			hostNameByLanguage[language] = make(map[string]string)
		}
		hostNameByLanguage[language][code] = translation.Name
	}
	s.hostNameByLanguage = hostNameByLanguage
	s.hostCacheLoaded = true

	return nil
}

func (s *championshipService) resolveHosts(hostCodes []string, language string) []domain.Host {
	hosts := make([]domain.Host, 0, len(hostCodes))
	for _, hostCode := range hostCodes {
		code := strings.ToUpper(strings.TrimSpace(hostCode))
		if code == "" {
			continue
		}
		hosts = append(hosts, domain.Host{
			Code: code,
			Name: s.resolveHostName(code, language),
		})
	}
	return hosts
}

func (s *championshipService) resolveChampion(championCode *string, language string) *domain.ChampionshipChampion {
	if championCode == nil {
		return nil
	}
	code := strings.ToUpper(strings.TrimSpace(*championCode))
	if code == "" {
		return nil
	}
	return &domain.ChampionshipChampion{
		Code: code,
		Name: s.resolveTeamName(code, language),
	}
}

func (s *championshipService) hydrateChampionshipStats(stats *domain.ChampionshipsStats, language string) {
	if stats == nil {
		return
	}
	stats.RunnerUp = s.resolvePodiumTeam(stats.RunnerUpCode, language)
	stats.ThirdPlace = s.resolvePodiumTeam(stats.ThirdPlaceCode, language)
	stats.FourthPlace = s.resolvePodiumTeam(stats.FourthPlaceCode, language)
}

func (s *championshipService) resolvePodiumTeam(teamCode string, language string) *domain.PodiumTeam {
	code := strings.ToUpper(strings.TrimSpace(teamCode))
	if code == "" {
		return nil
	}
	return &domain.PodiumTeam{
		Code: code,
		Name: s.resolveTeamName(code, language),
	}
}

func (s *championshipService) resolveHostName(code string, language string) string {
	return s.resolveTeamName(code, language)
}

func (s *championshipService) resolveTeamName(code string, language string) string {
	s.hostCacheMu.RLock()
	defer s.hostCacheMu.RUnlock()

	if namesByCode := s.hostNameByLanguage[normalizeHostLanguage(language)]; namesByCode != nil {
		if name := namesByCode[code]; name != "" {
			return name
		}
	}
	if namesByCode := s.hostNameByLanguage[defaultHostLanguage]; namesByCode != nil {
		if name := namesByCode[code]; name != "" {
			return name
		}
	}
	return code
}

func normalizeHostLanguage(language string) string {
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		return defaultHostLanguage
	}
	return language
}
