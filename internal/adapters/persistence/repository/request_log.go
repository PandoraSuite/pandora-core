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
	ctx context.Context, requestLog *models.RequestLog,
) error {
	if err := requestLog.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO request_log (environment_id, service_id, api_key, request_time, execution_status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		requestLog.EnvironmentID,
		requestLog.ServiceID,
		requestLog.APIKey,
		requestLog.RequestTime,
		requestLog.ExecutionStatus,
	).Scan(&requestLog.ID, &requestLog.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the request log: %w", err)
	}

	return nil
}
