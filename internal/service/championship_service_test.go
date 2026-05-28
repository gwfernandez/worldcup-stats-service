package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// MockChampionshipRepository is a mock implementation of ChampionshipRepository.
type MockChampionshipRepository struct {
	mock.Mock
}

func (m *MockChampionshipRepository) List(ctx context.Context, filter domain.ChampionshipFilter) ([]domain.Championship, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Championship), args.Get(1).(int64), args.Error(2)
}

func (m *MockChampionshipRepository) GetByYear(ctx context.Context, year int) (*domain.Championship, error) {
	args := m.Called(ctx, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Championship), args.Error(1)
}

func TestChampionshipService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipFilter{Page: 1, Size: 20}

		expected := []domain.Championship{{Year: 1930, HostCodes: []string{"URU"}}, {Year: 1934, HostCodes: []string{"ITA"}}}
		mockRepo.On("List", ctx, filter).Return(expected, int64(22), nil)

		res, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1, res.Pagination.Page)
		assert.Equal(t, 20, res.Pagination.Size)
		assert.Equal(t, int64(22), res.Pagination.TotalElements)
		assert.Len(t, res.Data, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid pagination page", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipFilter{Page: 0, Size: 20}

		res, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too small", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipFilter{Page: 1, Size: 0}

		res, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too large", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipFilter{Page: 1, Size: 101}

		res, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipFilter{Page: 1, Size: 20}

		mockRepo.On("List", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		res, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestChampionshipService_GetByYear(t *testing.T) {
	ctx := context.Background()

	t.Run("success with stats", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)

		expectedStats := &domain.ChampionshipsStats{
			TotalTeams: 13,
			TotalGoals: 70,
		}
		expected := &domain.Championship{Year: 1930, Stats: expectedStats}
		mockRepo.On("GetByYear", ctx, 1930).Return(expected, nil)

		res, err := svc.GetByYear(ctx, 1930)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1930, res.Year)
		require.NotNil(t, res.Stats)
		assert.Equal(t, int32(13), res.Stats.TotalTeams)
		assert.Equal(t, int32(70), res.Stats.TotalGoals)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success fallback to default stats", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)

		expected := &domain.Championship{Year: 2026, Stats: nil}
		mockRepo.On("GetByYear", ctx, 2026).Return(expected, nil)

		res, err := svc.GetByYear(ctx, 2026)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 2026, res.Year)
		require.NotNil(t, res.Stats)
		assert.Equal(t, int32(0), res.Stats.TotalTeams)
		assert.Equal(t, int32(0), res.Stats.TotalMatches)
		assert.Equal(t, "", res.Stats.RunnerUpCode)
		assert.Len(t, res.Stats.TopScorers, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)

		mockRepo.On("GetByYear", ctx, 1999).Return(nil, nil)

		res, err := svc.GetByYear(ctx, 1999)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)

		mockRepo.On("GetByYear", ctx, 2022).Return(nil, errors.New("db error"))

		res, err := svc.GetByYear(ctx, 2022)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}
