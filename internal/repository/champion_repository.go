package repository

import (
	"context"
	"strings"

	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// championRepository implements ChampionRepository using sqlc-generated queries.
type championRepository struct {
	queries *sqlc.Queries
}

// NewChampionRepository creates a new repository backed by sqlc.
func NewChampionRepository(db sqlc.DBTX) ChampionRepository {
	return &championRepository{queries: sqlc.New(db)}
}

func (r *championRepository) List(ctx context.Context, filter domain.ChampionFilter) ([]domain.Champion, int64, error) {
	total, err := r.queries.CountChampions(ctx)
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)
	rows, err := r.queries.ListChampions(ctx, sqlc.ListChampionsParams{
		LimitValue:  limit,
		OffsetValue: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	champions := make([]domain.Champion, len(rows))
	for i, row := range rows {
		champions[i] = toChampionDomain(row)
	}

	return champions, total, nil
}

func toChampionDomain(row sqlc.ListChampionsRow) domain.Champion {
	return domain.Champion{
		Team: domain.SimpleTeam{
			Code: strings.ToUpper(row.TeamCode),
		},
		Wins:  row.Wins,
		Years: row.Years,
	}
}
