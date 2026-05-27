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

// MockChampionshipService mocks the ChampionshipService interface.
type MockChampionshipService struct {
	mock.Mock
}

func (m *MockChampionshipService) List(ctx context.Context, filter domain.ChampionshipFilter) (*domain.ChampionshipListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChampionshipListResponse), args.Error(1)
}

func (m *MockChampionshipService) GetByYear(ctx context.Context, year int) (*domain.Championship, error) {
	args := m.Called(ctx, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Championship), args.Error(1)
}

func setupChampionshipRouter(svc *MockChampionshipService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewChampionshipHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestChampionshipHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		expected := &domain.ChampionshipListResponse{
			Data: []domain.Championship{{Year: 1930, StartDate: "1930-07-13", EndDate: "1930-07-30", HostNationCodes: []string{"URU"}}},
			Pagination: domain.PaginationInfo{
				Page:          1,
				Size:          20,
				TotalElements: 1,
				TotalPages:    1,
				HasNext:       false,
				HasPrevious:   false,
			},
		}

		svc.On("List", mock.Anything, domain.ChampionshipFilter{
			Year:              1930,
			Host:              "uru",
			ConfederationCode: "CONMEBOL",
			Page:              1,
			Size:              20,
		}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships?year=1930&host=uru&confederation_code=CONMEBOL", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("invalid page", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid size", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships?size=200", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid year param", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships?year=not-a-number", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service invalid input error", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		svc.On("List", mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidInput)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service internal error", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		svc.On("List", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/championships", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestChampionshipHandler_GetByYear(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		expected := &domain.Championship{
			Year: 1930,
			Stats: &domain.ChampionshipStats{
				TotalTeams: 13,
			},
		}

		svc.On("GetByYear", mock.Anything, 1930).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("invalid path year", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/invalid-year", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		svc.On("GetByYear", mock.Anything, 1999).Return(nil, domain.ErrNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		svc.On("GetByYear", mock.Anything, 2022).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/2022", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}
