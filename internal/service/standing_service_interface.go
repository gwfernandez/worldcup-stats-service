package service

import (
	"context"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// StandingService defines the contract for historical standings business logic.
type StandingService interface {
	List(ctx context.Context, filter domain.StandingFilter) (*domain.StandingListResponse, error)
}
