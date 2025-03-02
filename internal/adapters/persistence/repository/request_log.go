package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestLogRepository struct {
	pool *pgxpool.Pool
}

func (r *RequestLogRepository) Save(
	ctx context.Context, requestLog *entities.RequestLog,
) (*entities.RequestLog, error) {
	model := models.RequestLogFromEntity(requestLog)
	if err := r.save(ctx, model); err != nil {
		return nil, err
	}

	return model.ToEntity()
}

func (r *RequestLogRepository) save(
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
		return persistence.ConvertPgxError(err)
	}

	return nil
}

func NewRequestLogRepository(pool *pgxpool.Pool) *RequestLogRepository {
	return &RequestLogRepository{pool: pool}
}
