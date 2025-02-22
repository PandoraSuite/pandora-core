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
	ctx context.Context, newProject *models.Project,
) error {
	if err := newProject.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO client (client_id, name, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		newProject.ClientID,
		newProject.Name,
		newProject.Status,
	).Scan(&newProject.ID, &newProject.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the project: %w", err)
	}

	return nil
}
