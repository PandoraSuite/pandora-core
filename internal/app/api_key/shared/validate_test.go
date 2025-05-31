package shared

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/app/api_key/shared/mock"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type UseCaseSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	apiKeyRepo      *mock.MockValidateAPIKeyRepository
	projectRepo     *mock.MockValidateProjectRepository
	serviceRepo     *mock.MockValidateServiceRepository
	environmentRepo *mock.MockValidateEnvironmentRepository

	deps *ValidateDependencies

	ctx context.Context
}

func (s *UseCaseSuite) SetupTest() {
	time.Local = time.UTC
	s.ctrl = gomock.NewController(s.T())

	s.apiKeyRepo = mock.NewMockValidateAPIKeyRepository(s.ctrl)
	s.projectRepo = mock.NewMockValidateProjectRepository(s.ctrl)
	s.serviceRepo = mock.NewMockValidateServiceRepository(s.ctrl)
	s.environmentRepo = mock.NewMockValidateEnvironmentRepository(s.ctrl)

	s.deps = NewValidationDependencies(
		s.apiKeyRepo,
		s.serviceRepo,
		s.projectRepo,
		s.environmentRepo,
	)

	s.ctx = context.Background()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *UseCaseSuite) TestSuccess() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
		Services: []*entities.EnvironmentService{
			{
				ID:               service.ID,
				Name:             service.Name,
				Version:          service.Version,
				MaxRequest:       -1,
				AvailableRequest: -1,
				AssignedAt:       reqTime,
			},
		},
	}
	projectClient := &dto.ProjectClientInfoResponse{
		ProjectID:   1000,
		ProjectName: "TestProject",
		ClientID:    2000,
		ClientName:  "TestClient",
	}
	consumerInfo := &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(projectClient, nil).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.True(validateResponse.Valid)
	s.Empty(validateResponse.FailureCode)
	s.Equal(consumerInfo, validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Equal(consumerInfo.ProjectID, request.Project.ID)
	s.Equal(consumerInfo.ProjectName, request.Project.Name)
}

func (s *UseCaseSuite) TestAPIKeyNotFound() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "invalid-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(nil, errors.NewEntityNotFound(
			"APIKey",
			"api key not found",
			map[string]any{"key": req.APIKey},
			nil,
		)).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, gomock.Any()).
		Times(0)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, gomock.Any()).
		Times(0)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.False(validateResponse.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeAPIKeyInvalid, validateResponse.FailureCode)
	s.Nil(validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Zero(request.APIKey.ID)
	s.Zero(request.Environment.ID)
	s.Zero(request.Project.ID)
	s.Empty(request.Project.Name)
}

func (s *UseCaseSuite) TestServiceMismatch() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "NonExistentService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectClient := &dto.ProjectClientInfoResponse{
		ProjectID:   1000,
		ProjectName: "TestProject",
		ClientID:    2000,
		ClientName:  "TestClient",
	}
	consumerInfo := &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(nil, errors.NewEntityNotFound(
			"Service",
			"service not found",
			map[string]any{
				"name":    req.ServiceName,
				"version": req.ServiceVersion,
			},
			nil,
		)).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(projectClient, nil).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.False(validateResponse.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeServiceMismatch, validateResponse.FailureCode)
	s.Equal(consumerInfo, validateResponse.ConsumerInfo)

	s.Zero(request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Equal(consumerInfo.ProjectID, request.Project.ID)
	s.Equal(consumerInfo.ProjectName, request.Project.Name)
}

func (s *UseCaseSuite) TestServiceRepoInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	internalErr := errors.NewInternal("service repo error", nil)

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(nil, internalErr).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, gomock.Any()).
		Times(0)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, gomock.Any()).
		Times(0)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, gomock.Any()).
		Times(0)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().ErrorIs(err, internalErr)

	s.False(validateResponse.Valid)
	s.Empty(validateResponse.FailureCode)
	s.Nil(validateResponse.ConsumerInfo)

	s.Zero(request.Service.ID)
	s.Zero(request.APIKey.ID)
	s.Zero(request.Environment.ID)
	s.Zero(request.Project.ID)
	s.Empty(request.Project.Name)
}

func (s *UseCaseSuite) TestAPIKeyRepoInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}
	internalErr := errors.NewInternal("api key repo error", nil)

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(nil, internalErr).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, gomock.Any()).
		Times(0)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, gomock.Any()).
		Times(0)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().ErrorIs(err, internalErr)

	s.False(validateResponse.Valid)
	s.Empty(validateResponse.FailureCode)
	s.Nil(validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Zero(request.APIKey.ID)
	s.Zero(request.Environment.ID)
	s.Zero(request.Project.ID)
	s.Empty(request.Project.Name)
}

func (s *UseCaseSuite) TestEnvironmentRepoGetByIDInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	internalErr := errors.NewInternal("environment repo error", nil)

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(nil, internalErr).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, gomock.Any()).
		Times(0)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().ErrorIs(err, internalErr)

	s.False(validateResponse.Valid)
	s.Empty(validateResponse.FailureCode)
	s.Nil(validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Zero(request.Environment.ID)
	s.Zero(request.Project.ID)
	s.Empty(request.Project.Name)
}

func (s *UseCaseSuite) TestProjectRepoInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
		Services: []*entities.EnvironmentService{
			{
				ID:               service.ID,
				Name:             service.Name,
				Version:          service.Version,
				MaxRequest:       -1,
				AvailableRequest: -1,
				AssignedAt:       reqTime,
			},
		},
	}
	internalErr := errors.NewInternal("project repo error", nil)

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(nil, internalErr).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().ErrorIs(err, internalErr)

	s.False(validateResponse.Valid)
	s.Empty(validateResponse.FailureCode)
	s.Nil(validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Zero(request.Project.ID)
	s.Empty(request.Project.Name)
}

func (s *UseCaseSuite) TestEnvironmentDisabled() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusDisabled,
		ProjectID: 1000,
		Services: []*entities.EnvironmentService{
			{
				ID:               service.ID,
				Name:             service.Name,
				Version:          service.Version,
				MaxRequest:       -1,
				AvailableRequest: -1,
				AssignedAt:       reqTime,
			},
		},
	}
	projectClient := &dto.ProjectClientInfoResponse{
		ProjectID:   1000,
		ProjectName: "TestProject",
		ClientID:    2000,
		ClientName:  "TestClient",
	}
	consumerInfo := &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(projectClient, nil).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.False(validateResponse.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeEnvironmentDisabled, validateResponse.FailureCode)
	s.Equal(consumerInfo, validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Equal(consumerInfo.ProjectID, request.Project.ID)
	s.Equal(consumerInfo.ProjectName, request.Project.Name)
}

func (s *UseCaseSuite) TestServiceNotAssignedToEnvironment() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
		Services: []*entities.EnvironmentService{
			{
				ID:               10,
				Name:             "test",
				Version:          "1.0.0",
				MaxRequest:       -1,
				AvailableRequest: -1,
				AssignedAt:       reqTime,
			},
		},
	}
	projectClient := &dto.ProjectClientInfoResponse{
		ProjectID:   1000,
		ProjectName: "TestProject",
		ClientID:    2000,
		ClientName:  "TestClient",
	}
	consumerInfo := &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(projectClient, nil).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.False(validateResponse.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeServiceNotAssigned, validateResponse.FailureCode)
	s.Equal(consumerInfo, validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Equal(consumerInfo.ProjectID, request.Project.ID)
	s.Equal(consumerInfo.ProjectName, request.Project.Name)
}

func (s *UseCaseSuite) TestServiceDisabled() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "DisabledService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusDisabled,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
		Services: []*entities.EnvironmentService{
			{
				ID:               service.ID,
				Name:             service.Name,
				Version:          service.Version,
				MaxRequest:       -1,
				AvailableRequest: -1,
				AssignedAt:       reqTime,
			},
		},
	}
	projectClient := &dto.ProjectClientInfoResponse{
		ProjectID:   1000,
		ProjectName: "TestProject",
		ClientID:    2000,
		ClientName:  "TestClient",
	}
	consumerInfo := &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(projectClient, nil).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.False(validateResponse.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeServiceDisabled, validateResponse.FailureCode)
	s.Equal(consumerInfo, validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Equal(consumerInfo.ProjectID, request.Project.ID)
	s.Equal(consumerInfo.ProjectName, request.Project.Name)
}

func (s *UseCaseSuite) TestServiceDeprecated() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:         "valid-api-key",
		ServiceName:    "DeprecatedService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusDeprecated,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
		Services: []*entities.EnvironmentService{
			{
				ID:               service.ID,
				Name:             service.Name,
				Version:          service.Version,
				MaxRequest:       -1,
				AvailableRequest: -1,
				AssignedAt:       reqTime,
			},
		},
	}
	projectClient := &dto.ProjectClientInfoResponse{
		ProjectID:   1000,
		ProjectName: "TestProject",
		ClientID:    2000,
		ClientName:  "TestClient",
	}
	consumerInfo := &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(projectClient, nil).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.False(validateResponse.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeServiceDeprecated, validateResponse.FailureCode)
	s.Equal(consumerInfo, validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Equal(consumerInfo.ProjectID, request.Project.ID)
	s.Equal(consumerInfo.ProjectName, request.Project.Name)
}

func (s *UseCaseSuite) TestAPIKeyExpired() {
	reqTime := time.Now()
	pastTime := time.Now().Add(-24 * time.Hour)
	req := &dto.APIKeyValidate{
		APIKey:         "expired-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
		ExpiresAt:     pastTime,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
		Services: []*entities.EnvironmentService{
			{
				ID:               service.ID,
				Name:             service.Name,
				Version:          service.Version,
				MaxRequest:       -1,
				AvailableRequest: -1,
				AssignedAt:       reqTime,
			},
		},
	}
	projectClient := &dto.ProjectClientInfoResponse{
		ProjectID:   1000,
		ProjectName: "TestProject",
		ClientID:    2000,
		ClientName:  "TestClient",
	}
	consumerInfo := &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(projectClient, nil).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.False(validateResponse.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeAPIKeyExpired, validateResponse.FailureCode)
	s.Equal(consumerInfo, validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Equal(consumerInfo.ProjectID, request.Project.ID)
	s.Equal(consumerInfo.ProjectName, request.Project.Name)
}

func (s *UseCaseSuite) TestAPIKeyDisabled() {
	reqTime := time.Now()
	pastTime := time.Now().Add(-24 * time.Hour)
	req := &dto.APIKeyValidate{
		APIKey:         "disabled-api-key",
		ServiceName:    "TestService",
		ServiceVersion: "1.0.0",
		Request: &dto.RequestIncoming{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	service := &entities.Service{
		ID:      1,
		Name:    req.ServiceName,
		Version: req.ServiceVersion,
		Status:  enums.ServiceStatusEnabled,
	}
	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusDisabled,
		EnvironmentID: 100,
		ExpiresAt:     pastTime,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      "production",
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
		Services: []*entities.EnvironmentService{
			{
				ID:               service.ID,
				Name:             service.Name,
				Version:          service.Version,
				MaxRequest:       -1,
				AvailableRequest: -1,
				AssignedAt:       reqTime,
			},
		},
	}
	projectClient := &dto.ProjectClientInfoResponse{
		ProjectID:   1000,
		ProjectName: "TestProject",
		ClientID:    2000,
		ClientName:  "TestClient",
	}
	consumerInfo := &dto.ConsumerInfo{
		ClientID:        projectClient.ClientID,
		ClientName:      projectClient.ClientName,
		ProjectID:       projectClient.ProjectID,
		ProjectName:     projectClient.ProjectName,
		EnvironmentID:   environment.ID,
		EnvironmentName: environment.Name,
	}

	s.serviceRepo.EXPECT().
		GetByNameAndVersion(s.ctx, req.ServiceName, req.ServiceVersion).
		Return(service, nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByKey(s.ctx, req.APIKey).
		Return(apiKey, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, apiKey.EnvironmentID).
		Return(environment, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectClientInfoByID(s.ctx, environment.ProjectID).
		Return(projectClient, nil).
		Times(1)

	validateResponse := dto.APIKeyValidateResponse{}

	request := entities.Request{
		Path:        req.Request.Path,
		Method:      req.Request.Method,
		IPAddress:   req.Request.IPAddress,
		RequestTime: req.Request.RequestTime,
		APIKey:      &entities.RequestAPIKey{Key: req.APIKey},
		Service: &entities.RequestService{
			Name:    req.ServiceName,
			Version: req.ServiceVersion,
		},
		Environment: &entities.RequestEnvironment{},
		Project:     &entities.RequestProject{},
	}

	err := ValidateAPIKey(s.ctx, s.deps, req, &request, &validateResponse)

	s.Require().NoError(err)

	s.False(validateResponse.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeAPIKeyDisabled, validateResponse.FailureCode)
	s.Equal(consumerInfo, validateResponse.ConsumerInfo)

	s.Equal(service.ID, request.Service.ID)
	s.Equal(apiKey.ID, request.APIKey.ID)
	s.Equal(environment.ID, request.Environment.ID)
	s.Equal(consumerInfo.ProjectID, request.Project.ID)
	s.Equal(consumerInfo.ProjectName, request.Project.Name)
}

func TestValidateOnlySuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}
