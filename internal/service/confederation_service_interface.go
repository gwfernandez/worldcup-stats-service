package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ConfederationService defines the contract for confederation business logic.
type ConfederationService interface {
	List(ctx context.Context, language string) ([]domain.Confederation, error)
	GetByCode(ctx context.Context, code, language string) (*domain.Confederation, error)
}
