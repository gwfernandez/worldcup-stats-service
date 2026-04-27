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

func (r *confederationRepository) Create(ctx context.Context, code, name string) (*domain.Confederation, error) {
	row, err := r.queries.CreateConfederation(ctx, sqlc.CreateConfederationParams{
		Code: code,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	c := toDomain(row)
	return &c, nil
}

func (r *confederationRepository) Update(ctx context.Context, id int64, code, name string) (*domain.Confederation, error) {
	row, err := r.queries.UpdateConfederation(ctx, sqlc.UpdateConfederationParams{
		ID:   id,
		Code: code,
		Name: name,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	c := toDomain(row)
	return &c, nil
}

func (r *confederationRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteConfederation(ctx, id)
}

// toDomain converts a sqlc model to a domain entity.
func toDomain(row sqlc.Confederation) domain.Confederation {
	return domain.Confederation{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}
}
