package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ProjectRepository) RemoveServiceFromProject(
	ctx context.Context, id, serviceID int,
) (int64, *errors.Error) {
	query := `
		DELETE FROM project_service
		WHERE project_id = $1 AND service_id = $2;
	`

	result, err := r.pool.Exec(ctx, query, id, serviceID)
	if err != nil {
		return 0, r.handlerErr(err)
	}

	return result.RowsAffected(), nil
}

func (r *ProjectRepository) UpdateStatus(
	ctx context.Context, id int, status enums.ProjectStatus,
) *errors.Error {
	if status == enums.ProjectStatusNull {
		return errors.ErrProjectInvalidStatus
	}

	query := `
		UPDATE project
		SET status = $1
		WHERE id = $2;
	`

	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrProjectNotFound
	}

	return nil
}

func (r *ProjectRepository) Update(
	ctx context.Context, id int, update *dto.ProjectUpdate,
) *errors.Error {
	if update == nil {
		return nil
	}

	var updates []string
	args := []any{id}
	argIndex := 2

	if update.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, update.Name)
		argIndex++
	}

	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		`
			UPDATE project
			SET %s
			WHERE id = $1;
		`,
		strings.Join(updates, ", "),
	)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrProjectNotFound
	}

	return nil
}

func (r *ProjectRepository) Exists(
	ctx context.Context, id int,
) (bool, *errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM project
			WHERE id = $1
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, r.handlerErr(err)
	}

	return exists, nil
}

func (r *ProjectRepository) GetProjectServiceQuotaUsage(
	ctx context.Context, id, serviceID int,
) (*dto.QuotaUsage, *errors.Error) {
	query := `
		SELECT COALESCE(ps.max_request, -1), COALESCE(SUM(es.max_request), 0)
		FROM project_service ps
			LEFT JOIN environment e
				ON e.project_id = ps.project_id
			LEFT JOIN environment_service es
				ON es.environment_id = e.id
				AND es.service_id = ps.service_id
		WHERE ps.project_id = $1 AND ps.service_id = $2
		GROUP BY ps.max_request;
	`

	quota := new(dto.QuotaUsage)
	err := r.pool.QueryRow(ctx, query, id, serviceID).Scan(
		&quota.MaxAllowed,
		&quota.CurrentAllocated,
	)
	return quota, r.handlerErr(err)
}

func (r *ProjectRepository) FindByID(
	ctx context.Context, id int,
) (*entities.Project, *errors.Error) {
	query := `
		SELECT p.id, p.name, p.status, p.client_id, p.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'nextReset', ps.next_reset,
						'maxRequest', COALESCE(ps.max_request, -1),
						'resetFrequency', ps.reset_frequency,
						'assignedAt', ps.created_at
					)
				), '[]'
			)
		FROM project p
			LEFT JOIN project_service ps
				ON ps.project_id = p.id
			LEFT JOIN service s
				ON s.id = ps.service_id
		WHERE p.id = $1
		GROUP BY p.id;
	`

	project := new(entities.Project)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Status,
		&project.ClientID,
		&project.CreatedAt,
		&project.Services,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return project, nil
}

func (r *ProjectRepository) FindByClient(
	ctx context.Context, clientID int,
) ([]*entities.Project, *errors.Error) {
	query := `
		SELECT p.id, p.name, p.status, p.client_id, p.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'nextReset', ps.next_reset,
						'maxRequest', COALESCE(ps.max_request, -1),
						'resetFrequency', ps.reset_frequency,
						'assignedAt', ps.created_at
					)
				), '[]'
			)
		FROM project p
			JOIN client c
				ON c.id = p.client_id
			LEFT JOIN project_service ps
				ON ps.project_id = p.id
			LEFT JOIN service s
				ON s.id = ps.service_id
		WHERE c.id = $1
		GROUP BY p.id;
	`

	rows, err := r.pool.Query(ctx, query, clientID)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var projects []*entities.Project
	for rows.Next() {
		project := new(entities.Project)

		err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Status,
			&project.ClientID,
			&project.CreatedAt,
			&project.Services,
		)
		if err != nil {
			return nil, r.handlerErr(err)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return projects, nil
}

func (r *ProjectRepository) AddService(
	ctx context.Context, id int, service *entities.ProjectService,
) *errors.Error {
	query := `
		WITH inserted AS (
			INSERT INTO project_service (project_id, service_id, max_request, reset_frequency, next_reset)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING service_id
		)
		SELECT s.name, s.version
		FROM inserted i
			JOIN service s
				ON i.service_id = s.id;
	`

	var resetFrequency any
	if s := service.ResetFrequency.String(); s != "" {
		resetFrequency = s
	}

	var maxRequest any
	if service.MaxRequest != -1 {
		maxRequest = service.MaxRequest
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		service.ID,
		maxRequest,
		resetFrequency,
		service.NextReset,
	).Scan(&service.Name, &service.Version)

	return r.handlerErr(err)
}

func (r *ProjectRepository) Save(
	ctx context.Context, project *entities.Project,
) *errors.Error {
	tx, txErr := r.pool.Begin(ctx)
	if txErr != nil {
		return r.handlerErr(txErr)
	}

	if err := r.saveProject(ctx, tx, project); err != nil {
		tx.Rollback(ctx)
		return err
	}

	services, err := r.saveProjectServices(
		ctx, tx, project.ID, project.Services,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	project.Services = services
	return r.handlerErr(tx.Commit(ctx))
}

func (r *ProjectRepository) saveProject(
	ctx context.Context, tx pgx.Tx, project *entities.Project,
) *errors.Error {
	query := `
		INSERT INTO project (client_id, name, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := tx.QueryRow(
		ctx,
		query,
		project.ClientID,
		project.Name,
		project.Status,
	).Scan(&project.ID, &project.CreatedAt)

	return r.handlerErr(err)
}

func (r *ProjectRepository) saveProjectServices(
	ctx context.Context,
	tx pgx.Tx,
	projectID int,
	newServices []*entities.ProjectService,
) ([]*entities.ProjectService, *errors.Error) {
	if len(newServices) == 0 {
		return nil, nil
	}

	values := []string{}
	args := []any{}
	argIndex := 1

	for _, service := range newServices {
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

		var resetFrequency any
		if s := service.ResetFrequency.String(); s != "" {
			resetFrequency = s
		}

		var maxRequest any
		if service.MaxRequest != -1 {
			maxRequest = service.MaxRequest
		}

		args = append(
			args,
			projectID,
			service.ID,
			maxRequest,
			resetFrequency,
			service.NextReset,
		)
		argIndex += 5
	}

	query := fmt.Sprintf(
		`
			WITH inserted AS (
				INSERT INTO project_service (project_id, service_id, max_request, reset_frequency, next_reset)
				VALUES %s
				RETURNING *
			)
			SELECT s.id, s.name, s.version, COALESCE(i.max_request, -1), i.reset_frequency, i.next_reset, i.created_at
			FROM inserted i
				JOIN service s
					ON i.service_id = s.id;
		`,
		strings.Join(values, ", "),
	)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var services []*entities.ProjectService
	for rows.Next() {
		service := new(entities.ProjectService)

		err = rows.Scan(
			&service.ID,
			&service.Name,
			&service.Version,
			&service.MaxRequest,
			&service.ResetFrequency,
			&service.NextReset,
			&service.AssignedAt,
		)
		if err != nil {
			return nil, r.handlerErr(err)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return services, r.handlerErr(err)
}

func NewProjectRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *ProjectRepository {
	return &ProjectRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
