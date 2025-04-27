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

type ClientSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	clientRepo  *mock.MockClientPort
	projectRepo *mock.MockProjectPort

	useCase *ClientUseCase

	ctx context.Context
}

func (s *ClientSuite) SetupTest() {
	time.Local = time.UTC

	s.ctrl = gomock.NewController(s.T())

	s.clientRepo = mock.NewMockClientPort(s.ctrl)
	s.projectRepo = mock.NewMockProjectPort(s.ctrl)

	s.useCase = NewClientUseCase(s.clientRepo, s.projectRepo)

	s.ctx = context.Background()
}

func (s *ClientSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ClientSuite) TestUpdate_Success() {
	req := &dto.ClientUpdate{
		Type:  enums.ClientDeveloper,
		Name:  "Name",
		Email: "updated@test.com",
	}

	id := 1
	now := time.Now()

	s.clientRepo.EXPECT().
		Update(s.ctx, id, req).
		DoAndReturn(
			func(
				ctx context.Context, id int, req *dto.ClientUpdate,
			) (*entities.Client, *errors.Error) {
				return &entities.Client{
					ID:        id,
					Type:      req.Type,
					Name:      req.Name,
					Email:     req.Email,
					CreatedAt: now.Add(-24 * time.Hour),
					UpdatedAt: now,
				}, nil
			},
		).
		Times(1)

	resp, err := s.useCase.Update(s.ctx, id, req)

	s.Require().Nil(err)

	s.Equal(id, resp.ID)
	s.Equal(req.Type, resp.Type)
	s.Equal(req.Name, resp.Name)
	s.Equal(req.Email, resp.Email)
	s.Equal(now.Add(-24*time.Hour), resp.CreatedAt)
}

func (s *ClientSuite) TestUpdate_ClientRepoError() {
	req := &dto.ClientUpdate{}

	id := 1

	s.clientRepo.EXPECT().
		Update(s.ctx, id, req).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.Update(s.ctx, id, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func (s *ClientSuite) TestGetByID_Success() {
	id := 1
	now := time.Now()

	mockClient := &entities.Client{
		ID:        id,
		Type:      enums.ClientDeveloper,
		Name:      "Name",
		Email:     "client@test.com",
		CreatedAt: now.Add(-24 * time.Hour),
		UpdatedAt: now.Add(-24 * time.Hour),
	}

	s.clientRepo.EXPECT().
		FindByID(s.ctx, id).
		Return(mockClient, nil).
		Times(1)

	resp, err := s.useCase.GetByID(s.ctx, id)

	s.Require().Nil(err)

	s.Equal(id, resp.ID)
	s.Equal(mockClient.Type, resp.Type)
	s.Equal(mockClient.Name, resp.Name)
	s.Equal(mockClient.Email, resp.Email)
	s.Equal(mockClient.CreatedAt, resp.CreatedAt)
}

func (s *ClientSuite) TestGetByID_ClientRepoError() {
	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "ErrNotFound",
			mockErr:     errors.ErrNotFound,
			expectedErr: errors.ErrClientNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			id := 1

			s.clientRepo.EXPECT().
				FindByID(s.ctx, id).
				Return(nil, test.mockErr).
				Times(1)

			resp, err := s.useCase.GetByID(s.ctx, id)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ClientSuite) TestGetProjects_Successes() {
	id := 1
	now := time.Now()

	tests := []struct {
		name         string
		mockProjects []*entities.Project
	}{
		{
			name:         "ProjectEmpty",
			mockProjects: []*entities.Project{},
		},
		{
			name: "Projects",
			mockProjects: []*entities.Project{
				{
					ID:        1,
					Name:      "Project 1",
					Status:    enums.ProjectInProduction,
					ClientID:  id,
					CreatedAt: now.Add(-24 * time.Hour),
					Services: []*entities.ProjectService{
						{
							ID:             1,
							Name:           "Service 1",
							Version:        "1.0.0",
							NextReset:      now.Add(24 * time.Hour),
							MaxRequest:     100,
							ResetFrequency: enums.ProjectServiceDaily,
							AssignedAt:     now.Add(-24 * time.Hour),
						},
					},
				},
				{
					ID:        2,
					Name:      "Project 2",
					Status:    enums.ProjectInDevelopment,
					ClientID:  id,
					CreatedAt: now.Add(-24 * time.Hour),
					Services:  []*entities.ProjectService{},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.clientRepo.EXPECT().
				Exists(s.ctx, id).
				Return(true, nil).
				Times(1)

			s.projectRepo.EXPECT().
				FindByClient(s.ctx, id).
				Return(test.mockProjects, nil).
				Times(1)

			resp, err := s.useCase.GetProjects(s.ctx, id)

			s.Require().Nil(err)
			s.Len(resp, len(test.mockProjects))

			for i, mockProject := range test.mockProjects {
				s.Equal(mockProject.ID, resp[i].ID)
				s.Equal(mockProject.Name, resp[i].Name)
				s.Equal(mockProject.Status, resp[i].Status)
				s.Equal(mockProject.ClientID, resp[i].ClientID)
				s.Equal(mockProject.CreatedAt, resp[i].CreatedAt)

				for j, mockService := range mockProject.Services {
					s.Equal(mockService.ID, resp[i].Services[j].ID)
					s.Equal(mockService.Name, resp[i].Services[j].Name)
					s.Equal(mockService.Version, resp[i].Services[j].Version)
					s.Equal(mockService.NextReset, resp[i].Services[j].NextReset)
					s.Equal(mockService.MaxRequest, resp[i].Services[j].MaxRequest)
					s.Equal(mockService.ResetFrequency, resp[i].Services[j].ResetFrequency)
					s.Equal(mockService.AssignedAt, resp[i].Services[j].AssignedAt)
				}
			}
		})
	}
}

func (s *ClientSuite) TestGetProjects_ClientRepoErrors() {
	tests := []struct {
		name        string
		mockErr     *errors.Error
		mockExists  bool
		expectedErr *errors.Error
	}{
		{
			name:        "DoesNotExist",
			mockErr:     nil,
			mockExists:  false,
			expectedErr: errors.ErrClientNotFound,
		},
		{
			name:        "ErrPersistence",
			mockErr:     errors.ErrPersistence,
			mockExists:  false,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			id := 1

			s.clientRepo.EXPECT().
				Exists(s.ctx, id).
				Return(test.mockExists, test.mockErr).
				Times(1)

			s.projectRepo.EXPECT().
				FindByClient(s.ctx, id).
				Times(0)

			resp, err := s.useCase.GetProjects(s.ctx, id)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ClientSuite) TestGetProjects_ProjectRepoError() {
	id := 1

	s.clientRepo.EXPECT().
		Exists(s.ctx, id).
		Return(true, nil).
		Times(1)

	s.projectRepo.EXPECT().
		FindByClient(s.ctx, id).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.GetProjects(s.ctx, id)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func (s *ClientSuite) TestGetAll_Success() {
	now := time.Now()

	tests := []struct {
		name        string
		req         *dto.ClientFilter
		mockClients []*entities.Client
	}{
		{
			name: "WithoutFilter",
			req:  &dto.ClientFilter{},
			mockClients: []*entities.Client{
				{
					ID:        1,
					Type:      enums.ClientDeveloper,
					Name:      "Client 1",
					Email:     "client1@test.com",
					CreatedAt: now.Add(-24 * time.Hour),
					UpdatedAt: now.Add(-24 * time.Hour),
				},
				{
					ID:        2,
					Type:      enums.ClientOrganization,
					Name:      "Client 2",
					Email:     "client2@test.com",
					CreatedAt: now.Add(-24 * time.Hour),
					UpdatedAt: now.Add(-24 * time.Hour),
				},
			},
		},
		{
			name: "WithFilter",
			req:  &dto.ClientFilter{Type: enums.ClientDeveloper},
			mockClients: []*entities.Client{
				{
					ID:        1,
					Type:      enums.ClientDeveloper,
					Name:      "Client 1",
					Email:     "client1@test.com",
					CreatedAt: now.Add(-24 * time.Hour),
					UpdatedAt: now.Add(-24 * time.Hour),
				},
				{
					ID:        2,
					Type:      enums.ClientDeveloper,
					Name:      "Client 2",
					Email:     "client2@test.com",
					CreatedAt: now.Add(-24 * time.Hour),
					UpdatedAt: now.Add(-24 * time.Hour),
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.clientRepo.EXPECT().
				FindAll(s.ctx, test.req).
				Return(test.mockClients, nil).
				Times(1)

			resp, err := s.useCase.GetAll(s.ctx, test.req)

			s.Require().Nil(err)
			s.Len(resp, len(test.mockClients))

			for i, mockClient := range test.mockClients {
				s.Equal(mockClient.ID, resp[i].ID)
				s.Equal(mockClient.Type, resp[i].Type)
				s.Equal(mockClient.Name, resp[i].Name)
				s.Equal(mockClient.Email, resp[i].Email)
				s.Equal(mockClient.CreatedAt, resp[i].CreatedAt)
			}
		})
	}
}

func (s *ClientSuite) TestGetAll_ClientRepoError() {
	req := &dto.ClientFilter{}

	s.clientRepo.EXPECT().
		FindAll(s.ctx, req).
		Return(nil, errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.GetAll(s.ctx, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func (s *ClientSuite) TestCreate_Success() {
	req := &dto.ClientCreate{
		Type:  enums.ClientDeveloper,
		Name:  "Name",
		Email: "create@test.com",
	}

	now := time.Now()

	s.clientRepo.EXPECT().
		Save(s.ctx, gomock.AssignableToTypeOf(&entities.Client{})).
		DoAndReturn(
			func(
				ctx context.Context, client *entities.Client,
			) *errors.Error {
				client.ID = 1
				client.CreatedAt = now
				client.UpdatedAt = now
				return nil
			},
		).
		Times(1)

	resp, err := s.useCase.Create(s.ctx, req)

	s.Require().Nil(err)

	s.Equal(1, resp.ID)
	s.Equal(req.Type, resp.Type)
	s.Equal(req.Name, resp.Name)
	s.Equal(req.Email, resp.Email)
	s.Equal(now, resp.CreatedAt)
}

func (s *ClientSuite) TestCreate_ValidationErrors() {
	tests := []struct {
		name        string
		req         *dto.ClientCreate
		expectedErr *errors.Error
	}{
		{
			name: "EmptyName",
			req: &dto.ClientCreate{
				Type:  enums.ClientDeveloper,
				Name:  "",
				Email: "create@test.com",
			},
			expectedErr: errors.ErrNameCannotBeEmpty,
		},
		{
			name: "InvalidEmail",
			req: &dto.ClientCreate{
				Type:  enums.ClientDeveloper,
				Name:  "Name",
				Email: "invalid-email",
			},
			expectedErr: errors.ErrInvalidEmailFormat,
		},
		{
			name: "NullType",
			req: &dto.ClientCreate{
				Type:  enums.ClientTypeNull,
				Name:  "Name",
				Email: "create@test.com",
			},
			expectedErr: errors.ErrClientTypeCannotBeNull,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.clientRepo.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Times(0)

			resp, err := s.useCase.Create(s.ctx, test.req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *ClientSuite) TestCreate_ClientRepoError() {
	req := &dto.ClientCreate{
		Type:  enums.ClientDeveloper,
		Name:  "Name",
		Email: "create@test.com",
	}

	s.clientRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(errors.ErrPersistence).
		Times(1)

	resp, err := s.useCase.Create(s.ctx, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrPersistence, err)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
