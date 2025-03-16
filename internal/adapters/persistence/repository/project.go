package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ProjectRepository) FindByClient(
	ctx context.Context, clientID int,
) ([]*entities.Project, *errors.Error) {
	query := "SELECT * FROM project WHERE client_id = $1;"
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
			&project.ClientID,
			&project.Name,
			&project.Status,
			&project.CreatedAt,
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
