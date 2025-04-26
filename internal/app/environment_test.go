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

type EnvironmentSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	environmentRepo *mock.MockEnvironmentPort
	projectRepo     *mock.MockProjectPort

	useCase *EnvironmentUseCase

	ctx context.Context
}

func (s *EnvironmentSuite) SetupTest() {
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
			AssignedAt:       time.Now().UTC().Add(-24 * time.Hour),
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
		AssignedAt:       time.Now().UTC().Add(-24 * time.Hour),
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
		AssignedAt:       time.Now().UTC().Add(-24 * time.Hour),
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
	now := time.Now().UTC()

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
		AssignedAt:       time.Now().UTC().Add(-24 * time.Hour),
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

func TestEnvironmentSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentSuite))
}
