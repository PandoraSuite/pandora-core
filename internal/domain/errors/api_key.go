package errors

var (
	ErrAPIKeyNotFound = NewError(CodeNotFound, "API Key not found")

	ErrAPIKeyNotActive = NewError(CodeUnauthorized, "API Key is not active")

	ErrAPIKeyExpired          = NewError(CodeValidationError, "API Key expired")
	ErrAPIKeyAlreadyExists    = NewError(CodeValidationError, "API Key already exists")
	ErrAPIKeyInvalidStatus    = NewError(CodeValidationError, "Invalid API Key status")
	ErrAPIKeyInvalidExpiresAt = NewError(CodeValidationError, "Expires at cannot be in the past")

	ErrAPIKeyGenerationFailed = NewError(CodeInternalError, "API Key generation failed")
)
