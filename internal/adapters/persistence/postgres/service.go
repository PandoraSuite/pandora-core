package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceRepository struct {
	*Driver

	tableName string
}

func (r *ServiceRepository) Exists(
	ctx context.Context, id int,
) (bool, errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM service
			WHERE id = $1
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, r.errorMapper(err, r.tableName)
	}

	return exists, nil
}

func (r *ServiceRepository) Delete(
	ctx context.Context, id int,
) errors.Error {
	query := `
		DELETE FROM service
		WHERE id = $1;
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return r.errorMapper(err, r.tableName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(r.tableName, map[string]any{"id": id})
	}

	return nil
}

func (r *ServiceRepository) UpdateStatus(
	ctx context.Context, id int, status enums.ServiceStatus,
) (*entities.Service, errors.Error) {
	query := `
		UPDATE service
		SET status = $1
		WHERE id = $2
		RETURNING id, name, version, status, created_at;
	`

	service := new(entities.Service)
	err := r.pool.QueryRow(ctx, query, status, id).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.Status,
		&service.CreatedAt,
	)
	if err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	return service, nil
}

func (r *ServiceRepository) GetByNameAndVersion(
	ctx context.Context, name, version string,
) (*entities.Service, errors.Error) {
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
		return nil, r.errorMapper(err, r.tableName)
	}

	return service, nil
}

func (r *ServiceRepository) List(
	ctx context.Context, filter *dto.ServiceFilter,
) ([]*entities.Service, errors.Error) {
	query := `
		SELECT id, name, version, status, created_at
		FROM service
		ORDER BY created_at DESC;
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
		return nil, r.errorMapper(err, r.tableName)
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
			return nil, r.errorMapper(err, r.tableName)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	return services, nil
}

func (r *ServiceRepository) Create(
	ctx context.Context, service *entities.Service,
) errors.Error {
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

	return r.errorMapper(err, r.tableName)
}

func NewServiceRepository(driver *Driver) *ServiceRepository {
	return &ServiceRepository{Driver: driver, tableName: "service"}
}
