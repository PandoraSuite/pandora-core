package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func (r *ProjectRepository) Create(
	ctx context.Context, project *models.Project,
) error {
	if err := project.ValidateModel(); err != nil {
		return err
	}

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

	if err != nil {
		return fmt.Errorf("error when inserting the project: %w", err)
	}

	return nil
}
