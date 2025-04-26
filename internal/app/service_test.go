package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
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
			s.serviceRepo.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Times(0)

			resp, err := s.useCase.Create(s.ctx, test.req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ServiceSuite) TestCreate_ServiceRepoError() {
	req := &dto.ServiceCreate{Name: "Service", Version: "1.0.0"}

	s.serviceRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.Create(s.ctx, req)

	s.Require().Nil(resp)
	s.Error(err)
}

func (s *ServiceSuite) TestGetServices_Success() {
	now := time.Now().UTC()

	tests := []struct {
		name         string
		req          *dto.ServiceFilter
		mockServices []*entities.Service
	}{
		{
			name: "WithoutFilter",
			req:  &dto.ServiceFilter{},
			mockServices: []*entities.Service{
				{
					ID:        1,
					Name:      "Service 1",
					Status:    enums.ServiceActive,
					Version:   "1.0.0",
					CreatedAt: now.Add(-24 * time.Hour),
					UpdatedAt: now.Add(-24 * time.Hour),
				},
				{
					ID:        2,
					Name:      "Service 2",
					Status:    enums.ServiceDeactivated,
					Version:   "2.0.0",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
		{
			name: "WithFilter",
			req:  &dto.ServiceFilter{Status: enums.ServiceActive},
			mockServices: []*entities.Service{
				{
					ID:        1,
					Name:      "Service 1",
					Status:    enums.ServiceActive,
					Version:   "1.0.0",
					CreatedAt: now.Add(-24 * time.Hour),
					UpdatedAt: now.Add(-24 * time.Hour),
				},
				{
					ID:        2,
					Name:      "Service 2",
					Status:    enums.ServiceActive,
					Version:   "2.0.0",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.serviceRepo.EXPECT().
				FindAll(gomock.Any(), gomock.Any()).
				Return(test.mockServices, (*errors.Error)(nil)).
				Times(1)

			resp, err := s.useCase.GetServices(s.ctx, test.req)

			s.Require().Nil(err)
			s.Require().Len(resp, len(test.mockServices))

			for i, mockService := range test.mockServices {
				s.Equal(mockService.ID, resp[i].ID)
				s.Equal(mockService.Name, resp[i].Name)
				s.Equal(mockService.Status, resp[i].Status)
				s.Equal(mockService.Version, resp[i].Version)
				s.Equal(mockService.CreatedAt, resp[i].CreatedAt)
			}
		})
	}
}

func (s *ServiceSuite) TestGetServices_ServiceRepoError() {
	req := &dto.ServiceFilter{}

	s.serviceRepo.EXPECT().
		FindAll(gomock.Any(), gomock.Any()).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.GetServices(s.ctx, req)

	s.Require().Nil(resp)
	s.Error(err)
}

func (s *ServiceSuite) TestUpdateStatus_Success() {
	status := enums.ServiceDeprecated
	id := 42

	now := time.Now().UTC()
	mockService := &entities.Service{
		ID:        id,
		Name:      "Service",
		Status:    enums.ServiceActive,
		Version:   "1.0.0",
		CreatedAt: now.Add(-24 * time.Hour),
		UpdatedAt: now.Add(-24 * time.Hour),
	}

	s.serviceRepo.EXPECT().
		UpdateStatus(s.ctx, id, status).
		DoAndReturn(
			func(
				_ context.Context, id int, status enums.ServiceStatus,
			) (*entities.Service, *errors.Error) {
				mockService.Status = status
				mockService.UpdatedAt = now
				return mockService, nil
			},
		).
		Times(1)

	resp, err := s.useCase.UpdateStatus(s.ctx, id, status)

	s.Require().Nil(err)

	s.Equal(status, resp.Status)
	s.Equal(mockService.ID, resp.ID)
	s.Equal(mockService.Name, resp.Name)
	s.Equal(mockService.Version, resp.Version)
	s.Equal(mockService.CreatedAt, resp.CreatedAt)
}

func (s *ServiceSuite) TestUpdateStatus_ServiceRepoErrors() {
	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrServiceNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.serviceRepo.EXPECT().
				UpdateStatus(s.ctx, 42, enums.ServiceDeprecated).
				Return(nil, test.mockErr).
				Times(1)

			resp, err := s.useCase.UpdateStatus(
				s.ctx, 42, enums.ServiceDeprecated,
			)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ServiceSuite) TestDelete_Success() {
	id := 42

	gomock.InOrder(
		s.projectRepo.EXPECT().
			ExistsServiceIn(s.ctx, id).
			Return(false, (*errors.Error)(nil)).
			Times(1),

		s.serviceRepo.EXPECT().
			Delete(s.ctx, id).
			Return((*errors.Error)(nil)).
			Times(1),

		s.requestLogRepo.EXPECT().
			DeleteByService(s.ctx, id).
			Return((*errors.Error)(nil)).
			Times(1),
	)

	err := s.useCase.Delete(s.ctx, id)

	s.Require().Nil(err)
}

func (s *ServiceSuite) TestDelete_AssignedToProjects() {
	id := 42

	s.projectRepo.EXPECT().
		ExistsServiceIn(s.ctx, id).
		Return(true, (*errors.Error)(nil)).
		Times(1)

	s.serviceRepo.EXPECT().
		Delete(s.ctx, id).
		Times(0)

	s.requestLogRepo.EXPECT().
		DeleteByService(s.ctx, id).
		Times(0)

	err := s.useCase.Delete(s.ctx, id)

	s.Equal(errors.ErrServiceAssignedToProjects, err)
}

func (s *ServiceSuite) TestDelete_ProjectRepoError() {
	id := 42

	s.projectRepo.EXPECT().
		ExistsServiceIn(s.ctx, id).
		Return(false, errors.ErrPersistence).
		Times(1)

	s.serviceRepo.EXPECT().
		Delete(s.ctx, id).
		Times(0)

	s.requestLogRepo.EXPECT().
		DeleteByService(s.ctx, id).
		Times(0)

	err := s.useCase.Delete(s.ctx, id)

	s.Equal(errors.ErrPersistence, err)
}

func (s *ServiceSuite) TestDelete_ServiceRepoError() {
	id := 42

	s.projectRepo.EXPECT().
		ExistsServiceIn(s.ctx, id).
		Return(false, (*errors.Error)(nil)).
		Times(1)

	s.serviceRepo.EXPECT().
		Delete(s.ctx, id).
		Return(errors.ErrPersistence).
		Times(1)

	s.requestLogRepo.EXPECT().
		DeleteByService(s.ctx, id).
		Times(0)

	err := s.useCase.Delete(s.ctx, id)

	s.Equal(errors.ErrPersistence, err)
}

func (s *ServiceSuite) TestDelete_RequestLogRepoError() {
	id := 42

	s.projectRepo.EXPECT().
		ExistsServiceIn(s.ctx, id).
		Return(false, (*errors.Error)(nil)).
		Times(1)

	s.serviceRepo.EXPECT().
		Delete(s.ctx, id).
		Return((*errors.Error)(nil)).
		Times(1)

	s.requestLogRepo.EXPECT().
		DeleteByService(s.ctx, id).
		Return(errors.ErrPersistence).
		Times(1)

	err := s.useCase.Delete(s.ctx, id)

	s.Equal(errors.ErrPersistence, err)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
