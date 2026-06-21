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

func TestGoalRepository_ListByPlayer(t *testing.T) {
	t.Run("success with year filter and nullable fields", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewGoalRepository(mock)
		filter := domain.GoalFilter{PlayerID: 1524, Year: 2018, Page: 2, Size: 10}
		playerID := pgtype.Int8{Int64: 1524, Valid: true}

		mock.ExpectQuery(`^-- name: CountGoalsByPlayer :one.*`).
			WithArgs(playerID, int32(2018)).
			WillReturnRows(mock.NewRows([]string{"count"}).AddRow(int64(11)))

		matchDate := pgtype.Date{Time: time.Date(2018, time.June, 16, 0, 0, 0, 0, time.UTC), Valid: true}
		rows := mock.NewRows([]string{"year", "host_codes", "match_date", "opponent_team_code", "minute_regular", "penalty", "stage"}).
			AddRow(int32(2018), []string{"rus"}, matchDate, "isl", int32(64), pgtype.Bool{Bool: true, Valid: true}, "group_stage").
			AddRow(int32(2018), []string{"kor", "jpn"}, pgtype.Date{}, "cro", int32(80), pgtype.Bool{}, nil)
		mock.ExpectQuery(`^-- name: ListGoalsByPlayer :many.*`).
			WithArgs(playerID, int32(2018), int32(10), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.ListByPlayer(context.Background(), filter)
		require.NoError(t, err)
		assert.Equal(t, int64(11), total)
		require.Len(t, result, 2)
		assert.Equal(t, int32(2018), result[0].Year)
		assert.Equal(t, []domain.SimpleTeam{{Code: "RUS"}}, result[0].Hosts)
		assert.Equal(t, "2018-06-16", *result[0].MatchDate)
		assert.Equal(t, "ISL", result[0].OpponentTeam.Code)
		assert.Empty(t, result[0].OpponentTeam.Name)
		assert.Equal(t, int32(64), result[0].MinuteRegular)
		assert.True(t, *result[0].Penalty)
		assert.Equal(t, "group_stage", *result[0].Stage)
		assert.Equal(t, []domain.SimpleTeam{{Code: "KOR"}, {Code: "JPN"}}, result[1].Hosts)
		assert.Nil(t, result[1].MatchDate)
		assert.Nil(t, result[1].Penalty)
		assert.Nil(t, result[1].Stage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without year filter", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewGoalRepository(mock)
		playerID := pgtype.Int8{Int64: 10, Valid: true}
		mock.ExpectQuery(`^-- name: CountGoalsByPlayer :one.*`).
			WithArgs(playerID, int32(0)).
			WillReturnRows(mock.NewRows([]string{"count"}).AddRow(int64(0)))
		mock.ExpectQuery(`^-- name: ListGoalsByPlayer :many.*`).
			WithArgs(playerID, int32(0), int32(0), int32(20)).
			WillReturnRows(mock.NewRows([]string{"year", "host_codes", "match_date", "opponent_team_code", "minute_regular", "penalty", "stage"}))

		result, total, err := repo.ListByPlayer(context.Background(), domain.GoalFilter{PlayerID: 10, Page: 1, Size: 20})
		require.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("count error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewGoalRepository(mock)
		playerID := pgtype.Int8{Int64: 10, Valid: true}
		mock.ExpectQuery(`^-- name: CountGoalsByPlayer :one.*`).
			WithArgs(playerID, int32(0)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListByPlayer(context.Background(), domain.GoalFilter{PlayerID: 10, Page: 1, Size: 20})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("list error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewGoalRepository(mock)
		playerID := pgtype.Int8{Int64: 10, Valid: true}
		mock.ExpectQuery(`^-- name: CountGoalsByPlayer :one.*`).
			WithArgs(playerID, int32(0)).
			WillReturnRows(mock.NewRows([]string{"count"}).AddRow(int64(1)))
		mock.ExpectQuery(`^-- name: ListGoalsByPlayer :many.*`).
			WithArgs(playerID, int32(0), int32(0), int32(20)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.ListByPlayer(context.Background(), domain.GoalFilter{PlayerID: 10, Page: 1, Size: 20})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
