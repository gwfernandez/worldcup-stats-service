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

func (m *MockScorerRepository) GetByID(ctx context.Context, playerID int64) (*domain.ScorerDetail, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScorerDetail), args.Error(1)
}

func TestScorerService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success normalizes filters", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		resolver := new(mockTeamNameResolver)
		svc := service.NewScorerService(mockRepo, resolver)
		input := domain.ScorerFilter{Name: "messi", TeamCode: "arg", ConfederationCode: "conmebol", Page: 1, Size: 10}
		expectedFilter := domain.ScorerFilter{Name: "messi", TeamCode: "ARG", ConfederationCode: "CONMEBOL", Page: 1, Size: 10}
		expected := []domain.Scorer{{
			PlayerID:          10,
			FullName:          "Lionel Messi",
			Team:              domain.SimpleTeam{Code: "ARG"},
			Goals:             13,
			ListTeams:         []string{"ARG"},
			ConfederationCode: "CONMEBOL",
		}}
		mockRepo.On("List", ctx, expectedFilter).Return(expected, int64(11), nil)
		resolver.On("Resolve", ctx, "ARG", "").Return("Argentina", nil).Once()

		result, err := svc.List(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, "Argentina", result.Data[0].Team.Name)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 10, result.Pagination.Size)
		assert.Equal(t, int64(11), result.Pagination.TotalElements)
		assert.Equal(t, 2, result.Pagination.TotalPages)
		assert.True(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
		resolver.AssertExpectations(t)
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

func TestScorerService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success hydrates teams hosts and opponent", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		resolver := new(mockTeamNameResolver)
		svc := service.NewScorerService(mockRepo, resolver)
		position := "FW"
		expected := &domain.ScorerDetail{
			ID:            1524,
			FirstName:     "Lionel",
			LastName:      "Messi",
			Position:      &position,
			Championships: []int32{2006, 2022},
			Teams:         []domain.SimpleTeam{{Code: "ARG"}},
			Goals: []domain.Goal{{
				Year:         2022,
				Hosts:        []domain.SimpleTeam{{Code: "QAT"}},
				OpponentTeam: domain.SimpleTeam{Code: "FRA"},
			}},
		}
		mockRepo.On("GetByID", ctx, int64(1524)).Return(expected, nil)
		resolver.On("Resolve", ctx, "ARG", "es").Return("Argentina", nil).Once()
		resolver.On("Resolve", ctx, "QAT", "es").Return("Catar", nil).Once()
		resolver.On("Resolve", ctx, "FRA", "es").Return("Francia", nil).Once()

		result, err := svc.GetByID(ctx, 1524, "es")
		require.NoError(t, err)
		assert.Equal(t, "Argentina", result.Teams[0].Name)
		assert.Equal(t, "Catar", result.Goals[0].Hosts[0].Name)
		assert.Equal(t, "Francia", result.Goals[0].OpponentTeam.Name)
		mockRepo.AssertExpectations(t)
		resolver.AssertExpectations(t)
	})

	t.Run("nil arrays become empty arrays", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)
		expected := &domain.ScorerDetail{ID: 7, FirstName: "Juan", LastName: "Pérez"}
		mockRepo.On("GetByID", ctx, int64(7)).Return(expected, nil)

		result, err := svc.GetByID(ctx, 7, "es")
		require.NoError(t, err)
		assert.NotNil(t, result.Championships)
		assert.NotNil(t, result.Teams)
		assert.NotNil(t, result.Goals)
		assert.Empty(t, result.Championships)
		assert.Empty(t, result.Teams)
		assert.Empty(t, result.Goals)
		mockRepo.AssertExpectations(t)
	})

	t.Run("nil goal hosts become empty array", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		resolver := new(mockTeamNameResolver)
		svc := service.NewScorerService(mockRepo, resolver)
		expected := &domain.ScorerDetail{
			ID:    7,
			Goals: []domain.Goal{{OpponentTeam: domain.SimpleTeam{Code: "URU"}}},
		}
		mockRepo.On("GetByID", ctx, int64(7)).Return(expected, nil)
		resolver.On("Resolve", ctx, "URU", "en").Return("Uruguay", nil).Once()

		result, err := svc.GetByID(ctx, 7, "en")
		require.NoError(t, err)
		assert.NotNil(t, result.Goals[0].Hosts)
		assert.Empty(t, result.Goals[0].Hosts)
		resolver.AssertExpectations(t)
	})

	t.Run("invalid player id", func(t *testing.T) {
		svc := service.NewScorerService(new(MockScorerRepository))

		result, err := svc.GetByID(ctx, 0, "es")
		assert.ErrorIs(t, err, domain.ErrInvalidInput)
		assert.Nil(t, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

		result, err := svc.GetByID(ctx, 999, "es")
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		svc := service.NewScorerService(mockRepo)
		mockRepo.On("GetByID", ctx, int64(1524)).Return(nil, errors.New("db error"))

		result, err := svc.GetByID(ctx, 1524, "es")
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("resolver error", func(t *testing.T) {
		mockRepo := new(MockScorerRepository)
		resolver := new(mockTeamNameResolver)
		svc := service.NewScorerService(mockRepo, resolver)
		expected := &domain.ScorerDetail{
			ID:    1524,
			Teams: []domain.SimpleTeam{{Code: "ARG"}},
		}
		mockRepo.On("GetByID", ctx, int64(1524)).Return(expected, nil)
		resolver.On("Resolve", ctx, "ARG", "es").Return("", errors.New("cache error")).Once()

		result, err := svc.GetByID(ctx, 1524, "es")
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
		resolver.AssertExpectations(t)
	})
}
