package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ConfederationService defines the contract for confederation business logic.
type ConfederationService interface {
	List(ctx context.Context) ([]domain.Confederation, error)
	GetByID(ctx context.Context, id int64) (*domain.Confederation, error)
}
