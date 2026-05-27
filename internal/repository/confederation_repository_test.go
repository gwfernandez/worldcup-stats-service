package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

func TestConfederationRepository_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewConfederationRepository(mock)

		rows := mock.NewRows([]string{"code", "name"}).
			AddRow("CONMEBOL", "Confederación Sudamericana de Fútbol").
			AddRow("UEFA", "Union of European Football Associations")

		mock.ExpectQuery(`^-- name: ListConfederations :many.*`).WillReturnRows(rows)

		ctx := context.Background()
		result, err := repo.List(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewConfederationRepository(mock)

		mock.ExpectQuery(`^-- name: ListConfederations :many.*`).WillReturnError(errors.New("db error"))

		ctx := context.Background()
		result, err := repo.List(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestConfederationRepository_GetByCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewConfederationRepository(mock)

		rows := mock.NewRows([]string{"code", "name"}).
			AddRow( "CONMEBOL", "Confederación Sudamericana de Fútbol")

		mock.ExpectQuery(`^-- name: GetConfederationByCode :one.*`).WithArgs("CONMEBOL").WillReturnRows(rows)

		ctx := context.Background()
		result, err := repo.GetByCode(ctx, "CONMEBOL")

		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "CONMEBOL", result.Code)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewConfederationRepository(mock)

		mock.ExpectQuery(`^-- name: GetConfederationByCode :one.*`).WithArgs("ANTARCTICA").WillReturnError(pgx.ErrNoRows)

		ctx := context.Background()
		result, err := repo.GetByCode(ctx, "ANTARCTICA")

		assert.NoError(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewConfederationRepository(mock)

		mock.ExpectQuery(`^-- name: GetConfederationByCode :one.*`).WithArgs("ANTARCTICA").WillReturnError(errors.New("db error"))

		ctx := context.Background()
		result, err := repo.GetByCode(ctx, "ANTARCTICA")

		assert.Error(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
