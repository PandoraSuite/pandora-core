package errors

import "errors"

var (
	ErrInvalidToken       = errors.New("invalid authentication token")
	ErrInvalidTokenData   = errors.New("failed to validate authentication data")
	ErrInvalidTokenType   = errors.New("invalid token type, expected 'Bearer'")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrTokenSigningFailed = errors.New("failed to sign authentication token")
)
