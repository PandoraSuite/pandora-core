package validateonly

import (
	"context"
	"testing"
	"time"

	"github.com/MAD-py/pandora-core/internal/app/api_key/validate_only/mock"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	mockvalidator "github.com/MAD-py/pandora-core/internal/validator/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
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
	now := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: now,
			Metadata: &dto.RequestMetadata{
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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(true, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusForwarded, r.ExecutionStatus)
			s.Require().Equal(service.ID, r.ServiceID)
			s.Require().Equal(apiKey.ID, r.APIKeyID)
			s.Require().Equal(environment.ID, r.EnvironmentID)
			s.Require().Equal(projectCtx.ID, r.ProjectID)
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
	s.Equal(projectCtx, resp.ConsumerInfo)
}

func (s *UseCaseSuite) TestSuccessWithAPIKeyLastUsedErr() {
	now := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: now,
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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(true, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusForwarded, r.ExecutionStatus)
			s.Require().Equal(service.ID, r.ServiceID)
			s.Require().Equal(apiKey.ID, r.APIKeyID)
			s.Require().Equal(environment.ID, r.EnvironmentID)
			s.Require().Equal(projectCtx.ID, r.ProjectID)
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
	s.Equal(projectCtx, resp.ConsumerInfo)
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

func (s *UseCaseSuite) TestAPIKeyNotFound() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "invalid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
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
		Return(nil, errors.NewEntityNotFound(
			"APIKey",
			"api key not found",
			map[string]any{"key": req.APIKey},
			nil,
		)).
		Times(1)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			s.Require().Equal(0, r.APIKeyID)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeAPIKeyInvalid, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func (s *UseCaseSuite) TestServiceMismatch() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "NonExistentService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	wantRequestID := "req-id-789"

	apiKey := &entities.APIKey{
		ID:            10,
		Key:           req.APIKey,
		Status:        enums.APIKeyStatusEnabled,
		EnvironmentID: 100,
	}
	environment := &entities.Environment{
		ID:        100,
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

	s.validator.EXPECT().
		ValidateStruct(req, gomock.Any()).
		Return(nil).
		Times(1)

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			s.Require().Equal(0, r.ServiceID)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeServiceMismatch, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func (s *UseCaseSuite) TestServiceRepoInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
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
		GetProjectContextByID(s.ctx, gomock.Any()).
		Times(0)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.Any()).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Require().Equal(internalErr, err)
	s.Require().Equal(errors.CodeInternal, err.Code())
}

func (s *UseCaseSuite) TestAPIKeyRepoInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
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
		Return(nil, internalErr).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, gomock.Any()).
		Times(0)

	s.projectRepo.EXPECT().
		GetProjectContextByID(s.ctx, gomock.Any()).
		Times(0)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.Any()).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Require().Equal(internalErr, err)
	s.Require().Equal(errors.CodeInternal, err.Code())
}

func (s *UseCaseSuite) TestEnvironmentRepoGetByIDInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
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
		Return(nil, internalErr).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectContextByID(s.ctx, gomock.Any()).
		Times(0)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.Any()).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Require().Equal(internalErr, err)
	s.Require().Equal(errors.CodeInternal, err.Code())
}

func (s *UseCaseSuite) TestProjectRepoInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	internalErr := errors.NewInternal("project repo error", nil)

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(nil, internalErr).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.Any()).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Require().Equal(internalErr, err)
	s.Require().Equal(errors.CodeInternal, err.Code())
}

func (s *UseCaseSuite) TestEnvironmentRepoExistsServiceInInternalError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}
	internalErr := errors.NewInternal("environment repo exists service in error", nil)

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(false, internalErr).
		Times(1)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.Any()).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Require().Equal(internalErr, err)
	s.Require().Equal(errors.CodeInternal, err.Code())
}

func (s *UseCaseSuite) TestEnvironmentMismatch() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "staging",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	wantRequestID := "req-id-101"

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
		GetProjectContextByID(s.ctx, gomock.Any()).
		Times(0)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeEnvironmentMismatch, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func (s *UseCaseSuite) TestEnvironmentDisabled() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	wantRequestID := "req-id-101"

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
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(true, nil).
		Times(1)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeEnvironmentDisabled, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func (s *UseCaseSuite) TestServiceNotAssignedToEnvironment() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	wantRequestID := "req-id-112"

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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(false, nil).
		Times(1)

	s.requestRepo.EXPECT().Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeServiceNotAssigned, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func (s *UseCaseSuite) TestRequestCreationError() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(true, nil).
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

func (s *UseCaseSuite) TestServiceDisabled() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "DisabledService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	wantRequestID := "req-id-disabled-service"

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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(true, nil).
		Times(1)

	s.requestRepo.EXPECT().Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeServiceDisabled, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func (s *UseCaseSuite) TestServiceDeprecated() {
	reqTime := time.Now()
	req := &dto.APIKeyValidate{
		APIKey:          "valid-api-key",
		ServiceName:     "DeprecatedService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	wantRequestID := "req-id-deprecated-service"

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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(true, nil).
		Times(1)

	s.requestRepo.EXPECT().Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeServiceDeprecated, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func (s *UseCaseSuite) TestAPIKeyExpired() {
	reqTime := time.Now()
	pastTime := time.Now().Add(-24 * time.Hour)
	req := &dto.APIKeyValidate{
		APIKey:          "expired-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	wantRequestID := "req-id-expired-key"

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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(true, nil).
		Times(1)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeAPIKeyExpired, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func (s *UseCaseSuite) TestAPIKeyDisabled() {
	reqTime := time.Now()
	pastTime := time.Now().Add(-24 * time.Hour)
	req := &dto.APIKeyValidate{
		APIKey:          "disabled-api-key",
		ServiceName:     "TestService",
		ServiceVersion:  "1.0.0",
		EnvironmentName: "production",
		Request: &dto.RequestCreate{
			Path:        "/test",
			Method:      "GET",
			IPAddress:   "127.0.0.1",
			RequestTime: reqTime,
		},
	}

	wantRequestID := "req-id-disabled-key"

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
		Name:      req.EnvironmentName,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: 1000,
	}
	projectCtx := &dto.ProjectContextResponse{ID: 1000, Name: "TestProject"}

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
		GetProjectContextByID(s.ctx, environment.ProjectID).
		Return(projectCtx, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, environment.ID, service.ID).
		Return(true, nil).
		Times(1)

	s.requestRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Request{})).
		DoAndReturn(func(_ context.Context, r *entities.Request) errors.Error {
			s.Require().Equal(enums.RequestExecutionStatusUnauthorized, r.ExecutionStatus)
			r.ID = wantRequestID
			return nil
		}).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateLastUsed(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.False(resp.Valid)
	s.Equal(enums.APIKeyValidationFailureCodeAPIKeyDisabled, resp.FailureCode)
	s.Equal(wantRequestID, resp.RequestID)
}

func TestValidateOnlySuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}
