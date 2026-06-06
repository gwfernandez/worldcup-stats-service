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

// MockTeamService mocks the TeamService interface.
type MockTeamService struct {
	mock.Mock
}

func (m *MockTeamService) List(ctx context.Context, filter domain.TeamFilter) (*domain.TeamListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TeamListResponse), args.Error(1)
}

func (m *MockTeamService) GetByCode(ctx context.Context, code string) (*domain.Team, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func setupTeamRouter(svc *MockTeamService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewTeamHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestTeamHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		dissolutionDate := "1991-12-26"
		expected := &domain.TeamListResponse{
			Data: []domain.Team{{
				Code:              "URS",
				Name:              "Soviet Union",
				IsDissolved:       true,
				ConfederationCode: "UEFA",
				FederationName:    "Football Federation of the USSR",
				FederationCode:    "FFUSSR",
				DissolutionDate:   &dissolutionDate,
			}},
			Pagination: domain.PaginationInfo{Page: 1, Size: 20, TotalElements: 1, TotalPages: 1, HasNext: false, HasPrevious: false},
		}
		svc.On("List", mock.Anything, domain.TeamFilter{
			Name:           "argen",
			FederationName: "",
			FederationCode: "",
			Page:           1,
			Size:           20,
		}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/teams?name=argen", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [{
				"code": "URS",
				"name": "Soviet Union",
				"isDissolved": true,
				"confederationCode": "UEFA",
				"federationName": "Football Federation of the USSR",
				"federationCode": "FFUSSR",
				"dissolutionDate": "1991-12-26"
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
	})

	t.Run("bad request invalid page", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/teams?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request invalid size", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/teams?size=500", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("success with camelCase filters", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		confederationCode := "CONMEBOL"
		expected := &domain.TeamListResponse{Data: []domain.Team{}, Pagination: domain.PaginationInfo{Page: 1, Size: 20, TotalElements: 1, TotalPages: 1, HasNext: false, HasPrevious: false}}
		svc.On("List", mock.Anything, domain.TeamFilter{
			Name:              "",
			ConfederationCode: &confederationCode,
			FederationName:    "Asociación",
			FederationCode:    "AFA",
			IncludeDissolved:  true,
			Page:              1,
			Size:              20,
		}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/teams?confederationCode=CONMEBOL&federationName=Asociaci%C3%B3n&federationCode=AFA&includeDissolved=true", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": [],
			"pagination": {
				"page": 1,
				"size": 20,
				"totalElements": 1,
				"totalPages": 1,
				"hasNext": false,
				"hasPrevious": false
			}
		}`, w.Body.String())
	})

	t.Run("bad request invalid includeDissolved", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/teams?includeDissolved=nope", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid includeDissolved parameter"}`, w.Body.String())
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		svc.On("List", mock.Anything, domain.TeamFilter{
			Name:           "",
			FederationName: "",
			FederationCode: "",
			Page:           1,
			Size:           20,
		}).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/teams", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTeamHandler_GetByCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		expected := &domain.Team{
			Code:              "URS",
			Name:              "Soviet Union",
			IsDissolved:       true,
			ConfederationCode: "UEFA",
			FederationName:    "Football Federation of the USSR",
			FederationCode:    "FFUSSR",
		}
		svc.On("GetByCode", mock.Anything, "urs").Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/teams/urs", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"code": "URS",
			"name": "Soviet Union",
			"isDissolved": true,
			"confederationCode": "UEFA",
			"federationName": "Football Federation of the USSR",
			"federationCode": "FFUSSR",
			"dissolutionDate": null
		}`, w.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		svc.On("GetByCode", mock.Anything, "zzz").Return(nil, domain.ErrNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/teams/zzz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockTeamService)
		r := setupTeamRouter(svc)

		svc.On("GetByCode", mock.Anything, "zzz").Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/teams/zzz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
