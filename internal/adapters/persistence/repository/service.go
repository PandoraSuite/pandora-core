package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ServiceRepository) UpdateStatus(
	ctx context.Context, id int, status enums.ServiceStatus,
) *errors.Error {
	if status == enums.ServiceStatusNull {
		return errors.ErrServiceInvalidStatus
	}

	query := `
		UPDATE service
		SET status = $1
		WHERE id = $2;
	`

	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrServiceNotFound
	}

	return nil
}

func (r *ServiceRepository) FindByNameAndVersion(
	ctx context.Context, name, version string,
) (*entities.Service, *errors.Error) {
	query := `
		SELECT id, name, version, status, created_at
		FROM service
		WHERE name = $1 AND version = $2;
	`

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

func (r *ServiceRepository) FindAll(
	ctx context.Context, filter *dto.ServiceFilter,
) ([]*entities.Service, *errors.Error) {
	query := `
		SELECT id, name, version, status, created_at
		FROM service
	`

	var args []any
	if filter != nil {
		var where []string
		argIndex := 1

		if filter.Status != enums.ServiceStatusNull {
			where = append(where, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, filter.Status)
			argIndex++
		}

		if len(where) > 0 {
			query = fmt.Sprintf(
				"%s WHERE %s", query, strings.Join(where, " AND "),
			)
		}
	}

	query += ";"

	rows, err := r.pool.Query(ctx, query, args...)
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
