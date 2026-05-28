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

// teamRepository implements TeamRepository using sqlc-generated queries.
type teamRepository struct {
	queries *sqlc.Queries
}

// NewTeamRepository creates a new repository backed by sqlc.
func NewTeamRepository(db sqlc.DBTX) TeamRepository {
	return &teamRepository{queries: sqlc.New(db)}
}

func (r *teamRepository) List(ctx context.Context, filter domain.TeamFilter) ([]domain.Team, int64, error) {
	confederationCode := ""
	if filter.ConfederationCode != nil {
		confederationCode = *filter.ConfederationCode
	}

	total, err := r.queries.CountTeams(ctx, sqlc.CountTeamsParams{
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
	rows, err := r.queries.ListTeams(ctx, sqlc.ListTeamsParams{
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

	teams := make([]domain.Team, len(rows))
	for i, row := range rows {
		teams[i] = toTeamDomain(row)
	}

	return teams, total, nil
}

func (r *teamRepository) GetByCode(ctx context.Context, code string) (*domain.Team, error) {
	row, err := r.queries.GetTeamByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	team := toTeamDomain(row)
	return &team, nil
}

func toTeamDomain(row sqlc.Team) domain.Team {
	dissolutionDate := dateToStringPtr(row.DissolutionDate)

	return domain.Team{
		Name:              row.Name,
		Code:              strings.ToUpper(row.Code),
		DissolutionDate:   dissolutionDate,
		IsDissolved:       dissolutionDate != nil,
		ConfederationCode: strings.ToUpper(row.ConfederationCode),
		FederationName:    row.FederationName,
		FederationCode:    strings.ToUpper(row.FederationCode),
	}
}

func dateToStringPtr(date pgtype.Date) *string {
	if !date.Valid {
		return nil
	}
	formatted := date.Time.Format("2006-01-02")
	return &formatted
}
