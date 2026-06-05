package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

func TestMatchRepository_ListByYear(t *testing.T) {
	t.Run("success with matches", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewMatchRepository(mock)
		matchDate := time.Date(1938, 6, 5, 0, 0, 0, 0, time.UTC)
		matchTime := pgtype.Time{Microseconds: 17 * 60 * 60 * 1000 * 1000, Valid: true}

		rows := mock.NewRows([]string{
			"id", "year", "stage", "stage_type", "group_code", "replayed", "replay_of", "match_date", "match_time",
			"stadium_id", "home_team_code", "away_team_code", "home_team_score", "away_team_score", "extra_time",
			"penalty_shootout", "home_team_score_penalties", "away_team_score_penalties", "home_team_win",
			"away_team_win", "draw", "ref_id",
		}).AddRow(
			int64(10), int32(1938), "round_of_16", "knockout", nil, true, pgtype.Int8{Int64: 9, Valid: true},
			matchDate, matchTime, pgtype.Int8{Int64: 3, Valid: true}, "bra", "pol",
			pgtype.Int4{Int32: 6, Valid: true}, pgtype.Int4{Int32: 5, Valid: true}, true, false,
			pgtype.Int4{Valid: false}, pgtype.Int4{Valid: false}, pgtype.Bool{Bool: true, Valid: true},
			pgtype.Bool{Bool: false, Valid: true}, pgtype.Bool{Bool: false, Valid: true}, pgtype.Text{String: "M-1938-10", Valid: true},
		)

		mock.ExpectQuery(`^-- name: ListMatchesByYear :many.*`).
			WithArgs(int32(1938)).
			WillReturnRows(rows)

		result, err := repo.ListByYear(context.Background(), 1938)
		assert.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "round_of_16", result[0].Stage)
		assert.Equal(t, "", result[0].GroupCode)
		assert.Equal(t, int64(10), result[0].Match.ID)
		assert.Equal(t, "knockout", result[0].Match.StageType)
		assert.True(t, result[0].Match.Replayed)
		require.NotNil(t, result[0].Match.ReplayOf)
		assert.Equal(t, int64(9), *result[0].Match.ReplayOf)
		require.NotNil(t, result[0].Match.MatchDate)
		assert.Equal(t, "1938-06-05", *result[0].Match.MatchDate)
		require.NotNil(t, result[0].Match.MatchTime)
		assert.Equal(t, "17:00:00", *result[0].Match.MatchTime)
		assert.Equal(t, "BRA", result[0].Match.HomeTeamCode)
		assert.Equal(t, "POL", result[0].Match.AwayTeamCode)
		require.NotNil(t, result[0].Match.HomeTeamWin)
		assert.True(t, *result[0].Match.HomeTeamWin)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mock.Close()

		repo := repository.NewMatchRepository(mock)
		rows := mock.NewRows([]string{
			"id", "year", "stage", "stage_type", "group_code", "replayed", "replay_of", "match_date", "match_time",
			"stadium_id", "home_team_code", "away_team_code", "home_team_score", "away_team_score", "extra_time",
			"penalty_shootout", "home_team_score_penalties", "away_team_score_penalties", "home_team_win",
			"away_team_win", "draw", "ref_id",
		})

		mock.ExpectQuery(`^-- name: ListMatchesByYear :many.*`).
			WithArgs(int32(2023)).
			WillReturnRows(rows)

		result, err := repo.ListByYear(context.Background(), 2023)
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
