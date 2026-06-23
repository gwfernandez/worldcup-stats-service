package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

func TestChampionshipRepository_List(t *testing.T) {
	t.Run("success without host filter uses optimized query", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipFilter{
			Year:              1978,
			ConfederationCode: "CONMEBOL",
			Page:              1,
			Size:              10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(1))
		mock.ExpectQuery(`^-- name: CountChampionshipsWithoutHostFilter :one.*`).
			WithArgs(int32(1978), "CONMEBOL").
			WillReturnRows(countRows)

		startDate := time.Date(1978, 6, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(1978, 6, 25, 0, 0, 0, 0, time.UTC)
		rows := mock.NewRows([]string{"year", "start_date", "end_date", "host_codes", "confederation_codes", "champion_code"}).
			AddRow(int32(1978), startDate, endDate, []string{"arg"}, []string{"conmebol"}, pgtype.Text{String: "ARG", Valid: true})

		mock.ExpectQuery(`^-- name: ListChampionshipsWithoutHostFilter :many.*`).
			WithArgs(int32(1978), "CONMEBOL", int32(10), int32(0)).
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		require.Len(t, result, 1)
		assert.Equal(t, []string{"ARG"}, result[0].HostCodes)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success with filters", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipFilter{
			Year:              1978,
			Host:              "arg",
			Language:          "en",
			ConfederationCode: "CONMEBOL",
			Page:              1,
			Size:              10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(1))
		mock.ExpectQuery(`^-- name: CountChampionships :one.*`).
			WithArgs(int32(1978), "arg", "CONMEBOL", "en").
			WillReturnRows(countRows)

		startDate := time.Date(1978, 6, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(1978, 6, 25, 0, 0, 0, 0, time.UTC)
		rows := mock.NewRows([]string{"year", "start_date", "end_date", "host_codes", "confederation_codes", "champion_code"}).
			AddRow(int32(1978), startDate, endDate, []string{"arg"}, []string{"conmebol"}, pgtype.Text{String: "ARG", Valid: true})

		mock.ExpectQuery(`^-- name: ListChampionships :many.*`).
			WithArgs(int32(1978), "arg", "CONMEBOL", int32(10), int32(0), "en").
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		require.Len(t, result, 1)
		assert.Equal(t, 1978, result[0].Year)
		assert.Equal(t, "1978-06-01", result[0].StartDate)
		assert.Equal(t, "1978-06-25", result[0].EndDate)
		assert.Equal(t, []string{"ARG"}, result[0].HostCodes)
		assert.Equal(t, []string{"CONMEBOL"}, result[0].ConfederationCodes)
		require.NotNil(t, result[0].ChampionCode)
		assert.Equal(t, "ARG", *result[0].ChampionCode)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipFilter{Page: 1, Size: 10}

		mock.ExpectQuery(`^-- name: CountChampionshipsWithoutHostFilter :one.*`).
			WithArgs(int32(0), "").
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipFilter{Page: 1, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(5))
		mock.ExpectQuery(`^-- name: CountChampionshipsWithoutHostFilter :one.*`).
			WithArgs(int32(0), "").
			WillReturnRows(countRows)

		mock.ExpectQuery(`^-- name: ListChampionshipsWithoutHostFilter :many.*`).
			WithArgs(int32(0), "", int32(10), int32(0)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionshipRepository_ListTeamTranslations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)

		rows := mock.NewRows([]string{"team_code", "language", "name"}).
			AddRow("arg", "es", "Argentina").
			AddRow("KOR", "en", "South Korea")
		mock.ExpectQuery(`^-- name: ListTeamTranslations :many.*`).
			WillReturnRows(rows)

		result, err := repo.ListTeamTranslations(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []domain.TeamTranslation{
			{TeamCode: "ARG", Language: "es", Name: "Argentina"},
			{TeamCode: "KOR", Language: "en", Name: "South Korea"},
		}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)

		mock.ExpectQuery(`^-- name: ListTeamTranslations :many.*`).
			WillReturnError(errors.New("db error"))

		result, err := repo.ListTeamTranslations(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionshipRepository_ListTeamsByYear(t *testing.T) {
	t.Run("success without name filter uses optimized query", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipTeamFilter{
			Year:              1930,
			ConfederationCode: "CONMEBOL",
			GroupCode:         "1",
			Page:              1,
			Size:              10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(1))
		mock.ExpectQuery(`^-- name: CountChampionshipTeamsByYearWithoutNameFilter :one.*`).
			WithArgs(int32(1930), "CONMEBOL", "1").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"year", "team_code", "confederation_code", "group_code", "stage_reached", "managers"}).
			AddRow(int32(1930), "arg", "conmebol", pgtype.Text{String: "1", Valid: true}, "runner_up", "Francisco Olazar")

		mock.ExpectQuery(`^-- name: ListChampionshipTeamsByYearWithoutNameFilter :many.*`).
			WithArgs(int32(1930), "CONMEBOL", "1", int32(10), int32(0)).
			WillReturnRows(rows)

		result, total, err := repo.ListTeamsByYear(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		require.Len(t, result, 1)
		assert.Equal(t, "ARG", result[0].Team.Code)
		assert.Empty(t, result[0].Team.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success with filters", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipTeamFilter{
			Year:              1930,
			Name:              "argentina",
			Language:          "en",
			ConfederationCode: "CONMEBOL",
			GroupCode:         "1",
			Page:              1,
			Size:              10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(2))
		mock.ExpectQuery(`^-- name: CountChampionshipTeamsByYear :one.*`).
			WithArgs(int32(1930), "argentina", "CONMEBOL", "1", "en").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"year", "team_code", "confederation_code", "group_code", "stage_reached", "managers"}).
			AddRow(int32(1930), "arg", "conmebol", pgtype.Text{String: "1", Valid: true}, "runner_up", "Francisco Olazar").
			AddRow(int32(1930), "uru", "conmebol", pgtype.Text{Valid: false}, "champion", "")

		mock.ExpectQuery(`^-- name: ListChampionshipTeamsByYear :many.*`).
			WithArgs(int32(1930), "argentina", "CONMEBOL", "1", int32(10), int32(0), "en").
			WillReturnRows(rows)

		result, total, err := repo.ListTeamsByYear(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, result, 2)
		assert.Equal(t, domain.ChampionshipTeam{
			Year:              1930,
			Team:              domain.SimpleTeam{Code: "ARG"},
			ConfederationCode: "CONMEBOL",
			GroupCode:         "1",
			StageReached:      "runner_up",
			Managers:          "Francisco Olazar",
		}, result[0])
		assert.Equal(t, "", result[1].GroupCode)
		assert.Equal(t, "", result[1].Managers)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipTeamFilter{Year: 1930, Page: 1, Size: 10}

		mock.ExpectQuery(`^-- name: CountChampionshipTeamsByYearWithoutNameFilter :one.*`).
			WithArgs(int32(1930), "", "").
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListTeamsByYear(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipTeamFilter{Year: 1930, Page: 1, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(5))
		mock.ExpectQuery(`^-- name: CountChampionshipTeamsByYearWithoutNameFilter :one.*`).
			WithArgs(int32(1930), "", "").
			WillReturnRows(countRows)

		mock.ExpectQuery(`^-- name: ListChampionshipTeamsByYearWithoutNameFilter :many.*`).
			WithArgs(int32(1930), "", "", int32(10), int32(0)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListTeamsByYear(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionshipRepository_ListStadiumsByYear(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipStadiumFilter{
			Year: 1930,
			Name: "centenario",
			Page: 1,
			Size: 10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(2))
		mock.ExpectQuery(`^-- name: CountChampionshipStadiumsByYear :one.*`).
			WithArgs(int32(1930), "centenario").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"year", "id", "name", "city_name", "capacity", "matches_played"}).
			AddRow(int32(1930), int64(1), "Estadio Centenario", "Montevideo", int32(90000), int32(10)).
			AddRow(int32(1930), int64(2), "Estadio Pocitos", "", int32(0), int32(2))

		mock.ExpectQuery(`^-- name: ListChampionshipStadiumsByYear :many.*`).
			WithArgs(int32(1930), "centenario", int32(10), int32(0)).
			WillReturnRows(rows)

		result, total, err := repo.ListStadiumsByYear(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, result, 2)
		assert.Equal(t, domain.ChampionshipStadium{
			Year:          1930,
			ID:            1,
			Name:          "Estadio Centenario",
			CityName:      "Montevideo",
			Capacity:      90000,
			MatchesPlayed: 10,
		}, result[0])
		assert.Equal(t, "", result[1].CityName)
		assert.Equal(t, int32(0), result[1].Capacity)
		assert.Equal(t, int32(2), result[1].MatchesPlayed)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipStadiumFilter{Year: 1930, Page: 1, Size: 10}

		mock.ExpectQuery(`^-- name: CountChampionshipStadiumsByYear :one.*`).
			WithArgs(int32(1930), "").
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListStadiumsByYear(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipStadiumFilter{Year: 1930, Page: 2, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(5))
		mock.ExpectQuery(`^-- name: CountChampionshipStadiumsByYear :one.*`).
			WithArgs(int32(1930), "").
			WillReturnRows(countRows)

		mock.ExpectQuery(`^-- name: ListChampionshipStadiumsByYear :many.*`).
			WithArgs(int32(1930), "", int32(10), int32(10)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListStadiumsByYear(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionshipRepository_ListStandingsByYear(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipStandingFilter{
			Year:     1930,
			Language: "en",
			Page:     2,
			Size:     10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(13))
		mock.ExpectQuery(`^-- name: CountChampionshipStandingsByYear :one.*`).
			WithArgs(int32(1930)).
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{
			"team_code", "group_code", "matches_played", "wins", "draws", "losses",
			"goals_for", "goals_against", "goal_difference", "points", "unified_points",
			"position", "performance",
		}).
			AddRow("uru", "3", int32(4), int32(4), int32(0), int32(0), int32(15), int32(3), pgtype.Int4{Int32: 12, Valid: true}, int32(8), int32(12), pgtype.Int4{Int32: 1, Valid: true}, "champion").
			AddRow("arg", "1", int32(5), int32(4), int32(0), int32(1), int32(18), int32(9), pgtype.Int4{Int32: 9, Valid: true}, int32(8), int32(12), pgtype.Int4{Int32: 2, Valid: true}, "runner_up")

		mock.ExpectQuery(`^-- name: ListChampionshipStandingsByYear :many.*`).
			WithArgs(int32(1930), int32(10), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.ListStandingsByYear(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(13), total)
		require.Len(t, result, 2)
		assert.Equal(t, domain.ChampionshipStanding{
			Team:           domain.SimpleTeam{Code: "URU"},
			GroupCode:      "3",
			MatchesPlayed:  4,
			Wins:           4,
			Draws:          0,
			Losses:         0,
			GoalsFor:       15,
			GoalsAgainst:   3,
			GoalDifference: 12,
			Points:         8,
			UnifiedPoints:  12,
			Position:       1,
			Performance:    "champion",
		}, result[0])
		assert.Equal(t, "ARG", result[1].Team.Code)
		assert.Equal(t, "runner_up", result[1].Performance)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipStandingFilter{Year: 1930, Page: 1, Size: 10}

		mock.ExpectQuery(`^-- name: CountChampionshipStandingsByYear :one.*`).
			WithArgs(int32(1930)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListStandingsByYear(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipStandingFilter{Year: 1930, Page: 1, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(13))
		mock.ExpectQuery(`^-- name: CountChampionshipStandingsByYear :one.*`).
			WithArgs(int32(1930)).
			WillReturnRows(countRows)

		mock.ExpectQuery(`^-- name: ListChampionshipStandingsByYear :many.*`).
			WithArgs(int32(1930), int32(10), int32(0)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListStandingsByYear(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionshipRepository_ListScorersByYear(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipScorerFilter{
			Year:     1930,
			Name:     "guille",
			Language: "en",
			TeamCode: "ARG",
			Page:     2,
			Size:     10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(2))
		mock.ExpectQuery(`^-- name: CountChampionshipScorersByYear :one.*`).
			WithArgs(int32(1930), "guille", "ARG").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"player_id", "full_name", "team_code", "goals"}).
			AddRow(int64(10), "Guillermo Stabile", "arg", int32(8)).
			AddRow(int64(11), "Carlos Peucelle", "arg", int32(3))

		mock.ExpectQuery(`^-- name: ListChampionshipScorersByYear :many.*`).
			WithArgs(int32(1930), "guille", "ARG", int32(10), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.ListScorersByYear(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, result, 2)
		assert.Equal(t, domain.ChampionshipScorer{
			PlayerID: 10,
			FullName: "Guillermo Stabile",
			Team:     domain.SimpleTeam{Code: "ARG"},
			Goals:    8,
		}, result[0])
		assert.Equal(t, int64(11), result[1].PlayerID)
		assert.Equal(t, "ARG", result[1].Team.Code)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipScorerFilter{Year: 1930, Page: 1, Size: 10}

		mock.ExpectQuery(`^-- name: CountChampionshipScorersByYear :one.*`).
			WithArgs(int32(1930), "", "").
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListScorersByYear(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipScorerFilter{Year: 1930, Page: 1, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(5))
		mock.ExpectQuery(`^-- name: CountChampionshipScorersByYear :one.*`).
			WithArgs(int32(1930), "", "").
			WillReturnRows(countRows)

		mock.ExpectQuery(`^-- name: ListChampionshipScorersByYear :many.*`).
			WithArgs(int32(1930), "", "", int32(0), int32(10)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListScorersByYear(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionshipRepository_ListSquadByYearAndTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipSquadFilter{
			Year:     1930,
			TeamCode: "ARG",
			Page:     2,
			Size:     10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(13))
		mock.ExpectQuery(`^-- name: CountChampionshipSquadByYearAndTeam :one.*`).
			WithArgs(int32(1930), "ARG").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"player_id", "first_name", "last_name", "position", "shirt_number"}).
			AddRow(int64(10), "Guillermo", "Stabile", pgtype.Text{String: "forward", Valid: true}, pgtype.Int4{Int32: 10, Valid: true}).
			AddRow(int64(11), "Carlos", "Peucelle", pgtype.Text{Valid: false}, pgtype.Int4{Valid: false})

		mock.ExpectQuery(`^-- name: ListChampionshipSquadByYearAndTeam :many.*`).
			WithArgs(int32(1930), "ARG", int32(10), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.ListSquadByYearAndTeam(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(13), total)
		require.Len(t, result, 2)
		assert.Equal(t, int64(10), result[0].PlayerID)
		assert.Equal(t, "Guillermo", result[0].FirstName)
		assert.Equal(t, "Stabile", result[0].LastName)
		require.NotNil(t, result[0].Position)
		assert.Equal(t, "forward", *result[0].Position)
		require.NotNil(t, result[0].ShirtNumber)
		assert.Equal(t, int32(10), *result[0].ShirtNumber)
		assert.Nil(t, result[1].Position)
		assert.Nil(t, result[1].ShirtNumber)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipSquadFilter{Year: 1930, TeamCode: "ARG", Page: 1, Size: 10}

		mock.ExpectQuery(`^-- name: CountChampionshipSquadByYearAndTeam :one.*`).
			WithArgs(int32(1930), "ARG").
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListSquadByYearAndTeam(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipSquadFilter{Year: 1930, TeamCode: "ARG", Page: 1, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(13))
		mock.ExpectQuery(`^-- name: CountChampionshipSquadByYearAndTeam :one.*`).
			WithArgs(int32(1930), "ARG").
			WillReturnRows(countRows)

		mock.ExpectQuery(`^-- name: ListChampionshipSquadByYearAndTeam :many.*`).
			WithArgs(int32(1930), "ARG", int32(0), int32(10)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListSquadByYearAndTeam(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionshipRepository_GetByYear(t *testing.T) {
	t.Run("success with stats", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		startDate := time.Date(1930, 7, 13, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(1930, 7, 30, 0, 0, 0, 0, time.UTC)

		rows := mock.NewRows([]string{
			"year", "start_date", "end_date", "host_codes", "confederation_codes", "champion_code",
			"total_teams", "total_matches", "total_stadiums", "total_players", "total_goals",
			"stats_champion_code", "stats_runner_up_code", "stats_third_place_code", "stats_fourth_place_code",
			"top_scorer_ids", "top_scorer_goals",
		}).AddRow(
			int32(1930), startDate, endDate, []string{"uru"}, []string{"conmebol"}, pgtype.Text{String: "uru", Valid: true},
			pgtype.Int4{Int32: 13, Valid: true}, pgtype.Int4{Int32: 18, Valid: true}, pgtype.Int4{Int32: 3, Valid: true},
			pgtype.Int4{Int32: 189, Valid: true}, pgtype.Int4{Int32: 70, Valid: true},
			pgtype.Text{String: "uru", Valid: true}, pgtype.Text{String: "arg", Valid: true},
			pgtype.Text{String: "usa", Valid: true}, pgtype.Text{String: "yug", Valid: true},
			[]int64{}, pgtype.Int4{Int32: 8, Valid: true},
		)

		mock.ExpectQuery(`^-- name: GetChampionshipByYear :one.*`).
			WithArgs(int32(1930)).
			WillReturnRows(rows)

		result, err := repo.GetByYear(context.Background(), 1930)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1930, result.Year)
		assert.Equal(t, "1930-07-13", result.StartDate)
		assert.Equal(t, "1930-07-30", result.EndDate)
		assert.Equal(t, []string{"URU"}, result.HostCodes)
		assert.Equal(t, []string{"CONMEBOL"}, result.ConfederationCodes)
		require.NotNil(t, result.ChampionCode)
		assert.Equal(t, "URU", *result.ChampionCode)
		require.NotNil(t, result.Stats)
		assert.Equal(t, int32(13), result.Stats.TotalTeams)
		assert.Equal(t, int32(18), result.Stats.TotalMatches)
		assert.Equal(t, int32(3), result.Stats.TotalStadiums)
		assert.Equal(t, int32(189), result.Stats.TotalPlayers)
		assert.Equal(t, int32(70), result.Stats.TotalGoals)
		assert.Equal(t, "ARG", result.Stats.RunnerUpCode)
		assert.Equal(t, "USA", result.Stats.ThirdPlaceCode)
		assert.Equal(t, "YUG", result.Stats.FourthPlaceCode)
		assert.Equal(t, int32(8), result.Stats.TopScorerGoals)
		assert.Len(t, result.Stats.TopScorers, 0)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without stats", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		startDate := time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC)

		rows := mock.NewRows([]string{
			"year", "start_date", "end_date", "host_codes", "confederation_codes", "champion_code",
			"total_teams", "total_matches", "total_stadiums", "total_players", "total_goals",
			"stats_champion_code", "stats_runner_up_code", "stats_third_place_code", "stats_fourth_place_code",
			"top_scorer_ids", "top_scorer_goals",
		}).AddRow(
			int32(2026), startDate, endDate, []string{"usa", "can", "mex"}, []string{"concacaf", "uefa"}, pgtype.Text{Valid: false},
			pgtype.Int4{Valid: false}, pgtype.Int4{Valid: false}, pgtype.Int4{Valid: false},
			pgtype.Int4{Valid: false}, pgtype.Int4{Valid: false},
			pgtype.Text{Valid: false}, pgtype.Text{Valid: false},
			pgtype.Text{Valid: false}, pgtype.Text{Valid: false},
			nil, pgtype.Int4{Valid: false},
		)

		mock.ExpectQuery(`^-- name: GetChampionshipByYear :one.*`).
			WithArgs(int32(2026)).
			WillReturnRows(rows)

		result, err := repo.GetByYear(context.Background(), 2026)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2026, result.Year)
		assert.Equal(t, "2026-06-11", result.StartDate)
		assert.Equal(t, "2026-07-19", result.EndDate)
		assert.Equal(t, []string{"USA", "CAN", "MEX"}, result.HostCodes)
		assert.Equal(t, []string{"CONCACAF", "UEFA"}, result.ConfederationCodes)
		assert.Nil(t, result.ChampionCode)
		assert.Nil(t, result.Stats) // Stats should be nil here, to be filled by service

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		mock.ExpectQuery(`^-- name: GetChampionshipByYear :one.*`).
			WithArgs(int32(1999)).
			WillReturnError(pgx.ErrNoRows)

		result, err := repo.GetByYear(context.Background(), 1999)
		assert.NoError(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		mock.ExpectQuery(`^-- name: GetChampionshipByYear :one.*`).
			WithArgs(int32(2022)).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetByYear(context.Background(), 2022)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
