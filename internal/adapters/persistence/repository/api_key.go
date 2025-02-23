package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyRepository struct {
	pool *pgxpool.Pool
}

func (r *APIKeyRepository) Save(
	ctx context.Context, apiKey *entities.APIKey,
) (*entities.APIKey, error) {
	key := models.APIKeyFromEntity(apiKey)
	if err := r.save(ctx, key); err != nil {
		return nil, err
	}
	return key.ToEntity(), nil
}

func (r *APIKeyRepository) save(ctx context.Context, apiKey *models.APIKey) error {
	if err := apiKey.ValidateModel(); err != nil {
		return err
	}

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

	if err != nil {
		return fmt.Errorf("error when inserting the api key: %w", err)
	}

	return nil
}

func NewAPIKeyRepository(pool *pgxpool.Pool) outbound.APIKeyRepositoryPort {
	return &APIKeyRepository{pool: pool}
}
