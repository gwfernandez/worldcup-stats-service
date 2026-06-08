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
	repo repository.ChampionRepository
}

// NewChampionService creates a new ChampionService.
func NewChampionService(repo repository.ChampionRepository) ChampionService {
	return &championService{repo: repo}
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
