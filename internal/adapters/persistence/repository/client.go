package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func (r *ClientRepository) FindAll(
	ctx context.Context, clientType enums.ClientType,
) ([]*entities.Client, error) {
	query := "SELECT * FROM client"

	var args []any
	if clientType != enums.ClientTypeNull {
		query += " WHERE type = $1;"
		args = append(args, clientType)
	} else {
		query += ";"
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var clients []*models.Client
	for rows.Next() {
		client := new(models.Client)

		err = rows.Scan(
			&client.ID,
			&client.Type,
			&client.Name,
			&client.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return models.ClientsToEntity(clients)
}

func (r *ClientRepository) Save(
	ctx context.Context, client *entities.Client,
) (*entities.Client, error) {
	model := models.ClientFromEntity(client)
	if err := r.save(ctx, model); err != nil {
		return nil, err
	}

	return model.ToEntity()
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
