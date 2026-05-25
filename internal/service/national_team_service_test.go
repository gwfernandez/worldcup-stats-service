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

// MockNationalTeamRepository is a mock implementation of NationalTeamRepository.
type MockNationalTeamRepository struct {
	mock.Mock
}

func (m *MockNationalTeamRepository) List(ctx context.Context, filter domain.NationalTeamFilter) ([]domain.NationalTeam, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.NationalTeam), args.Get(1).(int64), args.Error(2)
}

func (m *MockNationalTeamRepository) GetByID(ctx context.Context, id int64) (*domain.NationalTeam, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NationalTeam), args.Error(1)
}

func (m *MockNationalTeamRepository) GetByCode(ctx context.Context, code string) (*domain.NationalTeam, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NationalTeam), args.Error(1)
}

func TestNationalTeamService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		filter := domain.NationalTeamFilter{Page: 1, Size: 20}

		expectedTeams := []domain.NationalTeam{{ID: 1, Name: "Argentina"}, {ID: 2, Name: "Brazil"}}
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
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		filter := domain.NationalTeamFilter{Page: 1, Size: 20}

		mockRepo.On("List", ctx, filter).Return([]domain.NationalTeam{}, int64(0), nil)

		result, err := svc.List(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, 0, result.Pagination.TotalPages)
		assert.False(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrevious)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		filter := domain.NationalTeamFilter{Page: 1, Size: 20}

		mockRepo.On("List", ctx, filter).Return(nil, int64(0), errors.New("db error"))

		result, err := svc.List(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestNationalTeamService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		expected := &domain.NationalTeam{ID: 1, Code: "arg", FederationCode: "afa"}
		mockRepo.On("GetByID", ctx, int64(1)).Return(expected, nil)

		result, err := svc.GetByID(ctx, 1)
		assert.NoError(t, err)
		require := assert.New(t)
		require.NotNil(result)
		assert.Equal(t, "ARG", result.Code)
		assert.Equal(t, "AFA", result.FederationCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		mockRepo.On("GetByID", ctx, int64(99)).Return(nil, nil)

		result, err := svc.GetByID(ctx, 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		mockRepo.On("GetByID", ctx, int64(1)).Return(nil, errors.New("db error"))

		result, err := svc.GetByID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestNationalTeamService_GetByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		expected := &domain.NationalTeam{ID: 5, Code: "urs", FederationCode: "ffsu"}
		mockRepo.On("GetByCode", ctx, "urs").Return(expected, nil)

		result, err := svc.GetByCode(ctx, "urs")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "URS", result.Code)
		assert.Equal(t, "FFSU", result.FederationCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		mockRepo.On("GetByCode", ctx, "zzz").Return(nil, nil)

		result, err := svc.GetByCode(ctx, "zzz")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockNationalTeamRepository)
		svc := service.NewNationalTeamService(mockRepo)
		mockRepo.On("GetByCode", ctx, "urs").Return(nil, errors.New("db error"))

		result, err := svc.GetByCode(ctx, "urs")
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
