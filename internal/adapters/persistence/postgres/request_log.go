package postgres

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type RequestLogRepository struct {
	*Driver

	tableName string
}

func (r *RequestLogRepository) DeleteByService(
	ctx context.Context, serviceID int,
) errors.PersistenceError {
	query := `
		DELETE FROM request_log
		WHERE service_id = $1;
	`

	_, err := r.pool.Exec(ctx, query, serviceID)
	if err != nil {
		return r.errorMapper(err, r.tableName)
	}

	return nil
}

func (r *RequestLogRepository) UpdateExecutionStatus(
	ctx context.Context, id string, executionStatus enums.RequestLogExecutionStatus,
) errors.PersistenceError {
	query := `
		UPDATE request_log
		SET execution_status = $1
		WHERE id = $2;
	`

	result, err := r.pool.Exec(ctx, query, executionStatus, id)
	if err != nil {
		return r.errorMapper(err, r.tableName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(r.tableName)
	}

	return nil
}

func (r *RequestLogRepository) Create(
	ctx context.Context, requestLog *entities.RequestLog,
) errors.PersistenceError {
	query := `
		INSERT INTO request_log (environment_id, service_id, api_key, start_point, request_time, execution_status, message)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at;
	`
	var startPoint any
	if requestLog.StartPoint != "" {
		startPoint = requestLog.StartPoint
	}

	var serviceID any
	if requestLog.ServiceID != 0 {
		serviceID = requestLog.ServiceID
	}

	var environmentID any
	if requestLog.EnvironmentID != 0 {
		environmentID = requestLog.EnvironmentID
	}
	err := r.pool.QueryRow(
		ctx,
		query,
		environmentID,
		serviceID,
		requestLog.APIKey,
		startPoint,
		requestLog.RequestTime,
		requestLog.ExecutionStatus,
		requestLog.Message,
	).Scan(&requestLog.ID, &requestLog.CreatedAt)

	return r.errorMapper(err, r.tableName)
}

func (r *RequestLogRepository) CreateAsInitialPoint(
	ctx context.Context, requestLog *entities.RequestLog,
) errors.PersistenceError {
	query := `
		WITH temp_table AS (
			SELECT gen_random_uuid() AS uuid
		)
		INSERT INTO request_log (id, environment_id, service_id, api_key, start_point, request_time, execution_status, message) 
		SELECT uuid, $1, $2, $3, uuid, $4, $5, $6
		FROM temp_table RETURNING id;
	`
	var serviceID any
	if requestLog.ServiceID != 0 {
		serviceID = requestLog.ServiceID
	}

	var environmentID any
	if requestLog.EnvironmentID != 0 {
		environmentID = requestLog.EnvironmentID
	}
	err := r.pool.QueryRow(
		ctx,
		query,
		environmentID,
		serviceID,
		requestLog.APIKey,
		requestLog.RequestTime,
		requestLog.ExecutionStatus,
		requestLog.Message,
	).Scan(&requestLog.ID)

	return r.errorMapper(err, r.tableName)
}

func NewRequestLogRepository(driver *Driver) *RequestLogRepository {
	return &RequestLogRepository{Driver: driver, tableName: "request_log"}
}
