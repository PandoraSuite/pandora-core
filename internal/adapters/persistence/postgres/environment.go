package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentRepository struct {
	*Driver

	tableName           string
	auxServiceTableName string
}

func (r *EnvironmentRepository) ExistsServiceWithInfiniteMaxRequest(
	ctx context.Context, projectID, serviceID int,
) (bool, errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM project_service ps
				JOIN environment e
					ON e.project_id = ps.project_id
				JOIN environment_service es
					ON es.environment_id = e.id
						AND es.service_id = ps.service_id
			WHERE ps.project_id = $1 AND ps.service_id = $2
				AND es.max_request IS NULL
		);
	`

	var hasInfinite bool
	err := r.pool.QueryRow(ctx, query, projectID, serviceID).Scan(&hasInfinite)
	if err != nil {
		return false, r.errorMapper(err, r.tableName)
	}

	return hasInfinite, nil
}

func (r *EnvironmentRepository) ResetAvailableRequests(
	ctx context.Context, id, serviceID int,
) (*entities.EnvironmentService, errors.Error) {
	query := `
		WITH updated AS (
			UPDATE environment_service
			SET available_request = max_request
			WHERE environment_id = $1 AND service_id = $2
			RETURNING *
		)
		SELECT s.id, s.name, s.version, COALESCE(u.max_request, -1), COALESCE(u.available_request, -1), u.created_at
		FROM updated u
			JOIN service s
				ON u.service_id = s.id;
	`

	service := new(entities.EnvironmentService)
	err := r.pool.QueryRow(ctx, query, id, serviceID).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.MaxRequest,
		&service.AvailableRequest,
		&service.AssignedAt,
	)
	return service, r.errorMapper(err, r.auxServiceTableName)
}

func (r *EnvironmentRepository) RemoveService(
	ctx context.Context, id, serviceID int,
) (int64, errors.Error) {
	query := `
		DELETE FROM environment_service
		WHERE environment_id = $1 AND service_id = $2;
	`

	result, err := r.pool.Exec(ctx, query, id, serviceID)
	if err != nil {
		return 0, r.errorMapper(err, r.auxServiceTableName)
	}

	return result.RowsAffected(), nil
}

func (r *EnvironmentRepository) RemoveServiceFromProjectEnvironments(
	ctx context.Context, projectID, serviceID int,
) (int64, errors.Error) {
	query := `
		DELETE FROM environment_service
		WHERE service_id = $2
			AND environment_id IN (
				SELECT id
				FROM environment
				WHERE project_id = $1
			);
	`

	result, err := r.pool.Exec(ctx, query, projectID, serviceID)
	if err != nil {
		return 0, r.errorMapper(err, r.auxServiceTableName)
	}

	return result.RowsAffected(), nil
}

func (r *EnvironmentRepository) UpdateStatus(
	ctx context.Context, id int, status enums.EnvironmentStatus,
) errors.Error {
	query := `
		UPDATE environment
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

func (r *EnvironmentRepository) UpdateService(
	ctx context.Context,
	id, serviceID int,
	update *dto.EnvironmentServiceUpdate,
) (*entities.EnvironmentService, errors.Error) {
	if update == nil {
		return r.GetServiceByID(ctx, id, serviceID)
	}

	query := `
		WITH updated AS (
			UPDATE environment_service
			SET max_request = $3, available_request = $4
			WHERE environment_id = $1 AND service_id = $2
			RETURNING *
		)
		SELECT s.id, s.name, s.version, COALESCE(u.max_request, -1), COALESCE(u.available_request, -1), u.created_at
		FROM updated u
			JOIN service s
				ON s.id = u.service_id;
	`

	var maxRequest any
	if update.MaxRequest != -1 {
		maxRequest = update.MaxRequest
	}

	var availableRequest any
	if update.AvailableRequest != -1 {
		availableRequest = update.AvailableRequest
	}

	service := new(entities.EnvironmentService)
	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		serviceID,
		maxRequest,
		availableRequest,
	).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.MaxRequest,
		&service.AvailableRequest,
		&service.AssignedAt,
	)
	if err != nil {
		return nil, r.errorMapper(err, r.auxServiceTableName)
	}

	return service, nil
}

func (r *EnvironmentRepository) Update(
	ctx context.Context, id int, update *dto.EnvironmentUpdate,
) (*entities.Environment, errors.Error) {
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
			UPDATE environment
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

func (r *EnvironmentRepository) Exists(
	ctx context.Context, id int,
) (bool, errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM environment
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

func (r *EnvironmentRepository) IsActive(
	ctx context.Context, id int,
) (bool, errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM environment
			WHERE id = $1
			AND status = $2
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id, enums.EnvironmentActive).Scan(&exists)
	if err != nil {
		return false, r.errorMapper(err, r.tableName)
	}

	return exists, nil
}

func (r *EnvironmentRepository) MissingResourceDiagnosis(
	ctx context.Context, id int, serviceID int,
) (bool, bool, errors.Error) {
	query := `
		SELECT
		EXISTS (
			SELECT 1 FROM environment_service 
			WHERE environment_id = $1 AND service_id = $2
		) 
		AS environment_service_found,

		EXISTS (
			SELECT 1 FROM environment_service 
			WHERE environment_id = $1 AND service_id = $2
			AND (available_request IS NULL OR available_request > 0)
		) 
		AS has_available_requests;
	`

	var environment_service_found, has_available_requests bool
	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		serviceID,
	).Scan(
		&environment_service_found,
		&has_available_requests,
	)
	if err != nil {
		return false, false, r.errorMapper(err, r.auxServiceTableName)
	}

	return environment_service_found, has_available_requests, nil
}

func (r *EnvironmentRepository) GetProjectServiceQuotaUsage(
	ctx context.Context, id, serviceID int,
) (*dto.QuotaUsage, errors.Error) {
	query := `
		SELECT COALESCE(ps.max_request, -1), COALESCE(SUM(es.max_request), 0)
		FROM environment e_target
			JOIN project_service ps
				ON ps.project_id = e_target.project_id
				AND ps.service_id = $2
			LEFT JOIN environment e
				ON e.project_id = ps.project_id
			LEFT JOIN environment_service es
				ON es.environment_id = e.id
				AND es.service_id = ps.service_id
		WHERE e_target.id = $1
		GROUP BY ps.max_request;
	`

	quota := new(dto.QuotaUsage)
	err := r.pool.QueryRow(ctx, query, id, serviceID).Scan(
		&quota.MaxAllowed,
		&quota.CurrentAllocated,
	)
	return quota, r.errorMapper(err, r.tableName)
}

func (r *EnvironmentRepository) IncreaseAvailableRequest(
	ctx context.Context, id, serviceID int,
) errors.Error {
	query := `
		UPDATE environment_service
		SET available_request =
			CASE
				WHEN available_request IS NOT NULL AND available_request > 0
				THEN available_request + 1
				ELSE available_request
			END
		WHERE environment_id = $1 AND service_id = $2
			AND (available_request IS NULL OR available_request > 0)
		RETURNING COALESCE(max_request, -1), COALESCE(available_request, -1);
	`

	result, err := r.pool.Exec(
		ctx, query, id, serviceID)
	if err != nil {
		return r.errorMapper(err, r.auxServiceTableName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(
			r.auxServiceTableName,
			map[string]any{"environment_id": id, "service_id": serviceID},
		)
	}

	return nil
}

func (r *EnvironmentRepository) DecrementAvailableRequest(
	ctx context.Context, id, serviceID int,
) (*dto.DecrementAvailableRequest, errors.Error) {
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
		RETURNING COALESCE(max_request, -1), COALESCE(available_request, -1);
	`

	result := new(dto.DecrementAvailableRequest)
	err := r.pool.QueryRow(ctx, query, id, serviceID).
		Scan(
			&result.MaxRequest,
			&result.AvailableRequest,
		)

	if err != nil {
		return nil, r.errorMapper(err, r.auxServiceTableName)
	}

	return result, nil
}

func (r *EnvironmentRepository) ExistsServiceIn(
	ctx context.Context, id, serviceID int,
) (bool, errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM environment_service
			WHERE environment_id = $1 AND service_id = $2
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id, serviceID).Scan(&exists)
	if err != nil {
		return false, r.errorMapper(err, r.auxServiceTableName)
	}

	return exists, nil
}

func (r *EnvironmentRepository) GetServiceByID(
	ctx context.Context, id, serviceID int,
) (*entities.EnvironmentService, errors.Error) {
	query := `
		SELECT s.id, s.name, s.version, COALESCE(es.max_request, -1),
			COALESCE(es.available_request, -1), es.created_at
		FROM environment_service es
			JOIN service s
				ON s.id = es.service_id
		WHERE es.environment_id = $1 AND es.service_id = $2;
	`

	service := new(entities.EnvironmentService)
	err := r.pool.QueryRow(ctx, query, id, serviceID).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.MaxRequest,
		&service.AvailableRequest,
		&service.AssignedAt,
	)
	if err != nil {
		return nil, r.errorMapper(err, r.auxServiceTableName)
	}

	return service, nil
}

func (r *EnvironmentRepository) GetByID(
	ctx context.Context, id int,
) (*entities.Environment, errors.Error) {
	query := `
		SELECT e.id, e.name, e.status, e.project_id, e.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'maxRequest', COALESCE(es.max_request, -1),
						'availableRequest', COALESCE(es.available_request, -1),
						'assignedAt', es.created_at
					)
				) FILTER (WHERE s.id IS NOT NULL), '[]'
			)
		FROM environment e
			LEFT JOIN environment_service es
				ON es.environment_id = e.id
			LEFT JOIN service s
				ON s.id = es.service_id
		WHERE e.id = $1
		GROUP BY e.id;
	`

	environment := new(entities.Environment)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&environment.ID,
		&environment.Name,
		&environment.Status,
		&environment.ProjectID,
		&environment.CreatedAt,
		&environment.Services,
	)
	if err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	return environment, nil
}

func (r *EnvironmentRepository) ListByProject(
	ctx context.Context, projectID int,
) ([]*entities.Environment, errors.Error) {
	query := `
		SELECT e.id, e.name, e.status, e.project_id, e.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'maxRequest', COALESCE(es.max_request, -1),
						'availableRequest', COALESCE(es.available_request, -1),
						'assignedAt', es.created_at
					)
				) FILTER (WHERE s.id IS NOT NULL), '[]'
			)
		FROM environment e
			JOIN project p
				ON p.id = e.project_id
			LEFT JOIN environment_service es
				ON es.environment_id = e.id
			LEFT JOIN service s
				ON s.id = es.service_id
		WHERE p.id = $1
		GROUP BY e.id;
	`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	defer rows.Close()

	var environments []*entities.Environment
	for rows.Next() {
		environment := new(entities.Environment)

		err = rows.Scan(
			&environment.ID,
			&environment.Name,
			&environment.Status,
			&environment.ProjectID,
			&environment.CreatedAt,
			&environment.Services,
		)
		if err != nil {
			return nil, r.errorMapper(err, r.tableName)
		}

		environments = append(environments, environment)
	}

	if err := rows.Err(); err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	return environments, nil
}

func (r *EnvironmentRepository) AddService(
	ctx context.Context, id int, service *entities.EnvironmentService,
) errors.Error {
	query := `
		WITH inserted AS (
			INSERT INTO environment_service (environment_id, service_id, max_request, available_request)
			VALUES ($1, $2, $3, $4)
			RETURNING service_id, created_at
		)
		SELECT s.name, s.version, i.created_at
		FROM inserted i
			JOIN service s
				ON i.service_id = s.id;
	`

	var maxRequest any
	if service.MaxRequest != -1 {
		maxRequest = service.MaxRequest
	}

	var availableRequest any
	if service.AvailableRequest != -1 {
		availableRequest = service.AvailableRequest
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		service.ID,
		maxRequest,
		availableRequest,
	).Scan(&service.Name, &service.Version, &service.AssignedAt)

	return r.errorMapper(err, r.auxServiceTableName)
}

func (r *EnvironmentRepository) Create(
	ctx context.Context, environment *entities.Environment,
) errors.Error {
	tx, txErr := r.pool.Begin(ctx)
	if txErr != nil {
		return r.errorMapper(txErr, r.tableName)
	}

	if err := r.createEnvironment(ctx, tx, environment); err != nil {
		tx.Rollback(ctx)
		return err
	}

	services, err := r.createEnvironmentServices(
		ctx, tx, environment.ID, environment.Services,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	environment.Services = services
	return r.errorMapper(tx.Commit(ctx), r.tableName)
}

func (r *EnvironmentRepository) createEnvironment(
	ctx context.Context, tx pgx.Tx, environment *entities.Environment,
) errors.Error {
	query := `
		INSERT INTO environment (project_id, name, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := tx.QueryRow(
		ctx,
		query,
		environment.ProjectID,
		environment.Name,
		environment.Status,
	).Scan(&environment.ID, &environment.CreatedAt)

	return r.errorMapper(err, r.tableName)
}

func (r *EnvironmentRepository) createEnvironmentServices(
	ctx context.Context,
	tx pgx.Tx,
	environmentID int,
	newServices []*entities.EnvironmentService,
) ([]*entities.EnvironmentService, errors.Error) {
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
				"($%d, $%d, $%d, $%d)",
				argIndex,
				argIndex+1,
				argIndex+2,
				argIndex+3,
			),
		)

		var maxRequest any
		if service.MaxRequest != -1 {
			maxRequest = service.MaxRequest
		}

		var availableRequest any
		if service.AvailableRequest != -1 {
			availableRequest = service.AvailableRequest
		}

		args = append(
			args,
			environmentID,
			service.ID,
			maxRequest,
			availableRequest,
		)
		argIndex += 4
	}

	query := fmt.Sprintf(
		`
			WITH inserted AS (
				INSERT INTO environment_service (environment_id, service_id, max_request, available_request)
				VALUES %s
				RETURNING *
			)
			SELECT s.id, s.name, s.version, COALESCE(i.max_request, -1), COALESCE(i.available_request, -1), i.created_at
			FROM inserted i
				JOIN service s
					ON i.service_id = s.id;
		`,
		strings.Join(values, ", "),
	)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, r.errorMapper(err, r.auxServiceTableName)
	}

	defer rows.Close()

	var services []*entities.EnvironmentService
	for rows.Next() {
		service := new(entities.EnvironmentService)

		err = rows.Scan(
			&service.ID,
			&service.Name,
			&service.Version,
			&service.MaxRequest,
			&service.AvailableRequest,
			&service.AssignedAt,
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

func NewEnvironmentRepository(driver *Driver) *EnvironmentRepository {
	return &EnvironmentRepository{
		Driver:              driver,
		tableName:           "environment",
		auxServiceTableName: "environment_service",
	}
}
