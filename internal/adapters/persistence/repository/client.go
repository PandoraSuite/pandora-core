package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func (r *ClientRepository) Save(
	ctx context.Context, client *entities.Client,
) (*entities.Client, error) {
	c := models.ClientFromEntity(client)
	if err := r.save(ctx, c); err != nil {
		return nil, err
	}
	return c.ToEntity(), nil
}

func (r *ClientRepository) save(ctx context.Context, client *models.Client) error {
	if err := client.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO client (type, name, email)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		client.Type,
		client.Name,
		client.Email,
	).Scan(&client.ID, &client.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the client: %w", err)
	}

	return nil
}

func NewClientRepository(pool *pgxpool.Pool) outbound.ClientRepositoryPort {
	return &ClientRepository{pool: pool}
}
