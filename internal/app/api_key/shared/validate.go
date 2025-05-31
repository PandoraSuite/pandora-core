package shared

import (
	"context"
	"slices"

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
}

type ValidateProjectRepository interface {
	GetProjectClientInfoByID(ctx context.Context, id int) (*dto.ProjectClientInfoResponse, errors.Error)
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
		request.Service.ID = service.ID

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

	request.APIKey.ID = apiKey.ID

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

	request.Environment.ID = environment.ID
	request.Environment.Name = environment.Name

	if !environment.IsEnabled() {
		setFailureIfEmpty(
			validateResponse,
			enums.APIKeyValidationFailureCodeEnvironmentDisabled,
		)
	}

	projectClient, err := deps.projectRepo.GetProjectClientInfoByID(
		ctx, environment.ProjectID,
	)
	if err != nil {
		return err
	}

	request.Project.ID = projectClient.ProjectID
	request.Project.Name = projectClient.ProjectName

	validateResponse.ConsumerInfo = &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	if service != nil && validateResponse.FailureCode == "" {
		index := slices.IndexFunc(
			environment.Services,
			func(s *entities.EnvironmentService) bool { return s.ID == service.ID },
		)

		if index == -1 {
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
