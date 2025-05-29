package enums

type APIKeyStatus string

const (
	APIKeyStatusNull     APIKeyStatus = ""
	APIKeyStatusEnabled  APIKeyStatus = "enabled"
	APIKeyStatusDisabled APIKeyStatus = "disabled"
)

func ParseAPIKeyStatus(status string) (APIKeyStatus, bool) {
	switch s := APIKeyStatus(status); s {
	case APIKeyStatusNull, APIKeyStatusEnabled, APIKeyStatusDisabled:
		return s, true
	default:
		return APIKeyStatusNull, false
	}
}

type APIKeyValidationFailureCode string

const (
	APIKeyValidationFailureCodeAPIKeyInvalid       APIKeyValidationFailureCode = "API_KEY_INVALID"
	APIKeyValidationFailureCodeQuotaExceeded       APIKeyValidationFailureCode = "QUOTA_EXCEEDED"
	APIKeyValidationFailureCodeAPIKeyExpired       APIKeyValidationFailureCode = "API_KEY_EXPIRED"
	APIKeyValidationFailureCodeAPIKeyDisabled      APIKeyValidationFailureCode = "API_KEY_DISABLED"
	APIKeyValidationFailureCodeServiceMismatch     APIKeyValidationFailureCode = "SERVICE_MISMATCH"
	APIKeyValidationFailureCodeServiceDisabled     APIKeyValidationFailureCode = "SERVICE_DISABLED"
	APIKeyValidationFailureCodeServiceDeprecated   APIKeyValidationFailureCode = "SERVICE_DEPRECATED"
	APIKeyValidationFailureCodeServiceNotAssigned  APIKeyValidationFailureCode = "SERVICE_NOT_ASSIGNED"
	APIKeyValidationFailureCodeEnvironmentDisabled APIKeyValidationFailureCode = "ENVIRONMENT_DISABLED"
)

func ParseAPIKeyValidationFailureCode(code string) (APIKeyValidationFailureCode, bool) {
	switch c := APIKeyValidationFailureCode(code); c {
	case APIKeyValidationFailureCodeAPIKeyInvalid,
		APIKeyValidationFailureCodeQuotaExceeded,
		APIKeyValidationFailureCodeAPIKeyExpired,
		APIKeyValidationFailureCodeAPIKeyDisabled,
		APIKeyValidationFailureCodeServiceMismatch,
		APIKeyValidationFailureCodeServiceNotAssigned,
		APIKeyValidationFailureCodeEnvironmentDisabled:
		return c, true
	default:
		return "", false
	}
}
