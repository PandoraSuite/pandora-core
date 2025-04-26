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
	now := time.Now().UTC()

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

func (s *ClientSuite) TestUpdate_ClientError() {
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

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
