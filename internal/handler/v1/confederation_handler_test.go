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

func (m *MockConfederationService) List(ctx context.Context, language string) ([]domain.Confederation, error) {
	args := m.Called(ctx, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Confederation), args.Error(1)
}

func (m *MockConfederationService) GetByCode(ctx context.Context, code, language string) (*domain.Confederation, error) {
	args := m.Called(ctx, code, language)
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
			{Code: "CONMEBOL", Name: "South America"},
		}

		svc.On("List", mock.Anything, "es").Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"code":"CONMEBOL","name":"South America"}]`, w.Body.String())
	})

	t.Run("uses english accept language", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		expected := []domain.Confederation{
			{Code: "CONMEBOL", Name: "South American Football Confederation"},
		}

		svc.On("List", mock.Anything, "en").Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations", nil)
		req.Header.Set("Accept-Language", "en-US")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"code":"CONMEBOL","name":"South American Football Confederation"}]`, w.Body.String())
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("List", mock.Anything, "es").Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestConfederationHandler_GetByCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		expected := &domain.Confederation{Name: "South America", Code: "CONMEBOL"}
		svc.On("GetByCode", mock.Anything, "CONMEBOL", "es").Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations/CONMEBOL", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"code":"CONMEBOL","name":"South America"}`, w.Body.String())
	})

	t.Run("uses english accept language", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		expected := &domain.Confederation{Name: "South American Football Confederation", Code: "CONMEBOL"}
		svc.On("GetByCode", mock.Anything, "CONMEBOL", "en").Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations/CONMEBOL", nil)
		req.Header.Set("Accept-Language", "en")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"code":"CONMEBOL","name":"South American Football Confederation"}`, w.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("GetByCode", mock.Anything, "ANTARCTICA", "es").Return(nil, domain.ErrNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations/ANTARCTICA", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("GetByCode", mock.Anything, "ANTARCTICA", "es").Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/confederations/ANTARCTICA", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
