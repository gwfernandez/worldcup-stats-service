package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

func TestStandingRepository_List(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewStandingRepository(mock)
		filter := domain.StandingFilter{Name: "arg", ConfederationCode: "CONMEBOL", Page: 2, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(11))
		mock.ExpectQuery(`^-- name: CountStandings :one.*`).
			WithArgs("arg", "CONMEBOL").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{
			"team_code",
			"name",
			"matches_played",
			"wins",
			"draws",
			"losses",
			"goals_for",
			"goals_against",
			"goal_difference",
			"points",
			"unified_points",
			"position",
			"unified_position",
		}).AddRow("arg", "Argentina", int32(88), int32(53), int32(10), int32(25), int32(152), int32(101), int32(51), int32(133), int32(159), int32(3), int32(3))

		mock.ExpectQuery(`^-- name: ListStandings :many.*`).
			WithArgs("arg", "CONMEBOL", int32(10), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(11), total)
		require.Len(t, result, 1)
		assert.Equal(t, "ARG", result[0].TeamCode)
		assert.Equal(t, "Argentina", result[0].Name)
		assert.Equal(t, int32(88), result[0].MatchesPlayed)
		assert.Equal(t, int32(53), result[0].Wins)
		assert.Equal(t, int32(10), result[0].Draws)
		assert.Equal(t, int32(25), result[0].Losses)
		assert.Equal(t, int32(152), result[0].GoalsFor)
		assert.Equal(t, int32(101), result[0].GoalsAgainst)
		assert.Equal(t, int32(51), result[0].GoalDifference)
		assert.Equal(t, int32(133), result[0].Points)
		assert.Equal(t, int32(159), result[0].UnifiedPoints)
		assert.Equal(t, int32(3), result[0].Position)
		assert.Equal(t, int32(3), result[0].UnifiedPosition)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewStandingRepository(mock)
		mock.ExpectQuery(`^-- name: CountStandings :one.*`).
			WithArgs("", "").
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), domain.StandingFilter{Page: 1, Size: 20})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewStandingRepository(mock)
		countRows := mock.NewRows([]string{"count"}).AddRow(int64(1))
		mock.ExpectQuery(`^-- name: CountStandings :one.*`).
			WithArgs("", "").
			WillReturnRows(countRows)
		mock.ExpectQuery(`^-- name: ListStandings :many.*`).
			WithArgs("", "", int32(20), int32(20)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), domain.StandingFilter{Page: 2, Size: 20})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
