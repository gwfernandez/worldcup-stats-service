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

// MockTeamRepository is a mock implementation of TeamRepository.
type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) List(ctx context.Context, filter domain.TeamFilter) ([]domain.Team, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Team), args.Get(1).(int64), args.Error(2)
}

func (m *MockTeamRepository) GetByCode(ctx context.Context, code, language string) (*domain.Team, error) {
	args := m.Called(ctx, code, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func TestTeamService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		svc := service.NewTeamService(mockRepo)
		filter := domain.TeamFilter{Page: 1, Size: 20}

		expectedTeams := []domain.Team{{Code: "ARG", Name: "Argentina"}, {Code: "BRA", Name: "Brazil"}}
		mockRepo.On("List", ctx, filter).Return(expectedTeams, int64(35), nil)

		result, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 20, result.Pagination.Size)
		assert.Equal(t, int64(35), result.Pagination.TotalElements)
		assert.Equal(t, 2, result.Pagination.TotalPages)
		assert.True(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrevious)
		assert.Len(t, result.Data, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty results", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		svc := service.NewTeamService(mockRepo)
		filter := domain.TeamFilter{Page: 1, Size: 20}

		mockRepo.On("List", ctx, filter).Return([]domain.Team{}, int64(0), nil)

		result, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, 0, result.Pagination.TotalPages)
		assert.False(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		svc := service.NewTeamService(mockRepo)
		filter := domain.TeamFilter{Page: 1, Size: 20}

		mockRepo.On("List", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		result, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTeamService_GetByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		svc := service.NewTeamService(mockRepo)
		expected := &domain.Team{Code: "urs", FederationCode: "ffsu"}
		mockRepo.On("GetByCode", ctx, "urs", "en").Return(expected, nil)

		result, err := svc.GetByCode(ctx, "urs", "en")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "URS", result.Code)
		assert.Equal(t, "FFSU", result.FederationCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		svc := service.NewTeamService(mockRepo)
		mockRepo.On("GetByCode", ctx, "zzz", "es").Return(nil, nil)

		result, err := svc.GetByCode(ctx, "zzz", "es")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		svc := service.NewTeamService(mockRepo)
		mockRepo.On("GetByCode", ctx, "urs", "es").Return(nil, errors.New("db error"))

		result, err := svc.GetByCode(ctx, "urs", "es")
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
