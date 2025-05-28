package shared

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ValidateServiceRepository interface {
	GetByNameAndVersion(ctx context.Context, name, version string) (*entities.Service, errors.Error)
}

type ValidateAPIKeyRepository interface {
	GetByKey(ctx context.Context, key string) (*entities.APIKey, errors.Error)
}

type ValidateEnvironmentRepository interface {
	GetByID(ctx context.Context, id int) (*entities.Environment, errors.Error)
	ExistsServiceIn(ctx context.Context, environmentID, serviceID int) (bool, errors.Error)
}

type ValidateProjectRepository interface {
	GetProjectContextByID(ctx context.Context, id int) (*dto.ProjectContextResponse, errors.Error)
}

type ValidateDependencies struct {
	apiKeyRepo      ValidateAPIKeyRepository
	serviceRepo     ValidateServiceRepository
	projectRepo     ValidateProjectRepository
	environmentRepo ValidateEnvironmentRepository
}

func NewValidationDependencies(
	apiKeyRepo ValidateAPIKeyRepository,
	serviceRepo ValidateServiceRepository,
	projectRepo ValidateProjectRepository,
	environmentRepo ValidateEnvironmentRepository,
) *ValidateDependencies {
	return &ValidateDependencies{
		apiKeyRepo:      apiKeyRepo,
		serviceRepo:     serviceRepo,
		projectRepo:     projectRepo,
		environmentRepo: environmentRepo,
	}
}

func ValidateAPIKey(
	ctx context.Context,
	deps *ValidateDependencies,
	req *dto.APIKeyValidate,
	request *entities.Request,
	validateResponse *dto.APIKeyValidateResponse,
) errors.Error {
	service, err := deps.serviceRepo.GetByNameAndVersion(
		ctx, req.ServiceName, req.ServiceVersion,
	)
	if err != nil {
		if err.Code() != errors.CodeNotFound {
			return err
		}
		setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeServiceMismatch,
		)
	} else {
		request.ServiceID = service.ID

		if service.IsDisabled() {
			setFailureIfEmpty(
				validateResponse,
				enums.APIKeyValidationFailureCodeServiceDisabled,
			)
		}

		if service.IsDeprecated() {
			setFailureIfEmpty(
				validateResponse,
				enums.APIKeyValidationFailureCodeServiceDeprecated,
			)
		}
	}

	apiKey, err := deps.apiKeyRepo.GetByKey(ctx, req.APIKey)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			setFailureIfEmpty(
				validateResponse,
				enums.APIKeyValidationFailureCodeAPIKeyInvalid,
			)
			return nil
		}
		return err
	}

	request.APIKeyID = apiKey.ID

	if !apiKey.IsEnabled() {
		setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeAPIKeyDisabled,
		)
	}

	if apiKey.IsExpired() {
		setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeAPIKeyExpired,
		)
	}

	environment, err := deps.environmentRepo.GetByID(ctx, apiKey.EnvironmentID)
	if err != nil {
		return err
	}

	if !environment.Is(req.EnvironmentName) {
		setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeEnvironmentMismatch,
		)
		return nil
	}

	request.EnvironmentID = environment.ID

	if !environment.IsEnabled() {
		setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeEnvironmentDisabled,
		)
	}

	consumer, err := deps.projectRepo.GetProjectContextByID(
		ctx, environment.ProjectID,
	)
	if err != nil {
		return err
	}

	request.ProjectID = consumer.ID
	request.ProjectName = consumer.Name
	validateResponse.ConsumerInfo = consumer

	if service != nil && validateResponse.FailureCode == "" {
		exists, err := deps.environmentRepo.ExistsServiceIn(
			ctx, environment.ID, service.ID,
		)
		if err != nil {
			return err
		}

		if !exists {
			setFailureIfEmpty(
				validateResponse,
				enums.APIKeyValidationFailureCodeServiceNotAssigned,
			)
		}
	}

	validateResponse.Valid = validateResponse.FailureCode == ""
	return nil
}

func setFailureIfEmpty(
	validateResponse *dto.APIKeyValidateResponse,
	failureCode enums.APIKeyValidationFailureCode,
) {
	if validateResponse.FailureCode == "" {
		validateResponse.FailureCode = failureCode
	}
}
