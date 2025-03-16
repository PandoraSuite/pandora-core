package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvironmentRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *EnvironmentRepository) FindByID(
	ctx context.Context, id int,
) (*entities.Environment, *errors.Error) {
	query := "SELECT * FROM environment WHERE id = $1;"

	environment := new(entities.Environment)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&environment.ID,
		&environment.ProjectID,
		&environment.Name,
		&environment.Status,
		&environment.CreatedAt,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return environment, nil
}

func (r *EnvironmentRepository) FindByProject(
	ctx context.Context, projectID int,
) ([]*entities.Environment, *errors.Error) {
	query := "SELECT * FROM environment WHERE project_id = $1;"

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
			&environment.ProjectID,
			&environment.Name,
			&environment.Status,
			&environment.CreatedAt,
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
