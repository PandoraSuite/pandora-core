package postgres

import (
	"context"

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
	ctx context.Context, serviceID int,
) ([]*dto.RequestResponse, errors.Error) {
	query := `
		SELECT id, start_point, api_key, api_key_id, project_name, project_id,
			environment_name, environment_id, service_name, service_version,
			service_id, status_code, execution_status, request_time, path,
			method, ip_address, created_at
		FROM request
		WHERE service_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query, serviceID)
	if err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	defer rows.Close()

	var requests []*dto.RequestResponse
	for rows.Next() {
		request := new(dto.RequestResponse)

		err := rows.Scan(
			&request.ID,
			&request.StartPoint,
			&request.APIKey,
			&request.APIKeyID,
			&request.ProjectName,
			&request.ProjectID,
			&request.EnvironmentName,
			&request.EnvironmentID,
			&request.ServiceName,
			&request.ServiceVersion,
			&request.ServiceID,
			&request.StatusCode,
			&request.ExecutionStatus,
			&request.RequestTime,
			&request.Path,
			&request.Method,
			&request.IPAddress,
			&request.CreateAt,
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
			method, ip_address, body, body_content_type, headers, query_params
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20
		) RETURNING id, created_at;
	`

	var startPoint any
	if request.StartPoint != "" {
		startPoint = request.StartPoint
	}

	var apiKeyID any
	if request.APIKeyID != 0 {
		apiKeyID = request.APIKeyID
	}

	var projectID any
	if request.ProjectID != 0 {
		projectID = request.ProjectID
	}

	var environmentID any
	if request.EnvironmentID != 0 {
		environmentID = request.EnvironmentID
	}

	var statusCode any
	if request.StatusCode != 0 {
		statusCode = request.StatusCode
	}

	var method any
	if request.Method != "" {
		method = request.Method
	}

	var headers any
	if request.Headers != "" {
		headers = request.Headers
	}

	var queryParams any
	if request.QueryParams != "" {
		queryParams = request.QueryParams
	}

	var body any
	if request.Body != "" {
		body = request.Body
	}

	var bodyContentType any
	if request.BodyContentType != enums.RequestBodyContentTypeNull {
		bodyContentType = request.BodyContentType
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		startPoint,
		request.APIKey,
		apiKeyID,
		request.ProjectName,
		projectID,
		request.EnvironmentName,
		environmentID,
		request.ServiceName,
		request.ServiceVersion,
		request.ServiceID,
		statusCode,
		request.ExecutionStatus,
		request.RequestTime,
		request.Path,
		method,
		request.IPAddress,
		body,
		bodyContentType,
		headers,
		queryParams,
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
			method, ip_address, body, body_content_type, headers, query_params
		)
		SELECT uuid, uuid, $1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19
		) RETURNING id, created_at;
		FROM temp_table RETURNING id;
	`

	var apiKeyID any
	if request.APIKeyID != 0 {
		apiKeyID = request.APIKeyID
	}

	var projectID any
	if request.ProjectID != 0 {
		projectID = request.ProjectID
	}

	var environmentID any
	if request.EnvironmentID != 0 {
		environmentID = request.EnvironmentID
	}

	var statusCode any
	if request.StatusCode != 0 {
		statusCode = request.StatusCode
	}

	var method any
	if request.Method != "" {
		method = request.Method
	}

	var headers any
	if request.Headers != "" {
		headers = request.Headers
	}

	var queryParams any
	if request.QueryParams != "" {
		queryParams = request.QueryParams
	}

	var body any
	if request.Body != "" {
		body = request.Body
	}

	var bodyContentType any
	if request.BodyContentType != enums.RequestBodyContentTypeNull {
		bodyContentType = request.BodyContentType
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		request.APIKey,
		apiKeyID,
		request.ProjectName,
		projectID,
		request.EnvironmentName,
		environmentID,
		request.ServiceName,
		request.ServiceVersion,
		request.ServiceID,
		statusCode,
		request.ExecutionStatus,
		request.RequestTime,
		request.Path,
		method,
		request.IPAddress,
		body,
		bodyContentType,
		headers,
		queryParams,
	).Scan(&request.ID, &request.CreatedAt)

	return r.errorMapper(err, r.tableName)
}

func NewRequestRepository(driver *Driver) *RequestRepository {
	return &RequestRepository{Driver: driver, tableName: "request"}
}
