package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

func TestNationalTeamRepository_List(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewNationalTeamRepository(mock)
		confederationCode := "CONMEBOL"
		filter := domain.NationalTeamFilter{
			Name:              "argen",
			ConfederationCode: &confederationCode,
			FederationName:    "futbol",
			FederationCode:    "afa",
			IncludeDissolved:  true,
			Page:              1,
			Size:              20,
		}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(2))
		mock.ExpectQuery(`^-- name: CountNationalTeams :one.*`).
			WithArgs("argen", "CONMEBOL", "futbol", "afa", true).
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"id", "name", "code", "dissolution_date", "confederation_code", "federation_name", "federation_code"}).
			AddRow(int64(1), "Argentina", "arg", nil, "CONMEBOL", "Asociación del Fútbol Argentino", "afa").
			AddRow(int64(5), "Soviet Union", "urs", time.Date(1991, 12, 26, 0, 0, 0, 0, time.UTC), "UEFA", "Football Federation of the Soviet Union", "ffsu")

		mock.ExpectQuery(`^-- name: ListNationalTeams :many.*`).
			WithArgs("argen", "CONMEBOL", "futbol", "afa", true, int32(20), int32(0)).
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, result, 2)
		assert.Equal(t, "ARG", result[0].Code)
		assert.Equal(t, "AFA", result[0].FederationCode)
		assert.Equal(t, "CONMEBOL", result[0].ConfederationCode)
		assert.False(t, result[0].IsDissolved)
		assert.True(t, result[1].IsDissolved)
		require.NotNil(t, result[1].DissolutionDate)
		assert.Equal(t, "1991-12-26", *result[1].DissolutionDate)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error on count", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewNationalTeamRepository(mock)
		filter := domain.NationalTeamFilter{Page: 1, Size: 20}

		mock.ExpectQuery(`^-- name: CountNationalTeams :one.*`).
			WithArgs("", "", "", "", false).
			WillReturnError(errors.New("db error"))

		result, total, err := repo.List(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNationalTeamRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewNationalTeamRepository(mock)
		rows := mock.NewRows([]string{"id", "name", "code", "dissolution_date", "confederation_code", "federation_name", "federation_code"}).
			AddRow(int64(1), "Argentina", "arg", nil, "CONMEBOL", "Asociación del Fútbol Argentino", "afa")
		mock.ExpectQuery(`^-- name: GetNationalTeamByID :one.*`).WithArgs(int64(1)).WillReturnRows(rows)

		result, err := repo.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, "ARG", result.Code)
		assert.Equal(t, "AFA", result.FederationCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewNationalTeamRepository(mock)
		mock.ExpectQuery(`^-- name: GetNationalTeamByID :one.*`).WithArgs(int64(999)).WillReturnError(pgx.ErrNoRows)

		result, err := repo.GetByID(context.Background(), 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNationalTeamRepository_GetByCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewNationalTeamRepository(mock)
		rows := mock.NewRows([]string{"id", "name", "code", "dissolution_date", "confederation_code", "federation_name", "federation_code"}).
			AddRow(int64(5), "Soviet Union", "urs", time.Date(1991, 12, 26, 0, 0, 0, 0, time.UTC), "UEFA", "Football Federation of the Soviet Union", "ffsu")
		mock.ExpectQuery(`^-- name: GetNationalTeamByCode :one.*`).WithArgs("urs").WillReturnRows(rows)

		result, err := repo.GetByCode(context.Background(), "urs")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "URS", result.Code)
		assert.Equal(t, "FFSU", result.FederationCode)
		assert.True(t, result.IsDissolved)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewNationalTeamRepository(mock)
		mock.ExpectQuery(`^-- name: GetNationalTeamByCode :one.*`).WithArgs("zzz").WillReturnError(pgx.ErrNoRows)

		result, err := repo.GetByCode(context.Background(), "zzz")
		assert.NoError(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
