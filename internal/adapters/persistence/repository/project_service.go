package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectServiceRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) error
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
		return nil, r.handlerErr(err)
	}

	return projectService.ToEntity()
}

func (r *ProjectServiceRepository) BulkSave(
	ctx context.Context, projectServices []*entities.ProjectService,
) error {
	if len(projectServices) == 0 {
		return nil
	}

	values := []string{}
	args := []any{}
	argIndex := 1

	for _, projectService := range projectServices {
		values = append(
			values,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d)",
				argIndex,
				argIndex+1,
				argIndex+2,
				argIndex+3,
				argIndex+4,
			),
		)
		args = append(
			args,
			projectService.ProjectID,
			projectService.ServiceID,
			projectService.MaxRequest,
			projectService.ResetFrequency,
			projectService.NextReset,
		)
		argIndex += 5
	}

	query := fmt.Sprintf(`
			INSERT INTO project_service (project_id, service_id, max_request, reset_frequency, next_reset)
			VALUES %s
		`,
		strings.Join(values, ", "),
	)

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

func (r *ProjectServiceRepository) Save(
	ctx context.Context, projectService *entities.ProjectService,
) error {
	model := models.ProjectServiceFromEntity(projectService)
	if err := r.save(ctx, &model); err != nil {
		return err
	}

	projectService.CreatedAt = model.EntityCreatedAt()
	return nil
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
		return r.handlerErr(err)
	}

	return nil
}

func NewProjectServiceRepository(
	pool *pgxpool.Pool, handlerErr func(error) error,
) *ProjectServiceRepository {
	return &ProjectServiceRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
