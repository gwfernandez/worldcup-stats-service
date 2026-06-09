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

// MockStandingService mocks the StandingService interface.
type MockStandingService struct {
	mock.Mock
}

func (m *MockStandingService) List(ctx context.Context, filter domain.StandingFilter) (*domain.StandingListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StandingListResponse), args.Error(1)
}

func setupStandingRouter(svc *MockStandingService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewStandingHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestStandingHandler_List(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		svc := new(MockStandingService)
		r := setupStandingRouter(svc)
		expected := &domain.StandingListResponse{
			Data: []domain.Standing{{
				TeamCode:        "BRA",
				TeamName:        "Brasil",
				MatchesPlayed:   114,
				Wins:            79,
				Draws:           14,
				Losses:          21,
				GoalsFor:        237,
				GoalsAgainst:    108,
				GoalDifference:  129,
				Points:          193,
				UnifiedPoints:   237,
				Position:        1,
				UnifiedPosition: 1,
			}},
			Pagination: domain.PaginationInfo{
				Page:          1,
				Size:          10,
				TotalElements: 1,
				TotalPages:    1,
				HasNext:       false,
				HasPrevious:   false,
			},
		}
		filter := domain.StandingFilter{Name: "bra", Language: "es", ConfederationCode: "CONMEBOL", Page: 1, Size: 10}
		svc.On("List", mock.Anything, filter).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/standings?page=1&size=10&name=bra&confederationCode=CONMEBOL", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"teamCode": "BRA",
				"teamName": "Brasil",
				"matchesPlayed": 114,
				"wins": 79,
				"draws": 14,
				"losses": 21,
				"goalsFor": 237,
				"goalsAgainst": 108,
				"goalDifference": 129,
				"points": 193,
				"unifiedPoints": 237,
				"position": 1,
				"unifiedPosition": 1
			}],
			"pagination": {
				"page": 1,
				"size": 10,
				"totalElements": 1,
				"totalPages": 1,
				"hasNext": false,
				"hasPrevious": false
			}
		}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("success with defaults and empty data", func(t *testing.T) {
		svc := new(MockStandingService)
		r := setupStandingRouter(svc)
		expected := &domain.StandingListResponse{
			Data: []domain.Standing{},
			Pagination: domain.PaginationInfo{
				Page:          1,
				Size:          20,
				TotalElements: 0,
				TotalPages:    0,
				HasNext:       false,
				HasPrevious:   false,
			},
		}
		svc.On("List", mock.Anything, domain.StandingFilter{Language: "es", Page: 1, Size: 20}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/standings", nil)
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
		svc := new(MockStandingService)
		r := setupStandingRouter(svc)
		expected := &domain.StandingListResponse{
			Data:       []domain.Standing{{TeamCode: "GER", TeamName: "Germany", Position: 2}},
			Pagination: domain.PaginationInfo{Page: 1, Size: 20, TotalElements: 1, TotalPages: 1},
		}
		svc.On("List", mock.Anything, domain.StandingFilter{Name: "ger", Language: "en", Page: 1, Size: 20}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/standings?name=ger", nil)
		req.Header.Set("Accept-Language", "en")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"teamCode": "GER",
				"teamName": "Germany",
				"matchesPlayed": 0,
				"wins": 0,
				"draws": 0,
				"losses": 0,
				"goalsFor": 0,
				"goalsAgainst": 0,
				"goalDifference": 0,
				"points": 0,
				"unifiedPoints": 0,
				"position": 2,
				"unifiedPosition": 0
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
		svc := new(MockStandingService)
		r := setupStandingRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/standings?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid page parameter"}`, w.Body.String())
	})

	t.Run("bad request non numeric page", func(t *testing.T) {
		svc := new(MockStandingService)
		r := setupStandingRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/standings?page=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid page parameter"}`, w.Body.String())
	})

	t.Run("bad request invalid size", func(t *testing.T) {
		svc := new(MockStandingService)
		r := setupStandingRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/standings?size=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid size parameter"}`, w.Body.String())
	})

	t.Run("bad request size greater than max", func(t *testing.T) {
		svc := new(MockStandingService)
		r := setupStandingRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/standings?size=101", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid size parameter"}`, w.Body.String())
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockStandingService)
		r := setupStandingRouter(svc)
		svc.On("List", mock.Anything, domain.StandingFilter{Language: "es", Page: 1, Size: 20}).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/standings", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to retrieve standings"}`, w.Body.String())
		svc.AssertExpectations(t)
	})
}
