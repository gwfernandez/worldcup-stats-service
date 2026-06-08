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

func TestChampionRepository_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewChampionRepository(mock)
		filter := domain.ChampionFilter{Page: 1, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(2))
		mock.ExpectQuery(`^-- name: CountChampions :one.*`).
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"team_code", "name", "wins", "years"}).
			AddRow("bra", "Brasil", int64(5), []int32{1958, 1962, 1970, 1994, 2002}).
			AddRow("arg", "Argentina", int64(3), []int32{1978, 1986, 2022})
		mock.ExpectQuery(`^-- name: ListChampions :many.*`).
			WithArgs(int32(10), int32(0)).
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, result, 2)
		assert.Equal(t, "BRA", result[0].TeamCode)
		assert.Equal(t, "Brasil", result[0].Name)
		assert.Equal(t, int64(5), result[0].Wins)
		assert.Equal(t, []int32{1958, 1962, 1970, 1994, 2002}, result[0].Years)
		assert.Equal(t, "ARG", result[1].TeamCode)
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
