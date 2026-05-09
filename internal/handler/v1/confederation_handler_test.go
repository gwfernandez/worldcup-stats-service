package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func (m *MockConfederationService) Create(ctx context.Context, req domain.CreateConfederationRequest) (*domain.Confederation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Confederation), args.Error(1)
}

func (m *MockConfederationService) Update(ctx context.Context, id int64, req domain.UpdateConfederationRequest) (*domain.Confederation, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Confederation), args.Error(1)
}

func (m *MockConfederationService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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

		svc.On("GetByID", mock.Anything, int64(99)).Return(nil, errors.New("not found"))

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

func TestConfederationHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		createReq := domain.CreateConfederationRequest{Code: "C", Name: "S"}
		expected := &domain.Confederation{ID: 1, Code: "C"}

		svc.On("Create", mock.Anything, createReq).Return(expected, nil)

		body, _ := json.Marshal(createReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/confederations", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		req, _ := http.NewRequest(http.MethodPost, "/api/confederations", bytes.NewBufferString("invalid"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("conflict", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		createReq := domain.CreateConfederationRequest{Code: "C", Name: "S"}
		svc.On("Create", mock.Anything, createReq).Return(nil, errors.New("duplicate key"))

		body, _ := json.Marshal(createReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/confederations", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		createReq := domain.CreateConfederationRequest{Code: "C", Name: "S"}
		svc.On("Create", mock.Anything, createReq).Return(nil, errors.New("db error"))

		body, _ := json.Marshal(createReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/confederations", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestConfederationHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		updateReq := domain.UpdateConfederationRequest{Code: "C", Name: "S"}
		expected := &domain.Confederation{ID: 1, Code: "C", Name: "S"}

		svc.On("Update", mock.Anything, int64(1), updateReq).Return(expected, nil)

		body, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPut, "/api/confederations/1", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		req, _ := http.NewRequest(http.MethodPut, "/api/confederations/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request body", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		req, _ := http.NewRequest(http.MethodPut, "/api/confederations/1", bytes.NewBufferString("invalid"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		updateReq := domain.UpdateConfederationRequest{Code: "C", Name: "S"}
		svc.On("Update", mock.Anything, int64(99), updateReq).Return(nil, errors.New("not found"))

		body, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPut, "/api/confederations/99", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("conflict", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		updateReq := domain.UpdateConfederationRequest{Code: "C", Name: "S"}
		svc.On("Update", mock.Anything, int64(1), updateReq).Return(nil, errors.New("duplicate key"))

		body, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPut, "/api/confederations/1", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		updateReq := domain.UpdateConfederationRequest{Code: "C", Name: "S"}
		svc.On("Update", mock.Anything, int64(1), updateReq).Return(nil, errors.New("db error"))

		body, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPut, "/api/confederations/1", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestConfederationHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("Delete", mock.Anything, int64(1)).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/api/confederations/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		req, _ := http.NewRequest(http.MethodDelete, "/api/confederations/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("Delete", mock.Anything, int64(99)).Return(errors.New("not found"))

		req, _ := http.NewRequest(http.MethodDelete, "/api/confederations/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockConfederationService)
		r := setupRouter(svc)

		svc.On("Delete", mock.Anything, int64(1)).Return(errors.New("db error"))

		req, _ := http.NewRequest(http.MethodDelete, "/api/confederations/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
