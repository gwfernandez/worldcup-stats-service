package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// MockStandingRepository is a mock implementation of StandingRepository.
type MockStandingRepository struct {
	mock.Mock
}

func (m *MockStandingRepository) List(ctx context.Context, filter domain.StandingFilter) ([]domain.Standing, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Standing), args.Get(1).(int64), args.Error(2)
}

func TestStandingService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success normalizes confederation code", func(t *testing.T) {
		mockRepo := new(MockStandingRepository)
		svc := service.NewStandingService(mockRepo)
		input := domain.StandingFilter{Name: "argen", ConfederationCode: "conmebol", Page: 1, Size: 10}
		expectedFilter := domain.StandingFilter{Name: "argen", ConfederationCode: "CONMEBOL", Page: 1, Size: 10}
		expected := []domain.Standing{{
			Team: domain.SimpleTeam{Code: "ARG", Name: "Argentina"},

			MatchesPlayed:   88,
			Wins:            53,
			Draws:           10,
			Losses:          25,
			GoalsFor:        152,
			GoalsAgainst:    101,
			GoalDifference:  51,
			Points:          133,
			UnifiedPoints:   159,
			Position:        3,
			UnifiedPosition: 3,
		}}
		mockRepo.On("List", ctx, expectedFilter).Return(expected, int64(11), nil)

		result, err := svc.List(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result.Data)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 10, result.Pagination.Size)
		assert.Equal(t, int64(11), result.Pagination.TotalElements)
		assert.Equal(t, 2, result.Pagination.TotalPages)
		assert.True(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty results", func(t *testing.T) {
		mockRepo := new(MockStandingRepository)
		svc := service.NewStandingService(mockRepo)
		filter := domain.StandingFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return([]domain.Standing{}, int64(0), nil)

		result, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		assert.Empty(t, result.Data)
		assert.Equal(t, int64(0), result.Pagination.TotalElements)
		assert.Equal(t, 0, result.Pagination.TotalPages)
		assert.False(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("nil results become empty array", func(t *testing.T) {
		mockRepo := new(MockStandingRepository)
		svc := service.NewStandingService(mockRepo)
		filter := domain.StandingFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return(nil, int64(0), nil)

		result, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, result.Data)
		assert.Empty(t, result.Data)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid page", func(t *testing.T) {
		mockRepo := new(MockStandingRepository)
		svc := service.NewStandingService(mockRepo)

		result, err := svc.List(ctx, domain.StandingFilter{Page: 0, Size: 20})
		assert.ErrorIs(t, err, domain.ErrInvalidInput)
		assert.Nil(t, result)
	})

	t.Run("invalid size", func(t *testing.T) {
		mockRepo := new(MockStandingRepository)
		svc := service.NewStandingService(mockRepo)

		result, err := svc.List(ctx, domain.StandingFilter{Page: 1, Size: 101})
		assert.ErrorIs(t, err, domain.ErrInvalidInput)
		assert.Nil(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockStandingRepository)
		svc := service.NewStandingService(mockRepo)
		filter := domain.StandingFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		result, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
