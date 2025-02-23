package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func (r *ProjectRepository) Save(
	ctx context.Context, prioject *entities.Project,
) (*entities.Project, error) {
	model := models.ProjectFromEntity(prioject)
	if err := r.save(ctx, model); err != nil {
		return nil, err
	}

	return model.ToEntity()
}

func (r *ProjectRepository) save(
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
