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
	req := &dto.Authenticate{Username: "User", Password: "Password1234"}

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
						HashedPassword:     "$2a$12$mtqngRdUqtAmrc3TMSKoteeVe4lwMmq0FBJiArGse8WrzuQVm8wW.",
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
						HashedPassword:     "$2a$12$mtqngRdUqtAmrc3TMSKoteeVe4lwMmq0FBJiArGse8WrzuQVm8wW.",
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
			req := &dto.Authenticate{Username: "User", Password: "Password1234"}

			gomock.InOrder(
				s.credentialsRepo.EXPECT().
					FindCredentials(s.ctx, req.Username).
					Return(nil, test.mockErr).
					Times(1),

				s.tokenProvider.EXPECT().
					GenerateToken(gomock.Any(), gomock.Any()).
					Times(0),
			)

			resp, err := s.useCase.Authenticate(s.ctx, req)

			s.Require().Nil(resp)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *AuthSuite) TestAuthenticate_TokenProviderError() {
	req := &dto.Authenticate{Username: "User", Password: "Password1234"}

	gomock.InOrder(
		s.credentialsRepo.EXPECT().
			FindCredentials(s.ctx, req.Username).
			DoAndReturn(
				func(
					_ context.Context, username string,
				) (*entities.Credentials, *errors.Error) {
					return &entities.Credentials{
						Username:           username,
						HashedPassword:     "$2a$12$mtqngRdUqtAmrc3TMSKoteeVe4lwMmq0FBJiArGse8WrzuQVm8wW.",
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

func (s *AuthSuite) TestChangePassword_Success() {
	req := &dto.ChangePassword{
		Username:        "User",
		NewPassword:     "NewPassword123",
		ConfirmPassword: "NewPassword123",
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
						HashedPassword:     "$2a$12$mtqngRdUqtAmrc3TMSKoteeVe4lwMmq0FBJiArGse8WrzuQVm8wW.",
						ForcePasswordReset: false,
					}, nil
				},
			).
			Times(1),

		s.credentialsRepo.EXPECT().
			ChangePassword(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1),
	)

	err := s.useCase.ChangePassword(s.ctx, req)

	s.Require().Nil(err)
}

func (s *AuthSuite) TestChangePassword_ValidationErrors() {
	tests := []struct {
		name        string
		req         *dto.ChangePassword
		expectedErr *errors.Error
	}{
		{
			name: "PasswordTooShort",
			req: &dto.ChangePassword{
				Username:        "User",
				NewPassword:     "Short",
				ConfirmPassword: "Short",
			},
			expectedErr: errors.ErrPasswordTooShort,
		},
		{
			name: "PasswordMismatch",
			req: &dto.ChangePassword{
				Username:        "User",
				NewPassword:     "NewPassword123",
				ConfirmPassword: "DifferentPassword123",
			},
			expectedErr: errors.ErrPasswordMismatch,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			gomock.InOrder(
				s.credentialsRepo.EXPECT().
					FindCredentials(gomock.Any(), gomock.Any()).
					Return(nil, (*errors.Error)(nil)).
					Times(0),

				s.credentialsRepo.EXPECT().
					ChangePassword(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0),
			)

			err := s.useCase.ChangePassword(s.ctx, test.req)
			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *AuthSuite) TestChangePassword_CredentialsRepoFindErrors() {
	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "CredentialsNotFound",
			mockErr:     errors.ErrCredentialsNotFound,
			expectedErr: errors.ErrPasswordChangeFailed,
		},
		{
			name:        "PersistenceError",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := &dto.ChangePassword{
				Username:        "User",
				NewPassword:     "NewPassword123",
				ConfirmPassword: "NewPassword123",
			}

			gomock.InOrder(
				s.credentialsRepo.EXPECT().
					FindCredentials(s.ctx, req.Username).
					Return(nil, test.mockErr).
					Times(1),

				s.credentialsRepo.EXPECT().
					ChangePassword(gomock.Any(), gomock.Any()).
					Times(0),
			)

			err := s.useCase.ChangePassword(s.ctx, req)

			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *AuthSuite) TestChangePassword_VerifyPasswordError() {
	req := &dto.ChangePassword{
		Username:        "User",
		NewPassword:     "Password1234",
		ConfirmPassword: "Password1234",
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
						HashedPassword:     "$2a$12$mtqngRdUqtAmrc3TMSKoteeVe4lwMmq0FBJiArGse8WrzuQVm8wW.",
						ForcePasswordReset: false,
					}, nil
				},
			).
			Times(1),

		s.credentialsRepo.EXPECT().
			ChangePassword(gomock.Any(), gomock.Any()).
			Times(0),
	)

	err := s.useCase.ChangePassword(s.ctx, req)

	s.Equal(errors.ErrPasswordUnchanged, err)
}

func (s *AuthSuite) TestChangePassword_CredentialsRepoChangeErrors() {
	tests := []struct {
		name        string
		mockErr     *errors.Error
		expectedErr *errors.Error
	}{
		{
			name:        "CredentialsNotFound",
			mockErr:     errors.ErrCredentialsNotFound,
			expectedErr: errors.ErrPasswordChangeFailed,
		},
		{
			name:        "PersistenceError",
			mockErr:     errors.ErrPersistence,
			expectedErr: errors.ErrPersistence,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := &dto.ChangePassword{
				Username:        "User",
				NewPassword:     "NewPassword123",
				ConfirmPassword: "NewPassword123",
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
								HashedPassword:     "$2a$12$mtqngRdUqtAmrc3TMSKoteeVe4lwMmq0FBJiArGse8WrzuQVm8wW.",
								ForcePasswordReset: false,
							}, nil
						},
					).
					Times(1),

				s.credentialsRepo.EXPECT().
					ChangePassword(s.ctx, gomock.Any()).
					Return(test.mockErr).
					Times(1),
			)

			err := s.useCase.ChangePassword(s.ctx, req)

			s.Equal(test.expectedErr, err)
		})
	}
}

func (s *AuthSuite) TestValidateToken_Success() {
	req := &dto.TokenRequest{
		Key:  "dummy-jwt",
		Type: "Bearer",
	}

	mockUsername := "User"

	s.tokenProvider.EXPECT().
		ValidateToken(s.ctx, req).
		Return(mockUsername, nil).
		Times(1)

	subject, err := s.useCase.ValidateToken(s.ctx, req)

	s.Require().Nil(err)
	s.Equal(mockUsername, subject)
}

func (s *AuthSuite) TestValidateToken_TokenProviderError() {
	req := &dto.TokenRequest{
		Key:  "dummy-jwt",
		Type: "Bearer",
	}

	s.tokenProvider.EXPECT().
		ValidateToken(s.ctx, req).
		Return("", errors.ErrInvalidToken).
		Times(1)

	username, err := s.useCase.ValidateToken(s.ctx, req)

	s.Require().Empty(username)
	s.Equal(errors.ErrInvalidToken, err)
}

func (s *AuthSuite) TestIsPasswordResetRequired_Successes() {
	tests := []struct {
		name      string
		username  string
		mockReset bool
	}{
		{
			name:      "PasswordResetRequired",
			username:  "User",
			mockReset: true,
		},
		{
			name:      "PasswordResetNotRequired",
			username:  "User",
			mockReset: false,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.credentialsRepo.EXPECT().
				FindCredentials(s.ctx, test.username).
				DoAndReturn(
					func(
						_ context.Context, username string,
					) (*entities.Credentials, *errors.Error) {
						return &entities.Credentials{
							Username:           username,
							HashedPassword:     "$2a$12$mtqngRdUqtAmrc3TMSKoteeVe4lwMmq0FBJiArGse8WrzuQVm8wW.",
							ForcePasswordReset: test.mockReset,
						}, nil
					},
				).
				Times(1)

			resetRequired, err := s.useCase.IsPasswordResetRequired(
				s.ctx, test.username,
			)

			s.Require().Nil(err)
			s.Equal(test.mockReset, resetRequired)
		})
	}
}

func (s *AuthSuite) TestIsPasswordResetRequired_CredentialsRepoError() {
	username := "User"

	s.credentialsRepo.EXPECT().
		FindCredentials(s.ctx, username).
		Return(nil, errors.ErrCredentialsNotFound).
		Times(1)

	resetRequired, err := s.useCase.IsPasswordResetRequired(s.ctx, username)

	s.Require().False(resetRequired)
	s.Equal(errors.ErrCredentialsNotFound, err)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
