package auth

import (
	accesstokenvalidation "github.com/MAD-py/pandora-core/internal/app/auth/access_token_validation"
	"github.com/MAD-py/pandora-core/internal/app/auth/authenticate"
	passwordchange "github.com/MAD-py/pandora-core/internal/app/auth/password_change"
	"github.com/MAD-py/pandora-core/internal/app/auth/reauthenticate"
	resetcheck "github.com/MAD-py/pandora-core/internal/app/auth/reset_check"
	scopedtokenvalidation "github.com/MAD-py/pandora-core/internal/app/auth/scoped_token_validation"
)

// ... Autenticate Use Case ...

type CredentialsGetRepository = authenticate.CredentialsRepository
type TokenGenerateProvider = authenticate.TokenProvider

// ... Password Change Use Case ...

type CredentialsPasswordChangeRepository = passwordchange.CredentialsRepository

// ... Reset Password Use Case ...

type CredentialsResetPasswordRepository = resetcheck.CredentialsRepository

// ... Access Token Validation Use Case ...

type AccessTokenValidationProvider = accesstokenvalidation.TokenProvider

// ... Reauthenticate Use Case ...

type CredentialsReauthenticateRepository = reauthenticate.CredentialsRepository
type TokenReauthenticateProvider = reauthenticate.TokenProvider

// ... Scoped Token Validation Use Case ...

type ScopedTokenValidationProvider = scopedtokenvalidation.TokenProvider
