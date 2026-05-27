package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ConfederationRepository defines the contract for confederation data access.
type ConfederationRepository interface {
	List(ctx context.Context) ([]domain.Confederation, error)
	GetByCode(ctx context.Context, code string) (*domain.Confederation, error)
}
