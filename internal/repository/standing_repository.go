package repository

import (
	"context"
	"strings"

	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// standingRepository implements StandingRepository using sqlc-generated queries.
type standingRepository struct {
	queries *sqlc.Queries
}

// NewStandingRepository creates a new repository backed by sqlc.
func NewStandingRepository(db sqlc.DBTX) StandingRepository {
	return &standingRepository{queries: sqlc.New(db)}
}

func (r *standingRepository) List(ctx context.Context, filter domain.StandingFilter) ([]domain.Standing, int64, error) {
	total, err := r.queries.CountStandings(ctx, sqlc.CountStandingsParams{
		Language:          filter.Language,
		NameFilter:        filter.Name,
		ConfederationCode: filter.ConfederationCode,
	})
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)
	rows, err := r.queries.ListStandings(ctx, sqlc.ListStandingsParams{
		Language:          filter.Language,
		NameFilter:        filter.Name,
		ConfederationCode: filter.ConfederationCode,
		LimitValue:        limit,
		OffsetValue:       offset,
	})
	if err != nil {
		return nil, 0, err
	}

	standings := make([]domain.Standing, len(rows))
	for i, row := range rows {
		standings[i] = toStandingDomain(row)
	}

	return standings, total, nil
}

func toStandingDomain(row sqlc.ListStandingsRow) domain.Standing {
	return domain.Standing{
		Team:            domain.SimpleTeam{Code: strings.ToUpper(row.TeamCode), Name: row.Name},
		MatchesPlayed:   row.MatchesPlayed,
		Wins:            row.Wins,
		Draws:           row.Draws,
		Losses:          row.Losses,
		GoalsFor:        row.GoalsFor,
		GoalsAgainst:    row.GoalsAgainst,
		GoalDifference:  row.GoalDifference,
		Points:          row.Points,
		UnifiedPoints:   row.UnifiedPoints,
		Position:        row.Position,
		UnifiedPosition: row.UnifiedPosition,
	}
}
