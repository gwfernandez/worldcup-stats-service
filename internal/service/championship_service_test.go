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

func (m *MockChampionshipRepository) ListTeamsByYear(ctx context.Context, filter domain.ChampionshipTeamFilter) ([]domain.ChampionshipTeam, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.ChampionshipTeam), args.Get(1).(int64), args.Error(2)
}

func (m *MockChampionshipRepository) ListStadiumsByYear(ctx context.Context, filter domain.ChampionshipStadiumFilter) ([]domain.ChampionshipStadium, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.ChampionshipStadium), args.Get(1).(int64), args.Error(2)
}

func (m *MockChampionshipRepository) ListScorersByYear(ctx context.Context, filter domain.ChampionshipScorerFilter) ([]domain.ChampionshipScorer, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.ChampionshipScorer), args.Get(1).(int64), args.Error(2)
}

func (m *MockChampionshipRepository) ListStandingsByYear(ctx context.Context, filter domain.ChampionshipStandingFilter) ([]domain.ChampionshipStanding, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.ChampionshipStanding), args.Get(1).(int64), args.Error(2)
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

func TestChampionshipService_ListTeamsByYear(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipTeamFilter{Year: 1930, Page: 1, Size: 10}

		expected := []domain.ChampionshipTeam{
			{Year: 1930, TeamCode: "URU", ConfederationCode: "CONMEBOL", GroupCode: "3", StageReached: "champion", Managers: "Alberto Suppici"},
			{Year: 1930, TeamCode: "ARG", ConfederationCode: "CONMEBOL", GroupCode: "1", StageReached: "runner_up", Managers: ""},
		}
		mockRepo.On("ListTeamsByYear", ctx, filter).Return(expected, int64(13), nil)

		res, err := svc.ListTeamsByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, 1, res.Pagination.Page)
		assert.Equal(t, 10, res.Pagination.Size)
		assert.Equal(t, int64(13), res.Pagination.TotalElements)
		assert.Equal(t, 2, res.Pagination.TotalPages)
		assert.True(t, res.Pagination.HasNext)
		assert.False(t, res.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty response", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipTeamFilter{Year: 9999, Page: 1, Size: 20}

		mockRepo.On("ListTeamsByYear", ctx, filter).Return(nil, int64(0), nil)

		res, err := svc.ListTeamsByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Empty(t, res.Data)
		assert.NotNil(t, res.Data)
		assert.Equal(t, 0, res.Pagination.TotalPages)
		assert.False(t, res.Pagination.HasNext)
		assert.False(t, res.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid pagination page", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipTeamFilter{Year: 1930, Page: 0, Size: 20}

		res, err := svc.ListTeamsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too small", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipTeamFilter{Year: 1930, Page: 1, Size: 0}

		res, err := svc.ListTeamsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too large", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipTeamFilter{Year: 1930, Page: 1, Size: 101}

		res, err := svc.ListTeamsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipTeamFilter{Year: 1930, Page: 1, Size: 20}

		mockRepo.On("ListTeamsByYear", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		res, err := svc.ListTeamsByYear(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestChampionshipService_ListStandingsByYear(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStandingFilter{Year: 1930, Page: 1, Size: 10}

		expected := []domain.ChampionshipStanding{
			{TeamCode: "URU", GroupCode: "3", MatchesPlayed: 4, Wins: 4, GoalsFor: 15, GoalsAgainst: 3, GoalDifference: 12, Points: 8, UnifiedPoints: 12, Position: 1, Performance: "champion"},
			{TeamCode: "ARG", GroupCode: "1", MatchesPlayed: 5, Wins: 4, Losses: 1, GoalsFor: 18, GoalsAgainst: 9, GoalDifference: 9, Points: 8, UnifiedPoints: 12, Position: 2, Performance: "runner_up"},
		}
		mockRepo.On("ListStandingsByYear", ctx, filter).Return(expected, int64(13), nil)

		res, err := svc.ListStandingsByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, 1, res.Pagination.Page)
		assert.Equal(t, 10, res.Pagination.Size)
		assert.Equal(t, int64(13), res.Pagination.TotalElements)
		assert.Equal(t, 2, res.Pagination.TotalPages)
		assert.True(t, res.Pagination.HasNext)
		assert.False(t, res.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty response", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStandingFilter{Year: 9999, Page: 1, Size: 20}

		mockRepo.On("ListStandingsByYear", ctx, filter).Return(nil, int64(0), nil)

		res, err := svc.ListStandingsByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Empty(t, res.Data)
		assert.NotNil(t, res.Data)
		assert.Equal(t, 0, res.Pagination.TotalPages)
		assert.False(t, res.Pagination.HasNext)
		assert.False(t, res.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("size 100 is valid", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStandingFilter{Year: 2026, Page: 1, Size: 100}

		mockRepo.On("ListStandingsByYear", ctx, filter).Return([]domain.ChampionshipStanding{}, int64(0), nil)

		res, err := svc.ListStandingsByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 100, res.Pagination.Size)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid pagination page", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStandingFilter{Year: 1930, Page: 0, Size: 20}

		res, err := svc.ListStandingsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too small", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStandingFilter{Year: 1930, Page: 1, Size: 0}

		res, err := svc.ListStandingsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too large", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStandingFilter{Year: 1930, Page: 1, Size: 101}

		res, err := svc.ListStandingsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStandingFilter{Year: 1930, Page: 1, Size: 20}

		mockRepo.On("ListStandingsByYear", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		res, err := svc.ListStandingsByYear(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestChampionshipService_ListStadiumsByYear(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStadiumFilter{Year: 1930, Name: "centenario", Page: 1, Size: 10}

		expected := []domain.ChampionshipStadium{
			{Year: 1930, ID: 1, Name: "Estadio Centenario", CityName: "Montevideo", Capacity: 90000, MatchesPlayed: 10},
			{Year: 1930, ID: 2, Name: "Estadio Pocitos", CityName: "Montevideo", Capacity: 1000, MatchesPlayed: 2},
		}
		mockRepo.On("ListStadiumsByYear", ctx, filter).Return(expected, int64(13), nil)

		res, err := svc.ListStadiumsByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, 1, res.Pagination.Page)
		assert.Equal(t, 10, res.Pagination.Size)
		assert.Equal(t, int64(13), res.Pagination.TotalElements)
		assert.Equal(t, 2, res.Pagination.TotalPages)
		assert.True(t, res.Pagination.HasNext)
		assert.False(t, res.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty response", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStadiumFilter{Year: 9999, Page: 1, Size: 20}

		mockRepo.On("ListStadiumsByYear", ctx, filter).Return(nil, int64(0), nil)

		res, err := svc.ListStadiumsByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Empty(t, res.Data)
		assert.NotNil(t, res.Data)
		assert.Equal(t, 0, res.Pagination.TotalPages)
		assert.False(t, res.Pagination.HasNext)
		assert.False(t, res.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("size 100 is valid", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStadiumFilter{Year: 2026, Page: 1, Size: 100}

		mockRepo.On("ListStadiumsByYear", ctx, filter).Return([]domain.ChampionshipStadium{}, int64(0), nil)

		res, err := svc.ListStadiumsByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 100, res.Pagination.Size)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid pagination page", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStadiumFilter{Year: 1930, Page: 0, Size: 20}

		res, err := svc.ListStadiumsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too small", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStadiumFilter{Year: 1930, Page: 1, Size: 0}

		res, err := svc.ListStadiumsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too large", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStadiumFilter{Year: 1930, Page: 1, Size: 101}

		res, err := svc.ListStadiumsByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipStadiumFilter{Year: 1930, Page: 1, Size: 20}

		mockRepo.On("ListStadiumsByYear", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		res, err := svc.ListStadiumsByYear(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestChampionshipService_ListScorersByYear(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipScorerFilter{Year: 1930, TeamCode: "ARG", Page: 1, Size: 10}

		expected := []domain.ChampionshipScorer{
			{FullName: "Guillermo Stabile", TeamCode: "ARG", Goals: 8},
			{FullName: "Carlos Peucelle", TeamCode: "ARG", Goals: 3},
		}
		mockRepo.On("ListScorersByYear", ctx, filter).Return(expected, int64(13), nil)

		res, err := svc.ListScorersByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, 1, res.Pagination.Page)
		assert.Equal(t, 10, res.Pagination.Size)
		assert.Equal(t, int64(13), res.Pagination.TotalElements)
		assert.Equal(t, 2, res.Pagination.TotalPages)
		assert.True(t, res.Pagination.HasNext)
		assert.False(t, res.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty response", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipScorerFilter{Year: 9999, Page: 1, Size: 20}

		mockRepo.On("ListScorersByYear", ctx, filter).Return(nil, int64(0), nil)

		res, err := svc.ListScorersByYear(ctx, filter)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Empty(t, res.Data)
		assert.NotNil(t, res.Data)
		assert.Equal(t, 0, res.Pagination.TotalPages)
		assert.False(t, res.Pagination.HasNext)
		assert.False(t, res.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid pagination page", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipScorerFilter{Year: 1930, Page: 0, Size: 20}

		res, err := svc.ListScorersByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too small", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipScorerFilter{Year: 1930, Page: 1, Size: 0}

		res, err := svc.ListScorersByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("invalid pagination size too large", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipScorerFilter{Year: 1930, Page: 1, Size: 101}

		res, err := svc.ListScorersByYear(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		assert.Nil(t, res)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockChampionshipRepository)
		svc := service.NewChampionshipService(mockRepo)
		filter := domain.ChampionshipScorerFilter{Year: 1930, Page: 1, Size: 20}

		mockRepo.On("ListScorersByYear", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		res, err := svc.ListScorersByYear(ctx, filter)
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
