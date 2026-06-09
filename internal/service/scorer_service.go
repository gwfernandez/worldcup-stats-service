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
	repo repository.ScorerRepository
}

// NewScorerService creates a new ScorerService.
func NewScorerService(repo repository.ScorerRepository) ScorerService {
	return &scorerService{repo: repo}
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
