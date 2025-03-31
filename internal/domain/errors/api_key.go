package errors

var (
	ErrAPIKeyNotFound = NewError(CodeNotFound, "API Key not found")

	ErrAPIKeyNotActive = NewError(CodeUnauthorized, "API Key is not active")

	ErrAPIKeyExpired       = NewError(CodeValidationError, "API Key expired")
	ErrAPIKeyAlreadyExists = NewError(CodeValidationError, "API Key already exists")
	ErrAPIKeyInvalidStatus = NewError(CodeValidationError, "Invalid API Key status")

	ErrAPIKeyGenerationFailed = NewError(CodeInternalError, "API Key generation failed")
)
