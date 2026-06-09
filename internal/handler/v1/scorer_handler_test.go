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
	"github.com/jendrix/worldcup-stats-service/internal/middleware"
)

// MockScorerService mocks the ScorerService interface.
type MockScorerService struct {
	mock.Mock
}

func (m *MockScorerService) List(ctx context.Context, filter domain.ScorerFilter) (*domain.ScorerListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScorerListResponse), args.Error(1)
}

func setupScorerRouter(svc *MockScorerService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewScorerHandler(svc)
	rg := r.Group("/api", middleware.Versioning())
	h.RegisterRoutes(rg)
	return r
}

func TestScorerHandler_List(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		svc := new(MockScorerService)
		r := setupScorerRouter(svc)
		expected := &domain.ScorerListResponse{
			Data: []domain.Scorer{{
				FullName:          "Lionel Messi",
				TeamCode:          "ARG",
				TeamName:          "Argentina",
				Goals:             13,
				ListTeams:         []string{"ARG"},
				ConfederationCode: "CONMEBOL",
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
		filter := domain.ScorerFilter{Name: "messi", Language: "en", TeamCode: "arg", ConfederationCode: "conmebol", Page: 1, Size: 10}
		svc.On("List", mock.Anything, filter).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/scorers?page=1&size=10&name=messi&teamCode=arg&confederationCode=conmebol", nil)
		req.Header.Set("Accept-Language", "en")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "1", w.Header().Get("API-Version-Used"))
		assert.JSONEq(t, `{
			"data": [{
				"fullName": "Lionel Messi",
				"teamCode": "ARG",
				"teamName": "Argentina",
				"goals": 13,
				"listTeams": ["ARG"],
				"confederationCode": "CONMEBOL"
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
		svc := new(MockScorerService)
		r := setupScorerRouter(svc)
		expected := &domain.ScorerListResponse{
			Data: []domain.Scorer{},
			Pagination: domain.PaginationInfo{
				Page:          1,
				Size:          20,
				TotalElements: 0,
				TotalPages:    0,
				HasNext:       false,
				HasPrevious:   false,
			},
		}
		svc.On("List", mock.Anything, domain.ScorerFilter{Language: "es", Page: 1, Size: 20}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/scorers", nil)
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

	t.Run("bad request invalid page", func(t *testing.T) {
		svc := new(MockScorerService)
		r := setupScorerRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/scorers?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid page parameter"}`, w.Body.String())
	})

	t.Run("bad request non numeric page", func(t *testing.T) {
		svc := new(MockScorerService)
		r := setupScorerRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/scorers?page=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid page parameter"}`, w.Body.String())
	})

	t.Run("bad request invalid size", func(t *testing.T) {
		svc := new(MockScorerService)
		r := setupScorerRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/scorers?size=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid size parameter"}`, w.Body.String())
	})

	t.Run("bad request size greater than max", func(t *testing.T) {
		svc := new(MockScorerService)
		r := setupScorerRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/scorers?size=101", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid size parameter"}`, w.Body.String())
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockScorerService)
		r := setupScorerRouter(svc)
		svc.On("List", mock.Anything, domain.ScorerFilter{Language: "es", Page: 1, Size: 20}).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/scorers", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to retrieve scorers"}`, w.Body.String())
		svc.AssertExpectations(t)
	})
}
