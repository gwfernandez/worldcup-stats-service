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

// MockNationalTeamService mocks the NationalTeamService interface.
type MockNationalTeamService struct {
	mock.Mock
}

func (m *MockNationalTeamService) List(ctx context.Context, filter domain.NationalTeamFilter) (*domain.NationalTeamListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NationalTeamListResponse), args.Error(1)
}

func (m *MockNationalTeamService) GetByCode(ctx context.Context, code string) (*domain.NationalTeam, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NationalTeam), args.Error(1)
}

func setupNationalTeamRouter(svc *MockNationalTeamService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewNationalTeamHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestNationalTeamHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		expected := &domain.NationalTeamListResponse{Pagination: domain.PaginationInfo{Page: 1, Size: 20, TotalElements: 1, TotalPages: 1, HasNext: false, HasPrevious: false}}
		svc.On("List", mock.Anything, domain.NationalTeamFilter{
			Name:           "argen",
			FederationName: "",
			FederationCode: "",
			Page:           1,
			Size:           20,
		}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams?name=argen", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request invalid page", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request invalid size", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams?size=500", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("success with confederation filter", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		confederationCode := "CONMEBOL"
		expected := &domain.NationalTeamListResponse{Pagination: domain.PaginationInfo{Page: 1, Size: 20, TotalElements: 1, TotalPages: 1, HasNext: false, HasPrevious: false}}
		svc.On("List", mock.Anything, domain.NationalTeamFilter{
			Name:              "",
			ConfederationCode: &confederationCode,
			FederationName:    "",
			FederationCode:    "",
			Page:              1,
			Size:              20,
		}).Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams?confederation_code=CONMEBOL", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request invalid include_dissolved", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams?include_dissolved=nope", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		svc.On("List", mock.Anything, domain.NationalTeamFilter{
			Name:           "",
			FederationName: "",
			FederationCode: "",
			Page:           1,
			Size:           20,
		}).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNationalTeamHandler_GetByCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		expected := &domain.NationalTeam{Code: "URS", Name: "Soviet Union"}
		svc.On("GetByCode", mock.Anything, "urs").Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams/urs", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		svc.On("GetByCode", mock.Anything, "zzz").Return(nil, domain.ErrNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams/zzz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(MockNationalTeamService)
		r := setupNationalTeamRouter(svc)

		svc.On("GetByCode", mock.Anything, "zzz").Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/national-teams/zzz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
