package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func (r *ClientRepository) Create(
	ctx context.Context, newClient *models.Client,
) error {
	if err := newClient.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO client (type, name, email)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		newClient.Type,
		newClient.Name,
		newClient.Email,
	).Scan(&newClient.ID, &newClient.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the client: %w", err)
	}

	return nil
}
