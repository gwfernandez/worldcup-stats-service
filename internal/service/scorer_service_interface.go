package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// ScorerService defines the contract for historical scorers business logic.
type ScorerService interface {
	List(ctx context.Context, filter domain.ScorerFilter) (*domain.ScorerListResponse, error)
	GetByID(ctx context.Context, playerID int64, language string) (*domain.ScorerDetail, error)
}
