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

func TestScorerRepository_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		filter := domain.ScorerFilter{Name: "messi", TeamCode: "ARG", ConfederationCode: "CONMEBOL", Page: 2, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(11))
		mock.ExpectQuery(`^-- name: CountScorers :one.*`).
			WithArgs("messi", "ARG", "CONMEBOL").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"full_name", "team_code", "goals", "list_teams", "confederation_code"}).
			AddRow("Lionel Messi", "arg", int32(13), []string{"arg"}, "conmebol").
			AddRow("Gabriel Batistuta", "arg", int32(10), []string{"arg"}, "conmebol")
		mock.ExpectQuery(`^-- name: ListScorers :many.*`).
			WithArgs("messi", "ARG", "CONMEBOL", int32(10), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(11), total)
		require.Len(t, result, 2)
		assert.Equal(t, "Lionel Messi", result[0].FullName)
		assert.Equal(t, "ARG", result[0].TeamCode)
		assert.Equal(t, int32(13), result[0].Goals)
		assert.Equal(t, []string{"ARG"}, result[0].ListTeams)
		assert.Equal(t, "CONMEBOL", result[0].ConfederationCode)
		assert.Equal(t, "Gabriel Batistuta", result[1].FullName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		mock.ExpectQuery(`^-- name: CountScorers :one.*`).
			WithArgs("", "", "").
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), domain.ScorerFilter{Page: 1, Size: 20})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on list", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		countRows := mock.NewRows([]string{"count"}).AddRow(int64(1))
		mock.ExpectQuery(`^-- name: CountScorers :one.*`).
			WithArgs("", "", "").
			WillReturnRows(countRows)
		mock.ExpectQuery(`^-- name: ListScorers :many.*`).
			WithArgs("", "", "", int32(20), int32(20)).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), domain.ScorerFilter{Page: 2, Size: 20})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
