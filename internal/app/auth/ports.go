package auth

import (
	"github.com/MAD-py/pandora-core/internal/app/auth/authenticate"
	passwordchange "github.com/MAD-py/pandora-core/internal/app/auth/password_change"
	resetcheck "github.com/MAD-py/pandora-core/internal/app/auth/reset_check"
	tokenvalidation "github.com/MAD-py/pandora-core/internal/app/auth/token_validation"
)

// ... Autenticate Use Case ...

type CredentialsGetRepository = authenticate.CredentialsRepository
type TokenGenerateProvider = authenticate.TokenProvider

// ... Password Change Use Case ...

type CredentialsPasswordChangeRepository = passwordchange.CredentialsRepository

// ... Reset Password Use Case ...

type CredentialsResetPasswordRepository = resetcheck.CredentialsRepository

// ... Token Validation Use Case ...

type TokenValidationProvider = tokenvalidation.TokenProvider
