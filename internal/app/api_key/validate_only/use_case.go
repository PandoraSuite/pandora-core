package validateonly

import (
	"context"
	"log"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.APIKeyValidate) (*dto.APIKeyValidateResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	apiKeyRepo      APIKeyRepository
	projectRepo     ProjectRepository
	serviceRepo     ServiceRepository
	requestRepo     RequestRepository
	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.APIKeyValidate,
) (*dto.APIKeyValidateResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	validateResponse := dto.APIKeyValidateResponse{}

	var requestMetadata entities.RequestMetadata
	if req.Request.Metadata != nil {
		requestMetadata = entities.RequestMetadata{
			QueryParams:     req.Request.Metadata.QueryParams,
			Headers:         req.Request.Metadata.Headers,
			Body:            req.Request.Metadata.Body,
			BodyContentType: req.Request.Metadata.BodyContentType,
		}
	}

	request := entities.Request{
		Path:            req.Request.Path,
		Method:          req.Request.Method,
		IPAddress:       req.Request.IPAddress,
		RequestTime:     req.Request.RequestTime,
		Metadata:        &requestMetadata,
		APIKey:          req.APIKey,
		ServiceName:     req.ServiceName,
		ServiceVersion:  req.ServiceVersion,
		EnvironmentName: req.EnvironmentName,
	}

	err := uc.validate(ctx, req, &request, &validateResponse)
	if err != nil {
		return nil, err
	}

	if validateResponse.Valid {
		request.ExecutionStatus = enums.RequestExecutionStatusForwarded
	} else {
		request.ExecutionStatus = enums.RequestExecutionStatusUnauthorized
	}

	if err := uc.requestRepo.Create(ctx, &request); err != nil {
		return nil, err
	}

	validateResponse.RequestID = request.ID

	if validateResponse.Valid {
		if err := uc.apiKeyRepo.UpdateLastUsed(ctx, req.APIKey); err != nil {
			log.Printf("Failed to update last_used for API Key %s: %v", req.APIKey, err)
		}
	}
	return &validateResponse, nil
}

func (uc *useCase) validate(
	ctx context.Context,
	req *dto.APIKeyValidate,
	request *entities.Request,
	validateResponse *dto.APIKeyValidateResponse,
) errors.Error {
	service, err := uc.serviceRepo.GetByNameAndVersion(
		ctx, req.ServiceName, req.ServiceVersion,
	)
	if err != nil {
		if err.Code() != errors.CodeNotFound {
			return err
		}
		uc.setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeServiceMismatch,
		)
	} else {
		request.ServiceID = service.ID

		if service.IsDisabled() {
			uc.setFailureIfEmpty(
				validateResponse,
				enums.APIKeyValidationFailureCodeServiceDisabled,
			)
		}

		if service.IsDeprecated() {
			uc.setFailureIfEmpty(
				validateResponse,
				enums.APIKeyValidationFailureCodeServiceDeprecated,
			)
		}
	}

	apiKey, err := uc.apiKeyRepo.GetByKey(ctx, req.APIKey)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			uc.setFailureIfEmpty(
				validateResponse,
				enums.APIKeyValidationFailureCodeAPIKeyInvalid,
			)
			return nil
		}
		return err
	}

	request.APIKeyID = apiKey.ID

	if !apiKey.IsEnabled() {
		uc.setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeAPIKeyDisabled,
		)
	}

	if apiKey.IsExpired() {
		uc.setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeAPIKeyExpired,
		)
	}

	environment, err := uc.environmentRepo.GetByID(ctx, apiKey.EnvironmentID)
	if err != nil {
		return err
	}

	if !environment.Is(req.EnvironmentName) {
		uc.setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeEnvironmentMismatch,
		)
		return nil
	}

	request.EnvironmentID = environment.ID

	if !environment.IsEnabled() {
		uc.setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeEnvironmentDisabled,
		)
	}

	consumer, err := uc.projectRepo.GetProjectContextByID(
		ctx, environment.ProjectID,
	)
	if err != nil {
		return err
	}

	request.ProjectID = consumer.ID
	request.ProjectName = consumer.Name
	validateResponse.ConsumerInfo = consumer

	if service != nil {
		exists, err := uc.environmentRepo.ExistsServiceIn(
			ctx, environment.ID, service.ID,
		)
		if err != nil {
			return err
		}

		if !exists {
			uc.setFailureIfEmpty(
				validateResponse,
				enums.APIKeyValidationFailureCodeServiceNotAssigned,
			)
		}
	}

	validateResponse.Valid = validateResponse.FailureCode == ""
	return nil
}

func (uc *useCase) setFailureIfEmpty(
	validateResponse *dto.APIKeyValidateResponse,
	failureCode enums.APIKeyValidationFailureCode,
) {
	if validateResponse.FailureCode == "" {
		validateResponse.FailureCode = failureCode
	}
}

func (uc *useCase) validateReq(req *dto.APIKeyValidate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"api_key.required":         "api_key is required",
			"request.required":         "request is required",
			"service.required":         "service is required",
			"environment.required":     "environment is required",
			"service_version.required": "service_version is required",

			"request.path.required":         "request.path is required",
			"request.ip_address.ip":         "request.ip_address must be a valid IP address",
			"request.method.required":       "request.method is required",
			"request.request_time.utc":      "request.request_time must be in UTC format",
			"request.ip_address.required":   "request.ip_address is required",
			"request.request_time.required": "request.request_time is required",

			"request.metadata.body_content_type.enums": "request.metadata.body_content_type must be one of the following: application/xml, application/json, text/plain, text/html, multipart/form-data, application/x-www-form-urlencoded, application/octet-stream",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	apiKeyRepo APIKeyRepository,
	projectRepo ProjectRepository,
	serviceRepo ServiceRepository,
	requestRepo RequestRepository,
	environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		apiKeyRepo:      apiKeyRepo,
		projectRepo:     projectRepo,
		serviceRepo:     serviceRepo,
		requestRepo:     requestRepo,
		environmentRepo: environmentRepo,
	}
}
