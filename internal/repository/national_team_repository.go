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

// nationalTeamRepository implements NationalTeamRepository using sqlc-generated queries.
type nationalTeamRepository struct {
	queries *sqlc.Queries
}

// NewNationalTeamRepository creates a new repository backed by sqlc.
func NewNationalTeamRepository(db sqlc.DBTX) NationalTeamRepository {
	return &nationalTeamRepository{queries: sqlc.New(db)}
}

func (r *nationalTeamRepository) List(ctx context.Context, filter domain.NationalTeamFilter) ([]domain.NationalTeam, int64, error) {
	confederationCode := ""
	if filter.ConfederationCode != nil {
		confederationCode = *filter.ConfederationCode
	}

	total, err := r.queries.CountNationalTeams(ctx, sqlc.CountNationalTeamsParams{
		Column1: filter.Name,
		Column2: confederationCode,
		Column3: filter.FederationName,
		Column4: filter.FederationCode,
		Column5: filter.IncludeDissolved,
	})
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)
	rows, err := r.queries.ListNationalTeams(ctx, sqlc.ListNationalTeamsParams{
		Column1: filter.Name,
		Column2: confederationCode,
		Column3: filter.FederationName,
		Column4: filter.FederationCode,
		Column5: filter.IncludeDissolved,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	teams := make([]domain.NationalTeam, len(rows))
	for i, row := range rows {
		teams[i] = toNationalTeamDomain(row)
	}

	return teams, total, nil
}

func (r *nationalTeamRepository) GetByCode(ctx context.Context, code string) (*domain.NationalTeam, error) {
	row, err := r.queries.GetNationalTeamByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	team := toNationalTeamDomain(row)
	return &team, nil
}

func toNationalTeamDomain(row sqlc.NationalTeam) domain.NationalTeam {
	dissolutionDate := dateToStringPtr(row.DissolutionDate)

	return domain.NationalTeam{
		Name:            row.Name,
		Code:            strings.ToUpper(row.Code),
		DissolutionDate: dissolutionDate,
		IsDissolved:     dissolutionDate != nil,
		ConfederationCode: strings.ToUpper(row.ConfederationCode),
		FederationName:  row.FederationName,
		FederationCode:  strings.ToUpper(row.FederationCode),
	}
}

func dateToStringPtr(date pgtype.Date) *string {
	if !date.Valid {
		return nil
	}
	formatted := date.Time.Format("2006-01-02")
	return &formatted
}
