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

	now := time.Now().UTC()

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

	now := time.Now().UTC()

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

	now := time.Now().UTC()

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

	now := time.Now().UTC()

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

	now := time.Now().UTC()

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

	now := time.Now().UTC()

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

	now := time.Now().UTC()

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

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}
