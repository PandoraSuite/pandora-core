package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvironmentServiceRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *EnvironmentServiceRepository) FindByProjectAndService(
	ctx context.Context, projectID, serviceID int,
) ([]*entities.EnvironmentService, *errors.Error) {
	query := `
		SELECT *
		FROM environment_service
		WHERE environment_id in (
			SELECT id
			FROM environment
			WHERE project_id = $1 AND status = 'active' ;
		) AND service_id = $2;
	`

	rows, err := r.pool.Query(ctx, query, projectID, serviceID)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var environmentServices []*entities.EnvironmentService
	for rows.Next() {
		environmentService := new(entities.EnvironmentService)

		err = rows.Scan(
			&environmentService.EnvironmentID,
			&environmentService.ServiceID,
			&environmentService.MaxRequest,
			&environmentService.AvailableRequest,
			&environmentService.CreatedAt,
		)
		if err != nil {
			return nil, r.handlerErr(err)
		}

		environmentServices = append(environmentServices, environmentService)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return environmentServices, nil
}

func (r *EnvironmentServiceRepository) DecrementAvailableRequest(
	ctx context.Context, environmentID, serviceID int,
) (*entities.EnvironmentService, *errors.Error) {
	query := `
		UPDATE environment_service
		SET available_request =
			CASE
				WHEN available_request IS NOT NULL AND available_request > 0
				THEN available_request - 1
				ELSE available_request
			END
		WHERE environment_id = $1 AND service_id = $2
		AND (available_request IS NULL OR available_request > 0)
		RETURNING *;
	`

	environmentService := new(entities.EnvironmentService)
	err := r.pool.QueryRow(ctx, query, environmentID, serviceID).
		Scan(
			&environmentService.EnvironmentID,
			&environmentService.ServiceID,
			&environmentService.MaxRequest,
			&environmentService.AvailableRequest,
			&environmentService.CreatedAt,
		)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return environmentService, nil
}

func (r *EnvironmentServiceRepository) Save(
	ctx context.Context, environmentService *entities.EnvironmentService,
) *errors.Error {
	query := `
		INSERT INTO environment_service (environment_id, service_id, max_request, available_request)
		VALUES ($1, $2, $3, $4) RETURNING created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		environmentService.EnvironmentID,
		environmentService.ServiceID,
		environmentService.MaxRequest,
		environmentService.AvailableRequest,
	).Scan(&environmentService.CreatedAt)

	return r.handlerErr(err)
}

func NewEnvironmentServiceRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *EnvironmentServiceRepository {
	return &EnvironmentServiceRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
