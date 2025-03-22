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

type ProjectRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ProjectRepository) UpdateStatus(
	ctx context.Context, id int, status enums.ProjectStatus,
) *errors.Error {
	if status == enums.ProjectStatusNull {
		return errors.ErrProjectInvalidStatus
	}

	query := "UPDATE project SET status = $1 WHERE id = $2;"
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrAPIKeyNotFound
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
		"UPDATE project SET %s WHERE id = $1;",
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

func (r *ProjectRepository) FindByID(
	ctx context.Context, id int,
) (*dto.ProjectResponse, *errors.Error) {
	query := `
		SELECT p.id, p.name, p.status, p.client_id, p.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'next_reset', ps.next_reset,
						'max_request', ps.max_request,
						'reset_frequency', ps.reset_frequency,
						'assigned_at', ps.created_at
					)
				), '[]'
			) AS services
		FROM project p
		LEFT JOIN project_service ps ON ps.project_id = p.id
		LEFT JOIN service s ON s.id = ps.service_id
		WHERE p.id = $1
		GROUP BY p.id;
	`

	project := new(dto.ProjectResponse)
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
) ([]*dto.ProjectResponse, *errors.Error) {
	query := `
		SELECT p.id, p.name, p.status, p.client_id, p.created_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', s.id,
						'name', s.name,
						'next_reset', ps.next_reset,
						'max_request', ps.max_request,
						'reset_frequency', ps.reset_frequency,
						'assigned_at', ps.created_at
					)
				), '[]'
			) AS services
		FROM project p
		JOIN client c ON c.id = p.client_id
		LEFT JOIN project_service ps ON ps.project_id = p.id
		LEFT JOIN service s ON s.id = ps.service_id
		WHERE c.id = $1
		GROUP BY p.id;
	`

	rows, err := r.pool.Query(ctx, query, clientID)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var projects []*dto.ProjectResponse
	for rows.Next() {
		project := new(dto.ProjectResponse)

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

func (r *ProjectRepository) Save(
	ctx context.Context, project *entities.Project,
) *errors.Error {
	query := `
		INSERT INTO project (client_id, name, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		project.ClientID,
		project.Name,
		project.Status,
	).Scan(&project.ID, &project.CreatedAt)

	return r.handlerErr(err)
}

func NewProjectRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *ProjectRepository {
	return &ProjectRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
