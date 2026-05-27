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

func (m *MockConfederationRepository) GetByCode(ctx context.Context, code string) (*domain.Confederation, error) {
	args := m.Called(ctx, code)
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
		{Code: "CONMEBOL", Name: "South America"},
	}

	mockRepo.On("List", ctx).Return(expected, nil)

	result, err := svc.List(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestConfederationService_GetByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		expected := &domain.Confederation{Name: "South America", Code: "CONMEBOL"}
		mockRepo.On("GetByCode", ctx, "CONMEBOL").Return(expected, nil)

		result, err := svc.GetByCode(ctx, "CONMEBOL")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		mockRepo.On("GetByCode", ctx, "ZZZ").Return(nil, nil)

		result, err := svc.GetByCode(ctx, "ZZZ")

		assert.Error(t, err)
		assert.Equal(t, "resource not found: confederation with code ZZZ not found", err.Error())
		assert.True(t, errors.Is(err, domain.ErrNotFound))
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		mockRepo.On("GetByCode", ctx, "ZZZ").Return(nil, errors.New("db error"))

		result, err := svc.GetByCode(ctx, "ZZZ")

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
