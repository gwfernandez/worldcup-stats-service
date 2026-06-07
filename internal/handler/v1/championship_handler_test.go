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

func (m *MockChampionshipService) ListTeamsByYear(ctx context.Context, filter domain.ChampionshipTeamFilter) (*domain.ChampionshipTeamListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChampionshipTeamListResponse), args.Error(1)
}

func (m *MockChampionshipService) ListScorersByYear(ctx context.Context, filter domain.ChampionshipScorerFilter) (*domain.ChampionshipScorerListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChampionshipScorerListResponse), args.Error(1)
}

func setupChampionshipRouter(svc *MockChampionshipService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewChampionshipHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestChampionshipHandler_ListScorersByYear(t *testing.T) {
	t.Run("success with defaults and filters", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		expected := &domain.ChampionshipScorerListResponse{
			Data: []domain.ChampionshipScorer{{
				FullName: "Guillermo Stabile",
				TeamCode: "ARG",
				Goals:    8,
			}},
			Pagination: domain.PaginationInfo{
				Page:          1,
				Size:          20,
				TotalElements: 1,
				TotalPages:    1,
			},
		}

		svc.On("ListScorersByYear", mock.Anything, domain.ChampionshipScorerFilter{
			Year:     1930,
			Name:     "guille",
			TeamCode: "ARG",
			Page:     1,
			Size:     20,
		}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/scores?name=guille&teamCode=arg", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"fullName": "Guillermo Stabile",
				"teamCode": "ARG",
				"goals": 8
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

	t.Run("success with explicit pagination", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		expected := &domain.ChampionshipScorerListResponse{
			Data: []domain.ChampionshipScorer{},
			Pagination: domain.PaginationInfo{
				Page:          2,
				Size:          10,
				TotalElements: 0,
				TotalPages:    0,
			},
		}

		svc.On("ListScorersByYear", mock.Anything, domain.ChampionshipScorerFilter{
			Year: 1930,
			Page: 2,
			Size: 10,
		}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/scores?page=2&size=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("invalid year", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/abc/scores", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid year parameter"}`, w.Body.String())
	})

	t.Run("invalid page", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/scores?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid page parameter"}`, w.Body.String())
	})

	t.Run("invalid size", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/scores?size=101", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid size parameter"}`, w.Body.String())
	})

	t.Run("service invalid input error", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		svc.On("ListScorersByYear", mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidInput)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/scores", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service internal error", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		svc.On("ListScorersByYear", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/scores", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestChampionshipHandler_ListTeamsByYear(t *testing.T) {
	t.Run("success with defaults and filters", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		expected := &domain.ChampionshipTeamListResponse{
			Data: []domain.ChampionshipTeam{{
				Year:              1930,
				TeamCode:          "ARG",
				ConfederationCode: "CONMEBOL",
				GroupCode:         "1",
				StageReached:      "runner_up",
				Managers:          "Francisco Olazar",
			}},
			Pagination: domain.PaginationInfo{
				Page:          1,
				Size:          20,
				TotalElements: 1,
				TotalPages:    1,
			},
		}

		svc.On("ListTeamsByYear", mock.Anything, domain.ChampionshipTeamFilter{
			Year:              1930,
			Name:              "argentina",
			ConfederationCode: "CONMEBOL",
			GroupCode:         "A",
			Page:              1,
			Size:              20,
		}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/teams?name=argentina&confederationCode=conmebol&groupCode=a", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"year": 1930,
				"teamCode": "ARG",
				"confederationCode": "CONMEBOL",
				"groupCode": "1",
				"stageReached": "runner_up",
				"managers": "Francisco Olazar"
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

	t.Run("invalid year", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/abc/teams", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid year parameter"}`, w.Body.String())
	})

	t.Run("invalid page", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/teams?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid page parameter"}`, w.Body.String())
	})

	t.Run("invalid size", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/teams?size=101", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid size parameter"}`, w.Body.String())
	})

	t.Run("service invalid input error", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		svc.On("ListTeamsByYear", mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidInput)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/teams", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service internal error", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		svc.On("ListTeamsByYear", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930/teams", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestChampionshipHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockChampionshipService)
		r := setupChampionshipRouter(svc)

		championCode := "URU"
		expected := &domain.ChampionshipListResponse{
			Data: []domain.Championship{{Year: 1930, StartDate: "1930-07-13", EndDate: "1930-07-30", HostCodes: []string{"URU"}, ChampionCode: &championCode}},
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

		req, _ := http.NewRequest(http.MethodGet, "/api/championships?year=1930&host=uru&confederationCode=CONMEBOL", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"year": 1930,
				"startDate": "1930-07-13",
				"endDate": "1930-07-30",
				"hostCodes": ["URU"],
				"championCode": "URU"
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
			Year:      1930,
			StartDate: "1930-07-13",
			EndDate:   "1930-07-30",
			HostCodes: []string{"URU"},
			Stats: &domain.ChampionshipsStats{
				TotalTeams:      13,
				TotalMatches:    18,
				TotalStadiums:   3,
				TotalPlayers:    176,
				TotalGoals:      70,
				RunnerUpCode:    "ARG",
				ThirdPlaceCode:  "USA",
				FourthPlaceCode: "YUG",
				TopScorers: []domain.TopScorer{{
					ID:         1,
					Name:       "Guillermo Stabile",
					NationCode: "ARG",
				}},
				TopScorerGoals: 8,
			},
		}

		svc.On("GetByYear", mock.Anything, 1930).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1930", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"year": 1930,
			"startDate": "1930-07-13",
			"endDate": "1930-07-30",
			"hostCodes": ["URU"],
			"championCode": null,
			"stats": {
				"totalTeams": 13,
				"totalMatches": 18,
				"totalStadiums": 3,
				"totalPlayers": 176,
				"totalGoals": 70,
				"runnerUpCode": "ARG",
				"thirdPlaceCode": "USA",
				"fourthPlaceCode": "YUG",
				"topScorers": [{
					"id": 1,
					"name": "Guillermo Stabile",
					"nationCode": "ARG"
				}],
				"topScorerGoals": 8
			}
		}`, w.Body.String())
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
