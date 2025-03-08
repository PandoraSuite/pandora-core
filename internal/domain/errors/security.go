package errors

import "errors"

var (
	ErrInvalidToken       = errors.New("invalid authentication token")
	ErrInvalidTokenData   = errors.New("failed to validate authentication data")
	ErrInvalidTokenType   = errors.New("invalid token type, expected 'Bearer'")
	ErrTokenSigningFailed = errors.New("failed to sign authentication token")

	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrCredentialsNotFound = errors.New("credentials not found")

	ErrPasswordMismatch         = errors.New("password and confirmation do not match")
	ErrPasswordTooShort         = errors.New("password must be at least 12 characters long")
	ErrPasswordChangeFailed     = errors.New("failed to change password")
	ErrPasswordProcessingFailed = errors.New("unable to process the password")
)
