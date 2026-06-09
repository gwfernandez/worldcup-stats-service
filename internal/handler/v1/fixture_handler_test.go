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

type mockFixtureService struct {
	mock.Mock
}

func (m *mockFixtureService) GetByYear(ctx context.Context, year int, language string) (*domain.Fixture, error) {
	args := m.Called(ctx, year, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Fixture), args.Error(1)
}

func setupFixtureRouter(svc *mockFixtureService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := v1.NewFixtureHandler(svc)
	rg := r.Group("/api")
	h.RegisterRoutes(rg)
	return r
}

func TestFixtureHandler_GetByYear(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mockFixtureService)
		r := setupFixtureRouter(svc)
		expected := &domain.Fixture{Year: 1978, Stages: []domain.FixtureStage{{
			Stage: "group_stage",
			Groups: []domain.FixtureGroup{{
				GroupCode: "A",
				Matches: []domain.FixtureMatch{{
					ID:           1,
					HomeTeamCode: "ARG",
					HomeTeamName: "Argentina",
					AwayTeamCode: "FRA",
					AwayTeamName: "France",
				}},
				Standings: []domain.GroupStanding{{
					TeamCode: "ARG",
					TeamName: "Argentina",
				}},
			}},
		}}}

		svc.On("GetByYear", mock.Anything, 1978, "en").Return(expected, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1978/fixture", nil)
		req.Header.Set("Accept-Language", "en")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"data": {
				"year": 1978,
				"stages": [{
					"stage": "group_stage",
					"groups": [{
						"groupCode": "A",
						"matches": [{
							"id": 1,
							"stageType": "",
							"replayed": false,
							"replayOf": null,
							"matchDate": null,
							"matchTime": null,
							"stadiumId": null,
							"homeTeamCode": "ARG",
							"homeTeamName": "Argentina",
							"awayTeamCode": "FRA",
							"awayTeamName": "France",
							"homeTeamScore": null,
							"awayTeamScore": null,
							"extraTime": false,
							"penaltyShootout": false,
							"homeTeamScorePenalties": null,
							"awayTeamScorePenalties": null,
							"homeTeamWin": null,
							"awayTeamWin": null,
							"draw": null,
							"refId": null
						}],
						"standings": [{
							"teamCode": "ARG",
							"teamName": "Argentina",
							"matchesPlayed": 0,
							"wins": 0,
							"draws": 0,
							"losses": 0,
							"goalsFor": 0,
							"goalsAgainst": 0,
							"goalDifference": 0,
							"points": 0,
							"unifiedPoints": 0,
							"position": null
						}]
					}]
				}]
			}
		}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("invalid year", func(t *testing.T) {
		svc := new(mockFixtureService)
		r := setupFixtureRouter(svc)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/abc/fixture", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid year"}`, w.Body.String())
		svc.AssertNotCalled(t, "GetByYear")
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(mockFixtureService)
		r := setupFixtureRouter(svc)

		svc.On("GetByYear", mock.Anything, 2023, "es").Return(nil, domain.ErrNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/2023/fixture", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"championship not found"}`, w.Body.String())
		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mockFixtureService)
		r := setupFixtureRouter(svc)

		svc.On("GetByYear", mock.Anything, 1978, "es").Return(nil, errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/championships/1978/fixture", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to retrieve fixture"}`, w.Body.String())
		svc.AssertExpectations(t)
	})
}
