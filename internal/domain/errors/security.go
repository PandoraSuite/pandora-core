package errors

var (
	ErrInvalidToken       = NewError(CodeUnauthorized, "invalid authentication token")
	ErrInvalidTokenData   = NewError(CodeUnauthorized, "failed to validate authentication data")
	ErrInvalidTokenType   = NewError(CodeUnauthorized, "invalid token type, expected 'Bearer'")
	ErrTokenSigningFailed = NewError(CodeInternalError, "failed to sign authentication token")

	ErrInvalidCredentials  = NewError(CodeUnauthorized, "invalid username or password")
	ErrCredentialsNotFound = NewError(CodeUnauthorized, "credentials not found")

	ErrPasswordMismatch         = NewError(CodeValidationError, "password and confirmation do not match")
	ErrPasswordTooShort         = NewError(CodeValidationError, "password must be at least 12 characters long")
	ErrPasswordChangeFailed     = NewError(CodeInternalError, "failed to change password")
	ErrPasswordProcessingFailed = NewError(CodeInternalError, "unable to process the password")
)
