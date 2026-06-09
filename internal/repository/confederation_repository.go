package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// confederationRepository implements ConfederationRepository using sqlc-generated queries.
type confederationRepository struct {
	queries *sqlc.Queries
}

// NewConfederationRepository creates a new repository backed by sqlc.
func NewConfederationRepository(db sqlc.DBTX) ConfederationRepository {
	return &confederationRepository{
		queries: sqlc.New(db),
	}
}

func (r *confederationRepository) List(ctx context.Context, language string) ([]domain.Confederation, error) {
	rows, err := r.queries.ListConfederations(ctx, language)
	if err != nil {
		return nil, err
	}

	confederations := make([]domain.Confederation, len(rows))
	for i, row := range rows {
		confederations[i] = toDomain(row.Code, row.Name)
	}
	return confederations, nil
}

func (r *confederationRepository) GetByCode(ctx context.Context, code, language string) (*domain.Confederation, error) {
	row, err := r.queries.GetConfederationByCode(ctx, sqlc.GetConfederationByCodeParams{
		Language: language,
		Code:     code,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	c := toDomain(row.Code, row.Name)
	return &c, nil
}

// toDomain builds a domain entity from sqlc row values.
func toDomain(code, name string) domain.Confederation {
	return domain.Confederation{
		Code: code,
		Name: name,
	}
}
