package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestLogRepository struct {
	pool *pgxpool.Pool
}

func (r *RequestLogRepository) Create(
	ctx context.Context, newRequestLog *models.RequestLog,
) error {
	if err := newRequestLog.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO client (environment_id, service_id, api_key, request_time, execution_status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		newRequestLog.EnvironmentID,
		newRequestLog.ServiceID,
		newRequestLog.APIKey,
		newRequestLog.RequestTime,
		newRequestLog.ExecutionStatus,
	).Scan(&newRequestLog.ID, &newRequestLog.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the request log: %w", err)
	}

	return nil
}
