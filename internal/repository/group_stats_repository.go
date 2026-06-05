package repository

import (
	"context"
	"strings"

	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// groupStatsRepository implements GroupStatsRepository using sqlc-generated queries.
type groupStatsRepository struct {
	queries *sqlc.Queries
}

// NewGroupStatsRepository creates a new repository backed by sqlc.
func NewGroupStatsRepository(db sqlc.DBTX) GroupStatsRepository {
	return &groupStatsRepository{queries: sqlc.New(db)}
}

func (r *groupStatsRepository) ListByYear(ctx context.Context, year int) ([]domain.GroupStandingRecord, error) {
	rows, err := r.queries.ListGroupsStatsByYear(ctx, int32(year))
	if err != nil {
		return nil, err
	}

	standings := make([]domain.GroupStandingRecord, len(rows))
	for i, row := range rows {
		standings[i] = toGroupStandingRecord(row)
	}

	return standings, nil
}

func toGroupStandingRecord(row sqlc.ChampionshipsGroupsStat) domain.GroupStandingRecord {
	return domain.GroupStandingRecord{
		Stage:     enumString(row.Stage),
		GroupCode: strings.ToUpper(row.GroupCode),
		Standing: domain.GroupStanding{
			TeamCode:       strings.ToUpper(row.TeamCode),
			MatchesPlayed:  row.MatchesPlayed,
			Wins:           row.Wins,
			Draws:          row.Draws,
			Losses:         row.Losses,
			GoalsFor:       row.GoalsFor,
			GoalsAgainst:   row.GoalsAgainst,
			GoalDifference: row.GoalDifference.Int32,
			Points:         row.Points,
			UnifiedPoints:  row.UnifiedPoints,
			Position:       int4Ptr(row.Position),
		},
	}
}
