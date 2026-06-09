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

// MockChampionRepository is a mock implementation of ChampionRepository.
type MockChampionRepository struct {
	mock.Mock
}

func (m *MockChampionRepository) List(ctx context.Context, filter domain.ChampionFilter) ([]domain.Champion, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Champion), args.Get(1).(int64), args.Error(2)
}

func TestChampionService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockChampionRepository)
		svc := service.NewChampionService(mockRepo)
		filter := domain.ChampionFilter{Page: 1, Size: 10}
		expected := []domain.Champion{{
			TeamCode: "BRA",
			TeamName: "Brasil",
			Wins:     5,
			Years:    []int32{1958, 1962, 1970, 1994, 2002},
		}}
		mockRepo.On("List", ctx, filter).Return(expected, int64(11), nil)

		result, err := svc.List(ctx, filter)
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
		mockRepo := new(MockChampionRepository)
		svc := service.NewChampionService(mockRepo)
		filter := domain.ChampionFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return([]domain.Champion{}, int64(0), nil)

		result, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		assert.Empty(t, result.Data)
		assert.Equal(t, 0, result.Pagination.TotalPages)
		assert.False(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("nil results become empty array", func(t *testing.T) {
		mockRepo := new(MockChampionRepository)
		svc := service.NewChampionService(mockRepo)
		filter := domain.ChampionFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return(nil, int64(0), nil)

		result, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, result.Data)
		assert.Empty(t, result.Data)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid page", func(t *testing.T) {
		mockRepo := new(MockChampionRepository)
		svc := service.NewChampionService(mockRepo)

		result, err := svc.List(ctx, domain.ChampionFilter{Page: 0, Size: 20})
		assert.ErrorIs(t, err, domain.ErrInvalidInput)
		assert.Nil(t, result)
	})

	t.Run("invalid size", func(t *testing.T) {
		mockRepo := new(MockChampionRepository)
		svc := service.NewChampionService(mockRepo)

		result, err := svc.List(ctx, domain.ChampionFilter{Page: 1, Size: 101})
		assert.ErrorIs(t, err, domain.ErrInvalidInput)
		assert.Nil(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockChampionRepository)
		svc := service.NewChampionService(mockRepo)
		filter := domain.ChampionFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		result, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
