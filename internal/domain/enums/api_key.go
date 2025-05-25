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
