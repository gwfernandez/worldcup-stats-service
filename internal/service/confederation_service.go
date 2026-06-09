package service

import (
	"context"
	"fmt"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

// confederationService implements ConfederationService with business logic.
type confederationService struct {
	repo repository.ConfederationRepository
}

// NewConfederationService creates a new service that depends on the repository interface.
func NewConfederationService(repo repository.ConfederationRepository) ConfederationService {
	return &confederationService{repo: repo}
}

func (s *confederationService) List(ctx context.Context, language string) ([]domain.Confederation, error) {
	return s.repo.List(ctx, language)
}

func (s *confederationService) GetByCode(ctx context.Context, code, language string) (*domain.Confederation, error) {
	confederation, err := s.repo.GetByCode(ctx, code, language)
	if err != nil {
		return nil, err
	}
	if confederation == nil {
		return nil, fmt.Errorf("%w: confederation with code %s not found", domain.ErrNotFound, code)
	}
	return confederation, nil
}
