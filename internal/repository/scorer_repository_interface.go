package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ScorerRepository defines the contract for historical scorers data access.
type ScorerRepository interface {
	List(ctx context.Context, filter domain.ScorerFilter) ([]domain.Scorer, int64, error)
	GetByID(ctx context.Context, playerID int64) (*domain.ScorerDetail, error)
}
