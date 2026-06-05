package repository

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// GroupStatsRepository defines the contract for group standings data access.
type GroupStatsRepository interface {
	ListByYear(ctx context.Context, year int) ([]domain.GroupStandingRecord, error)
}
