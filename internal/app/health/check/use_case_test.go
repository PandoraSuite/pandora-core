package check

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/app/health/check/mock"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Suite struct {
	suite.Suite

	ctrl *gomock.Controller

	database *mock.MockDatabase

	useCase UseCase

	ctx context.Context
}

func (s *Suite) SetupTest() {
	time.Local = time.UTC

	s.ctrl = gomock.NewController(s.T())

	s.database = mock.NewMockDatabase(s.ctrl)

	s.useCase = NewUseCase(s.database)

	s.ctx = context.Background()
}

func (s *Suite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *Suite) TestSuccess() {
	latency := int64(32)

	s.database.EXPECT().
		Ping().
		Return(nil).
		Times(1)

	s.database.EXPECT().
		Latency().
		Return(latency, nil).
		Times(1)

	response := s.useCase.Execute()

	s.Equal(enums.HealthStatusOK, response.Status)
	s.NotNil(response.Timestamp)

	s.Require().NotNil(response.Check)
	s.Require().NotNil(response.Check.Database)

	s.Equal(enums.HealthStatusOK, response.Check.Database.Status)
	s.Equal("database is reachable", response.Check.Database.Message)
	s.Equal(latency, response.Check.Database.Latency)
}

func (s *Suite) TestDatabaseDown() {
	pingErr := errors.NewInternal("failed to ping database", nil)
	s.database.EXPECT().
		Ping().
		Return(pingErr).
		Times(1)

	s.database.EXPECT().
		Latency().
		Times(0)

	response := s.useCase.Execute()

	s.Equal(enums.HealthStatusDown, response.Status)
	s.NotNil(response.Timestamp)

	s.Require().NotNil(response.Check)
	s.Require().NotNil(response.Check.Database)

	s.Equal(enums.HealthStatusDown, response.Check.Database.Status)
	s.Equal(pingErr.Error(), response.Check.Database.Message)
	s.Equal(int64(0), response.Check.Database.Latency)
}

func (s *Suite) TestDatabaseDegraded() {
	latencyErr := errors.NewInternal("failed to measure DB latency", nil)

	s.database.EXPECT().
		Ping().
		Return(nil).
		Times(1)

	s.database.EXPECT().
		Latency().
		Return(int64(0), latencyErr).
		Times(1)

	response := s.useCase.Execute()

	s.Equal(enums.HealthStatusDegraded, response.Status)
	s.NotNil(response.Timestamp)

	s.Require().NotNil(response.Check)
	s.Require().NotNil(response.Check.Database)

	s.Equal(enums.HealthStatusDegraded, response.Check.Database.Status)
	s.Equal(latencyErr.Error(), response.Check.Database.Message)
	s.Equal(int64(0), response.Check.Database.Latency)
}

func TestServiceCreateSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
