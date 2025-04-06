package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *APIKeyRepository) UpdateStatus(
	ctx context.Context, id int, status enums.APIKeyStatus,
) *errors.Error {
	if status == enums.APIKeyStatusNull {
		return errors.ErrAPIKeyInvalidStatus
	}

	query := `
		UPDATE api_key
		SET status = $1
		WHERE id = $2;
	`

	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrAPIKeyNotFound
	}

	return nil
}

func (r *APIKeyRepository) Update(
	ctx context.Context, id int, update *dto.APIKeyUpdate,
) *errors.Error {
	if update == nil {
		return nil
	}

	var updates []string
	args := []any{id}
	argIndex := 2

	if !update.ExpiresAt.IsZero() {
		updates = append(updates, fmt.Sprintf("expires_at = $%d", argIndex))
		args = append(args, update.ExpiresAt)
		argIndex++
	}

	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		`
			UPDATE api_key
			SET %s
			WHERE id = $1;
		`,
		strings.Join(updates, ", "),
	)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrAPIKeyNotFound
	}

	return nil
}

func (r *APIKeyRepository) UpdateLastUsed(
	ctx context.Context, key string,
) *errors.Error {
	query := `
		UPDATE api_key
		SET last_used = NOW()
		WHERE key = $1;
	`

	result, err := r.pool.Exec(ctx, query, key)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrAPIKeyNotFound
	}

	return nil
}

func (r *APIKeyRepository) FindByEnvironment(
	ctx context.Context, environmentID int,
) ([]*entities.APIKey, *errors.Error) {
	query := `
		SELECT id, environment_id, key, status, created_at,
			COALESCE(expires_at, '0001-01-01 00:00:00.0+00'),
			COALESCE(last_used, '0001-01-01 00:00:00.0+00')
		FROM api_key
		WHERE environment_id = $1;
	`

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
			&apiKey.Status,
			&apiKey.CreatedAt,
			&apiKey.ExpiresAt,
			&apiKey.LastUsed,
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
	query := `
		SELECT id, environment_id, key, status, created_at,
			COALESCE(expires_at, '0001-01-01 00:00:00.0+00'),
			COALESCE(last_used, '0001-01-01 00:00:00.0+00')
		FROM api_key
		WHERE key = $1;
	`

	apiKey := new(entities.APIKey)
	err := r.pool.QueryRow(ctx, query, key).Scan(
		&apiKey.ID,
		&apiKey.EnvironmentID,
		&apiKey.Key,
		&apiKey.Status,
		&apiKey.CreatedAt,
		&apiKey.ExpiresAt,
		&apiKey.LastUsed,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return apiKey, nil
}

func (r *APIKeyRepository) FindByID(
	ctx context.Context, id int,
) (*entities.APIKey, *errors.Error) {
	query := `
		SELECT id, environment_id, key, status, created_at,
			COALESCE(expires_at, '0001-01-01 00:00:00.0+00'),
			COALESCE(last_used, '0001-01-01 00:00:00.0+00')
		FROM api_key
		WHERE id = $1;
	`

	apiKey := new(entities.APIKey)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&apiKey.ID,
		&apiKey.EnvironmentID,
		&apiKey.Key,
		&apiKey.Status,
		&apiKey.CreatedAt,
		&apiKey.ExpiresAt,
		&apiKey.LastUsed,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return apiKey, nil
}

func (r *APIKeyRepository) Exists(
	ctx context.Context, key string,
) (bool, *errors.Error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM api_key
			WHERE key = $1
		);
	`

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

	var expiresAt any
	if !apiKey.ExpiresAt.IsZero() {
		expiresAt = apiKey.ExpiresAt
	}

	var lastUsed any
	if !apiKey.LastUsed.IsZero() {
		lastUsed = apiKey.LastUsed
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		apiKey.EnvironmentID,
		apiKey.Key,
		expiresAt,
		lastUsed,
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
