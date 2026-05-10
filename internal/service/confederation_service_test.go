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

// MockConfederationRepository is a mock implementation of ConfederationRepository
type MockConfederationRepository struct {
	mock.Mock
}

func (m *MockConfederationRepository) List(ctx context.Context) ([]domain.Confederation, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Confederation), args.Error(1)
}

func (m *MockConfederationRepository) GetByID(ctx context.Context, id int64) (*domain.Confederation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Confederation), args.Error(1)
}

func TestConfederationService_List(t *testing.T) {
	mockRepo := new(MockConfederationRepository)
	svc := service.NewConfederationService(mockRepo)
	ctx := context.Background()

	expected := []domain.Confederation{
		{ID: 1, Code: "CONMEBOL", Name: "South America"},
	}

	mockRepo.On("List", ctx).Return(expected, nil)

	result, err := svc.List(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestConfederationService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		expected := &domain.Confederation{ID: 1, Code: "CONMEBOL"}
		mockRepo.On("GetByID", ctx, int64(1)).Return(expected, nil)

		result, err := svc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		mockRepo.On("GetByID", ctx, int64(99)).Return(nil, nil)

		result, err := svc.GetByID(ctx, 99)

		assert.Error(t, err)
		assert.Equal(t, "confederation with id 99 not found", err.Error())
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		mockRepo.On("GetByID", ctx, int64(1)).Return(nil, errors.New("db error"))

		result, err := svc.GetByID(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
