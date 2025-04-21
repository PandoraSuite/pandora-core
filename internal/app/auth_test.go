package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type AuthSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	tokenProvider   *mock.MockTokenPort
	credentialsRepo *mock.MockCredentialsPort

	useCase *AuthUseCase

	ctx context.Context
}

func (s *AuthSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.tokenProvider = mock.NewMockTokenPort(s.ctrl)
	s.credentialsRepo = mock.NewMockCredentialsPort(s.ctrl)

	s.useCase = NewAuthUseCase(s.tokenProvider, s.credentialsRepo)

	s.ctx = context.Background()
}

func (s *AuthSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
