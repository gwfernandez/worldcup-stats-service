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

func (s *confederationService) List(ctx context.Context) ([]domain.Confederation, error) {
	return s.repo.List(ctx)
}

func (s *confederationService) GetByID(ctx context.Context, id int64) (*domain.Confederation, error) {
	confederation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if confederation == nil {
		return nil, fmt.Errorf("confederation with id %d not found", id)
	}
	return confederation, nil
}

func (s *confederationService) Create(ctx context.Context, req domain.CreateConfederationRequest) (*domain.Confederation, error) {
	return s.repo.Create(ctx, req.Code, req.Name)
}

func (s *confederationService) Update(ctx context.Context, id int64, req domain.UpdateConfederationRequest) (*domain.Confederation, error) {
	confederation, err := s.repo.Update(ctx, id, req.Code, req.Name)
	if err != nil {
		return nil, err
	}
	if confederation == nil {
		return nil, fmt.Errorf("confederation with id %d not found", id)
	}
	return confederation, nil
}

func (s *confederationService) Delete(ctx context.Context, id int64) error {
	// Verify existence before deleting
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("confederation with id %d not found", id)
	}
	return s.repo.Delete(ctx, id)
}
