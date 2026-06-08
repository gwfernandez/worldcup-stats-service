package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// StandingRepository defines the contract for historical standings data access.
type StandingRepository interface {
	List(ctx context.Context, filter domain.StandingFilter) ([]domain.Standing, int64, error)
}
