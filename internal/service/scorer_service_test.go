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

// MockScorerRepository is a mock implementation of ScorerRepository.
type MockScorerRepository struct {
	mock.Mock
}

func (m *MockScorerRepository) List(ctx context.Context, filter domain.ScorerFilter) ([]domain.Scorer, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Scorer), args.Get(1).(int64), args.Error(2)
}

func TestScorerService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success normalizes filters", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)
		input := domain.ScorerFilter{Name: "messi", TeamCode: "arg", ConfederationCode: "conmebol", Page: 1, Size: 10}
		expectedFilter := domain.ScorerFilter{Name: "messi", TeamCode: "ARG", ConfederationCode: "CONMEBOL", Page: 1, Size: 10}
		expected := []domain.Scorer{{
			FullName:          "Lionel Messi",
			TeamCode:          "ARG",
			Goals:             13,
			ListTeams:         []string{"ARG"},
			ConfederationCode: "CONMEBOL",
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
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)
		filter := domain.ScorerFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return([]domain.Scorer{}, int64(0), nil)

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
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)
		filter := domain.ScorerFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return(nil, int64(0), nil)

		result, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, result.Data)
		assert.Empty(t, result.Data)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid page", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)

		result, err := svc.List(ctx, domain.ScorerFilter{Page: 0, Size: 20})
		assert.ErrorIs(t, err, domain.ErrInvalidInput)
		assert.Nil(t, result)
	})

	t.Run("invalid size", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)

		result, err := svc.List(ctx, domain.ScorerFilter{Page: 1, Size: 101})
		assert.ErrorIs(t, err, domain.ErrInvalidInput)
		assert.Nil(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)
		filter := domain.ScorerFilter{Page: 1, Size: 20}
		mockRepo.On("List", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		result, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
