package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvironmentServiceRepository struct {
	pool *pgxpool.Pool
}

func (r *EnvironmentServiceRepository) Save(
	ctx context.Context, environmentService *entities.EnvironmentService,
) (*entities.EnvironmentService, error) {
	es := models.EnvironmentServiceFromEntity(environmentService)
	if err := r.save(ctx, es); err != nil {
		return nil, err
	}
	return es.ToEntity(), nil
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
