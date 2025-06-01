package auth

import (
	accesstokenvalidation "github.com/MAD-py/pandora-core/internal/app/auth/access_token_validation"
	"github.com/MAD-py/pandora-core/internal/app/auth/authenticate"
	passwordchange "github.com/MAD-py/pandora-core/internal/app/auth/password_change"
	"github.com/MAD-py/pandora-core/internal/app/auth/reauthenticate"
	resetcheck "github.com/MAD-py/pandora-core/internal/app/auth/reset_check"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Autenticate Use Case ...

type AutenticateUseCase = authenticate.UseCase

func NewAutenticateUseCase(
	validator validator.Validator,
	tokenProvider TokenGenerateProvider,
	credentialsRepo CredentialsGetRepository,
) AutenticateUseCase {
	return authenticate.NewUseCase(validator, tokenProvider, credentialsRepo)
}

// ... Password Change Use Case ...

type PasswordChangeUseCase = passwordchange.UseCase

func NewPasswordChangeUseCase(
	validator validator.Validator,
	credentialsRepo CredentialsPasswordChangeRepository,
) PasswordChangeUseCase {
	return passwordchange.NewUseCase(validator, credentialsRepo)
}

// ... Reset Password Use Case ...

type ResetPasswordUseCase = resetcheck.UseCase

func NewResetPasswordUseCase(
	validator validator.Validator,
	credentialsRepo CredentialsResetPasswordRepository,
) ResetPasswordUseCase {
	return resetcheck.NewUseCase(validator, credentialsRepo)
}

// ... Token Validation Use Case ...

type TokenValidationUseCase = accesstokenvalidation.UseCase

func NewTokenValidationUseCase(
	validator validator.Validator,
	tokenProvider TokenValidationProvider,
) TokenValidationUseCase {
	return accesstokenvalidation.NewUseCase(validator, tokenProvider)
}

// ... Reauthenticate Use Case ...

type ReauthenticateUseCase = reauthenticate.UseCase

func NewReauthenticateUseCase(
	validator validator.Validator,
	tokenProvider TokenReauthenticateProvider,
	credentialsRepo CredentialsReauthenticateRepository,
) ReauthenticateUseCase {
	return reauthenticate.NewUseCase(validator, tokenProvider, credentialsRepo)
}
