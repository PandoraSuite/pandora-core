package app

import (
	"context"
	"testing"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ServiceSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	serviceRepo    *mock.MockServicePort
	projectRepo    *mock.MockProjectPort
	requestLogRepo *mock.MockRequestLogPort

	useCase *ServiceUseCase

	ctx context.Context
}

func (s *ServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.serviceRepo = mock.NewMockServicePort(s.ctrl)
	s.projectRepo = mock.NewMockProjectPort(s.ctrl)
	s.requestLogRepo = mock.NewMockRequestLogPort(s.ctrl)

	s.useCase = NewServiceUseCase(
		s.serviceRepo, s.projectRepo, s.requestLogRepo,
	)

	s.ctx = context.Background()
}

func (s *ServiceSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ServiceSuite) TestCreate_Success() {
	req := &dto.ServiceCreate{Name: "Service", Version: "1.0.0"}

	now := time.Now().UTC()

	s.serviceRepo.EXPECT().
		Save(s.ctx, gomock.AssignableToTypeOf(&entities.Service{})).
		DoAndReturn(
			func(_ context.Context, service *entities.Service) *errors.Error {
				service.ID = 42
				service.CreatedAt = now
				service.UpdatedAt = now
				return nil
			},
		).
		Times(1)

	resp, err := s.useCase.Create(s.ctx, req)

	s.Require().Nil(err)

	s.Equal(42, resp.ID)
	s.Equal(req.Name, resp.Name)
	s.Equal(now, resp.CreatedAt)
	s.Equal(req.Version, resp.Version)
	s.Equal(enums.ServiceActive, resp.Status)
}

func (s *ServiceSuite) TestCreate_ValidationErrors() {
	tests := []struct {
		name        string
		req         *dto.ServiceCreate
		expectedErr *errors.Error
	}{
		{
			name:        "EmptyName",
			req:         &dto.ServiceCreate{Name: "", Version: "1.0.0"},
			expectedErr: errors.ErrNameCannotBeEmpty,
		},
		{
			name:        "EmptyVersion",
			req:         &dto.ServiceCreate{Name: "Service", Version: ""},
			expectedErr: errors.ErrVersionCannotBeEmpty,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			ctrl := gomock.NewController(s.T())

			serviceRepo := mock.NewMockServicePort(s.ctrl)
			projectRepo := mock.NewMockProjectPort(s.ctrl)
			requestLogRepo := mock.NewMockRequestLogPort(s.ctrl)

			uc := NewServiceUseCase(
				serviceRepo, projectRepo, requestLogRepo,
			)

			serviceRepo.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Times(0)

			resp, err := uc.Create(s.ctx, test.req)

			s.Nil(resp)
			s.Equal(test.expectedErr, err)

			ctrl.Finish()
		})
	}
}

func (s *ServiceSuite) TestCreate_RepoError() {
	req := &dto.ServiceCreate{Name: "Service", Version: "1.0.0"}

	s.serviceRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(errors.ErrPersistence)

	resp, err := s.useCase.Create(s.ctx, req)

	s.Nil(resp)
	s.Error(err)
}

func (s *ServiceSuite) TestGetServices_Success() {
	tests := []struct {
		name string
		req  *dto.ServiceFilter
	}{
		{
			name: "WithoutFilter",
			req:  &dto.ServiceFilter{},
		},
		{
			name: "WithFilter",
			req:  &dto.ServiceFilter{Status: enums.ServiceActive},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			ctrl := gomock.NewController(s.T())

			serviceRepo := mock.NewMockServicePort(s.ctrl)
			projectRepo := mock.NewMockProjectPort(s.ctrl)
			requestLogRepo := mock.NewMockRequestLogPort(s.ctrl)

			uc := NewServiceUseCase(
				serviceRepo, projectRepo, requestLogRepo,
			)

			serviceRepo.EXPECT().
				FindAll(s.ctx, test.req).
				Return().
				Times(1)

			resp, err := uc.GetServices(s.ctx, test.req)

			s.Nil(resp)
			s.Equal(test.expectedErr, err)

			ctrl.Finish()
		})
	}
}

func TestCreateServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
