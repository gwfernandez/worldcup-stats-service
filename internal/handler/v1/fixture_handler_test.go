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

type mockFixtureService struct {
	mock.Mock
}

func (m *mockFixtureService) GetByYear(ctx context.Context, year int) (*domain.Fixture, error) {
	args := m.Called(ctx, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Fixture), args.Error(1)
}

func setupFixtureRouter(svc *mockFixtureService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewFixtureHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestFixtureHandler_GetByYear(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mockFixtureService)
		r := setupFixtureRouter(svc)
		expected := &domain.Fixture{Year: 1978, Stages: []domain.FixtureStage{{Stage: "group_stage"}}}

		svc.On("GetByYear", mock.Anything, 1978).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1978/fixture", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"data":{"year":1978,"stages":[{"stage":"group_stage"}]}}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("invalid year", func(t *testing.T) {
		svc := new(mockFixtureService)
		r := setupFixtureRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/abc/fixture", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid year"}`, w.Body.String())
		svc.AssertNotCalled(t, "GetByYear")
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(mockFixtureService)
		r := setupFixtureRouter(svc)

		svc.On("GetByYear", mock.Anything, 2023).Return(nil, domain.ErrNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/2023/fixture", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"championship not found"}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mockFixtureService)
		r := setupFixtureRouter(svc)

		svc.On("GetByYear", mock.Anything, 1978).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1978/fixture", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to retrieve fixture"}`, w.Body.String())
		svc.AssertExpectations(t)
	})
}
