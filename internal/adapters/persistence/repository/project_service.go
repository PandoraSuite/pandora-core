package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectServiceRepository struct {
	pool *pgxpool.Pool
}

func (r *ProjectServiceRepository) Save(
	ctx context.Context, priojectService *entities.ProjectService,
) (*entities.ProjectService, error) {
	ps := models.ProjectServiceFromEntity(priojectService)
	if err := r.save(ctx, ps); err != nil {
		return nil, err
	}
	return ps.ToEntity(), nil
}

func (r *ProjectServiceRepository) save(
	ctx context.Context, projectService *models.ProjectService,
) error {
	if err := projectService.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO project_service (project_id, service_id, max_request, reset_frequency, next_reset)
		VALUES ($1, $2, $3, $4, $5) RETURNING created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		projectService.ProjectID,
		projectService.ServiceID,
		projectService.MaxRequest,
		projectService.ResetFrequency,
		projectService.NextReset,
	).Scan(&projectService.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the project service: %w", err)
	}

	return nil
}
