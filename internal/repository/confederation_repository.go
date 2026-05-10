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

func (r *confederationRepository) List(ctx context.Context) ([]domain.Confederation, error) {
	rows, err := r.queries.ListConfederations(ctx)
	if err != nil {
		return nil, err
	}

	confederations := make([]domain.Confederation, len(rows))
	for i, row := range rows {
		confederations[i] = toDomain(row)
	}
	return confederations, nil
}

func (r *confederationRepository) GetByID(ctx context.Context, id int64) (*domain.Confederation, error) {
	row, err := r.queries.GetConfederation(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	c := toDomain(row)
	return &c, nil
}

// toDomain converts a sqlc model to a domain entity.
func toDomain(row sqlc.Confederation) domain.Confederation {
	return domain.Confederation{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}
}
