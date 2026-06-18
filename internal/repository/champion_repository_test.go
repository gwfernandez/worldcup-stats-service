package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

func TestChampionRepository_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionRepository(mock)
		filter := domain.ChampionFilter{Language: "en", Page: 1, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(2))
		mock.ExpectQuery(`^-- name: CountChampions :one.*`).
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"team_code", "wins", "years", "confederation_code"}).
			AddRow("bra", int64(5), []int32{1958, 1962, 1970, 1994, 2002}, "conmebol").
			AddRow("arg", int64(3), []int32{1978, 1986, 2022}, "conmebol")
		mock.ExpectQuery(`^-- name: ListChampions :many.*`).
			WithArgs(int32(0), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, result, 2)
		assert.Equal(t, "BRA", result[0].Team.Code)
		assert.Empty(t, result[0].Team.Name)
		assert.Equal(t, int64(5), result[0].Wins)
		assert.Equal(t, []int32{1958, 1962, 1970, 1994, 2002}, result[0].Years)
		assert.Equal(t, "CONMEBOL", result[0].ConfederationCode)
		assert.Equal(t, "ARG", result[1].Team.Code)
		assert.Equal(t, "CONMEBOL", result[1].ConfederationCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionRepository(mock)
		mock.ExpectQuery(`^-- name: CountChampions :one.*`).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), domain.ChampionFilter{Page: 1, Size: 20})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionRepository(mock)
		countRows := mock.NewRows([]string{"count"}).AddRow(int64(1))
		mock.ExpectQuery(`^-- name: CountChampions :one.*`).
			WillReturnRows(countRows)
		mock.ExpectQuery(`^-- name: ListChampions :many.*`).
			WithArgs(int32(20), int32(20)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), domain.ChampionFilter{Page: 2, Size: 20})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChampionRepository_ListFinalsWonByTeam(t *testing.T) {
	t.Run("success with regular final and 1950 deciding match", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionRepository(mock)
		filter := domain.ChampionFinalFilter{TeamCode: "URU", Page: 1, Size: 10}

		mock.ExpectQuery(`^-- name: CountFinalsWonByTeam :one.*`).
			WithArgs("URU").
			WillReturnRows(mock.NewRows([]string{"count"}).AddRow(int64(2)))

		rows := mock.NewRows([]string{
			"year", "match_date", "match_time",
			"home_team_code", "home_team_score", "home_team_score_penalties",
			"away_team_code", "away_team_score", "away_team_score_penalties",
		}).
			AddRow(
				int32(1930),
				pgtype.Date{Time: time.Date(1930, 7, 30, 0, 0, 0, 0, time.UTC), Valid: true},
				pgtype.Time{Microseconds: 14 * 60 * 60 * 1_000_000, Valid: true},
				"uru",
				pgtype.Int4{Int32: 4, Valid: true},
				pgtype.Int4{},
				"arg",
				pgtype.Int4{Int32: 2, Valid: true},
				pgtype.Int4{},
			).
			AddRow(
				int32(1950),
				pgtype.Date{Time: time.Date(1950, 7, 16, 0, 0, 0, 0, time.UTC), Valid: true},
				pgtype.Time{},
				"bra",
				pgtype.Int4{Int32: 1, Valid: true},
				pgtype.Int4{},
				"uru",
				pgtype.Int4{Int32: 2, Valid: true},
				pgtype.Int4{},
			)
		mock.ExpectQuery(`^-- name: ListFinalsWonByTeam :many.*`).
			WithArgs("URU", int32(0), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.ListFinalsWonByTeam(context.Background(), filter)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, result, 2)
		assert.Equal(t, int32(1930), result[0].Year)
		assert.Equal(t, "1930-07-30", *result[0].MatchDate)
		assert.Equal(t, "14:00:00", *result[0].MatchTime)
		assert.Equal(t, "URU", result[0].HomeTeam.Code)
		assert.Equal(t, int32(4), *result[0].HomeTeamScore)
		assert.Nil(t, result[0].HomeTeamScorePenalties)
		assert.Equal(t, "ARG", result[0].AwayTeam.Code)
		assert.Equal(t, int32(2), *result[0].AwayTeamScore)
		assert.Nil(t, result[0].AwayTeamScorePenalties)
		assert.Equal(t, int32(1950), result[1].Year)
		assert.Nil(t, result[1].MatchTime)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionRepository(mock)
		mock.ExpectQuery(`^-- name: CountFinalsWonByTeam :one.*`).
			WithArgs("ARG").
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListFinalsWonByTeam(
			context.Background(),
			domain.ChampionFinalFilter{TeamCode: "ARG", Page: 1, Size: 20},
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionRepository(mock)
		mock.ExpectQuery(`^-- name: CountFinalsWonByTeam :one.*`).
			WithArgs("ARG").
			WillReturnRows(mock.NewRows([]string{"count"}).AddRow(int64(3)))
		mock.ExpectQuery(`^-- name: ListFinalsWonByTeam :many.*`).
			WithArgs("ARG", int32(20), int32(20)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListFinalsWonByTeam(
			context.Background(),
			domain.ChampionFinalFilter{TeamCode: "ARG", Page: 2, Size: 20},
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
