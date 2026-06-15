package repository

import (
	"context"
	"strings"

	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// scorerRepository implements ScorerRepository using sqlc-generated queries.
type scorerRepository struct {
	queries *sqlc.Queries
}

// NewScorerRepository creates a new repository backed by sqlc.
func NewScorerRepository(db sqlc.DBTX) ScorerRepository {
	return &scorerRepository{queries: sqlc.New(db)}
}

func (r *scorerRepository) List(ctx context.Context, filter domain.ScorerFilter) ([]domain.Scorer, int64, error) {
	total, err := r.queries.CountScorers(ctx, sqlc.CountScorersParams{
		Column1: filter.Name,
		Column2: filter.TeamCode,
		Column3: filter.ConfederationCode,
	})
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)
	rows, err := r.queries.ListScorers(ctx, sqlc.ListScorersParams{
		Column1:     filter.Name,
		Column2:     filter.TeamCode,
		Column3:     filter.ConfederationCode,
		Language:    filter.Language,
		LimitValue:  limit,
		OffsetValue: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	scorers := make([]domain.Scorer, len(rows))
	for i, row := range rows {
		scorers[i] = toScorerDomain(row)
	}

	return scorers, total, nil
}

func toScorerDomain(row sqlc.ListScorersRow) domain.Scorer {
	return domain.Scorer{
		FullName:          row.FullName,
		Team:              domain.SimpleTeam{Code: strings.ToUpper(row.TeamCode), Name: row.Name},
		Goals:             row.Goals,
		ListTeams:         uppercaseSlice(row.ListTeams),
		ConfederationCode: strings.ToUpper(row.ConfederationCode),
	}
}
