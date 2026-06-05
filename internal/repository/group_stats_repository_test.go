package repository_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

func TestGroupStatsRepository_ListByYear(t *testing.T) {
	t.Run("success with standings", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewGroupStatsRepository(mock)
		rows := mock.NewRows([]string{
			"year", "stage", "group_code", "team_code", "matches_played", "wins", "draws", "losses",
			"goals_for", "goals_against", "goal_difference", "points", "unified_points", "position",
		}).AddRow(
			int32(1978), "group_stage", "a", "arg", int32(3), int32(2), int32(1), int32(0),
			int32(4), int32(1), pgtype.Int4{Int32: 3, Valid: true}, int32(5), int32(7),
			pgtype.Int4{Int32: 1, Valid: true},
		)

		mock.ExpectQuery(`^-- name: ListGroupsStatsByYear :many.*`).
			WithArgs(int32(1978)).
			WillReturnRows(rows)

		result, err := repo.ListByYear(context.Background(), 1978)
		assert.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "group_stage", result[0].Stage)
		assert.Equal(t, "A", result[0].GroupCode)
		assert.Equal(t, "ARG", result[0].Standing.TeamCode)
		assert.Equal(t, int32(3), result[0].Standing.GoalDifference)
		require.NotNil(t, result[0].Standing.Position)
		assert.Equal(t, int32(1), *result[0].Standing.Position)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewGroupStatsRepository(mock)
		rows := mock.NewRows([]string{
			"year", "stage", "group_code", "team_code", "matches_played", "wins", "draws", "losses",
			"goals_for", "goals_against", "goal_difference", "points", "unified_points", "position",
		})

		mock.ExpectQuery(`^-- name: ListGroupsStatsByYear :many.*`).
			WithArgs(int32(1938)).
			WillReturnRows(rows)

		result, err := repo.ListByYear(context.Background(), 1938)
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
