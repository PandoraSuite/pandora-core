package disabled

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/app/api_key/disabled/mock"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	mockvalidator "github.com/MAD-py/pandora-core/internal/validator/mock"
)

type Suite struct {
	suite.Suite

	ctrl *gomock.Controller

	validator  *mockvalidator.MockValidator
	apiKeyRepo *mock.MockAPIKeyRepository

	useCase UseCase

	ctx context.Context
}

func (s *Suite) SetupTest() {
	time.Local = time.UTC

	s.ctrl = gomock.NewController(s.T())

	s.apiKeyRepo = mock.NewMockAPIKeyRepository(s.ctrl)
	s.validator = mockvalidator.NewMockValidator(s.ctrl)

	s.useCase = NewUseCase(s.validator, s.apiKeyRepo)

	s.ctx = context.Background()
}

func (s *Suite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *Suite) TestSuccess() {
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

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:        id,
				Status:    enums.APIKeyStatusEnabled,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			nil,
		).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
		Return(nil).
		Times(1)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().NoError(err)
}

func (s *Suite) TestAlreadyDisabled() {
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

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:     id,
				Status: enums.APIKeyStatusDisabled,
			},
			nil,
		).
		Times(1)

	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
		Times(0)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().NoError(err)
}

func (s *Suite) TestExpired() {
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

	s.apiKeyRepo.EXPECT().
		GetByID(s.ctx, id).
		Return(
			&entities.APIKey{
				ID:        id,
				Status:    enums.APIKeyStatusEnabled,
				ExpiresAt: time.Now().Add(-24 * time.Hour),
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

	s.Equal("APIKey", e.Entity())
	s.Equal("API key is expired and cannot be disabled", e.Message())
	s.Equal(map[string]any{"id": id}, e.Identifiers())
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
				Status:    enums.APIKeyStatusEnabled,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			nil,
		).
		Times(1)

	expectedErr := errors.NewInternal("Database error", nil)
	s.apiKeyRepo.EXPECT().
		UpdateStatus(s.ctx, id, enums.APIKeyStatusDisabled).
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
