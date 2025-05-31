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

type RequestRepository struct {
	*Driver

	tableName string
}

func (r *RequestRepository) DeleteByService(
	ctx context.Context, serviceID int,
) errors.Error {
	query := `
		DELETE FROM request
		WHERE service_id = $1;
	`

	_, err := r.pool.Exec(ctx, query, serviceID)
	if err != nil {
		return r.errorMapper(err, r.tableName)
	}

	return nil
}

func (r *RequestRepository) UpdateExecutionStatus(
	ctx context.Context, id string, executionStatus enums.RequestExecutionStatus,
) errors.Error {
	query := `
		UPDATE request
		SET execution_status = $1
		WHERE id = $2;
	`

	result, err := r.pool.Exec(ctx, query, executionStatus, id)
	if err != nil {
		return r.errorMapper(err, r.tableName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(r.tableName, map[string]any{"id": id})
	}

	return nil
}

func (r *RequestRepository) ListByService(
	ctx context.Context, serviceID int, filter *dto.RequestFilter,
) ([]*entities.Request, errors.Error) {
	query := `
		SELECT id, COALESCE(start_point::text, ''), api_key, COALESCE(api_key_id, 0),
			COALESCE(project_name, ''), COALESCE(project_id, 0),
			COALESCE(environment_name, ''), COALESCE(environment_id, 0),
			service_name, service_version, COALESCE(service_id, 0), COALESCE(status_code, 0),
			execution_status, request_time, path, method, ip_address,
			COALESCE(unauthorized_reason, ''), created_at
		FROM request
		WHERE service_id = $1
	`

	args := []any{serviceID}
	if filter != nil {
		var where []string
		argIndex := 2

		if !filter.RequestTimeTo.IsZero() {
			where = append(where, fmt.Sprintf("request_time <= $%d", argIndex))
			args = append(args, filter.RequestTimeTo)
			argIndex++
		}

		if !filter.RequestTimeFrom.IsZero() {
			where = append(where, fmt.Sprintf("request_time >= $%d", argIndex))
			args = append(args, filter.RequestTimeFrom)
			argIndex++
		}

		if filter.ExecutionStatus != enums.RequestExecutionStatusNull {
			where = append(where, fmt.Sprintf("execution_status = $%d", argIndex))
			args = append(args, filter.ExecutionStatus)
			argIndex++
		}

		if len(where) > 0 {
			query = fmt.Sprintf("%s AND %s", query, strings.Join(where, " AND "))
		}
	}

	query += " ORDER BY request_time DESC;"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	defer rows.Close()

	var requests []*entities.Request
	for rows.Next() {
		request := &entities.Request{
			APIKey:      &entities.RequestAPIKey{},
			Project:     &entities.RequestProject{},
			Environment: &entities.RequestEnvironment{},
			Service:     &entities.RequestService{},
		}

		err := rows.Scan(
			&request.ID,
			&request.StartPoint,
			&request.APIKey.Key,
			&request.APIKey.ID,
			&request.Project.Name,
			&request.Project.ID,
			&request.Environment.Name,
			&request.Environment.ID,
			&request.Service.Name,
			&request.Service.Version,
			&request.Service.ID,
			&request.StatusCode,
			&request.ExecutionStatus,
			&request.RequestTime,
			&request.Path,
			&request.Method,
			&request.IPAddress,
			&request.UnauthorizedReason,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, r.errorMapper(err, r.tableName)
		}

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}
	return requests, nil
}

func (r *RequestRepository) Create(
	ctx context.Context, request *entities.Request,
) errors.Error {
	query := `
		INSERT INTO request (
			start_point, api_key, api_key_id, project_name, project_id,
			environment_name, environment_id, service_name, service_version,
			service_id, status_code, execution_status, request_time, path,
			method, ip_address, metadata, unauthorized_reason
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16,
			$17, $18
		) RETURNING id, created_at;
	`

	var startPoint any
	if request.StartPoint != "" {
		startPoint = request.StartPoint
	}

	var apiKeyID any
	if request.APIKey.ID != 0 {
		apiKeyID = request.APIKey.ID
	}

	var projectID any
	if request.Project.ID != 0 {
		projectID = request.Project.ID
	}

	var projectName any
	if request.Project.Name != "" {
		projectName = request.Project.Name
	}

	var environmentID any
	if request.Environment.ID != 0 {
		environmentID = request.Environment.ID
	}

	var environmentName any
	if request.Environment.Name != "" {
		environmentName = request.Environment.Name
	}

	var statusCode any
	if request.StatusCode != 0 {
		statusCode = request.StatusCode
	}

	var method any
	if request.Method != "" {
		method = request.Method
	}

	var UnauthorizedReason any
	if request.UnauthorizedReason != "" {
		UnauthorizedReason = request.UnauthorizedReason
	}

	metadata := make(map[string]any)
	if request.Metadata != nil {
		metadata["body"] = request.Metadata.Body
		metadata["headers"] = request.Metadata.Headers
		metadata["queryParams"] = request.Metadata.QueryParams
		metadata["bodyContentType"] = request.Metadata.BodyContentType
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		startPoint,
		request.APIKey.Key,
		apiKeyID,
		projectName,
		projectID,
		environmentName,
		environmentID,
		request.Service.Name,
		request.Service.Version,
		request.Service.ID,
		statusCode,
		request.ExecutionStatus,
		request.RequestTime,
		request.Path,
		method,
		request.IPAddress,
		metadata,
		UnauthorizedReason,
	).Scan(&request.ID, &request.CreatedAt)

	return r.errorMapper(err, r.tableName)
}

func (r *RequestRepository) CreateAsInitialPoint(
	ctx context.Context, request *entities.Request,
) errors.Error {
	query := `
		WITH temp_table AS (
			SELECT gen_random_uuid() AS uuid
		)
		INSERT INTO request (
			id, start_point, api_key, api_key_id, project_name, project_id,
			environment_name, environment_id, service_name, service_version,
			service_id, status_code, execution_status, request_time, path,
			method, ip_address, metadata, unauthorized_reason
		)
		SELECT uuid, uuid, $1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16, $17
		) RETURNING id, created_at;
		FROM temp_table RETURNING id;
	`

	var startPoint any
	if request.StartPoint != "" {
		startPoint = request.StartPoint
	}

	var apiKeyID any
	if request.APIKey.ID != 0 {
		apiKeyID = request.APIKey.ID
	}

	var projectID any
	if request.Project.ID != 0 {
		projectID = request.Project.ID
	}

	var projectName any
	if request.Project.Name != "" {
		projectName = request.Project.Name
	}

	var environmentID any
	if request.Environment.ID != 0 {
		environmentID = request.Environment.ID
	}

	var environmentName any
	if request.Environment.Name != "" {
		environmentName = request.Environment.Name
	}

	var statusCode any
	if request.StatusCode != 0 {
		statusCode = request.StatusCode
	}

	var method any
	if request.Method != "" {
		method = request.Method
	}

	var UnauthorizedReason any
	if request.UnauthorizedReason != "" {
		UnauthorizedReason = request.UnauthorizedReason
	}

	metadata := map[string]any{
		"body":            "",
		"headers":         "",
		"queryParams":     "",
		"bodyContentType": enums.RequestBodyContentTypeNull,
	}
	if request.Metadata != nil {
		metadata["body"] = request.Metadata.Body
		metadata["headers"] = request.Metadata.Headers
		metadata["queryParams"] = request.Metadata.QueryParams
		metadata["bodyContentType"] = request.Metadata.BodyContentType
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		startPoint,
		request.APIKey,
		apiKeyID,
		projectName,
		projectID,
		environmentName,
		environmentID,
		request.Service.Name,
		request.Service.Version,
		request.Service.ID,
		statusCode,
		request.ExecutionStatus,
		request.RequestTime,
		request.Path,
		method,
		request.IPAddress,
		metadata,
		UnauthorizedReason,
	).Scan(&request.ID, &request.CreatedAt)

	return r.errorMapper(err, r.tableName)
}

func NewRequestRepository(driver *Driver) *RequestRepository {
	return &RequestRepository{Driver: driver, tableName: "request"}
}
