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

func TestScorerRepository_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		filter := domain.ScorerFilter{Name: "messi", Language: "en", TeamCode: "ARG", ConfederationCode: "CONMEBOL", Page: 2, Size: 10}

		countRows := mock.NewRows([]string{"count"}).AddRow(int64(11))
		mock.ExpectQuery(`^-- name: CountScorers :one.*`).
			WithArgs("messi", "ARG", "CONMEBOL").
			WillReturnRows(countRows)

		rows := mock.NewRows([]string{"player_id", "full_name", "team_code", "goals", "list_teams", "confederation_code"}).
			AddRow(int64(10), "Lionel Messi", "arg", int32(13), []string{"arg"}, "conmebol").
			AddRow(int64(11), "Gabriel Batistuta", "arg", int32(10), []string{"arg"}, "conmebol")
		mock.ExpectQuery(`^-- name: ListScorers :many.*`).
			WithArgs("messi", "ARG", "CONMEBOL", int32(10), int32(10)).
			WillReturnRows(rows)

		result, total, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(11), total)
		require.Len(t, result, 2)
		assert.Equal(t, int64(10), result[0].PlayerID)
		assert.Equal(t, "Lionel Messi", result[0].FullName)
		assert.Equal(t, "ARG", result[0].Team.Code)
		assert.Empty(t, result[0].Team.Name)
		assert.Equal(t, int32(13), result[0].Goals)
		assert.Equal(t, []string{"ARG"}, result[0].ListTeams)
		assert.Equal(t, "CONMEBOL", result[0].ConfederationCode)
		assert.Equal(t, int64(11), result[1].PlayerID)
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

func TestScorerRepository_GetByID(t *testing.T) {
	t.Run("success preserves order and maps goals", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		matchDate := pgtype.Date{Time: time.Date(2022, 12, 18, 0, 0, 0, 0, time.UTC), Valid: true}

		playerRows := mock.NewRows([]string{"id", "first_name", "last_name", "position", "list_championships", "list_teams"}).
			AddRow(int64(1524), "Lionel", "Messi", "FW", []int32{2006, 2010, 2022}, []string{"arg", "fcb"})
		mock.ExpectQuery(`^-- name: GetScorerByPlayerID :one.*`).
			WithArgs(int64(1524)).
			WillReturnRows(playerRows)

		goalRows := mock.NewRows([]string{
			"year", "host_codes", "match_date", "opponent_team_code",
			"minute_regular", "penalty", "stage",
		}).AddRow(
			int32(2022), []string{"qat"}, matchDate, "fra",
			int32(23), pgtype.Bool{Bool: true, Valid: true}, "final",
		)
		mock.ExpectQuery(`^-- name: ListScorerGoalsByPlayer :many.*`).
			WithArgs(pgtype.Int8{Int64: 1524, Valid: true}).
			WillReturnRows(goalRows)

		result, err := repo.GetByID(context.Background(), 1524)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, int64(1524), result.ID)
		assert.Equal(t, "Lionel", result.FirstName)
		assert.Equal(t, "Messi", result.LastName)
		require.NotNil(t, result.Position)
		assert.Equal(t, "FW", *result.Position)
		assert.Equal(t, []int32{2006, 2010, 2022}, result.Championships)
		assert.Equal(t, []domain.SimpleTeam{{Code: "ARG"}, {Code: "FCB"}}, result.Teams)
		require.Len(t, result.Goals, 1)
		assert.Equal(t, int32(2022), result.Goals[0].Year)
		assert.Equal(t, []domain.SimpleTeam{{Code: "QAT"}}, result.Goals[0].Hosts)
		assert.Equal(t, "FRA", result.Goals[0].OpponentTeam.Code)
		assert.Equal(t, "2022-12-18", *result.Goals[0].MatchDate)
		assert.Equal(t, "final", *result.Goals[0].Stage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("nullable position and no goals", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		playerRows := mock.NewRows([]string{"id", "first_name", "last_name", "position", "list_championships", "list_teams"}).
			AddRow(int64(7), "Juan", "Pérez", "", []int32{}, []string{})
		mock.ExpectQuery(`^-- name: GetScorerByPlayerID :one.*`).
			WithArgs(int64(7)).
			WillReturnRows(playerRows)
		mock.ExpectQuery(`^-- name: ListScorerGoalsByPlayer :many.*`).
			WithArgs(pgtype.Int8{Int64: 7, Valid: true}).
			WillReturnRows(mock.NewRows([]string{
				"year", "host_codes", "match_date", "opponent_team_code",
				"minute_regular", "penalty", "stage",
			}))

		result, err := repo.GetByID(context.Background(), 7)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Nil(t, result.Position)
		assert.Empty(t, result.Goals)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("player not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		mock.ExpectQuery(`^-- name: GetScorerByPlayerID :one.*`).
			WithArgs(int64(999)).
			WillReturnError(pgx.ErrNoRows)

		result, err := repo.GetByID(context.Background(), 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error retrieving player", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		mock.ExpectQuery(`^-- name: GetScorerByPlayerID :one.*`).
			WithArgs(int64(1524)).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetByID(context.Background(), 1524)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error retrieving goals", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewScorerRepository(mock)
		playerRows := mock.NewRows([]string{"id", "first_name", "last_name", "position", "list_championships", "list_teams"}).
			AddRow(int64(1524), "Lionel", "Messi", "FW", []int32{2022}, []string{"arg"})
		mock.ExpectQuery(`^-- name: GetScorerByPlayerID :one.*`).
			WithArgs(int64(1524)).
			WillReturnRows(playerRows)
		mock.ExpectQuery(`^-- name: ListScorerGoalsByPlayer :many.*`).
			WithArgs(pgtype.Int8{Int64: 1524, Valid: true}).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetByID(context.Background(), 1524)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
