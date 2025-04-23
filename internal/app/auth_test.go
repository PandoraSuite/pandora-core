package app

import (
	"context"
	"testing"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type AuthSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	tokenProvider   *mock.MockTokenPort
	credentialsRepo *mock.MockCredentialsPort

	useCase *AuthUseCase

	ctx context.Context
}

func (s *AuthSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.tokenProvider = mock.NewMockTokenPort(s.ctrl)
	s.credentialsRepo = mock.NewMockCredentialsPort(s.ctrl)

	s.useCase = NewAuthUseCase(s.tokenProvider, s.credentialsRepo)

	s.ctx = context.Background()
}

func (s *AuthSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *AuthSuite) TestAuthenticate_Success() {
	req := &dto.Authenticate{Username: "User", Password: "Password"}

	mockToken := &dto.TokenResponse{
		Token:     "dummy-jwt",
		TokenType: "Bearer",
		ExpiresIn: time.Now().Add(time.Hour),
	}

	gomock.InOrder(
		s.credentialsRepo.EXPECT().
			FindCredentials(s.ctx, req.Username).
			DoAndReturn(
				func(
					_ context.Context, username string,
				) (*entities.Credentials, *errors.Error) {
					return &entities.Credentials{
						Username:           username,
						HashedPassword:     "$2a$12$PeI5jw3hPZvOh5IQ7fTElO/iCatcyLcaw3d1lQjrSwSnhDEmLPn9q",
						ForcePasswordReset: false,
					}, nil
				},
			).
			Times(1),

		s.tokenProvider.EXPECT().
			GenerateToken(s.ctx, req.Username).
			Return(mockToken, nil).
			Times(1),
	)

	resp, err := s.useCase.Authenticate(s.ctx, req)

	s.Require().Nil(err)
	s.Equal(mockToken.Token, resp.Token)
	s.Equal(mockToken.TokenType, resp.TokenType)
	s.Equal(mockToken.ExpiresIn, resp.ExpiresIn)
	s.False(resp.ForcePasswordReset)
}

func (s *AuthSuite) TestAuthenticate_VerifyPasswordError() {
	req := &dto.Authenticate{Username: "User", Password: "WrongPassword"}

	gomock.InOrder(
		s.credentialsRepo.EXPECT().
			FindCredentials(s.ctx, req.Username).
			DoAndReturn(
				func(
					_ context.Context, username string,
				) (*entities.Credentials, *errors.Error) {
					return &entities.Credentials{
						Username:           username,
						HashedPassword:     "$2a$12$PeI5jw3hPZvOh5IQ7fTElO/iCatcyLcaw3d1lQjrSwSnhDEmLPn9q",
						ForcePasswordReset: false,
					}, nil
				},
			).
			Times(1),

		s.tokenProvider.EXPECT().
			GenerateToken(gomock.Any(), gomock.Any()).
			Times(0),
	)

	resp, err := s.useCase.Authenticate(s.ctx, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrInvalidCredentials, err)
}

func (s *AuthSuite) TestAuthenticate_CredentialsRepoErrors() {
	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "CredentialsNotFound",
			mockErr:     errors.ErrCredentialsNotFound,
			expectedErr: errors.ErrInvalidCredentials,
		},
		{
			name:        "PersistenceError",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			ctrl := gomock.NewController(s.T())
			defer ctrl.Finish()

			tokenProvider := mock.NewMockTokenPort(ctrl)
			credentialsRepo := mock.NewMockCredentialsPort(ctrl)

			uc := NewAuthUseCase(tokenProvider, credentialsRepo)

			req := &dto.Authenticate{Username: "User", Password: "Password"}

			gomock.InOrder(
				credentialsRepo.EXPECT().
					FindCredentials(s.ctx, req.Username).
					Return(nil, test.mockErr).
					Times(1),

				tokenProvider.EXPECT().
					GenerateToken(gomock.Any(), gomock.Any()).
					Times(0),
			)

			resp, err := uc.Authenticate(s.ctx, req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *AuthSuite) TestAuthenticate_TokenProviderError() {
	req := &dto.Authenticate{Username: "User", Password: "Password"}

	gomock.InOrder(
		s.credentialsRepo.EXPECT().
			FindCredentials(s.ctx, req.Username).
			DoAndReturn(
				func(
					_ context.Context, username string,
				) (*entities.Credentials, *errors.Error) {
					return &entities.Credentials{
						Username:           username,
						HashedPassword:     "$2a$12$PeI5jw3hPZvOh5IQ7fTElO/iCatcyLcaw3d1lQjrSwSnhDEmLPn9q",
						ForcePasswordReset: false,
					}, nil
				},
			).
			Times(1),

		s.tokenProvider.EXPECT().
			GenerateToken(s.ctx, req.Username).
			Return(nil, errors.ErrTokenSigningFailed).
			Times(1),
	)

	resp, err := s.useCase.Authenticate(s.ctx, req)

	s.Require().Nil(resp)
	s.Equal(errors.ErrTokenSigningFailed, err)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
