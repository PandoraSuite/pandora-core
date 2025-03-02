package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectServiceRepository struct {
	pool *pgxpool.Pool
}

func (r *ProjectServiceRepository) FindByProjectAndService(
	ctx context.Context, projectID, serviceID int,
) (*entities.ProjectService, error) {
	query := `
		SELECT *
		FROM project_service
		WHERE project_id = $1 AND service_id = $2;
	`

	var projectService models.ProjectService
	err := r.pool.QueryRow(ctx, query, projectID, serviceID).Scan(
		&projectService.ProjectID,
		&projectService.ServiceID,
		&projectService.MaxRequest,
		&projectService.ResetFrequency,
		&projectService.NextReset,
		&projectService.CreatedAt,
	)
	if err != nil {
		return nil, persistence.ConvertPgxError(err)
	}

	return projectService.ToEntity()
}

func (r *ProjectServiceRepository) Save(
	ctx context.Context, projectService *entities.ProjectService,
) (*entities.ProjectService, error) {
	model := models.ProjectServiceFromEntity(projectService)
	if err := r.save(ctx, model); err != nil {
		return nil, err
	}

	return model.ToEntity()
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
		return persistence.ConvertPgxError(err)
	}

	return nil
}

func NewProjectServiceRepository(pool *pgxpool.Pool) *ProjectServiceRepository {
	return &ProjectServiceRepository{pool: pool}
}
