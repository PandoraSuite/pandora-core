package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *APIKeyRepository) FindByEnvironment(
	ctx context.Context, environmentID int,
) ([]*entities.APIKey, *errors.Error) {
	query := "SELECT * FROM api_key WHERE environment_id = $1;"
	rows, err := r.pool.Query(ctx, query, environmentID)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var apiKeys []*entities.APIKey
	for rows.Next() {
		apiKey := new(entities.APIKey)

		err = rows.Scan(
			&apiKey.ID,
			&apiKey.EnvironmentID,
			&apiKey.Key,
			&apiKey.ExpiresAt,
			&apiKey.LastUsed,
			&apiKey.Status,
			&apiKey.CreatedAt,
		)
		if err != nil {
			return nil, r.handlerErr(err)
		}

		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return apiKeys, nil
}

func (r *APIKeyRepository) FindByKey(
	ctx context.Context, key string,
) (*entities.APIKey, *errors.Error) {
	query := "SELECT * FROM api_key WHERE key = $1;"

	apiKey := new(entities.APIKey)
	err := r.pool.QueryRow(ctx, query, key).Scan(
		&apiKey.ID,
		&apiKey.EnvironmentID,
		&apiKey.Key,
		&apiKey.ExpiresAt,
		&apiKey.LastUsed,
		&apiKey.Status,
		&apiKey.CreatedAt,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return apiKey, nil
}

func (r *APIKeyRepository) Exists(ctx context.Context, key string) (bool, *errors.Error) {
	query := "SELECT EXISTS (SELECT 1 FROM api_key WHERE key = $1;);"

	var exists bool
	err := r.pool.QueryRow(ctx, query, key).Scan(&exists)
	if err != nil {
		return false, r.handlerErr(err)
	}

	return exists, nil
}

func (r *APIKeyRepository) Save(
	ctx context.Context, apiKey *entities.APIKey,
) *errors.Error {
	query := `
		INSERT INTO api_key (environment_id, key, expires_at, last_used, status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		apiKey.EnvironmentID,
		apiKey.Key,
		apiKey.ExpiresAt,
		apiKey.LastUsed,
		apiKey.Status,
	).Scan(&apiKey.ID, &apiKey.CreatedAt)

	return r.handlerErr(err)
}

func NewAPIKeyRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *APIKeyRepository {
	return &APIKeyRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
