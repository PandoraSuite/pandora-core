package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RequestLogSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	requestLogRepo *mock.MockRequestLogPort

	useCase *RequestLogUseCase

	ctx context.Context
}

func (s *RequestLogSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.requestLogRepo = mock.NewMockRequestLogPort(s.ctrl)

	s.useCase = NewRequestLogUseCase(s.requestLogRepo)

	s.ctx = context.Background()
}

func (s *RequestLogSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *RequestLogSuite) TestUpdateExecutionStatus_Success() {
	id := "test-id"

	s.requestLogRepo.EXPECT().
		UpdateExecutionStatus(s.ctx, id, enums.RequestLogSuccess).
		Return(nil)

	err := s.useCase.UpdateExecutionStatus(
		s.ctx, id, enums.RequestLogSuccess,
	)

	s.Require().Nil(err)
}

func (s *RequestLogSuite) TestUpdateExecutionStatus_ValidationErrors() {
	id := "test-id"

	tests := []struct {
		name            string
		executionStatus enums.RequestLogExecutionStatus
		expectedError   *errors.Error
	}{
		{
			name:            "NullExecutionStatus",
			executionStatus: enums.RequestLogExecutionStatusNull,
			expectedError:   errors.ErrCannotUpdateToNullExecutionStatus,
		},
		{
			name:            "PendingExecutionStatus",
			executionStatus: enums.RequestLogPending,
			expectedError:   errors.ErrCannotUpdateToPendingExecutionStatus,
		},
		{
			name:            "UnauthorizedExecutionStatus",
			executionStatus: enums.RequestLogUnauthorized,
			expectedError:   errors.ErrCannotUpdateToUnauthorizedExecutionStatus,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := s.useCase.UpdateExecutionStatus(
				s.ctx, id, test.executionStatus,
			)

			s.Equal(test.expectedError, err)
		})
	}
}

func (s *RequestLogSuite) TestUpdateExecutionStatus_RequestLogRepoError() {
	id := "test-id"

	s.requestLogRepo.EXPECT().
		UpdateExecutionStatus(s.ctx, id, enums.RequestLogSuccess).
		Return(errors.ErrPersistence).
		Times(1)

	err := s.useCase.UpdateExecutionStatus(
		s.ctx, id, enums.RequestLogSuccess,
	)

	s.Equal(errors.ErrPersistence, err)
}

func TestRequestLogSuite(t *testing.T) {
	suite.Run(t, new(RequestLogSuite))
}
