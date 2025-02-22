package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	pool *pgxpool.Pool
}

func (r *ServiceRepository) Create(
	ctx context.Context, newService *models.Service,
) error {
	query := `
		INSERT INTO service (name, version, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		newService.Name,
		newService.Version,
		newService.Status,
	).Scan(&newService.ID, &newService.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the service: %w", err)
	}

	return nil
}
