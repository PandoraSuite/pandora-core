package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ServiceRepository) FindByNameAndVersion(
	ctx context.Context, name, version string,
) (*entities.Service, *errors.Error) {
	query := `SELECT * FROM service WHERE name = $1 AND version = $2;`

	service := new(entities.Service)
	err := r.pool.QueryRow(ctx, query, name, version).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.Status,
		&service.CreatedAt,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return service, nil
}

func (r *ServiceRepository) FindActiveServices(ctx context.Context) ([]*entities.Service, *errors.Error) {
	query := "SELECT * FROM service WHERE status = 'active';"
	return r.find(ctx, query)
}

func (r *ServiceRepository) FindAll(ctx context.Context) ([]*entities.Service, *errors.Error) {
	query := "SELECT * FROM service;"
	return r.find(ctx, query)
}

func (r *ServiceRepository) find(
	ctx context.Context, query string,
) ([]*entities.Service, *errors.Error) {
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var services []*entities.Service
	for rows.Next() {
		service := new(entities.Service)

		err = rows.Scan(
			&service.ID,
			&service.Name,
			&service.Version,
			&service.Status,
			&service.CreatedAt,
		)
		if err != nil {
			return nil, r.handlerErr(err)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return services, nil
}

func (r *ServiceRepository) Save(
	ctx context.Context, service *entities.Service,
) *errors.Error {
	query := `
		INSERT INTO service (name, version, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		service.Name,
		service.Version,
		service.Status,
	).Scan(&service.ID, &service.CreatedAt)

	return r.handlerErr(err)
}

func NewServiceRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *ServiceRepository {
	return &ServiceRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
