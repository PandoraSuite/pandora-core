package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectServiceRepository struct {
	pool *pgxpool.Pool
}

func (r *ProjectServiceRepository) Create(
	ctx context.Context, newProjectService *models.ProjectService,
) error {
	if err := newProjectService.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO client (project_id, service_id, max_request, reset_frequency, next_reset)
		VALUES ($1, $2, $3, $4, $5) RETURNING created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		newProjectService.ProjectID,
		newProjectService.ServiceID,
		newProjectService.MaxRequest,
		newProjectService.ResetFrequency,
		newProjectService.NextReset,
	).Scan(&newProjectService.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the project service: %w", err)
	}

	return nil
}
