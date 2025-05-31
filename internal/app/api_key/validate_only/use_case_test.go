package validateonly

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/app/api_key/validate_only/mock"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	mockvalidator "github.com/MAD-py/pandora-core/internal/validator/mock"
)

type UseCaseSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	validator       *mockvalidator.MockValidator
	apiKeyRepo      *mock.MockAPIKeyRepository
	projectRepo     *mock.MockProjectRepository
	serviceRepo     *mock.MockServiceRepository
	requestRepo     *mock.MockRequestRepository
	environmentRepo *mock.MockEnvironmentRepository

	useCase UseCase

	ctx context.Context
}

func (s *UseCaseSuite) SetupTest() {
	time.Local = time.UTC
	s.ctrl = gomock.NewController(s.T())

	s.validator = mockvalidator.NewMockValidator(s.ctrl)
	s.apiKeyRepo = mock.NewMockAPIKeyRepository(s.ctrl)
	s.projectRepo = mock.NewMockProjectRepository(s.ctrl)
	s.serviceRepo = mock.NewMockServiceRepository(s.ctrl)
	s.requestRepo = mock.NewMockRequestRepository(s.ctrl)
	s.environmentRepo = mock.NewMockEnvironmentRepository(s.ctrl)

	s.useCase = NewUseCase(
		s.validator,
		s.apiKeyRepo,
		s.projectRepo,
		s.serviceRepo,
		s.requestRepo,
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
			Metadata: &dto.RequestIncomingMetadata{
				QueryParams:     `{"key": "value"}`,
				Headers:         `{"key": "value"}`,
				Body:            `{"key": "value"}`,
				BodyContentType: enums.RequestBodyContentTypeJSON,
			},
		},
	}

	wantRequestID := "req-id-123"

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

	s.validator.EXPECT().
		ValidateStruct(req, gomock.Any()).
		Return(nil).
		Times(1)

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

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusForwarded, r.ExecutionStatus)
			s.Require().Zero(r.UnauthorizedReason)
			s.Require().Equal(service.ID, r.Service.ID)
			s.Require().Equal(apiKey.ID, r.APIKey.ID)
			s.Require().Equal(environment.ID, r.Environment.ID)
			s.Require().Equal(projectClient.ProjectID, r.Project.ID)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, req.APIKey).
		Return(nil).
		Times(1)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.True(resp.Valid)
	s.Empty(resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
	s.Equal(projectClient.ProjectID, resp.Project.ID)
	s.Equal(projectClient.ProjectName, resp.Project.Name)
	s.Equal(projectClient.ClientID, resp.Client.ID)
	s.Equal(projectClient.ClientName, resp.Client.Name)
	s.Equal(environment.ID, resp.Environment.ID)
	s.Equal(environment.Name, resp.Environment.Name)
}

func (s *UseCaseSuite) TestSuccessUnauthorized() {
	reqTime := time.Now()
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

	wantRequestID := "req-id-123"

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

	s.validator.EXPECT().
		ValidateStruct(req, gomock.Any()).
		Return(nil).
		Times(1)

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

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			s.Require().Equal(enums.APIKeyValidationFailureCodeAPIKeyDisabled, r.UnauthorizedReason)
			s.Require().Equal(service.ID, r.Service.ID)
			s.Require().Equal(apiKey.ID, r.APIKey.ID)
			s.Require().Equal(environment.ID, r.Environment.ID)
			s.Require().Equal(projectClient.ProjectID, r.Project.ID)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, req.APIKey).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeAPIKeyDisabled, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
	s.Equal(projectClient.ProjectID, resp.Project.ID)
	s.Equal(projectClient.ProjectName, resp.Project.Name)
	s.Equal(projectClient.ClientID, resp.Client.ID)
	s.Equal(projectClient.ClientName, resp.Client.Name)
	s.Equal(environment.ID, resp.Environment.ID)
	s.Equal(environment.Name, resp.Environment.Name)
}

func (s *UseCaseSuite) TestValidateInternalError() {
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

	s.validator.EXPECT().
		ValidateStruct(req, gomock.Any()).
		Return(nil).
		Times(1)

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

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, req.APIKey).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().ErrorIs(err, internalErr)
	s.Require().Nil(resp)
}

func (s *UseCaseSuite) TestSuccessWithAPIKeyLastUsedErr() {
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

	wantRequestID := "req-id-123"

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

	s.validator.EXPECT().
		ValidateStruct(req, gomock.Any()).
		Return(nil).
		Times(1)

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

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusForwarded, r.ExecutionStatus)
			s.Require().Zero(r.UnauthorizedReason)
			s.Require().Equal(service.ID, r.Service.ID)
			s.Require().Equal(apiKey.ID, r.APIKey.ID)
			s.Require().Equal(environment.ID, r.Environment.ID)
			s.Require().Equal(projectClient.ProjectID, r.Project.ID)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, req.APIKey).
		Return(errors.NewInternal("Internal Error", nil)).
		Times(1)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.True(resp.Valid)
	s.Empty(resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
	s.Equal(projectClient.ProjectID, resp.Project.ID)
	s.Equal(projectClient.ProjectName, resp.Project.Name)
	s.Equal(projectClient.ClientID, resp.Client.ID)
	s.Equal(projectClient.ClientName, resp.Client.Name)
	s.Equal(environment.ID, resp.Environment.ID)
	s.Equal(environment.Name, resp.Environment.Name)
}

func (s *UseCaseSuite) TestValidationError() {
	req := &dto.APIKeyValidate{}

	validationErr := errors.NewValidationFailed("validation error", nil)
	s.validator.EXPECT().
		ValidateStruct(req, gomock.Any()).
		Return(validationErr).
		Times(1)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Equal(errors.CodeValidationFailed, err.Code())
	s.Equal(validationErr, err)
}

func (s *UseCaseSuite) TestRequestCreationError() {
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

	repoErr := errors.NewInternal("database error", nil)

	s.validator.EXPECT().
		ValidateStruct(req, gomock.Any()).
		Return(nil).
		Times(1)

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

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		Return(repoErr).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Equal(errors.CodeInternal, err.Code())
	s.Equal(repoErr, err)
}

func TestValidateOnlySuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}
