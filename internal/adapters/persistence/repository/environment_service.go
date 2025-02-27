package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvironmentServiceRepository struct {
	pool *pgxpool.Pool
}

func (r *EnvironmentServiceRepository) DecrementAvailableRequest(
	ctx context.Context, environmentID, serviceID int,
) (*entities.EnvironmentService, error) {
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

	environmentService := new(models.EnvironmentService)
	err := r.pool.QueryRow(ctx, query, environmentID, serviceID).
		Scan(
			&environmentService.EnvironmentID,
			&environmentService.ServiceID,
			&environmentService.MaxRequest,
			&environmentService.AvailableRequest,
			&environmentService.CreatedAt,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("no more requests available")
		}
		return nil, err
	}

	return environmentService.ToEntity(), nil
}

func (r *EnvironmentServiceRepository) Save(
	ctx context.Context, environmentService *entities.EnvironmentService,
) (*entities.EnvironmentService, error) {
	model := models.EnvironmentServiceFromEntity(environmentService)
	if err := r.save(ctx, model); err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *EnvironmentServiceRepository) save(
	ctx context.Context, environmentService *models.EnvironmentService,
) error {
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

	if err != nil {
		return fmt.Errorf(
			"error when inserting the environment service: %w",
			err,
		)
	}

	return nil
}

func NewEnvironmentServiceRepository(
	pool *pgxpool.Pool,
) outbound.EnvironmentServiceRepositoryPort {
	return &EnvironmentServiceRepository{pool: pool}
}
