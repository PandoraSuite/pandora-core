package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type RequestLogRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *RequestLogRepository) UpdateExecutionStatus(
	ctx context.Context, id int, executionStatus enums.RequestLogExecutionStatus,
) *errors.Error {
	if executionStatus == enums.RequestLogExecutionStatusNull {
		return errors.ErrAPIKeyInvalidStatus
	}

	query := `
		UPDATE request_log
		SET execution_status = $1
		WHERE id = $2;
	`

	result, err := r.pool.Exec(ctx, query, executionStatus, id)
	if err != nil {
		return r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrRequestLogNotFound
	}

	return nil
}

func (r *RequestLogRepository) Save(
	ctx context.Context, requestLog *entities.RequestLog,
) *errors.Error {
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

	return r.handlerErr(err)
}

func NewRequestLogRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *RequestLogRepository {
	return &RequestLogRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
