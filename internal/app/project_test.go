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

type ProjectSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	projectRepo     *mock.MockProjectPort
	environmentRepo *mock.MockEnvironmentPort

	useCase *ProjectUseCase

	ctx context.Context
}

func (s *ProjectSuite) SetupTest() {
	time.Local = time.UTC

	s.ctrl = gomock.NewController(s.T())

	s.projectRepo = mock.NewMockProjectPort(s.ctrl)
	s.environmentRepo = mock.NewMockEnvironmentPort(s.ctrl)

	s.useCase = NewProjectUseCase(s.projectRepo, s.environmentRepo)

	s.ctx = context.Background()
}

func (s *ProjectSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ProjectSuite) TestResetServiceAvailableRequests_Successes() {
	id := 1
	serviceID := 2

	now := time.Now()

	mockService := &entities.ProjectService{
		ID:             serviceID,
		Name:           "Service",
		Version:        "1.0.0",
		NextReset:      now.Add(24 * time.Hour),
		MaxRequest:     100,
		ResetFrequency: enums.ProjectServiceDaily,
		AssignedAt:     now.Add(-24 * time.Hour),
	}

	mockEnvServices := []*dto.EnvironmentServiceReset{
		{
			ID:     1,
			Name:   "Environment 1",
			Status: enums.EnvironmentActive,
			Service: &dto.EnvironmentServiceResponse{
				ID:               serviceID,
				Name:             "Service 1",
				Version:          "1.0.0",
				MaxRequest:       20,
				AvailableRequest: 20,
				AssignedAt:       now.Add(-24 * time.Hour),
			},
		},
		{
			ID:     2,
			Name:   "Environment 2",
			Status: enums.EnvironmentActive,
			Service: &dto.EnvironmentServiceResponse{
				ID:               serviceID,
				Name:             "Service 2",
				Version:          "1.0.0",
				MaxRequest:       30,
				AvailableRequest: 30,
				AssignedAt:       now.Add(-24 * time.Hour),
			},
		},
	}

	tests := []struct {
		name string
		req  *dto.ProjectServiceResetRequest
	}{
		{
			name: "RecalculateNextReset",
			req: &dto.ProjectServiceResetRequest{
				RecalculateNextReset: true,
			},
		},
		{
			name: "WithoutRecalculateNextReset",
			req: &dto.ProjectServiceResetRequest{
				RecalculateNextReset: false,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil)

			s.projectRepo.EXPECT().
				FindServiceByID(s.ctx, id, serviceID).
				Return(mockService, nil)

			if test.req.RecalculateNextReset {
				s.projectRepo.EXPECT().
					ResetProjectServiceUsage(
						s.ctx, id, serviceID, gomock.Any(),
					).
					Return(mockEnvServices, nil).
					Times(1)
			} else {
				s.projectRepo.EXPECT().
					ResetAvailableRequestsForEnvsService(
						s.ctx, id, serviceID,
					).
					Return(mockEnvServices, nil).
					Times(1)
			}

			resp, err := s.useCase.ResetServiceAvailableRequests(
				s.ctx, id, serviceID, test.req,
			)

			s.Require().Nil(err)

			s.Equal(resp.ProjectService.ID, mockService.ID)
			s.Equal(resp.ProjectService.Name, mockService.Name)
			s.Equal(resp.ProjectService.Version, mockService.Version)
			s.Equal(resp.ProjectService.MaxRequest, mockService.MaxRequest)
			s.Equal(resp.ProjectService.ResetFrequency, mockService.ResetFrequency)
			s.Equal(resp.ProjectService.NextReset, mockService.NextReset)
			s.Equal(resp.ProjectService.AssignedAt, mockService.AssignedAt)

			s.Equal(resp.ResetCount, len(mockEnvServices))

			for i, envService := range resp.EnvironmentServices {
				s.Equal(envService.ID, mockEnvServices[i].ID)
				s.Equal(envService.Name, mockEnvServices[i].Name)
				s.Equal(envService.Status, mockEnvServices[i].Status)
				s.Equal(envService.Service.ID, mockEnvServices[i].Service.ID)
				s.Equal(envService.Service.Name, mockEnvServices[i].Service.Name)
				s.Equal(envService.Service.Version, mockEnvServices[i].Service.Version)
				s.Equal(envService.Service.MaxRequest, mockEnvServices[i].Service.MaxRequest)
				s.Equal(envService.Service.AvailableRequest, mockEnvServices[i].Service.AvailableRequest)
				s.Equal(envService.Service.AssignedAt, mockEnvServices[i].Service.AssignedAt)
			}
		})
	}
}

func (s *ProjectSuite) TestResetServiceAvailableRequests_ExistsErrors() {
	id := 1
	serviceID := 2

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "DoesNotExist",
			mockErr:     nil,
			expectedErr: errors.ErrProjectNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(false, test.mockErr).
				Times(1)

			s.projectRepo.EXPECT().
				FindServiceByID(s.ctx, id, serviceID).
				Times(0)

			s.projectRepo.EXPECT().
				ResetProjectServiceUsage(
					s.ctx, id, serviceID, gomock.Any(),
				).
				Times(0)

			s.projectRepo.EXPECT().
				ResetAvailableRequestsForEnvsService(
					s.ctx, id, serviceID,
				).
				Times(0)

			resp, err := s.useCase.ResetServiceAvailableRequests(
				s.ctx, id, serviceID, &dto.ProjectServiceResetRequest{},
			)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestResetServiceAvailableRequests_FindServiceErrors() {
	id := 1
	serviceID := 2

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ServiceNotAssignedToProject",
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrServiceNotAssignedToProject,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.projectRepo.EXPECT().
				FindServiceByID(s.ctx, id, serviceID).
				Return(nil, test.mockErr).
				Times(1)

			s.projectRepo.EXPECT().
				ResetProjectServiceUsage(
					s.ctx, id, serviceID, gomock.Any(),
				).
				Times(0)

			s.projectRepo.EXPECT().
				ResetAvailableRequestsForEnvsService(
					s.ctx, id, serviceID,
				).
				Times(0)

			resp, err := s.useCase.ResetServiceAvailableRequests(
				s.ctx, id, serviceID, &dto.ProjectServiceResetRequest{},
			)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestResetServiceAvailableRequests_EnvServicesErrors() {
	id := 1
	serviceID := 2

	now := time.Now()

	mockService := &entities.ProjectService{
		ID:             serviceID,
		Name:           "Service",
		Version:        "1.0.0",
		NextReset:      now.Add(24 * time.Hour),
		MaxRequest:     100,
		ResetFrequency: enums.ProjectServiceDaily,
		AssignedAt:     now.Add(-24 * time.Hour),
	}

	tests := []struct {
		name string
		req  *dto.ProjectServiceResetRequest
	}{
		{
			name: "RecalculateNextReset",
			req: &dto.ProjectServiceResetRequest{
				RecalculateNextReset: true,
			},
		},
		{
			name: "WithoutRecalculateNextReset",
			req: &dto.ProjectServiceResetRequest{
				RecalculateNextReset: false,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.projectRepo.EXPECT().
				FindServiceByID(s.ctx, id, serviceID).
				Return(mockService, nil).
				Times(1)

			if test.req.RecalculateNextReset {
				s.projectRepo.EXPECT().
					ResetProjectServiceUsage(
						s.ctx, id, serviceID, gomock.Any(),
					).
					Return(nil, errors.ErrPersistence).
					Times(1)
			} else {
				s.projectRepo.EXPECT().
					ResetAvailableRequestsForEnvsService(
						s.ctx, id, serviceID,
					).
					Return(nil, errors.ErrPersistence).
					Times(1)
			}

			resp, err := s.useCase.ResetServiceAvailableRequests(
				s.ctx, id, serviceID, test.req,
			)

			s.Require().Nil(resp)
			s.Equal(errors.ErrPersistence, err)
		})
	}
}

func (s *ProjectSuite) TestUpdateService_Successes() {
	id := 1
	serviceID := 2

	now := time.Now()

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       -1,
		CurrentAllocated: 20,
	}

	tests := []struct {
		name string
		req  *dto.ProjectServiceUpdate
	}{
		{
			name: "NextResetIsZeroWithMaxRequest",
			req: &dto.ProjectServiceUpdate{
				MaxRequest:     100,
				ResetFrequency: enums.ProjectServiceDaily,
				NextReset:      time.Time{},
			},
		},
		{
			name: "NextResetIsNotZeroWithMaxRequest",
			req: &dto.ProjectServiceUpdate{
				MaxRequest:     100,
				ResetFrequency: enums.ProjectServiceDaily,
				NextReset:      now.Add(48 * time.Hour),
			},
		},
		{
			name: "NextResetIsZeroInfiniteMaxRequest",
			req: &dto.ProjectServiceUpdate{
				MaxRequest:     -1,
				ResetFrequency: enums.ProjectServiceNull,
				NextReset:      time.Time{},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.projectRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Return(mockQuota, nil).
				Times(1)

			if test.req.MaxRequest != -1 {
				s.environmentRepo.EXPECT().
					ExistsServiceWithInfiniteMaxRequest(s.ctx, id, serviceID).
					Return(false, nil).
					Times(1)
			}

			s.projectRepo.EXPECT().
				UpdateService(s.ctx, id, serviceID, test.req).
				DoAndReturn(
					func(
						ctx context.Context,
						id, serviceID int,
						update *dto.ProjectServiceUpdate,
					) (*entities.ProjectService, *errors.Error) {
						return &entities.ProjectService{
							ID:             serviceID,
							Name:           "Service",
							Version:        "1.0.0",
							NextReset:      update.NextReset,
							MaxRequest:     update.MaxRequest,
							ResetFrequency: update.ResetFrequency,
							AssignedAt:     now.Add(-24 * time.Hour),
						}, nil
					},
				).
				Times(1)

			resp, err := s.useCase.UpdateService(
				s.ctx, id, serviceID, test.req,
			)

			s.Require().Nil(err)

			s.Equal(serviceID, resp.ID)
			s.Equal("Service", resp.Name)
			s.Equal("1.0.0", resp.Version)
			s.Equal(test.req.MaxRequest, resp.MaxRequest)
			s.Equal(test.req.ResetFrequency, resp.ResetFrequency)
			s.Equal(now.Add(-24*time.Hour), resp.AssignedAt)

			if test.req.NextReset.IsZero() && test.req.MaxRequest != -1 {
				s.Equal(now.Add(24*time.Hour), resp.NextReset)
			} else {
				s.Equal(test.req.NextReset, resp.NextReset)
			}
		})
	}

}

func (s *ProjectSuite) TestUpdateService_ValidationErrors() {
	id := 1
	serviceID := 2

	now := time.Now()

	tests := []struct {
		name        string
		req         *dto.ProjectServiceUpdate
		expectedErr *errors.Error
	}{
		{
			name: "ErrInvalidMaxRequest",
			req: &dto.ProjectServiceUpdate{
				MaxRequest:     -2,
				ResetFrequency: enums.ProjectServiceNull,
				NextReset:      time.Time{},
			},
			expectedErr: errors.ErrInvalidMaxRequest,
		},
		{
			name: "ErrProjectServiceNextResetInPast",
			req: &dto.ProjectServiceUpdate{
				MaxRequest:     100,
				ResetFrequency: enums.ProjectServiceDaily,
				NextReset:      now.Add(-24 * time.Hour),
			},
			expectedErr: errors.ErrProjectServiceNextResetInPast,
		},
		{
			name: "ErrProjectServiceNextResetWithInfiniteQuota",
			req: &dto.ProjectServiceUpdate{
				MaxRequest:     -1,
				ResetFrequency: enums.ProjectServiceNull,
				NextReset:      now.Add(24 * time.Hour),
			},
			expectedErr: errors.ErrProjectServiceNextResetWithInfiniteQuota,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Times(0)

			s.projectRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Times(0)

			s.environmentRepo.EXPECT().
				ExistsServiceWithInfiniteMaxRequest(s.ctx, id, serviceID).
				Times(0)

			s.projectRepo.EXPECT().
				UpdateService(s.ctx, id, serviceID, test.req).
				Times(0)

			resp, err := s.useCase.UpdateService(
				s.ctx, id, serviceID, test.req,
			)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestUpdateService_ExistsErrors() {
	id := 1
	serviceID := 2

	now := time.Now()

	req := &dto.ProjectServiceUpdate{
		MaxRequest:     100,
		ResetFrequency: enums.ProjectServiceDaily,
		NextReset:      now.Add(48 * time.Hour),
	}

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "DoesNotExist",
			mockErr:     nil,
			expectedErr: errors.ErrProjectNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(false, test.mockErr).
				Times(1)

			s.projectRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Times(0)

			s.environmentRepo.EXPECT().
				ExistsServiceWithInfiniteMaxRequest(s.ctx, id, serviceID).
				Times(0)

			s.projectRepo.EXPECT().
				UpdateService(s.ctx, id, serviceID, req).
				Times(0)

			resp, err := s.useCase.UpdateService(
				s.ctx, id, serviceID, &dto.ProjectServiceUpdate{},
			)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestUpdateService_QuotaErrors() {
	id := 1
	serviceID := 2

	now := time.Now()

	req := &dto.ProjectServiceUpdate{
		MaxRequest:     100,
		ResetFrequency: enums.ProjectServiceDaily,
		NextReset:      now.Add(48 * time.Hour),
	}

	tests := []struct {
		name        string
		mockQuota   *dto.QuotaUsage
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockQuota:   nil,
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrServiceNotAssignedToProject,
		},
		{
			name:        "ErrPersistence",
			mockQuota:   nil,
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
		{
			name: "ServiceMaxRequestBelow",
			mockQuota: &dto.QuotaUsage{
				MaxAllowed:       1000,
				CurrentAllocated: 1000,
			},
			mockErr:     nil,
			expectedErr: errors.ErrProjectServiceMaxRequestBelow,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.projectRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Return(test.mockQuota, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				ExistsServiceWithInfiniteMaxRequest(s.ctx, id, serviceID).
				Times(0)

			s.projectRepo.EXPECT().
				UpdateService(s.ctx, id, serviceID, req).
				Times(0)

			resp, err := s.useCase.UpdateService(
				s.ctx, id, serviceID, req,
			)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestUpdateService_EnvWithInfiniteErrors() {
	id := 1
	serviceID := 2

	now := time.Now()

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       -1,
		CurrentAllocated: 20,
	}

	req := &dto.ProjectServiceUpdate{
		MaxRequest:     100,
		ResetFrequency: enums.ProjectServiceDaily,
		NextReset:      now.Add(48 * time.Hour),
	}

	tests := []struct {
		name            string
		mockHasInfinite bool
		mockErr         *errors.Error
		expectedErr     *errors.Error
	}{
		{
			name:            "ErrPersistence",
			mockHasInfinite: false,
			mockErr:         errors.ErrPersistence,
			expectedErr:     errors.ErrPersistence,
		},
		{
			name:            "ErrServiceWithInfiniteMaxRequest",
			mockHasInfinite: true,
			mockErr:         nil,
			expectedErr:     errors.ErrProjectServiceFiniteQuotaWithInfiniteEnvs,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.projectRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Return(mockQuota, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				ExistsServiceWithInfiniteMaxRequest(s.ctx, id, serviceID).
				Return(test.mockHasInfinite, test.mockErr).
				Times(1)

			s.projectRepo.EXPECT().
				UpdateService(s.ctx, id, serviceID, req).
				Times(0)

			resp, err := s.useCase.UpdateService(
				s.ctx, id, serviceID, req,
			)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestUpdateService_ValidateErrors() {
	id := 1
	serviceID := 2

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       100,
		CurrentAllocated: 20,
	}

	req := &dto.ProjectServiceUpdate{
		MaxRequest:     -1,
		ResetFrequency: enums.ProjectServiceDaily,
		NextReset:      time.Time{},
	}

	s.projectRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
		Return(mockQuota, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceWithInfiniteMaxRequest(s.ctx, id, serviceID).
		Times(0)

	s.projectRepo.EXPECT().
		UpdateService(s.ctx, id, serviceID, req).
		Times(0)

	resp, err := s.useCase.UpdateService(
		s.ctx, id, serviceID, req,
	)

	s.Require().Nil(resp)
	s.Equal(errors.ErrProjectServiceResetFrequencyNotPermitted, err)
}

func (s *ProjectSuite) TestUpdateService_UpdateServiceError() {
	id := 1
	serviceID := 2

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       100,
		CurrentAllocated: 20,
	}

	req := &dto.ProjectServiceUpdate{
		MaxRequest:     -1,
		ResetFrequency: enums.ProjectServiceNull,
		NextReset:      time.Time{},
	}

	s.projectRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
		Return(mockQuota, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceWithInfiniteMaxRequest(s.ctx, id, serviceID).
		Times(0)

	s.projectRepo.EXPECT().
		UpdateService(s.ctx, id, serviceID, req).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.UpdateService(
		s.ctx, id, serviceID, req,
	)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func (s *ProjectSuite) TestUpdate_Success() {
	id := 1

	now := time.Now()

	mockProject := &entities.Project{
		ID:       id,
		Name:     "Name",
		Status:   enums.ProjectInProduction,
		ClientID: 1,
		Services: []*entities.ProjectService{
			{
				ID:             1,
				Name:           "Service 1",
				Version:        "1.0.0",
				MaxRequest:     100,
				ResetFrequency: enums.ProjectServiceDaily,
				NextReset:      now.Add(24 * time.Hour),
				AssignedAt:     now.Add(-24 * time.Hour),
			},
		},
		CreatedAt: now.Add(-24 * time.Hour),
	}

	req := &dto.ProjectUpdate{
		Name: "Updated Project",
	}

	s.projectRepo.EXPECT().
		Update(s.ctx, id, req).
		DoAndReturn(
			func(
				ctx context.Context, id int, update *dto.ProjectUpdate,
			) (*entities.Project, *errors.Error) {
				mockProject.Name = update.Name
				return mockProject, nil
			},
		).
		Times(1)

	resp, err := s.useCase.Update(s.ctx, id, req)

	s.Require().Nil(err)

	s.Equal(mockProject.ID, resp.ID)
	s.Equal(req.Name, resp.Name)
	s.Equal(mockProject.Status, resp.Status)
	s.Equal(mockProject.ClientID, resp.ClientID)
	s.Equal(mockProject.CreatedAt, resp.CreatedAt)

	s.Equal(len(mockProject.Services), len(resp.Services))

	for i, service := range mockProject.Services {
		s.Equal(service.ID, resp.Services[i].ID)
		s.Equal(service.Name, resp.Services[i].Name)
		s.Equal(service.Version, resp.Services[i].Version)
		s.Equal(service.MaxRequest, resp.Services[i].MaxRequest)
		s.Equal(service.ResetFrequency, resp.Services[i].ResetFrequency)
		s.Equal(service.NextReset, resp.Services[i].NextReset)
		s.Equal(service.AssignedAt, resp.Services[i].AssignedAt)
	}
}

func (s *ProjectSuite) TestUpdate_ProjectRepoError() {
	id := 1

	req := &dto.ProjectUpdate{
		Name: "Updated Project",
	}

	s.projectRepo.EXPECT().
		Update(s.ctx, id, req).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.Update(s.ctx, id, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func (s *ProjectSuite) TestAssignService_Success() {
	id := 1

	now := time.Now()

	req := &dto.ProjectService{
		ID:             1,
		MaxRequest:     10,
		ResetFrequency: enums.ProjectServiceDaily,
	}

	s.projectRepo.EXPECT().
		AddService(
			s.ctx, id, gomock.AssignableToTypeOf(&entities.ProjectService{}),
		).
		DoAndReturn(
			func(
				ctx context.Context, id int, service *entities.ProjectService,
			) *errors.Error {
				service.Name = "Service"
				service.Version = "1.0.0"
				service.AssignedAt = now
				return nil
			},
		).
		Times(1)

	resp, err := s.useCase.AssignService(s.ctx, id, req)

	s.Require().Nil(err)

	s.Equal(req.ID, resp.ID)
	s.Equal("Service", resp.Name)
	s.Equal("1.0.0", resp.Version)
	s.Equal(req.MaxRequest, resp.MaxRequest)
	s.Equal(req.ResetFrequency, resp.ResetFrequency)
	s.Equal(now, resp.AssignedAt)
	s.True(resp.NextReset.Equal(now.Add(24 * time.Hour).Truncate(24 * time.Hour)))
}

func (s *ProjectSuite) TestAssignService_ValidateError() {
	id := 1

	req := &dto.ProjectService{
		ID:             1,
		MaxRequest:     -1,
		ResetFrequency: enums.ProjectServiceDaily,
	}

	s.projectRepo.EXPECT().
		AddService(
			s.ctx, id, gomock.AssignableToTypeOf(&entities.ProjectService{}),
		).
		Times(0)

	resp, err := s.useCase.AssignService(s.ctx, id, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrProjectServiceResetFrequencyNotPermitted, err)
}

func (s *ProjectSuite) TestAssignService_ProjectRepoErrors() {
	id := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrUniqueViolation",
			mockErr:     errors.ErrUniqueViolation,
			expectedErr: errors.ErrProjectServiceAlreadyExists,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := &dto.ProjectService{
				ID:             1,
				MaxRequest:     10,
				ResetFrequency: enums.ProjectServiceDaily,
			}

			s.projectRepo.EXPECT().
				AddService(
					s.ctx, id, gomock.AssignableToTypeOf(&entities.ProjectService{}),
				).
				Return(test.mockErr).
				Times(1)

			resp, err := s.useCase.AssignService(s.ctx, id, req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestRemoveService_Success() {
	id := 1
	serviceID := 1

	s.projectRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		RemoveServiceFromProjectEnvironments(s.ctx, id, serviceID).
		Return(int64(0), nil).
		Times(1)

	s.projectRepo.EXPECT().
		RemoveService(s.ctx, id, serviceID).
		Return(int64(1), nil).
		Times(1)

	err := s.useCase.RemoveService(s.ctx, id, serviceID)

	s.Require().Nil(err)
}

func (s *ProjectSuite) TestRemoveService_ExistsError() {
	id := 1
	serviceID := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "DoesNotExist",
			mockErr:     nil,
			expectedErr: errors.ErrProjectNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(false, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				RemoveServiceFromProjectEnvironments(s.ctx, id, serviceID).
				Times(0)

			s.projectRepo.EXPECT().
				RemoveService(s.ctx, id, serviceID).
				Times(0)

			err := s.useCase.RemoveService(s.ctx, id, serviceID)

			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestRemoveService_EnvironmentRepoError() {
	id := 1
	serviceID := 1

	s.projectRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		RemoveServiceFromProjectEnvironments(s.ctx, id, serviceID).
		Return(int64(0), errors.ErrPersistence).
		Times(1)

	s.projectRepo.EXPECT().
		RemoveService(s.ctx, id, serviceID).
		Times(0)

	err := s.useCase.RemoveService(s.ctx, id, serviceID)

	s.Equal(errors.ErrPersistence, err)
}

func (s *ProjectSuite) TestRemoveService_RemoveServiceErrors() {
	id := 1
	serviceID := 1

	tests := []struct {
		name        string
		mockN       int64
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockN:       0,
			mockErr:     nil,
			expectedErr: errors.ErrServiceNotFound,
		},
		{
			name:        "ErrPersistence",
			mockN:       0,
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				RemoveServiceFromProjectEnvironments(s.ctx, id, serviceID).
				Return(int64(0), nil).
				Times(1)

			s.projectRepo.EXPECT().
				RemoveService(s.ctx, id, serviceID).
				Return(test.mockN, test.mockErr).
				Times(1)

			err := s.useCase.RemoveService(s.ctx, id, serviceID)

			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestGetByID_Success() {
	id := 1

	now := time.Now()

	mockProject := &entities.Project{
		ID:       id,
		Name:     "Name",
		Status:   enums.ProjectInProduction,
		ClientID: 1,
		Services: []*entities.ProjectService{
			{
				ID:             1,
				Name:           "Service 1",
				Version:        "1.0.0",
				MaxRequest:     100,
				ResetFrequency: enums.ProjectServiceDaily,
				NextReset:      now.Add(24 * time.Hour),
				AssignedAt:     now.Add(-24 * time.Hour),
			},
		},
		CreatedAt: now.Add(-24 * time.Hour),
	}

	s.projectRepo.EXPECT().
		FindByID(s.ctx, id).
		Return(mockProject, nil).
		Times(1)

	resp, err := s.useCase.GetByID(s.ctx, id)

	s.Require().Nil(err)

	s.Equal(mockProject.ID, resp.ID)
	s.Equal(mockProject.Name, resp.Name)
	s.Equal(mockProject.Status, resp.Status)
	s.Equal(mockProject.ClientID, resp.ClientID)
	s.Equal(mockProject.CreatedAt, resp.CreatedAt)

	s.Equal(len(mockProject.Services), len(resp.Services))

	for i, service := range mockProject.Services {
		s.Equal(service.ID, resp.Services[i].ID)
		s.Equal(service.Name, resp.Services[i].Name)
		s.Equal(service.Version, resp.Services[i].Version)
		s.Equal(service.MaxRequest, resp.Services[i].MaxRequest)
		s.Equal(service.ResetFrequency, resp.Services[i].ResetFrequency)
		s.Equal(service.NextReset, resp.Services[i].NextReset)
		s.Equal(service.AssignedAt, resp.Services[i].AssignedAt)
	}
}

func (s *ProjectSuite) TestGetByID_ProjectRepoErrors() {
	id := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrProjectNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				FindByID(s.ctx, id).
				Return(nil, test.mockErr).
				Times(1)

			resp, err := s.useCase.GetByID(s.ctx, id)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestGetEnvironments_Success() {
	id := 1

	now := time.Now()

	mockProject := []*entities.Environment{
		{
			ID:     1,
			Name:   "Environment 1",
			Status: enums.EnvironmentActive,
			Services: []*entities.EnvironmentService{
				{
					ID:               1,
					Name:             "Service 1",
					Version:          "1.0.0",
					MaxRequest:       100,
					AvailableRequest: 100,
					AssignedAt:       now.Add(-24 * time.Hour),
				},
				{
					ID:               2,
					Name:             "Service 2",
					Version:          "1.0.0",
					MaxRequest:       -1,
					AvailableRequest: -1,
					AssignedAt:       now.Add(-24 * time.Hour),
				},
			},
			CreatedAt: now.Add(-24 * time.Hour),
		},
		{
			ID:        2,
			Name:      "Environment 2",
			Status:    enums.EnvironmentActive,
			Services:  []*entities.EnvironmentService{},
			CreatedAt: now.Add(-24 * time.Hour),
		},
	}

	s.projectRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		FindByProject(s.ctx, id).
		Return(mockProject, nil).
		Times(1)

	resp, err := s.useCase.GetEnvironments(s.ctx, id)

	s.Require().Nil(err)

	s.Equal(len(mockProject), len(resp))

	for i, env := range mockProject {
		s.Equal(env.ID, resp[i].ID)
		s.Equal(env.Name, resp[i].Name)
		s.Equal(env.Status, resp[i].Status)
		s.Equal(env.CreatedAt, resp[i].CreatedAt)

		s.Equal(len(env.Services), len(resp[i].Services))

		for j, service := range env.Services {
			s.Equal(service.ID, resp[i].Services[j].ID)
			s.Equal(service.Name, resp[i].Services[j].Name)
			s.Equal(service.Version, resp[i].Services[j].Version)
			s.Equal(service.MaxRequest, resp[i].Services[j].MaxRequest)
			s.Equal(service.AvailableRequest, resp[i].Services[j].AvailableRequest)
			s.Equal(service.AssignedAt, resp[i].Services[j].AssignedAt)
		}
	}
}

func (s *ProjectSuite) TestGetEnvironments_ExistsErrors() {
	id := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "DoesNotExist",
			mockErr:     nil,
			expectedErr: errors.ErrProjectNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.projectRepo.EXPECT().
				Exists(s.ctx, id).
				Return(false, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				FindByProject(s.ctx, id).
				Times(0)

			resp, err := s.useCase.GetEnvironments(s.ctx, id)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ProjectSuite) TestGetEnvironments_EnvironmentRepoError() {
	id := 1

	s.projectRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		FindByProject(s.ctx, id).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.GetEnvironments(s.ctx, id)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}
