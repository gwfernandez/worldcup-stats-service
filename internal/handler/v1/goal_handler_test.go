package v1_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	v1 "github.com/jendrix/worldcup-stats-service/internal/handler/v1"
	"github.com/jendrix/worldcup-stats-service/internal/middleware"
)

type mockGoalService struct {
	mock.Mock
}

func (m *mockGoalService) ListByPlayer(ctx context.Context, filter domain.GoalFilter) (*domain.GoalListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GoalListResponse), args.Error(1)
}

func setupGoalRouter(svc *mockGoalService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := v1.NewGoalHandler(svc)
	rg := r.Group("/api", middleware.Versioning())
	h.RegisterRoutes(rg)
	return r
}

func TestGoalHandler_List(t *testing.T) {
	t.Run("success with year and pagination", func(t *testing.T) {
		svc := new(mockGoalService)
		r := setupGoalRouter(svc)
		matchDate := "2018-06-16"
		stage := "group_stage"
		penalty := true
		expected := &domain.GoalListResponse{
			Data: []domain.Goal{{
				Year:          2018,
				Hosts:         []domain.SimpleTeam{{Code: "RUS", Name: "Russia"}},
				MatchDate:     &matchDate,
				OpponentTeam:  domain.SimpleTeam{Code: "ISL", Name: "Iceland"},
				MinuteRegular: 64,
				Penalty:       &penalty,
				Stage:         &stage,
			}},
			Pagination: domain.PaginationInfo{Page: 2, Size: 10, TotalElements: 11, TotalPages: 2, HasPrevious: true},
		}
		filter := domain.GoalFilter{PlayerID: 1524, Year: 2018, Language: "en", Page: 2, Size: 10}
		svc.On("ListByPlayer", mock.Anything, filter).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/players/1524/goals?year=2018&page=2&size=10", nil)
		req.Header.Set("Accept-Language", "en")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "1", w.Header().Get("API-Version-Used"))
		assert.JSONEq(t, `{
			"data": [{
				"year": 2018,
				"hosts": [{"code": "RUS", "name": "Russia"}],
				"matchDate": "2018-06-16",
				"opponentTeam": {"code": "ISL", "name": "Iceland"},
				"minuteRegular": 64,
				"penalty": true,
				"stage": "group_stage"
			}],
			"pagination": {
				"page": 2,
				"size": 10,
				"totalElements": 11,
				"totalPages": 2,
				"hasNext": false,
				"hasPrevious": true
			}
		}`, w.Body.String())
		body := w.Body.String()
		assert.Less(t, strings.Index(body, `"year"`), strings.Index(body, `"hosts"`))
		assert.Less(t, strings.Index(body, `"hosts"`), strings.Index(body, `"matchDate"`))
		svc.AssertExpectations(t)
	})

	t.Run("success with defaults and nullable fields", func(t *testing.T) {
		svc := new(mockGoalService)
		r := setupGoalRouter(svc)
		expected := &domain.GoalListResponse{
			Data: []domain.Goal{{
				Year:          1930,
				Hosts:         []domain.SimpleTeam{},
				OpponentTeam:  domain.SimpleTeam{Code: "USA", Name: "Estados Unidos"},
				MinuteRegular: 10,
			}},
			Pagination: domain.PaginationInfo{Page: 1, Size: 20, TotalElements: 1, TotalPages: 1},
		}
		filter := domain.GoalFilter{PlayerID: 1, Language: "es", Page: 1, Size: 20}
		svc.On("ListByPlayer", mock.Anything, filter).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/players/1/goals", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"year": 1930,
				"hosts": [],
				"matchDate": null,
				"opponentTeam": {"code": "USA", "name": "Estados Unidos"},
				"minuteRegular": 10,
				"penalty": null,
				"stage": null
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

	for _, tc := range []struct {
		name  string
		path  string
		error string
	}{
		{name: "non numeric playerId", path: "/api/players/abc/goals", error: "invalid playerId parameter"},
		{name: "non positive playerId", path: "/api/players/0/goals", error: "invalid playerId parameter"},
		{name: "non numeric year", path: "/api/players/1/goals?year=abc", error: "invalid year parameter"},
		{name: "non positive year", path: "/api/players/1/goals?year=0", error: "invalid year parameter"},
		{name: "invalid page", path: "/api/players/1/goals?page=0", error: "invalid page parameter"},
		{name: "non numeric page", path: "/api/players/1/goals?page=x", error: "invalid page parameter"},
		{name: "invalid size", path: "/api/players/1/goals?size=0", error: "invalid size parameter"},
		{name: "size over max", path: "/api/players/1/goals?size=101", error: "invalid size parameter"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mockGoalService)
			r := setupGoalRouter(svc)
			req, _ := http.NewRequest(http.MethodGet, tc.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, `{"error":"`+tc.error+`"}`, w.Body.String())
		})
	}

	t.Run("internal error", func(t *testing.T) {
		svc := new(mockGoalService)
		r := setupGoalRouter(svc)
		filter := domain.GoalFilter{PlayerID: 1, Language: "es", Page: 1, Size: 20}
		svc.On("ListByPlayer", mock.Anything, filter).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/players/1/goals", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to retrieve goals"}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("valid year without results", func(t *testing.T) {
		svc := new(mockGoalService)
		r := setupGoalRouter(svc)
		filter := domain.GoalFilter{PlayerID: 1, Year: 2024, Language: "es", Page: 1, Size: 20}
		expected := &domain.GoalListResponse{
			Data:       []domain.Goal{},
			Pagination: domain.PaginationInfo{Page: 1, Size: 20},
		}
		svc.On("ListByPlayer", mock.Anything, filter).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/players/1/goals?year=2024", nil)
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

	t.Run("old route is not registered", func(t *testing.T) {
		svc := new(mockGoalService)
		r := setupGoalRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/goals/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		svc.AssertNotCalled(t, "ListByPlayer", mock.Anything, mock.Anything)
	})
}
