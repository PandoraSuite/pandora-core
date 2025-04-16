package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ClientRepository) Update(
	ctx context.Context, id int, update *dto.ClientUpdate,
) (*entities.Client, *errors.Error) {
	if update == nil {
		return r.FindByID(ctx, id)
	}

	var updates []string
	args := []any{id}
	argIndex := 2

	if update.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, update.Name)
		argIndex++
	}

	if update.Email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, update.Email)
		argIndex++
	}

	if update.Type != enums.ClientTypeNull {
		updates = append(updates, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, update.Type)
		argIndex++
	}

	if len(updates) == 0 {
		return r.FindByID(ctx, id)
	}

	query := fmt.Sprintf(
		`
			UPDATE client
			SET %s
			WHERE id = $1
			RETURNING id, type, name, email, created_at;
		`,
		strings.Join(updates, ", "),
	)

	client := new(entities.Client)
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&client.ID,
		&client.Type,
		&client.Name,
		&client.Email,
		&client.CreatedAt,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return client, nil
}

func (r *ClientRepository) Exists(
	ctx context.Context, id int,
) (bool, *errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM client
			WHERE id = $1
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, r.handlerErr(err)
	}

	return exists, nil
}

func (r *ClientRepository) FindByID(
	ctx context.Context, id int,
) (*entities.Client, *errors.Error) {
	query := `
		SELECT id, type, name, email, created_at
		FROM client
		WHERE id = $1;
	`

	client := new(entities.Client)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&client.ID,
		&client.Type,
		&client.Name,
		&client.Email,
		&client.CreatedAt,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return client, nil
}

func (r *ClientRepository) FindAll(
	ctx context.Context, filter *dto.ClientFilter,
) ([]*entities.Client, *errors.Error) {
	query := `
		SELECT id, type, name, email, created_at
		FROM client
	`

	var args []any
	if filter != nil {
		var where []string
		argIndex := 1

		if filter.Type != enums.ClientTypeNull {
			where = append(where, fmt.Sprintf("type = $%d", argIndex))
			args = append(args, filter.Type)
			argIndex++
		}

		if len(where) > 0 {
			query = fmt.Sprintf(
				"%s WHERE %s", query, strings.Join(where, " AND "),
			)
		}
	}

	query += ";"

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
