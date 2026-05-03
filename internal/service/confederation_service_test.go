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

func (m *MockConfederationRepository) Create(ctx context.Context, code, name string) (*domain.Confederation, error) {
	args := m.Called(ctx, code, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Confederation), args.Error(1)
}

func (m *MockConfederationRepository) Update(ctx context.Context, id int64, code, name string) (*domain.Confederation, error) {
	args := m.Called(ctx, id, code, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Confederation), args.Error(1)
}

func (m *MockConfederationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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

func TestConfederationService_Create(t *testing.T) {
	mockRepo := new(MockConfederationRepository)
	svc := service.NewConfederationService(mockRepo)
	ctx := context.Background()

	req := domain.CreateConfederationRequest{Code: "CONMEBOL", Name: "South America"}
	expected := &domain.Confederation{ID: 1, Code: "CONMEBOL", Name: "South America"}

	mockRepo.On("Create", ctx, "CONMEBOL", "South America").Return(expected, nil)

	result, err := svc.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestConfederationService_Update(t *testing.T) {
	ctx := context.Background()
	req := domain.UpdateConfederationRequest{Code: "CONMEBOL", Name: "South America Mod"}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		expected := &domain.Confederation{ID: 1, Code: "CONMEBOL", Name: "South America Mod"}
		mockRepo.On("Update", ctx, int64(1), "CONMEBOL", "South America Mod").Return(expected, nil)

		result, err := svc.Update(ctx, 1, req)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		mockRepo.On("Update", ctx, int64(99), "CONMEBOL", "South America Mod").Return(nil, nil)

		result, err := svc.Update(ctx, 99, req)

		assert.Error(t, err)
		assert.Equal(t, "confederation with id 99 not found", err.Error())
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		mockRepo.On("Update", ctx, int64(1), "CONMEBOL", "South America Mod").Return(nil, errors.New("db error"))

		result, err := svc.Update(ctx, 1, req)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestConfederationService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		existing := &domain.Confederation{ID: 1, Code: "CONMEBOL"}
		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		mockRepo.On("Delete", ctx, int64(1)).Return(nil)

		err := svc.Delete(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		mockRepo.On("GetByID", ctx, int64(99)).Return(nil, nil)

		err := svc.Delete(ctx, 99)

		assert.Error(t, err)
		assert.Equal(t, "confederation with id 99 not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("db error on get", func(t *testing.T) {
		mockRepo := new(MockConfederationRepository)
		svc := service.NewConfederationService(mockRepo)

		mockRepo.On("GetByID", ctx, int64(1)).Return(nil, errors.New("db error"))

		err := svc.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
