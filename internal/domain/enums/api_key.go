package enums

type APIKeyStatus string

const (
	APIKeyStatusNull     APIKeyStatus = ""
	APIKeyStatusEnabled  APIKeyStatus = "enabled"
	APIKeyStatusDisabled APIKeyStatus = "disabled"
)
