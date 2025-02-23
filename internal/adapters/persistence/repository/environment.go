package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvironmentRepository struct {
	pool *pgxpool.Pool
}

func (r *EnvironmentRepository) Save(
	ctx context.Context, environment *entities.Environment,
) (*entities.Environment, error) {
	model := models.EnvironmentFromEntity(environment)
	if err := r.save(ctx, model); err != nil {
		return nil, err
	}

	return model.ToEntity()
}

func (r *EnvironmentRepository) save(
	ctx context.Context, environment *models.Environment,
) error {
	if err := environment.ValidateModel(); err != nil {
		return err
	}

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

	if err != nil {
		return fmt.Errorf("error when inserting the environment: %w", err)
	}

	return nil
}
