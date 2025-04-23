package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ProjectSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	projectRepo     *mock.MockProjectPort
	environmentRepo *mock.MockEnvironmentPort

	useCase *ProjectUseCase

	ctx context.Context
}

func (s *ProjectSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.projectRepo = mock.NewMockProjectPort(s.ctrl)
	s.environmentRepo = mock.NewMockEnvironmentPort(s.ctrl)

	s.useCase = NewProjectUseCase(s.projectRepo, s.environmentRepo)

	s.ctx = context.Background()
}

func (s *ProjectSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}
