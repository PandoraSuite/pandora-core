package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
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

func (r *EnvironmentRepository) FindByID(
	ctx context.Context, id int,
) (*dto.EnvironmentResponse, *errors.Error) {
	query := `
		SELECT e.id, e.name, e.status, e.project_id, e.createdAt,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
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

	environment := new(dto.EnvironmentResponse)
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
) ([]*dto.EnvironmentResponse, *errors.Error) {
	query := `
		SELECT e.id, e.name, e.status, e.project_id, e.createdAt,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
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

	var environments []*dto.EnvironmentResponse
	for rows.Next() {
		environment := new(dto.EnvironmentResponse)

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

func (r *EnvironmentRepository) Save(
	ctx context.Context, environment *entities.Environment,
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

func NewEnvironmentRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *EnvironmentRepository {
	return &EnvironmentRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
