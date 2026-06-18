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

func (r *championRepository) ListFinalsWonByTeam(ctx context.Context, filter domain.ChampionFinalFilter) ([]domain.ChampionFinal, int64, error) {
	total, err := r.queries.CountFinalsWonByTeam(ctx, filter.TeamCode)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.queries.ListFinalsWonByTeam(ctx, sqlc.ListFinalsWonByTeamParams{
		TeamCode:    filter.TeamCode,
		LimitValue:  int32(filter.Size),
		OffsetValue: int32((filter.Page - 1) * filter.Size),
	})
	if err != nil {
		return nil, 0, err
	}

	finals := make([]domain.ChampionFinal, len(rows))
	for i, row := range rows {
		finals[i] = toChampionFinalDomain(row)
	}

	return finals, total, nil
}

func toChampionDomain(row sqlc.ListChampionsRow) domain.Champion {
	return domain.Champion{
		Team: domain.SimpleTeam{
			Code: strings.ToUpper(row.TeamCode),
		},
		Wins:              row.Wins,
		Years:             row.Years,
		ConfederationCode: strings.ToUpper(row.ConfederationCode),
	}
}

func toChampionFinalDomain(row sqlc.ListFinalsWonByTeamRow) domain.ChampionFinal {
	return domain.ChampionFinal{
		Year:                   row.Year,
		MatchDate:              datePtr(row.MatchDate),
		MatchTime:              timePtr(row.MatchTime),
		HomeTeam:               domain.SimpleTeam{Code: strings.ToUpper(row.HomeTeamCode)},
		HomeTeamScore:          int4Ptr(row.HomeTeamScore),
		HomeTeamScorePenalties: int4Ptr(row.HomeTeamScorePenalties),
		AwayTeam:               domain.SimpleTeam{Code: strings.ToUpper(row.AwayTeamCode)},
		AwayTeamScore:          int4Ptr(row.AwayTeamScore),
		AwayTeamScorePenalties: int4Ptr(row.AwayTeamScorePenalties),
	}
}
