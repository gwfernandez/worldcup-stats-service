package v1_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	v1 "github.com/jendrix/worldcup-stats-service/internal/handler/v1"
)

// MockConfederationService mocks the ConfederationService interface
type MockConfederationService struct {
	mock.Mock
}

func (m *MockConfederationService) List(ctx context.Context) ([]domain.Confederation, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Confederation), args.Error(1)
}

func (m *MockConfederationService) GetByID(ctx context.Context, id int64) (*domain.Confederation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Confederation), args.Error(1)
}

func setupRouter(svc *MockConfederationService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewConfederationHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestConfederationHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		expected := []domain.Confederation{
			{ID: 1, Code: "CONMEBOL", Name: "South America"},
		}

		svc.On("List", mock.Anything).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("List", mock.Anything).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestConfederationHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		expected := &domain.Confederation{ID: 1, Code: "CONMEBOL"}
		svc.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("GetByID", mock.Anything, int64(99)).Return(nil, domain.ErrNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("GetByID", mock.Anything, int64(99)).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
