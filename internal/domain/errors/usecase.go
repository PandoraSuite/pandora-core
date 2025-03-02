package errors

import "errors"

var (
	ErrAPIKeyExpired  = errors.New("api key expired")
	ErrAPIKeyNotFound = errors.New("api key not found")

	ErrServiceNotFound    = errors.New("service not found")
	ErrServiceDeprecated  = errors.New("service is deprecated")
	ErrServiceDeactivated = errors.New("service is deactivated")

	ErrEnvironmentNotFound = errors.New("environment not found")

	ErrProjectServiceNotFound = errors.New("project service not found")

	ErrEnvironmentServiceNotFound = errors.New("environment service not found")

	ErrNameCannotBeEmpty                     = errors.New("name cannot be empty")
	ErrInvalidEmailFormat                    = errors.New("invalid email format")
	ErrNoAvailableRequests                   = errors.New("no available requests")
	ErrAPIKeyGenerationFailed                = errors.New("api key generation failed")
	ErrMaxRequestExceededForServiceInProyect = errors.New("max request exceeded for service in proyect")
)
