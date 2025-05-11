package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository struct {
	*Driver

	talbeName string
}

func (r *APIKeyRepository) UpdateStatus(
	ctx context.Context, id int, status enums.APIKeyStatus,
) errors.PersistenceError {
	query := `
		UPDATE api_key
		SET status = $1
		WHERE id = $2;
	`

	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return r.errorMapper(err, r.talbeName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(r.talbeName)
	}

	return nil
}

func (r *APIKeyRepository) Update(
	ctx context.Context, id int, update *dto.APIKeyUpdate,
) (*entities.APIKey, errors.PersistenceError) {
	if update == nil {
		return r.GetByID(ctx, id)
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
		return r.GetByID(ctx, id)
	}

	query := fmt.Sprintf(
		`
			UPDATE api_key
			SET %s
			WHERE id = $1
			RETURNING id, environment_id, key, status, created_at,
				COALESCE(expires_at, '0001-01-01 00:00:00.0+00'),
				COALESCE(last_used, '0001-01-01 00:00:00.0+00');
		`,
		strings.Join(updates, ", "),
	)

	apiKey := new(entities.APIKey)
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&apiKey.ID,
		&apiKey.EnvironmentID,
		&apiKey.Key,
		&apiKey.Status,
		&apiKey.CreatedAt,
		&apiKey.ExpiresAt,
		&apiKey.LastUsed,
	)
	if err != nil {
		return nil, r.errorMapper(err, r.talbeName)
	}

	return apiKey, nil
}

func (r *APIKeyRepository) UpdateLastUsed(
	ctx context.Context, key string,
) errors.PersistenceError {
	query := `
		UPDATE api_key
		SET last_used = NOW()
		WHERE key = $1;
	`

	result, err := r.pool.Exec(ctx, query, key)
	if err != nil {
		return r.errorMapper(err, r.talbeName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(r.talbeName)
	}

	return nil
}

func (r *APIKeyRepository) ListByEnvironment(
	ctx context.Context, environmentID int,
) ([]*entities.APIKey, errors.PersistenceError) {
	query := `
		SELECT id, environment_id, key, status, created_at,
			COALESCE(expires_at, '0001-01-01 00:00:00.0+00'),
			COALESCE(last_used, '0001-01-01 00:00:00.0+00')
		FROM api_key
		WHERE environment_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, environmentID)
	if err != nil {
		return nil, r.errorMapper(err, r.talbeName)
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
			return nil, r.errorMapper(err, r.talbeName)
		}

		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, r.errorMapper(err, r.talbeName)
	}

	return apiKeys, nil
}

func (r *APIKeyRepository) GetByKey(
	ctx context.Context, key string,
) (*entities.APIKey, errors.PersistenceError) {
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
		return nil, r.errorMapper(err, r.talbeName)
	}

	return apiKey, nil
}

func (r *APIKeyRepository) GetByID(
	ctx context.Context, id int,
) (*entities.APIKey, errors.PersistenceError) {
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
		return nil, r.errorMapper(err, r.talbeName)
	}

	return apiKey, nil
}

func (r *APIKeyRepository) Exists(
	ctx context.Context, key string,
) (bool, errors.PersistenceError) {
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
		return false, r.errorMapper(err, r.talbeName)
	}

	return exists, nil
}

func (r *APIKeyRepository) Create(
	ctx context.Context, apiKey *entities.APIKey,
) errors.PersistenceError {
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

	return r.errorMapper(err, r.talbeName)
}

func NewAPIKeyRepository(driver *Driver) *APIKeyRepository {
	return &APIKeyRepository{Driver: driver, talbeName: "api_key"}
}
