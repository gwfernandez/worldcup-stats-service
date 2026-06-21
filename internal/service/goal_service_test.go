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

type mockGoalRepository struct {
	mock.Mock
}

func (m *mockGoalRepository) ListByPlayer(ctx context.Context, filter domain.GoalFilter) ([]domain.Goal, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Goal), args.Get(1).(int64), args.Error(2)
}

func TestGoalService_ListByPlayer(t *testing.T) {
	ctx := context.Background()

	t.Run("success hydrates opponent and builds pagination", func(t *testing.T) {
		repo := new(mockGoalRepository)
		resolver := new(mockTeamNameResolver)
		svc := service.NewGoalService(repo, resolver)
		filter := domain.GoalFilter{PlayerID: 1524, Year: 2018, Language: "es", Page: 1, Size: 2}
		goals := []domain.Goal{
			{Year: 2018, OpponentTeam: domain.SimpleTeam{Code: "ISL"}, MinuteRegular: 64},
			{Year: 2018, OpponentTeam: domain.SimpleTeam{Code: "CRO"}, MinuteRegular: 80},
		}
		repo.On("ListByPlayer", ctx, filter).Return(goals, int64(3), nil)
		resolver.On("Resolve", ctx, "ISL", "es").Return("Islandia", nil).Once()
		resolver.On("Resolve", ctx, "CRO", "es").Return("Croacia", nil).Once()

		result, err := svc.ListByPlayer(ctx, filter)
		require.NoError(t, err)
		assert.Equal(t, "Islandia", result.Data[0].OpponentTeam.Name)
		assert.Equal(t, "Croacia", result.Data[1].OpponentTeam.Name)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 2, result.Pagination.Size)
		assert.Equal(t, int64(3), result.Pagination.TotalElements)
		assert.Equal(t, 2, result.Pagination.TotalPages)
		assert.True(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrevious)
		repo.AssertExpectations(t)
		resolver.AssertExpectations(t)
	})

	t.Run("nil results become empty data", func(t *testing.T) {
		repo := new(mockGoalRepository)
		resolver := new(mockTeamNameResolver)
		svc := service.NewGoalService(repo, resolver)
		filter := domain.GoalFilter{PlayerID: 10, Page: 2, Size: 20}
		repo.On("ListByPlayer", ctx, filter).Return(nil, int64(0), nil)

		result, err := svc.ListByPlayer(ctx, filter)
		require.NoError(t, err)
		assert.NotNil(t, result.Data)
		assert.Empty(t, result.Data)
		assert.Equal(t, 0, result.Pagination.TotalPages)
		assert.False(t, result.Pagination.HasNext)
		assert.True(t, result.Pagination.HasPrevious)
		repo.AssertExpectations(t)
	})

	for _, tc := range []struct {
		name   string
		filter domain.GoalFilter
	}{
		{name: "invalid playerId", filter: domain.GoalFilter{PlayerID: 0, Page: 1, Size: 20}},
		{name: "invalid year", filter: domain.GoalFilter{PlayerID: 1, Year: -1, Page: 1, Size: 20}},
		{name: "invalid page", filter: domain.GoalFilter{PlayerID: 1, Page: 0, Size: 20}},
		{name: "invalid size", filter: domain.GoalFilter{PlayerID: 1, Page: 1, Size: 101}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockGoalRepository)
			resolver := new(mockTeamNameResolver)
			svc := service.NewGoalService(repo, resolver)

			result, err := svc.ListByPlayer(ctx, tc.filter)
			assert.ErrorIs(t, err, domain.ErrInvalidInput)
			assert.Nil(t, result)
		})
	}

	t.Run("repository error", func(t *testing.T) {
		repo := new(mockGoalRepository)
		resolver := new(mockTeamNameResolver)
		svc := service.NewGoalService(repo, resolver)
		filter := domain.GoalFilter{PlayerID: 10, Page: 1, Size: 20}
		repo.On("ListByPlayer", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		result, err := svc.ListByPlayer(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("resolver error", func(t *testing.T) {
		repo := new(mockGoalRepository)
		resolver := new(mockTeamNameResolver)
		svc := service.NewGoalService(repo, resolver)
		filter := domain.GoalFilter{PlayerID: 10, Language: "en", Page: 1, Size: 20}
		repo.On("ListByPlayer", ctx, filter).Return([]domain.Goal{{OpponentTeam: domain.SimpleTeam{Code: "ARG"}}}, int64(1), nil)
		resolver.On("Resolve", ctx, "ARG", "en").Return("", errors.New("cache error")).Once()

		result, err := svc.ListByPlayer(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
		resolver.AssertExpectations(t)
	})
}
