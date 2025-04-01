package errors

var (
	ErrInvalidCredentials  = NewError(CodeUnauthorized, "Invalid Username or Password")
	ErrCredentialsNotFound = NewError(CodeUnauthorized, "Credentials not found")

	ErrPasswordMismatch = NewError(CodeValidationError, "Password and confirmation do not match")
	ErrPasswordTooShort = NewError(CodeValidationError, "Password must be at least 12 characters long")

	ErrPasswordChangeFailed     = NewError(CodeInternalError, "Failed to change Password")
	ErrPasswordProcessingFailed = NewError(CodeInternalError, "Unable to process the Password")
)
