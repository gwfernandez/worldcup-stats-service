package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ConfederationService defines the contract for confederation business logic.
type ConfederationService interface {
	List(ctx context.Context) ([]domain.Confederation, error)
	GetByID(ctx context.Context, id int64) (*domain.Confederation, error)
	Create(ctx context.Context, req domain.CreateConfederationRequest) (*domain.Confederation, error)
	Update(ctx context.Context, id int64, req domain.UpdateConfederationRequest) (*domain.Confederation, error)
	Delete(ctx context.Context, id int64) error
}
