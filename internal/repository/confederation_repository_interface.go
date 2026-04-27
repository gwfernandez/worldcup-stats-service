package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ConfederationRepository defines the contract for confederation data access.
type ConfederationRepository interface {
	List(ctx context.Context) ([]domain.Confederation, error)
	GetByID(ctx context.Context, id int64) (*domain.Confederation, error)
	Create(ctx context.Context, code, name string) (*domain.Confederation, error)
	Update(ctx context.Context, id int64, code, name string) (*domain.Confederation, error)
	Delete(ctx context.Context, id int64) error
}
