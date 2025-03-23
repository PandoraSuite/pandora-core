package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvironmentRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *EnvironmentRepository) UpdateStatus(
	ctx context.Context, id int, status enums.EnvironmentStatus,
) *errors.Error {
	if status == enums.EnvironmentStatusNull {
		return errors.ErrEnvironmentInvalidStatus
	}

	query := "UPDATE environment SET status = $1 WHERE id = $2;"
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrAPIKeyNotFound
	}

	return nil
}

func (r *EnvironmentRepository) Update(
	ctx context.Context, id int, update *dto.EnvironmentUpdate,
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
		"UPDATE environment SET %s WHERE id = $1;",
		strings.Join(updates, ", "),
	)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrAPIKeyNotFound
	}

	return nil
}

func (r *EnvironmentRepository) DecrementAvailableRequest(
	ctx context.Context, id, serviceID int,
) (*dto.DecrementAvailableRequest, *errors.Error) {
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
		RETURNING max_request, available_request;
	`

	result := new(dto.DecrementAvailableRequest)
	err := r.pool.QueryRow(ctx, query, id, serviceID).
		Scan(
			&result.MaxRequest,
			&result.AvailableRequest,
		)

	if err != nil {
		return nil, r.handlerErr(err)
	}

	return result, nil
}

func (r *EnvironmentRepository) ExistsEnvironmentService(
	ctx context.Context, id, serviceID int,
) (bool, *errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM environment_service
			WHERE environment_id = $1 AND service_id = $2;
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id, serviceID).Scan(&exists)
	if err != nil {
		return false, r.handlerErr(err)
	}

	return exists, nil
}

func (r *EnvironmentRepository) GetAllMaxRequestForServiceInEnvironments(
	ctx context.Context, id, serviceID int,
) ([]int, *errors.Error) {
	query := `
		SELECT es.max_request
		FROM environment_service es
		JOIN environment e ON e.id = es.environment_id
		WHERE e.project_id = (
			SELECT project_id FROM environment WHERE id = $1;
		)
		AND es.service_id = $2;
	`

	rows, err := r.pool.Query(ctx, query, id, serviceID)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var maxRequests []int
	for rows.Next() {
		var maxResuest int
		if err := rows.Scan(&maxResuest); err != nil {
			return nil, r.handlerErr(err)
		}

		maxRequests = append(maxRequests, maxResuest)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return maxRequests, nil
}

func (r *EnvironmentRepository) GetMaxRequestForServiceInProject(
	ctx context.Context, id, serviceID int,
) (int, *errors.Error) {
	query := `
		SELECT ps.max_request
		FROM environment e
		JOIN project_service ps ON ps.project_id = e.project_id
		WHERE e.id = $1 AND ps.service_id = $2;
	`

	var maxRequest int
	err := r.pool.QueryRow(ctx, query, id, serviceID).
		Scan(&maxRequest)
	if err != nil {
		return 0, r.handlerErr(err)
	}

	return maxRequest, nil
}

func (r *EnvironmentRepository) FindByID(
	ctx context.Context, id int,
) (*entities.Environment, *errors.Error) {
	query := `
		SELECT e.id, e.name, e.status, e.project_id, e.createdAt,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'max_request', es.max_request,
						'assigned_at', es.created_at
					)
				), '[]'
			) AS services
		FROM environment e
		LEFT JOIN environment_service es ON es.environment_id = e.id
		LEFT JOIN service s ON s.id = es.service_id
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
		return nil, r.handlerErr(err)
	}

	return environment, nil
}

func (r *EnvironmentRepository) FindByProject(
	ctx context.Context, projectID int,
) ([]*entities.Environment, *errors.Error) {
	query := `
		SELECT e.id, e.name, e.status, e.project_id, e.createdAt,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'version', s.version,
						'max_request', es.max_request,
						'assigned_at', es.created_at
					)
				), '[]'
			) AS services
		FROM environment e
		JOIN project p ON p.id = e.project_id
		LEFT JOIN environment_service es ON es.environment_id = e.id
		LEFT JOIN service s ON s.id = es.service_id
		WHERE p.id = $1
		GROUP BY e.id;
	`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, r.handlerErr(err)
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
			return nil, r.handlerErr(err)
		}

		environments = append(environments, environment)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return environments, nil
}

func (r *EnvironmentRepository) AddService(
	ctx context.Context, id int, service *entities.EnvironmentService,
) *errors.Error {
	query := `
		WITH inserted AS (
			INSERT INTO environment_service (environment_id, service_id, max_request, available_request)
			VALUES ($1, $2, $3, $4)
			RETURNING service_id
		)
		SELECT s.name, s.version
		FROM inserted i
		JOIN service s ON i.service_id = s.id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		service.ID,
		service.MaxRequest,
		service.AvailableRequest,
	).Scan(&service.Name, &service.Version)

	return r.handlerErr(err)
}

func (r *EnvironmentRepository) Save(
	ctx context.Context, environment *entities.Environment,
) *errors.Error {
	tx, txErr := r.pool.Begin(ctx)
	if txErr != nil {
		return r.handlerErr(txErr)
	}

	if err := r.saveEnvironment(ctx, tx, environment); err != nil {
		tx.Rollback(ctx)
		return err
	}

	services, err := r.saveEnvironmentServices(
		ctx, tx, environment.ID, environment.Services,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	environment.Services = services
	return r.handlerErr(tx.Commit(ctx))
}

func (r *EnvironmentRepository) saveEnvironment(
	ctx context.Context, tx pgx.Tx, environment *entities.Environment,
) *errors.Error {
	query := `
		INSERT INTO environment (project_id, name, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		environment.ProjectID,
		environment.Name,
		environment.Status,
	).Scan(&environment.ID, &environment.CreatedAt)

	return r.handlerErr(err)
}

func (r *EnvironmentRepository) saveEnvironmentServices(
	ctx context.Context,
	tx pgx.Tx,
	environmentID int,
	newServices []*entities.EnvironmentService,
) ([]*entities.EnvironmentService, *errors.Error) {
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
		args = append(
			args,
			environmentID,
			service.ID,
			service.MaxRequest,
			service.AvailableRequest,
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
			SELECT s.id, s.name, s.version, i.max_request, i.available_request, i.create_at
			FROM inserted i
			JOIN service s ON i.service_id = s.id
		`,
		strings.Join(values, ", "),
	)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, r.handlerErr(err)
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
			return nil, r.handlerErr(err)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return services, r.handlerErr(err)
}

func NewEnvironmentRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *EnvironmentRepository {
	return &EnvironmentRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
