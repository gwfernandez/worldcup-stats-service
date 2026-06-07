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
	t.Run("success with filters", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipFilter{
			Year:              1978,
			Host:              "arg",
			ConfederationCode: "CONMEBOL",
			Page:              1,
			Size:              10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(1))
		mock.ExpectQuery(`^-- name: CountChampionships :one.*`).
			WithArgs(int32(1978), "arg", "CONMEBOL").
			WillReturnRows(countRows)

		startDate := time.Date(1978, 6, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(1978, 6, 25, 0, 0, 0, 0, time.UTC)
		rows := mock.NewRows([]string{"year", "start_date", "end_date", "host_codes", "champion_code"}).
			AddRow(int32(1978), startDate, endDate, []string{"arg"}, pgtype.Text{String: "ARG", Valid: true})

		mock.ExpectQuery(`^-- name: ListChampionships :many.*`).
			WithArgs(int32(1978), "arg", "CONMEBOL", int32(10), int32(0)).
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		require.Len(t, result, 1)
		assert.Equal(t, 1978, result[0].Year)
		assert.Equal(t, "1978-06-01", result[0].StartDate)
		assert.Equal(t, "1978-06-25", result[0].EndDate)
		assert.Equal(t, []string{"ARG"}, result[0].HostCodes)
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

		mock.ExpectQuery(`^-- name: CountChampionships :one.*`).
			WithArgs(int32(0), "", "").
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
		mock.ExpectQuery(`^-- name: CountChampionships :one.*`).
			WithArgs(int32(0), "", "").
			WillReturnRows(countRows)

		mock.ExpectQuery(`^-- name: ListChampionships :many.*`).
			WithArgs(int32(0), "", "", int32(10), int32(0)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionshipRepository_ListTeamsByYear(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		filter := domain.ChampionshipTeamFilter{
			Year:              1930,
			Name:              "argentina",
			ConfederationCode: "CONMEBOL",
			GroupCode:         "1",
			Page:              1,
			Size:              10,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(2))
		mock.ExpectQuery(`^-- name: CountChampionshipTeamsByYear :one.*`).
			WithArgs(int32(1930), "argentina", "CONMEBOL", "1").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"year", "team_code", "confederation_code", "group_code", "stage_reached", "managers"}).
			AddRow(int32(1930), "arg", "conmebol", pgtype.Text{String: "1", Valid: true}, "runner_up", "Francisco Olazar").
			AddRow(int32(1930), "uru", "conmebol", pgtype.Text{Valid: false}, "champion", "")

		mock.ExpectQuery(`^-- name: ListChampionshipTeamsByYear :many.*`).
			WithArgs(int32(1930), "argentina", "CONMEBOL", "1", int32(10), int32(0)).
			WillReturnRows(rows)

		result, total, err := repo.ListTeamsByYear(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, result, 2)
		assert.Equal(t, domain.ChampionshipTeam{
			Year:              1930,
			TeamCode:          "ARG",
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

		mock.ExpectQuery(`^-- name: CountChampionshipTeamsByYear :one.*`).
			WithArgs(int32(1930), "", "", "").
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
		mock.ExpectQuery(`^-- name: CountChampionshipTeamsByYear :one.*`).
			WithArgs(int32(1930), "", "", "").
			WillReturnRows(countRows)

		mock.ExpectQuery(`^-- name: ListChampionshipTeamsByYear :many.*`).
			WithArgs(int32(1930), "", "", "", int32(10), int32(0)).
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

func TestChampionshipRepository_GetByYear(t *testing.T) {
	t.Run("success with stats", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionshipRepository(mock)
		startDate := time.Date(1930, 7, 13, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(1930, 7, 30, 0, 0, 0, 0, time.UTC)

		rows := mock.NewRows([]string{
			"year", "start_date", "end_date", "host_codes", "champion_code",
			"total_teams", "total_matches", "total_stadiums", "total_players", "total_goals",
			"stats_champion_code", "stats_runner_up_code", "stats_third_place_code", "stats_fourth_place_code",
			"top_scorer_ids", "top_scorer_goals",
		}).AddRow(
			int32(1930), startDate, endDate, []string{"uru"}, pgtype.Text{String: "uru", Valid: true},
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
			"year", "start_date", "end_date", "host_codes", "champion_code",
			"total_teams", "total_matches", "total_stadiums", "total_players", "total_goals",
			"stats_champion_code", "stats_runner_up_code", "stats_third_place_code", "stats_fourth_place_code",
			"top_scorer_ids", "top_scorer_goals",
		}).AddRow(
			int32(2026), startDate, endDate, []string{"usa", "can", "mex"}, pgtype.Text{Valid: false},
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
