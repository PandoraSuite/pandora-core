package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectRepository struct {
	*Driver

	tableName           string
	auxServiceTableName string
}

func (r *ProjectRepository) Delete(
	ctx context.Context, id int,
) errors.Error {
	query := `
		DELETE FROM project
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

func (r *ProjectRepository) GetProjectClientInfoByID(
	ctx context.Context, id int,
) (*dto.ProjectClientInfoResponse, errors.Error) {
	query := `
		SELECT p.id, p.name, c.id, c.name
		FROM project p
			JOIN client c ON c.id = p.client_id
		WHERE p.id = $1;
	`

	projectCxt := new(dto.ProjectClientInfoResponse)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&projectCxt.ProjectID,
		&projectCxt.ProjectName,
		&projectCxt.ClientID,
		&projectCxt.ClientName,
	)
	return projectCxt, r.errorMapper(err, r.auxServiceTableName)
}

func (r *ProjectRepository) ResetAvailableRequestsForEnvsService(
	ctx context.Context, id, serviceID int,
) ([]*dto.EnvironmentServiceReset, errors.Error) {
	return r.resetAvailableRequestsForEnvsService(
		ctx, nil, id, serviceID,
	)
}

func (r *ProjectRepository) ResetProjectServiceUsage(
	ctx context.Context, id, serviceID int, nextReset time.Time,
) ([]*dto.EnvironmentServiceReset, errors.Error) {
	tx, txErr := r.pool.Begin(ctx)
	if txErr != nil {
		return nil, r.errorMapper(txErr, r.tableName)
	}

	if err := r.updateNextReset(ctx, tx, id, serviceID, nextReset); err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	environmentsService, err := r.resetAvailableRequestsForEnvsService(
		ctx, tx, id, serviceID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	return environmentsService, r.errorMapper(tx.Commit(ctx), r.tableName)
}

func (r *ProjectRepository) updateNextReset(
	ctx context.Context, tx pgx.Tx, id, serviceID int, nextReset time.Time,
) errors.Error {
	query := `
		UPDATE project_service
		SET next_reset = $3
		WHERE project_id = $1 AND service_id = $2;
	`

	var internalNextReset any
	if !nextReset.IsZero() {
		internalNextReset = nextReset
	}

	result, err := tx.Exec(ctx, query, id, serviceID, internalNextReset)
	if err != nil {
		return r.errorMapper(err, r.auxServiceTableName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(
			r.auxServiceTableName,
			map[string]any{"project_id": id, "service_id": serviceID},
		)
	}

	return nil
}

func (r *ProjectRepository) resetAvailableRequestsForEnvsService(
	ctx context.Context, tx pgx.Tx, id, serviceID int,
) ([]*dto.EnvironmentServiceReset, errors.Error) {
	query := `
		WITH updated AS (
			UPDATE environment_service es
			SET available_request = max_requests
			FROM project p
				JOIN environment e ON e.project_id = p.id
			WHERE es.environment_id = e.id AND p.id = $1 AND es.service_id = $2
			RETURNING e.id, e.name, e.status, es.max_requests,
				es.available_request, es.created_at, es.service_id
		)
		SELECT u.id, u.name, u.status, JSON_BUILD_OBJECT(
			'id', s.id,
			'name', s.name,
			'version', s.version,
			'max_requests', COALESCE(u.max_requests, -1),
			'available_request', COALESCE(u.available_request, -1),
			'assigned_at', u.created_at
		)
		FROM updated u
			JOIN service s ON s.id = u.service_id;
	`

	var err error
	var rows pgx.Rows
	if tx != nil {
		rows, err = tx.Query(ctx, query, id, serviceID)
	} else {
		rows, err = r.pool.Query(ctx, query, id, serviceID)
	}

	if err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	defer rows.Close()

	var environmentsService []*dto.EnvironmentServiceReset
	for rows.Next() {
		environmentService := new(dto.EnvironmentServiceReset)

		err = rows.Scan(
			&environmentService.ID,
			&environmentService.Name,
			&environmentService.Status,
			&environmentService.Service,
		)
		if err != nil {
			return nil, r.errorMapper(err, r.tableName)
		}

		environmentsService = append(environmentsService, environmentService)
	}

	if err := rows.Err(); err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	return environmentsService, nil
}

func (r *ProjectRepository) GetServiceByID(
	ctx context.Context, id, serviceID int,
) (*entities.ProjectService, errors.Error) {
	query := `
		SELECT s.id, s.name, s.version,
			COALESCE(ps.max_requests, -1),
			COALESCE(ps.reset_frequency, ''),
			COALESCE(ps.next_reset, '0001-01-01 00:00:00.0+00'),
			ps.created_at
		FROM project_service ps
			JOIN service s
				ON s.id = ps.service_id
		WHERE ps.project_id = $1 AND ps.service_id = $2;
	`

	service := new(entities.ProjectService)
	err := r.pool.QueryRow(ctx, query, id, serviceID).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.MaxRequests,
		&service.ResetFrequency,
		&service.NextReset,
		&service.AssignedAt,
	)
	if err != nil {
		return nil, r.errorMapper(err, r.auxServiceTableName)
	}

	return service, nil
}

func (r *ProjectRepository) UpdateService(
	ctx context.Context, id, serviceID int, update *dto.ProjectServiceUpdate,
) (*entities.ProjectService, errors.Error) {
	if update == nil {
		return r.GetServiceByID(ctx, id, serviceID)
	}

	query := `
		WITH updated AS (
			UPDATE project_service
			SET max_requests = $3, reset_frequency = $4, next_reset = $5
			WHERE project_id = $1 AND service_id = $2
			RETURNING *
		)
		SELECT s.id, s.name, s.version,
			COALESCE(u.max_requests, -1),
			COALESCE(u.reset_frequency, ''),
			COALESCE(u.next_reset, '0001-01-01 00:00:00.0+00'),
			u.created_at
		FROM updated u
			JOIN service s ON s.id = u.service_id;
	`

	var resetFrequency any
	if update.ResetFrequency != enums.ProjectServiceResetFrequencyNull {
		resetFrequency = update.ResetFrequency
	}

	var maxRequests any
	if update.MaxRequests != -1 {
		maxRequests = update.MaxRequests
	}

	var nextReset any
	if !update.NextReset.IsZero() {
		nextReset = update.NextReset
	}

	service := new(entities.ProjectService)
	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		serviceID,
		maxRequests,
		resetFrequency,
		nextReset,
	).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.MaxRequests,
		&service.ResetFrequency,
		&service.NextReset,
		&service.AssignedAt,
	)
	if err != nil {
		return nil, r.errorMapper(err, r.auxServiceTableName)
	}

	return service, nil
}

func (r *ProjectRepository) ExistsServiceIn(
	ctx context.Context, serviceID int,
) (bool, errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM project_service
			WHERE service_id = $1
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, serviceID).Scan(&exists)
	if err != nil {
		return false, r.errorMapper(err, r.auxServiceTableName)
	}

	return exists, nil
}

func (r *ProjectRepository) RemoveService(
	ctx context.Context, id, serviceID int,
) (int64, errors.Error) {
	query := `
		DELETE FROM project_service
		WHERE project_id = $1 AND service_id = $2;
	`

	result, err := r.pool.Exec(ctx, query, id, serviceID)
	if err != nil {
		return 0, r.errorMapper(err, r.auxServiceTableName)
	}

	return result.RowsAffected(), nil
}

func (r *ProjectRepository) UpdateStatus(
	ctx context.Context, id int, status enums.ProjectStatus,
) errors.Error {
	query := `
		UPDATE project
		SET status = $1
		WHERE id = $2;
	`

	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return r.errorMapper(err, r.tableName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(r.tableName, map[string]any{"id": id})
	}

	return nil
}

func (r *ProjectRepository) Update(
	ctx context.Context, id int, update *dto.ProjectUpdate,
) (*entities.Project, errors.Error) {
	if update == nil {
		return r.GetByID(ctx, id)
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
		return r.GetByID(ctx, id)
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
		return nil, r.errorMapper(err, r.tableName)
	}

	if result.RowsAffected() == 0 {
		return nil, r.entityNotFoundError(r.tableName, map[string]any{"id": id})
	}

	return r.GetByID(ctx, id)
}

func (r *ProjectRepository) Exists(
	ctx context.Context, id int,
) (bool, errors.Error) {
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
		return false, r.errorMapper(err, r.tableName)
	}

	return exists, nil
}

func (r *ProjectRepository) GetProjectServiceQuotaUsage(
	ctx context.Context, id, serviceID int,
) (*dto.QuotaUsage, errors.Error) {
	query := `
		SELECT COALESCE(ps.max_requests, -1), COALESCE(SUM(es.max_requests), 0)
		FROM project_service ps
			LEFT JOIN environment e
				ON e.project_id = ps.project_id
			LEFT JOIN environment_service es
				ON es.environment_id = e.id
				AND es.service_id = ps.service_id
		WHERE ps.project_id = $1 AND ps.service_id = $2
		GROUP BY ps.max_requests;
	`

	quota := new(dto.QuotaUsage)
	err := r.pool.QueryRow(ctx, query, id, serviceID).Scan(
		&quota.MaxAllowed,
		&quota.CurrentAllocated,
	)
	return quota, r.errorMapper(err, r.auxServiceTableName)
}

func (r *ProjectRepository) GetByID(
	ctx context.Context, id int,
) (*entities.Project, errors.Error) {
	query := `
		SELECT p.id, p.name, p.status, p.client_id, p.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'nextReset', COALESCE(ps.next_reset, '0001-01-01 00:00:00.0+00'),
						'maxRequests', COALESCE(ps.max_requests, -1),
						'resetFrequency', COALESCE(ps.reset_frequency, ''),
						'assignedAt', ps.created_at
					)
				) FILTER (WHERE s.id IS NOT NULL), '[]'
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
		return nil, r.errorMapper(err, r.tableName)
	}

	return project, nil
}

func (r *ProjectRepository) List(ctx context.Context) ([]*entities.Project, errors.Error) {
	query := `
		SELECT p.id, p.name, p.status, p.client_id, p.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'nextReset', COALESCE(ps.next_reset, '0001-01-01 00:00:00.0+00'),
						'maxRequests', COALESCE(ps.max_requests, -1),
						'resetFrequency', COALESCE(ps.reset_frequency, ''),
						'assignedAt', ps.created_at
					)
				) FILTER (WHERE s.id IS NOT NULL), '[]'
			)
		FROM project p
			LEFT JOIN project_service ps
				ON ps.project_id = p.id
			LEFT JOIN service s
				ON s.id = ps.service_id
		GROUP BY p.id;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, r.errorMapper(err, r.tableName)
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
			return nil, r.errorMapper(err, r.tableName)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	return projects, nil
}

func (r *ProjectRepository) ListByClient(
	ctx context.Context, clientID int,
) ([]*entities.Project, errors.Error) {
	query := `
		SELECT p.id, p.name, p.status, p.client_id, p.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'nextReset', COALESCE(ps.next_reset, '0001-01-01 00:00:00.0+00'),
						'maxRequests', COALESCE(ps.max_requests, -1),
						'resetFrequency', ps.reset_frequency,
						'assignedAt', ps.created_at
					)
				) FILTER (WHERE s.id IS NOT NULL), '[]'
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
		return nil, r.errorMapper(err, r.tableName)
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
			return nil, r.errorMapper(err, r.tableName)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	return projects, nil
}

func (r *ProjectRepository) AddService(
	ctx context.Context, id int, service *entities.ProjectService,
) errors.Error {
	query := `
		WITH inserted AS (
			INSERT INTO project_service (project_id, service_id, max_requests, reset_frequency, next_reset)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING service_id, created_at
		)
		SELECT s.name, s.version, i.created_at
		FROM inserted i
			JOIN service s
				ON i.service_id = s.id;
	`

	var resetFrequency any
	if service.ResetFrequency != enums.ProjectServiceResetFrequencyNull {
		resetFrequency = service.ResetFrequency
	}

	var maxRequests any
	if service.MaxRequests != -1 {
		maxRequests = service.MaxRequests
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		service.ID,
		maxRequests,
		resetFrequency,
		service.NextReset,
	).Scan(&service.Name, &service.Version, &service.AssignedAt)

	return r.errorMapper(err, r.auxServiceTableName)
}

func (r *ProjectRepository) Create(
	ctx context.Context, project *entities.Project,
) errors.Error {
	tx, txErr := r.pool.Begin(ctx)
	if txErr != nil {
		return r.errorMapper(txErr, r.tableName)
	}

	if err := r.createProject(ctx, tx, project); err != nil {
		tx.Rollback(ctx)
		return err
	}

	services, err := r.createProjectServices(
		ctx, tx, project.ID, project.Services,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	project.Services = services
	return r.errorMapper(tx.Commit(ctx), r.tableName)
}

func (r *ProjectRepository) createProject(
	ctx context.Context, tx pgx.Tx, project *entities.Project,
) errors.Error {
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

	return r.errorMapper(err, r.tableName)
}

func (r *ProjectRepository) createProjectServices(
	ctx context.Context,
	tx pgx.Tx,
	projectID int,
	newServices []*entities.ProjectService,
) ([]*entities.ProjectService, errors.Error) {
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
		if service.ResetFrequency != enums.ProjectServiceResetFrequencyNull {
			resetFrequency = service.ResetFrequency
		}

		var maxRequests any
		if service.MaxRequests != -1 {
			maxRequests = service.MaxRequests
		}

		var nextReset any
		if !service.NextReset.IsZero() {
			nextReset = service.NextReset
		}

		args = append(
			args,
			projectID,
			service.ID,
			maxRequests,
			resetFrequency,
			nextReset,
		)
		argIndex += 5
	}

	query := fmt.Sprintf(
		`
			WITH inserted AS (
				INSERT INTO project_service (
					project_id, service_id, max_requests, reset_frequency, next_reset
				)
				VALUES %s
				RETURNING *
			)
			SELECT s.id, s.name, s.version, i.created_at,
				COALESCE(i.reset_frequency, ''),
				COALESCE(i.max_requests, -1),
				COALESCE(i.next_reset, '0001-01-01 00:00:00.0+00')
			FROM inserted i
				JOIN service s ON i.service_id = s.id;
		`,
		strings.Join(values, ", "),
	)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, r.errorMapper(err, r.auxServiceTableName)
	}

	defer rows.Close()

	var services []*entities.ProjectService
	for rows.Next() {
		service := new(entities.ProjectService)

		err = rows.Scan(
			&service.ID,
			&service.Name,
			&service.Version,
			&service.AssignedAt,
			&service.ResetFrequency,
			&service.MaxRequests,
			&service.NextReset,
		)
		if err != nil {
			return nil, r.errorMapper(err, r.auxServiceTableName)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, r.errorMapper(err, r.auxServiceTableName)
	}

	return services, r.errorMapper(err, r.auxServiceTableName)
}

func NewProjectRepository(driver *Driver) *ProjectRepository {
	return &ProjectRepository{
		Driver:              driver,
		tableName:           "project",
		auxServiceTableName: "project_service",
	}
}
