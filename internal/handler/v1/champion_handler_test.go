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

// MockChampionService mocks the ChampionService interface.
type MockChampionService struct {
	mock.Mock
}

func (m *MockChampionService) List(ctx context.Context, filter domain.ChampionFilter) (*domain.ChampionListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChampionListResponse), args.Error(1)
}

func setupChampionRouter(svc *MockChampionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewChampionHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestChampionHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockChampionService)
		r := setupChampionRouter(svc)
		expected := &domain.ChampionListResponse{
			Data: []domain.Champion{{
				TeamCode: "BRA",
				Name:     "Brasil",
				Wins:     5,
				Years:    []int32{1958, 1962, 1970, 1994, 2002},
			}},
			Pagination: domain.PaginationInfo{
				Page:          1,
				Size:          10,
				TotalElements: 8,
				TotalPages:    1,
				HasNext:       false,
				HasPrevious:   false,
			},
		}
		svc.On("List", mock.Anything, domain.ChampionFilter{Language: "es", Page: 1, Size: 10}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/champions?page=1&size=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"teamCode": "BRA",
				"name": "Brasil",
				"wins": 5,
				"years": [1958, 1962, 1970, 1994, 2002]
			}],
			"pagination": {
				"page": 1,
				"size": 10,
				"totalElements": 8,
				"totalPages": 1,
				"hasNext": false,
				"hasPrevious": false
			}
		}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("success with defaults and empty data", func(t *testing.T) {
		svc := new(MockChampionService)
		r := setupChampionRouter(svc)
		expected := &domain.ChampionListResponse{
			Data: []domain.Champion{},
			Pagination: domain.PaginationInfo{
				Page:          1,
				Size:          20,
				TotalElements: 0,
				TotalPages:    0,
				HasNext:       false,
				HasPrevious:   false,
			},
		}
		svc.On("List", mock.Anything, domain.ChampionFilter{Language: "es", Page: 1, Size: 20}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/champions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [],
			"pagination": {
				"page": 1,
				"size": 20,
				"totalElements": 0,
				"totalPages": 0,
				"hasNext": false,
				"hasPrevious": false
			}
		}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("uses english accept language", func(t *testing.T) {
		svc := new(MockChampionService)
		r := setupChampionRouter(svc)
		expected := &domain.ChampionListResponse{
			Data:       []domain.Champion{{TeamCode: "GER", Name: "Germany", Wins: 4, Years: []int32{1954, 1974, 1990, 2014}}},
			Pagination: domain.PaginationInfo{Page: 1, Size: 20, TotalElements: 1, TotalPages: 1},
		}
		svc.On("List", mock.Anything, domain.ChampionFilter{Language: "en", Page: 1, Size: 20}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/champions", nil)
		req.Header.Set("Accept-Language", "en")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"teamCode": "GER",
				"name": "Germany",
				"wins": 4,
				"years": [1954, 1974, 1990, 2014]
			}],
			"pagination": {
				"page": 1,
				"size": 20,
				"totalElements": 1,
				"totalPages": 1,
				"hasNext": false,
				"hasPrevious": false
			}
		}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("bad request invalid page", func(t *testing.T) {
		svc := new(MockChampionService)
		r := setupChampionRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/champions?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid page parameter"}`, w.Body.String())
	})

	t.Run("bad request non numeric page", func(t *testing.T) {
		svc := new(MockChampionService)
		r := setupChampionRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/champions?page=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid page parameter"}`, w.Body.String())
	})

	t.Run("bad request invalid size", func(t *testing.T) {
		svc := new(MockChampionService)
		r := setupChampionRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/champions?size=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid size parameter"}`, w.Body.String())
	})

	t.Run("bad request size greater than max", func(t *testing.T) {
		svc := new(MockChampionService)
		r := setupChampionRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/champions?size=101", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid size parameter"}`, w.Body.String())
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockChampionService)
		r := setupChampionRouter(svc)
		svc.On("List", mock.Anything, domain.ChampionFilter{Language: "es", Page: 1, Size: 20}).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/champions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to retrieve champions"}`, w.Body.String())
		svc.AssertExpectations(t)
	})
}
