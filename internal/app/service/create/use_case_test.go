package create

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/app/service/create/mock"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	mockvalidator "github.com/MAD-py/pandora-core/internal/validator/mock"
)

type Suite struct {
	suite.Suite

	ctrl *gomock.Controller

	validator   *mockvalidator.MockValidator
	serviceRepo *mock.MockServiceRepository

	useCase UseCase

	ctx context.Context
}

func (s *Suite) SetupTest() {
	time.Local = time.UTC

	s.ctrl = gomock.NewController(s.T())

	s.serviceRepo = mock.NewMockServiceRepository(s.ctrl)
	s.validator = mockvalidator.NewMockValidator(s.ctrl)

	s.useCase = NewUseCase(s.validator, s.serviceRepo)

	s.ctx = context.Background()
}

func (s *Suite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *Suite) TestSuccess() {
	req := dto.ServiceCreate{Name: "Service", Version: "1.0.0"}

	now := time.Now()

	s.validator.EXPECT().
		ValidateStruct(&req, gomock.Any()).
		Return(nil).
		Times(1)

	s.serviceRepo.EXPECT().
		Create(s.ctx, gomock.Any()).
		DoAndReturn(
			func(_ context.Context, service *entities.Service) errors.Error {
				service.ID = 42
				service.CreatedAt = now
				return nil
			},
		).
		Times(1)

	resp, err := s.useCase.Execute(s.ctx, &req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Equal(42, resp.ID)
	s.Equal(req.Name, resp.Name)
	s.Equal(enums.ServiceStatusEnabled, resp.Status)
	s.Equal(req.Version, resp.Version)
	s.Equal(now, resp.CreatedAt)
}

func (s *Suite) TestAlreadyExists() {
	req := dto.ServiceCreate{Name: "Service", Version: "1.0.0"}

	s.validator.EXPECT().
		ValidateStruct(&req, gomock.Any()).
		Return(nil).
		Times(1)

	s.serviceRepo.EXPECT().
		Create(s.ctx, gomock.Any()).
		Return(errors.NewEntityAlreadyExists(
			"Service", "", map[string]any{}, nil,
		)).
		Times(1)

	resp, err := s.useCase.Execute(s.ctx, &req)

	s.Require().Nil(resp)
	s.Require().Error(err)

	s.Equal(errors.CodeAlreadyExists, err.Code())

	var e *errors.EntityError
	s.Require().ErrorAs(err, &e)

	s.Equal("Service", e.Entity())
	s.Equal("Service with this name and version already exists", e.Message())
	s.Equal(map[string]any{"name": req.Name, "version": req.Version}, e.Identifiers())
}

func (s *Suite) TestValidationError() {
	req := dto.ServiceCreate{Name: "", Version: ""}

	validationErr := errors.NewValidationFailed("Validation Error", nil)
	s.validator.EXPECT().
		ValidateStruct(&req, gomock.Any()).
		Return(validationErr).
		Times(1)

	s.serviceRepo.EXPECT().
		Create(s.ctx, gomock.Any()).
		Times(0)

	resp, err := s.useCase.Execute(s.ctx, &req)

	s.Require().Nil(resp)
	s.Require().Error(err)

	s.Equal(errors.CodeValidationFailed, err.Code())
	s.Equal(validationErr, err)
}

func (s *Suite) TestRepositoryError() {
	req := dto.ServiceCreate{Name: "Service", Version: "1.0.0"}

	s.validator.EXPECT().
		ValidateStruct(&req, gomock.Any()).
		Return(nil).
		Times(1)

	repositoryErr := errors.NewInternal("Repository Error", nil)
	s.serviceRepo.EXPECT().
		Create(s.ctx, gomock.AssignableToTypeOf(&entities.Service{})).
		Return(repositoryErr).
		Times(1)

	resp, err := s.useCase.Execute(s.ctx, &req)

	s.Require().Nil(resp)
	s.Require().Error(err)

	s.Equal(errors.CodeInternal, err.Code())
	s.Equal(repositoryErr, err)
}

func TestServiceCreateSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
