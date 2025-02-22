package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyRepository struct {
	pool *pgxpool.Pool
}

func (r *APIKeyRepository) Create(
	ctx context.Context, newAPIKey *models.APIKey,
) error {
	if err := newAPIKey.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO api_key (environment_id, key, expires_at, last_used, status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		newAPIKey.EnvironmentID,
		newAPIKey.Key,
		newAPIKey.ExpiresAt,
		newAPIKey.LastUsed,
		newAPIKey.Status,
	).Scan(&newAPIKey.ID, &newAPIKey.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the api key: %w", err)
	}

	return nil
}
