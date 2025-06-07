package delete

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/app/environment/delete/mock"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	mockvalidator "github.com/MAD-py/pandora-core/internal/validator/mock"
)

type Suite struct {
	suite.Suite

	ctrl *gomock.Controller

	validator       *mockvalidator.MockValidator
	environmentRepo *mock.MockEnvironmentRepository

	useCase UseCase

	ctx context.Context
}

func (s *Suite) SetupTest() {
	time.Local = time.UTC

	s.ctrl = gomock.NewController(s.T())

	s.environmentRepo = mock.NewMockEnvironmentRepository(s.ctrl)
	s.validator = mockvalidator.NewMockValidator(s.ctrl)

	s.useCase = NewUseCase(s.validator, s.environmentRepo)

	s.ctx = context.Background()
}

func (s *Suite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *Suite) TestSuccess() {
	id := 42

	s.validator.EXPECT().
		ValidateVariable(id, "id", gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	s.environmentRepo.EXPECT().
		Delete(s.ctx, id).
		Return(nil).
		Times(1)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().NoError(err)
}

func (s *Suite) TestValidationError() {
	id := -1

	validationErr := errors.NewValidationFailed("Validation Error", nil)
	s.validator.EXPECT().
		ValidateVariable(id, "id", gomock.Any(), gomock.Any()).
		Return(validationErr).
		Times(1)

	s.environmentRepo.EXPECT().
		Delete(s.ctx, id).
		Times(0)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().Error(err)

	s.Equal(errors.CodeValidationFailed, err.Code())
	s.Equal(validationErr, err)
}

func (s *Suite) TestRepositoryError() {
	id := 42

	s.validator.EXPECT().
		ValidateVariable(id, "id", gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	repositoryErr := errors.NewInternal("Repository Error", nil)
	s.environmentRepo.EXPECT().
		Delete(s.ctx, id).
		Return(repositoryErr).
		Times(1)

	err := s.useCase.Execute(s.ctx, id)

	s.Require().Error(err)

	s.Equal(errors.CodeInternal, err.Code())
	s.Equal(repositoryErr, err)
}

func TestServiceCreateSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
