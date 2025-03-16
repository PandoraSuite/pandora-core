package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ClientRepository) FindAll(
	ctx context.Context, filter *dto.ClientFilter,
) ([]*entities.Client, *errors.Error) {
	query := "SELECT * FROM client"

	var args []any
	if filter.Type != enums.ClientTypeNull {
		query += " WHERE type = $1;"
		args = append(args, filter.Type)
	} else {
		query += ";"
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var clients []*entities.Client
	for rows.Next() {
		client := new(entities.Client)

		err = rows.Scan(
			&client.ID,
			&client.Type,
			&client.Name,
			&client.Email,
			&client.CreatedAt,
		)
		if err != nil {
			return nil, r.handlerErr(err)
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return clients, nil
}

func (r *ClientRepository) Save(
	ctx context.Context, client *entities.Client,
) *errors.Error {
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

	return r.handlerErr(err)
}

func NewClientRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *ClientRepository {
	return &ClientRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
