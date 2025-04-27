package app

import (
	"context"
	"fmt"
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

type EnvironmentSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	environmentRepo *mock.MockEnvironmentPort
	projectRepo     *mock.MockProjectPort

	useCase *EnvironmentUseCase

	ctx context.Context
}

func (s *EnvironmentSuite) SetupTest() {
	time.Local = time.UTC

	s.ctrl = gomock.NewController(s.T())

	s.environmentRepo = mock.NewMockEnvironmentPort(s.ctrl)
	s.projectRepo = mock.NewMockProjectPort(s.ctrl)

	s.useCase = NewEnvironmentUseCase(s.environmentRepo, s.projectRepo)

	s.ctx = context.Background()
}

func (s *EnvironmentSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *EnvironmentSuite) TestUpdateService_Successes() {
	id := 1
	serviceID := 1
	MaxRequest := 10

	buildService := func(opt func(*entities.EnvironmentService)) *entities.EnvironmentService {
		s := &entities.EnvironmentService{
			ID:               serviceID,
			Name:             "Service",
			Version:          "1.0.0",
			MaxRequest:       20,
			AvailableRequest: 20,
			AssignedAt:       time.Now().Add(-24 * time.Hour),
		}
		if opt != nil {
			opt(s)
		}
		return s
	}

	buildQuota := func(opt func(*dto.QuotaUsage)) *dto.QuotaUsage {
		q := &dto.QuotaUsage{
			MaxAllowed:       -1,
			CurrentAllocated: 0,
		}

		if opt != nil {
			opt(q)
		}
		return q
	}

	tests := []struct {
		name                     string
		mockService              func(*entities.EnvironmentService)
		mockQuota                func(*dto.QuotaUsage)
		expectedAvailableRequest int
	}{
		{
			name: "AvailableRequestMinor",
			mockService: func(s *entities.EnvironmentService) {
				s.AvailableRequest = 1
			},
			mockQuota:                nil,
			expectedAvailableRequest: 1,
		},
		{
			name: "AvailableRequestEqual",
			mockService: func(s *entities.EnvironmentService) {
				s.AvailableRequest = 10
			},
			mockQuota:                nil,
			expectedAvailableRequest: MaxRequest,
		},
		{
			name:                     "AvailableRequestGreater",
			mockService:              nil,
			mockQuota:                nil,
			expectedAvailableRequest: MaxRequest,
		},
		{
			name: "MaxRequestInfinite",
			mockService: func(s *entities.EnvironmentService) {
				s.MaxRequest = -1
				s.AvailableRequest = -1
			},
			mockQuota:                nil,
			expectedAvailableRequest: MaxRequest,
		},
		{
			name:        "LimitedQuota",
			mockService: nil,
			mockQuota: func(q *dto.QuotaUsage) {
				q.MaxAllowed = 100
				q.CurrentAllocated = 30
			},
			expectedAvailableRequest: MaxRequest,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := &dto.EnvironmentServiceUpdate{MaxRequest: MaxRequest}

			mockService := buildService(test.mockService)
			mockQuota := buildQuota(test.mockQuota)

			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				FindServiceByID(s.ctx, id, serviceID).
				Return(mockService, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Return(mockQuota, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				UpdateService(s.ctx, id, serviceID, req).
				DoAndReturn(
					func(
						ctx context.Context,
						id, serviceID int,
						req *dto.EnvironmentServiceUpdate,
					) (*entities.EnvironmentService, *errors.Error) {
						return &entities.EnvironmentService{
							ID:               serviceID,
							Name:             mockService.Name,
							Version:          mockService.Version,
							MaxRequest:       req.MaxRequest,
							AvailableRequest: req.AvailableRequest,
							AssignedAt:       mockService.AssignedAt,
						}, nil
					},
				).
				Times(1)

			resp, err := s.useCase.UpdateService(s.ctx, id, serviceID, req)

			s.Require().Nil(err)

			s.Equal(mockService.ID, resp.ID)
			s.Equal(mockService.Name, resp.Name)
			s.Equal(req.MaxRequest, resp.MaxRequest)
			s.Equal(mockService.Version, resp.Version)
			s.Equal(mockService.AssignedAt, resp.AssignedAt)
			s.Equal(test.expectedAvailableRequest, resp.AvailableRequest)
		})
	}
}

func (s *EnvironmentSuite) TestUpdateService_ValidationError() {
	req := &dto.EnvironmentServiceUpdate{MaxRequest: -2}

	s.environmentRepo.EXPECT().
		Exists(s.ctx, gomock.Any()).
		Times(0)

	s.environmentRepo.EXPECT().
		FindServiceByID(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	s.environmentRepo.EXPECT().
		GetProjectServiceQuotaUsage(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	s.environmentRepo.EXPECT().
		UpdateService(s.ctx, gomock.Any(), gomock.Any(), req).
		Times(0)

	resp, err := s.useCase.UpdateService(s.ctx, 1, 1, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrInvalidMaxRequest, err)
}

func (s *EnvironmentSuite) TestUpdateService_ExistsErrors() {
	id := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "DoesNotExist",
			mockErr:     nil,
			expectedErr: errors.ErrEnvironmentNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := &dto.EnvironmentServiceUpdate{MaxRequest: 10}

			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(false, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				FindServiceByID(s.ctx, id, gomock.Any()).
				Times(0)

			s.environmentRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, gomock.Any()).
				Times(0)

			s.environmentRepo.EXPECT().
				UpdateService(s.ctx, id, gomock.Any(), req).
				Times(0)

			resp, err := s.useCase.UpdateService(s.ctx, id, 1, req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestUpdateService_FindServiceErrors() {
	id := 1
	serviceID := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrServiceNotAssignedToEnvironment,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := &dto.EnvironmentServiceUpdate{MaxRequest: 10}

			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				FindServiceByID(s.ctx, id, serviceID).
				Return(nil, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Times(0)

			s.environmentRepo.EXPECT().
				UpdateService(s.ctx, id, serviceID, req).
				Times(0)

			resp, err := s.useCase.UpdateService(s.ctx, id, serviceID, req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestUpdateService_QuotaErrors() {
	id := 1
	serviceID := 1

	mockService := &entities.EnvironmentService{
		ID:               serviceID,
		Name:             "Service",
		Version:          "1.0.0",
		MaxRequest:       20,
		AvailableRequest: 20,
		AssignedAt:       time.Now().Add(-24 * time.Hour),
	}

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       100,
		CurrentAllocated: 30,
	}

	tests := []struct {
		name        string
		req         *dto.EnvironmentServiceUpdate
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			req:         &dto.EnvironmentServiceUpdate{MaxRequest: 10},
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrServiceNotAssignedToProject,
		},
		{
			name:        "ErrPersistence",
			req:         &dto.EnvironmentServiceUpdate{MaxRequest: 10},
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
		{
			name:        "ErrInfiniteRequestsNotAllowed",
			req:         &dto.EnvironmentServiceUpdate{MaxRequest: -1},
			mockErr:     nil,
			expectedErr: errors.ErrInfiniteRequestsNotAllowed,
		},
		{
			name:        "ErrMaxRequestExceededForServiceInProyect",
			req:         &dto.EnvironmentServiceUpdate{MaxRequest: 1000},
			mockErr:     nil,
			expectedErr: errors.ErrMaxRequestExceededForServiceInProyect,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				FindServiceByID(s.ctx, id, serviceID).
				Return(mockService, nil).
				Times(1)

			if test.mockErr != nil {
				s.environmentRepo.EXPECT().
					GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
					Return(nil, test.mockErr).
					Times(1)
			} else {
				s.environmentRepo.EXPECT().
					GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
					Return(mockQuota, nil).
					Times(1)
			}

			s.environmentRepo.EXPECT().
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

func (s *EnvironmentSuite) TestUpdateService_UpdateError() {
	id := 1
	serviceID := 1

	req := &dto.EnvironmentServiceUpdate{MaxRequest: 10}

	mockService := &entities.EnvironmentService{
		ID:               serviceID,
		Name:             "Service",
		Version:          "1.0.0",
		MaxRequest:       20,
		AvailableRequest: 20,
		AssignedAt:       time.Now().Add(-24 * time.Hour),
	}

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       -1,
		CurrentAllocated: 0,
	}

	s.environmentRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		FindServiceByID(s.ctx, id, serviceID).
		Return(mockService, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
		Return(mockQuota, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		UpdateService(s.ctx, id, serviceID, req).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.UpdateService(s.ctx, id, serviceID, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func (s *EnvironmentSuite) TestUpdate_Success() {
	id := 1
	now := time.Now()

	req := &dto.EnvironmentUpdate{Name: "Updated"}

	mockEnvironment := &entities.Environment{
		ID:        id,
		Name:      req.Name,
		Status:    enums.EnvironmentActive,
		ProjectID: 1,
		Services: []*entities.EnvironmentService{
			{
				ID:               1,
				Name:             "Service 1",
				Version:          "1.0.0",
				MaxRequest:       20,
				AvailableRequest: 20,
				AssignedAt:       now.Add(-24 * time.Hour),
			},
			{
				ID:               2,
				Name:             "Service 2",
				Version:          "1.0.0",
				MaxRequest:       20,
				AvailableRequest: 20,
				AssignedAt:       now.Add(-24 * time.Hour),
			},
		},
		CreatedAt: now.Add(-24 * time.Hour),
	}

	s.environmentRepo.EXPECT().
		Update(s.ctx, id, req).
		DoAndReturn(
			func(
				ctx context.Context, id int, req *dto.EnvironmentUpdate,
			) (*entities.Environment, *errors.Error) {
				mockEnvironment.Name = req.Name
				return mockEnvironment, nil
			},
		).
		Times(1)

	resp, err := s.useCase.Update(s.ctx, id, req)

	s.Require().Nil(err)

	s.Equal(mockEnvironment.ID, resp.ID)
	s.Equal(mockEnvironment.Name, resp.Name)
	s.Equal(mockEnvironment.Status, resp.Status)
	s.Equal(mockEnvironment.ProjectID, resp.ProjectID)
	s.Equal(mockEnvironment.CreatedAt, resp.CreatedAt)

	s.Equal(len(mockEnvironment.Services), len(resp.Services))

	for i, service := range mockEnvironment.Services {
		s.Equal(service.ID, resp.Services[i].ID)
		s.Equal(service.Name, resp.Services[i].Name)
		s.Equal(service.Version, resp.Services[i].Version)
		s.Equal(service.MaxRequest, resp.Services[i].MaxRequest)
		s.Equal(service.AvailableRequest, resp.Services[i].AvailableRequest)
		s.Equal(service.AssignedAt, resp.Services[i].AssignedAt)
	}
}

func (s *EnvironmentSuite) TestUpdate_EnvironmentRepoError() {
	id := 1
	req := &dto.EnvironmentUpdate{Name: "Updated"}

	s.environmentRepo.EXPECT().
		Update(s.ctx, id, req).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.Update(s.ctx, id, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func (s *EnvironmentSuite) TestResetServiceRequests_Success() {
	id := 1
	serviceID := 1

	mockService := &entities.EnvironmentService{
		ID:               serviceID,
		Name:             "Service",
		Version:          "1.0.0",
		MaxRequest:       20,
		AvailableRequest: 20,
		AssignedAt:       time.Now().Add(-24 * time.Hour),
	}

	s.environmentRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ResetAvailableRequests(s.ctx, id, serviceID).
		Return(mockService, nil).
		Times(1)

	resp, err := s.useCase.ResetServiceRequests(s.ctx, id, serviceID)

	s.Require().Nil(err)

	s.Equal(mockService.ID, resp.ID)
	s.Equal(mockService.Name, resp.Name)
	s.Equal(mockService.Version, resp.Version)
	s.Equal(mockService.MaxRequest, resp.MaxRequest)
	s.Equal(mockService.AvailableRequest, resp.AvailableRequest)
	s.Equal(mockService.AssignedAt, resp.AssignedAt)
}

func (s *EnvironmentSuite) TestResetServiceRequests_ExistsErrors() {
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
			expectedErr: errors.ErrEnvironmentNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(false, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				ResetAvailableRequests(s.ctx, id, serviceID).
				Times(0)

			resp, err := s.useCase.ResetServiceRequests(s.ctx, id, serviceID)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestResetServiceRequests_EnvironmentRepoErrors() {
	id := 1
	serviceID := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrServiceNotAssignedToEnvironment,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				ResetAvailableRequests(s.ctx, id, serviceID).
				Return(nil, test.mockErr).
				Times(1)

			resp, err := s.useCase.ResetServiceRequests(s.ctx, id, serviceID)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestGetByID_Success() {
	id := 1

	now := time.Now()

	mockEnvironment := &entities.Environment{
		ID:        id,
		Name:      "Environment",
		Status:    enums.EnvironmentActive,
		ProjectID: 1,
		Services: []*entities.EnvironmentService{
			{
				ID:               1,
				Name:             "Service 1",
				Version:          "1.0.0",
				MaxRequest:       20,
				AvailableRequest: 20,
				AssignedAt:       now.Add(-24 * time.Hour),
			},
			{
				ID:               2,
				Name:             "Service 2",
				Version:          "1.0.0",
				MaxRequest:       20,
				AvailableRequest: 20,
				AssignedAt:       now.Add(-24 * time.Hour),
			},
		},
	}

	s.environmentRepo.EXPECT().
		FindByID(s.ctx, id).
		Return(mockEnvironment, nil).
		Times(1)

	resp, err := s.useCase.GetByID(s.ctx, id)

	s.Require().Nil(err)

	s.Equal(mockEnvironment.ID, resp.ID)
	s.Equal(mockEnvironment.Name, resp.Name)
	s.Equal(mockEnvironment.Status, resp.Status)
	s.Equal(mockEnvironment.ProjectID, resp.ProjectID)

	s.Equal(len(mockEnvironment.Services), len(resp.Services))

	for i, service := range mockEnvironment.Services {
		s.Equal(service.ID, resp.Services[i].ID)
		s.Equal(service.Name, resp.Services[i].Name)
		s.Equal(service.Version, resp.Services[i].Version)
		s.Equal(service.MaxRequest, resp.Services[i].MaxRequest)
		s.Equal(service.AvailableRequest, resp.Services[i].AvailableRequest)
		s.Equal(service.AssignedAt, resp.Services[i].AssignedAt)
	}
}

func (s *EnvironmentSuite) TestGetByID_EnvironmentRepoErrors() {
	id := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrEnvironmentNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				FindByID(s.ctx, id).
				Return(nil, test.mockErr).
				Times(1)

			resp, err := s.useCase.GetByID(s.ctx, id)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestRemoveService_Success() {
	id := 1
	serviceID := 1

	s.environmentRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		RemoveService(s.ctx, id, serviceID).
		Return(int64(1), nil).
		Times(1)

	err := s.useCase.RemoveService(s.ctx, id, serviceID)

	s.Require().Nil(err)
}

func (s *EnvironmentSuite) TestRemoveService_ExistsErrors() {
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
			expectedErr: errors.ErrEnvironmentNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(false, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				RemoveService(s.ctx, id, serviceID).
				Times(0)

			err := s.useCase.RemoveService(s.ctx, id, serviceID)

			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestRemoveService_EnvironmentRepoErrors() {
	id := 1
	serviceID := 1

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
		{
			name:        "ErrServiceNotFound",
			mockErr:     nil,
			expectedErr: errors.ErrServiceNotFound,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				RemoveService(s.ctx, id, serviceID).
				Return(int64(0), test.mockErr).
				Times(1)

			err := s.useCase.RemoveService(s.ctx, id, serviceID)

			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestAssignService_Successes() {
	id := 1
	serviceID := 1

	mockAssignedAt := time.Now().Add(-24 * time.Hour)

	req := &dto.EnvironmentService{
		ID:         serviceID,
		MaxRequest: 10,
	}

	tests := []struct {
		name      string
		mockQuota *dto.QuotaUsage
	}{
		{
			name: "QuotaWithInfiniteMaxAllowed",
			mockQuota: &dto.QuotaUsage{
				MaxAllowed:       -1,
				CurrentAllocated: 0,
			},
		},
		{
			name: "QuotaWithLimitedMaxAllowed",
			mockQuota: &dto.QuotaUsage{
				MaxAllowed:       100,
				CurrentAllocated: 20,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				ExistsServiceIn(s.ctx, id, serviceID).
				Return(false, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Return(test.mockQuota, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				AddService(
					s.ctx, id,
					gomock.AssignableToTypeOf(&entities.EnvironmentService{}),
				).
				DoAndReturn(
					func(
						ctx context.Context, id int,
						service *entities.EnvironmentService,
					) *errors.Error {
						service.Name = "Service"
						service.Version = "1.0.0"
						service.AssignedAt = mockAssignedAt
						return nil
					},
				).
				Times(1)

			resp, err := s.useCase.AssignService(s.ctx, id, req)

			s.Require().Nil(err)

			s.Equal(req.ID, resp.ID)
			s.Equal(req.MaxRequest, resp.MaxRequest)
			s.Equal(req.MaxRequest, resp.AvailableRequest)
			s.Equal("Service", resp.Name)
			s.Equal("1.0.0", resp.Version)
			s.Equal(mockAssignedAt, resp.AssignedAt)
		})
	}
}

func (s *EnvironmentSuite) TestAssignService_ValidationErrors() {
	id := 1

	req := &dto.EnvironmentService{
		ID:         1,
		MaxRequest: -2,
	}

	s.environmentRepo.EXPECT().
		Exists(s.ctx, id).
		Times(0)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, id, req.ID).
		Times(0)

	s.environmentRepo.EXPECT().
		GetProjectServiceQuotaUsage(s.ctx, id, req.ID).
		Times(0)

	s.environmentRepo.EXPECT().
		AddService(s.ctx, id, gomock.Any()).
		Times(0)

	resp, err := s.useCase.AssignService(s.ctx, id, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrInvalidMaxRequest, err)
}

func (s *EnvironmentSuite) TestAssignService_ExistsErrors() {
	id := 1
	serviceID := 1

	req := &dto.EnvironmentService{
		ID:         serviceID,
		MaxRequest: 10,
	}

	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "DoesNotExist",
			mockErr:     nil,
			expectedErr: errors.ErrEnvironmentNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(false, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				ExistsServiceIn(s.ctx, id, serviceID).
				Times(0)

			s.environmentRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Times(0)

			s.environmentRepo.EXPECT().
				AddService(s.ctx, id, gomock.Any()).
				Times(0)

			resp, err := s.useCase.AssignService(s.ctx, id, req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestAssignService_AlreadyExistsErrors() {
	id := 1
	serviceID := 1

	req := &dto.EnvironmentService{
		ID:         serviceID,
		MaxRequest: 10,
	}

	tests := []struct {
		name        string
		mockExists  bool
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrAlreadyExists",
			mockExists:  true,
			mockErr:     nil,
			expectedErr: errors.ErrEnvironmentServiceAlreadyExists,
		},
		{
			name:        "ErrPersistence",
			mockExists:  false,
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				ExistsServiceIn(s.ctx, id, serviceID).
				Return(test.mockExists, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
				Times(0)

			s.environmentRepo.EXPECT().
				AddService(s.ctx, id, gomock.Any()).
				Times(0)

			resp, err := s.useCase.AssignService(s.ctx, id, req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestAssignService_QuotaErrors() {
	id := 1

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       30,
		CurrentAllocated: 20,
	}

	tests := []struct {
		name        string
		req         *dto.EnvironmentService
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name: "ErrNotFound",
			req: &dto.EnvironmentService{
				ID:         1,
				MaxRequest: 10,
			},
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrServiceNotAssignedToProject,
		},
		{
			name: "ErrPersistence",
			req: &dto.EnvironmentService{
				ID:         1,
				MaxRequest: 10,
			},
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
		{
			name: "ErrInfiniteRequestsNotAllowed",
			req: &dto.EnvironmentService{
				ID:         1,
				MaxRequest: -1,
			},
			mockErr:     nil,
			expectedErr: errors.ErrInfiniteRequestsNotAllowed,
		},
		{
			name: "ErrMaxRequestExceededForServiceInProyect",
			req: &dto.EnvironmentService{
				ID:         1,
				MaxRequest: 1000,
			},
			mockErr:     nil,
			expectedErr: errors.ErrMaxRequestExceededForServiceInProyect,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.environmentRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				ExistsServiceIn(s.ctx, id, test.req.ID).
				Return(false, nil).
				Times(1)

			s.environmentRepo.EXPECT().
				GetProjectServiceQuotaUsage(s.ctx, id, test.req.ID).
				Return(mockQuota, test.mockErr).
				Times(1)

			s.environmentRepo.EXPECT().
				AddService(s.ctx, id, gomock.Any()).
				Times(0)

			resp, err := s.useCase.AssignService(s.ctx, id, test.req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestAssignService_AddServiceError() {
	id := 1
	serviceID := 1

	req := &dto.EnvironmentService{
		ID:         serviceID,
		MaxRequest: 10,
	}

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       -1,
		CurrentAllocated: 0,
	}

	s.environmentRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		ExistsServiceIn(s.ctx, id, serviceID).
		Return(false, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		GetProjectServiceQuotaUsage(s.ctx, id, serviceID).
		Return(mockQuota, nil).
		Times(1)

	s.environmentRepo.EXPECT().
		AddService(
			s.ctx, id,
			gomock.AssignableToTypeOf(&entities.EnvironmentService{}),
		).
		Return(errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.AssignService(s.ctx, id, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func (s *EnvironmentSuite) TestCreate_Success() {
	now := time.Now()

	req := &dto.EnvironmentCreate{
		Name:      "Name",
		ProjectID: 1,
		Services: []*dto.EnvironmentService{
			{
				ID:         1,
				MaxRequest: 20,
			},
			{
				ID:         2,
				MaxRequest: 10,
			},
		},
	}

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       -1,
		CurrentAllocated: 0,
	}

	s.projectRepo.EXPECT().
		Exists(s.ctx, req.ProjectID).
		Return(true, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectServiceQuotaUsage(
			s.ctx, req.ProjectID, gomock.AssignableToTypeOf(0),
		).
		Return(mockQuota, nil).
		Times(2)

	s.environmentRepo.EXPECT().
		Save(s.ctx, gomock.AssignableToTypeOf(&entities.Environment{})).
		DoAndReturn(
			func(
				ctx context.Context, env *entities.Environment,
			) *errors.Error {
				env.ID = 1
				env.CreatedAt = now

				for _, service := range env.Services {
					service.Name = fmt.Sprintf("Service %d", service.ID)
					service.Version = "1.0.0"
					service.AssignedAt = now
				}
				return nil
			},
		).
		Times(1)

	envirinment, err := s.useCase.Create(s.ctx, req)

	s.Require().Nil(err)

	s.Equal(1, envirinment.ID)
	s.Equal(req.Name, envirinment.Name)
	s.Equal(now, envirinment.CreatedAt)
	s.Equal(req.ProjectID, envirinment.ProjectID)
	s.Equal(enums.EnvironmentActive, envirinment.Status)

	s.Equal(len(req.Services), len(envirinment.Services))

	for i, service := range req.Services {
		s.Equal(service.ID, envirinment.Services[i].ID)
		s.Equal(service.MaxRequest, envirinment.Services[i].MaxRequest)
		s.Equal(service.MaxRequest, envirinment.Services[i].AvailableRequest)
		s.Equal(fmt.Sprintf("Service %d", service.ID), envirinment.Services[i].Name)
		s.Equal("1.0.0", envirinment.Services[i].Version)
		s.Equal(now, envirinment.Services[i].AssignedAt)
	}
}

func (s *EnvironmentSuite) TestCreate_ExistsErrors() {
	req := &dto.EnvironmentCreate{
		Name:      "Name",
		ProjectID: 1,
		Services: []*dto.EnvironmentService{
			{
				ID:         1,
				MaxRequest: 20,
			},
			{
				ID:         2,
				MaxRequest: 10,
			},
		},
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
				Exists(s.ctx, req.ProjectID).
				Return(false, test.mockErr).
				Times(1)

			s.projectRepo.EXPECT().
				GetProjectServiceQuotaUsage(
					s.ctx, req.ProjectID, gomock.Any(),
				).
				Times(0)

			s.environmentRepo.EXPECT().
				Save(s.ctx, gomock.Any()).
				Times(0)

			resp, err := s.useCase.Create(s.ctx, req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *EnvironmentSuite) TestCreate_QuotaErrors() {
	req := &dto.EnvironmentCreate{
		Name:      "Name",
		ProjectID: 1,
	}

	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       10,
		CurrentAllocated: 0,
	}

	buildService := func(opt func(*dto.EnvironmentService)) []*dto.EnvironmentService {
		s := &dto.EnvironmentService{
			ID:         1,
			MaxRequest: 20,
		}

		if opt != nil {
			opt(s)
		}
		return []*dto.EnvironmentService{s}
	}

	tests := []struct {
		name        string
		mockService func(*dto.EnvironmentService)
		mockErr     *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockService: nil,
			mockErr:     errors.ErrNotFound,
		},
		{
			name:        "ErrPersistence",
			mockService: nil,
			mockErr:     errors.ErrPersistence,
		},
		{
			name:        "ErrInfiniteRequestsNotAllowed",
			mockService: func(s *dto.EnvironmentService) { s.MaxRequest = -1 },
			mockErr:     nil,
		},
		{
			name:        "ErrMaxRequestExceededForServiceInProyect",
			mockService: func(s *dto.EnvironmentService) { s.MaxRequest = 1000 },
			mockErr:     nil,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req.Services = buildService(test.mockService)

			s.projectRepo.EXPECT().
				Exists(s.ctx, req.ProjectID).
				Return(true, nil).
				Times(1)

			if test.mockErr != nil {
				s.projectRepo.EXPECT().
					GetProjectServiceQuotaUsage(
						s.ctx, req.ProjectID, gomock.AssignableToTypeOf(0),
					).
					Return(nil, test.mockErr).
					Times(1)
			} else {
				s.projectRepo.EXPECT().
					GetProjectServiceQuotaUsage(
						s.ctx, req.ProjectID, gomock.AssignableToTypeOf(0),
					).
					Return(mockQuota, nil).
					Times(1)
			}

			s.environmentRepo.EXPECT().
				Save(s.ctx, gomock.Any()).
				Times(0)

			resp, err := s.useCase.Create(s.ctx, req)

			s.Require().Nil(resp)
			s.Error(err)
		})
	}
}

func (s *EnvironmentSuite) TestCreate_ValidationErrors() {
	mockQuota := &dto.QuotaUsage{
		MaxAllowed:       -1,
		CurrentAllocated: 0,
	}

	buildReq := func(opt func(*dto.EnvironmentCreate)) *dto.EnvironmentCreate {
		req := &dto.EnvironmentCreate{
			Name:      "Name",
			ProjectID: 1,
		}

		if opt != nil {
			opt(req)
		}
		return req
	}

	buildService := func(opt func(*dto.EnvironmentService)) []*dto.EnvironmentService {
		s := &dto.EnvironmentService{
			ID:         1,
			MaxRequest: 20,
		}

		if opt != nil {
			opt(s)
		}
		return []*dto.EnvironmentService{s}
	}

	tests := []struct {
		name        string
		mockReq     func(*dto.EnvironmentCreate)
		mockService func(*dto.EnvironmentService)
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNameCannotBeEmpty",
			mockReq:     func(e *dto.EnvironmentCreate) { e.Name = "" },
			mockService: nil,
			expectedErr: errors.ErrNameCannotBeEmpty,
		},
		{
			name:        "ErrInService",
			mockReq:     nil,
			mockService: func(s *dto.EnvironmentService) { s.MaxRequest = -2 },
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		req := buildReq(test.mockReq)
		req.Services = buildService(test.mockService)

		s.projectRepo.EXPECT().
			Exists(s.ctx, req.ProjectID).
			Return(true, nil).
			Times(1)

		s.projectRepo.EXPECT().
			GetProjectServiceQuotaUsage(
				s.ctx, req.ProjectID, gomock.AssignableToTypeOf(0),
			).
			Return(mockQuota, nil).
			Times(1)

		s.environmentRepo.EXPECT().
			Save(s.ctx, gomock.Any()).
			Times(0)

		resp, err := s.useCase.Create(s.ctx, req)

		s.Require().Nil(resp)

		if test.expectedErr != nil {
			s.Equal(test.expectedErr, err)
		} else {
			s.Error(err)
		}
	}
}

func (s *EnvironmentSuite) TestCreate_SaveError() {
	req := &dto.EnvironmentCreate{
		Name:      "Name",
		ProjectID: 1,
	}

	s.projectRepo.EXPECT().
		Exists(s.ctx, req.ProjectID).
		Return(true, nil).
		Times(1)

	s.projectRepo.EXPECT().
		GetProjectServiceQuotaUsage(
			s.ctx, req.ProjectID, gomock.Any(),
		).
		Times(0)

	s.environmentRepo.EXPECT().
		Save(s.ctx, gomock.Any()).
		Return(errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.Create(s.ctx, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func TestEnvironmentSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentSuite))
}
