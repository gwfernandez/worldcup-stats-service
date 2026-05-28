package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// championshipRepository implements ChampionshipRepository using sqlc.
type championshipRepository struct {
	queries *sqlc.Queries
}

// NewChampionshipRepository creates a new ChampionshipRepository.
func NewChampionshipRepository(db sqlc.DBTX) ChampionshipRepository {
	return &championshipRepository{queries: sqlc.New(db)}
}

// List retrieves a paginated list of championships based on the given filters.
func (r *championshipRepository) List(ctx context.Context, filter domain.ChampionshipFilter) ([]domain.Championship, int64, error) {
	total, err := r.queries.CountChampionships(ctx, sqlc.CountChampionshipsParams{
		Column1: int32(filter.Year),
		Column2: filter.Host,
		Column3: filter.ConfederationCode,
	})
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)

	rows, err := r.queries.ListChampionships(ctx, sqlc.ListChampionshipsParams{
		Column1: int32(filter.Year),
		Column2: filter.Host,
		Column3: filter.ConfederationCode,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	championships := make([]domain.Championship, len(rows))
	for i, row := range rows {
		championships[i] = toChampionshipDomain(row)
	}

	return championships, total, nil
}

// GetByYear retrieves a championship and its stats by year.
func (r *championshipRepository) GetByYear(ctx context.Context, year int) (*domain.Championship, error) {
	row, err := r.queries.GetChampionshipByYear(ctx, int32(year))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	championship := toChampionshipDetailDomain(row)
	return &championship, nil
}

func toChampionshipDomain(row sqlc.Championship) domain.Championship {
	var championCode *string
	if row.ChampionCode.Valid {
		val := strings.ToUpper(row.ChampionCode.String)
		championCode = &val
	}

	return domain.Championship{
		Year:         int(row.Year),
		StartDate:    dateToString(row.StartDate),
		EndDate:      dateToString(row.EndDate),
		HostCodes:    uppercaseSlice(row.HostCodes),
		ChampionCode: championCode,
	}
}

func toChampionshipDetailDomain(row sqlc.GetChampionshipByYearRow) domain.Championship {
	var championCode *string
	if row.ChampionCode.Valid {
		val := strings.ToUpper(row.ChampionCode.String)
		championCode = &val
	}

	c := domain.Championship{
		Year:         int(row.Year),
		StartDate:    dateToString(row.StartDate),
		EndDate:      dateToString(row.EndDate),
		HostCodes:    uppercaseSlice(row.HostCodes),
		ChampionCode: championCode,
	}

	// If TotalTeams is valid, it means we have stats in the DB
	if row.TotalTeams.Valid {
		var runnerUpCode string
		if row.StatsRunnerUpCode.Valid {
			runnerUpCode = strings.ToUpper(row.StatsRunnerUpCode.String)
		}

		var thirdPlaceCode string
		if row.StatsThirdPlaceCode.Valid {
			thirdPlaceCode = strings.ToUpper(row.StatsThirdPlaceCode.String)
		}

		var fourthPlaceCode string
		if row.StatsFourthPlaceCode.Valid {
			fourthPlaceCode = strings.ToUpper(row.StatsFourthPlaceCode.String)
		}

		c.Stats = &domain.ChampionshipsStats{
			TotalTeams:      row.TotalTeams.Int32,
			TotalMatches:    row.TotalMatches.Int32,
			TotalStadiums:   row.TotalStadiums.Int32,
			TotalPlayers:    row.TotalPlayers.Int32,
			TotalGoals:      row.TotalGoals.Int32,
			RunnerUpCode:    runnerUpCode,
			ThirdPlaceCode:  thirdPlaceCode,
			FourthPlaceCode: fourthPlaceCode,
			TopScorers:      make([]domain.TopScorer, 0), // Default empty slice until players table exists
			TopScorerGoals:  row.TopScorerGoals.Int32,
		}
	}

	return c
}

func dateToString(date pgtype.Date) string {
	if !date.Valid {
		return ""
	}
	return date.Time.Format("2006-01-02")
}

func uppercaseSlice(slice []string) []string {
	if slice == nil {
		return nil
	}
	res := make([]string, len(slice))
	for i, val := range slice {
		res[i] = strings.ToUpper(val)
	}
	return res
}
