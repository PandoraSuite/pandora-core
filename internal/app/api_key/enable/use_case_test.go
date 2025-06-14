package enable

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/app/api_key/enable/mock"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	mockvalidator "github.com/MAD-py/pandora-core/internal/validator/mock"
)

type Suite struct {
	suite.Suite

	ctrl *gomock.Controller

	validator       *mockvalidator.MockValidator
	apiKeyRepo      *mock.MockAPIKeyRepository
	environmentRepo *mock.MockEnvironmentRepository

	useCase UseCase

	ctx context.Context
}

func (s *Suite) SetupTest() {
	time.Local = time.UTC

	s.ctrl = gomock.NewController(s.T())

	s.apiKeyRepo = mock.NewMockAPIKeyRepository(s.ctrl)
	s.environmentRepo = mock.NewMockEnvironmentRepository(s.ctrl)
	s.validator = mockvalidator.NewMockValidator(s.ctrl)

	s.useCase = NewUseCase(s.validator, s.apiKeyRepo, s.environmentRepo)

	s.ctx = context.Background()
}

func (s *Suite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *Suite) TestSuccess() {
	id := 42
	environmentID := 42

	s.validator.EXPECT().
		ValidateVariable(
			id,
			"id",
			"required,gt=0",
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:            id,
				Status:        enums.APIKeyStatusDisabled,
				ExpiresAt:     time.Now().Add(24 * time.Hour),
				EnvironmentID: environmentID,
			},
			nil,
		).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, environmentID).
		Return(
			&entities.Environment{
				ID:        environmentID,
				Name:      "Test Environment",
				Status:    enums.EnvironmentStatusEnabled,
				ProjectID: 1,
				CreatedAt: time.Now(),
			},
			nil,
		).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusEnabled).
		Return(nil).
		Times(1)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().NoError(err)
}

func (s *Suite) TestAlreadyEnabled() {
	id := 42
	environmentID := 42

	s.validator.EXPECT().
		ValidateVariable(
			id,
			"id",
			"required,gt=0",
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:            id,
				Status:        enums.APIKeyStatusEnabled,
				EnvironmentID: environmentID,
			},
			nil,
		).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, environmentID).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
		Times(0)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().NoError(err)
}

func (s *Suite) TestExpired() {
	id := 42
	environmentID := 42

	s.validator.EXPECT().
		ValidateVariable(
			id,
			"id",
			"required,gt=0",
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:        id,
				Status:    enums.APIKeyStatusDisabled,
				ExpiresAt: time.Now().Add(-24 * time.Hour),
			},
			nil,
		).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, environmentID).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
		Times(0)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().Error(err)

	s.Equal(errors.CodeValidationFailed, err.Code())

	var e *errors.EntityError
	s.Require().ErrorAs(err, &e)

	s.Equal("APIKey", e.Entity())
	s.Equal("API key is expired and cannot be enabled", e.Message())
	s.Equal(map[string]any{"id": id}, e.Identifiers())
}

func (s *Suite) TestEnvironmentDisabled() {
	id := 42
	environmentID := 42

	s.validator.EXPECT().
		ValidateVariable(
			id,
			"id",
			"required,gt=0",
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:            id,
				Status:        enums.APIKeyStatusDisabled,
				EnvironmentID: environmentID,
			},
			nil,
		).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, environmentID).
		Return(
			&entities.Environment{
				ID:     environmentID,
				Status: enums.EnvironmentStatusDisabled,
			},
			nil,
		).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
		Times(0)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().Error(err)

	s.Equal(errors.CodeValidationFailed, err.Code())

	var e *errors.EntityError
	s.Require().ErrorAs(err, &e)

	s.Equal("Environment", e.Entity())
	s.Equal("Environment is disabled and cannot enable API key for it", e.Message())
	s.Equal(map[string]any{"id": environmentID}, e.Identifiers())
}

func (s *Suite) TestValidationError() {
	id := 0

	validationErr := errors.NewValidationFailed("Validation Error", nil)
	s.validator.EXPECT().
		ValidateVariable(
			id,
			"id",
			"required,gt=0",
			gomock.Any(),
		).
		Return(validationErr).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, gomock.Any()).
		Times(0)

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
		Times(0)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().Error(err)

	s.Equal(errors.CodeValidationFailed, err.Code())
	s.Equal(validationErr, err)
}

func (s *Suite) TestRepositoryGetByIDError() {
	id := 42

	s.validator.EXPECT().
		ValidateVariable(
			id,
			"id",
			"required,gt=0",
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	expectedErr := errors.NewInternal("Database error", nil)
	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(nil, expectedErr).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, gomock.Any()).
		Times(0)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
		Times(0)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().Error(err)

	s.Equal(errors.CodeInternal, err.Code())
	s.Equal(expectedErr, err)
}

func (s *Suite) TestRepositoryEnvGetByIDError() {
	id := 42
	environmentID := 42

	s.validator.EXPECT().
		ValidateVariable(
			id,
			"id",
			"required,gt=0",
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:            id,
				Status:        enums.APIKeyStatusDisabled,
				EnvironmentID: environmentID,
			},
			nil,
		).
		Times(1)

	expectedErr := errors.NewInternal("Database error", nil)
	s.environmentRepo.EXPECT().
		GetByID(s.ctx, environmentID).
		Return(nil, expectedErr).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
		Times(0)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().Error(err)

	s.Equal(errors.CodeInternal, err.Code())
	s.Equal(expectedErr, err)
}

func (s *Suite) TestRepositoryUpdateStatusError() {
	id := 42
	environmentID := 42

	s.validator.EXPECT().
		ValidateVariable(
			id,
			"id",
			"required,gt=0",
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:            id,
				Status:        enums.APIKeyStatusDisabled,
				EnvironmentID: environmentID,
			},
			nil,
		).
		Times(1)

	s.environmentRepo.EXPECT().
		GetByID(s.ctx, environmentID).
		Return(
			&entities.Environment{
				ID:     environmentID,
				Status: enums.EnvironmentStatusEnabled,
			},
			nil,
		).
		Times(1)

	expectedErr := errors.NewInternal("Update error", nil)
	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusEnabled).
		Return(expectedErr).
		Times(1)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().Error(err)

	s.Equal(errors.CodeInternal, err.Code())
	s.Equal(expectedErr, err)
}

func TestServiceCreateSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
