package repository

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// goalRepository implements GoalRepository using sqlc-generated queries.
type goalRepository struct {
	queries *sqlc.Queries
}

// NewGoalRepository creates a new repository backed by sqlc.
func NewGoalRepository(db sqlc.DBTX) GoalRepository {
	return &goalRepository{queries: sqlc.New(db)}
}

func (r *goalRepository) ListByPlayer(ctx context.Context, filter domain.GoalFilter) ([]domain.Goal, int64, error) {
	total, err := r.queries.CountGoalsByPlayer(ctx, sqlc.CountGoalsByPlayerParams{
		PlayerID: pgtype.Int8{Int64: filter.PlayerID, Valid: true},
		Year:     int32(filter.Year),
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.queries.ListGoalsByPlayer(ctx, sqlc.ListGoalsByPlayerParams{
		PlayerID:    pgtype.Int8{Int64: filter.PlayerID, Valid: true},
		Year:        int32(filter.Year),
		LimitValue:  int32(filter.Size),
		OffsetValue: int32((filter.Page - 1) * filter.Size),
	})
	if err != nil {
		return nil, 0, err
	}

	goals := make([]domain.Goal, len(rows))
	for i, row := range rows {
		goals[i] = toGoalDomain(row)
	}

	return goals, total, nil
}

func toGoalDomain(row sqlc.ListGoalsByPlayerRow) domain.Goal {
	return domain.Goal{
		Year:          row.Year,
		MatchDate:     datePtr(row.MatchDate),
		OpponentTeam:  domain.SimpleTeam{Code: strings.ToUpper(row.OpponentTeamCode)},
		MinuteRegular: row.MinuteRegular,
		Penalty:       boolPtr(row.Penalty),
		Stage:         nonEmptyStringPtr(enumString(row.Stage)),
	}
}

func nonEmptyStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
