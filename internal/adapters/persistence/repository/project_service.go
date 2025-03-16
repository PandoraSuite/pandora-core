package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectServiceRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ProjectServiceRepository) FindByProjectAndService(
	ctx context.Context, projectID, serviceID int,
) (*entities.ProjectService, *errors.Error) {
	query := `
		SELECT *
		FROM project_service
		WHERE project_id = $1 AND service_id = $2;
	`

	projectService := new(entities.ProjectService)
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

	return projectService, nil
}

func (r *ProjectServiceRepository) BulkSave(
	ctx context.Context, projectServices []*entities.ProjectService,
) *errors.Error {
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
	return r.handlerErr(err)
}

func (r *ProjectServiceRepository) Save(
	ctx context.Context, projectService *entities.ProjectService,
) *errors.Error {
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

	return r.handlerErr(err)
}

func NewProjectServiceRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *ProjectServiceRepository {
	return &ProjectServiceRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
