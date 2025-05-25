package enums

type APIKeyStatus string

const (
	APIKeyStatusNull APIKeyStatus = ""
	APIKeyEnabled    APIKeyStatus = "enabled"
	APIKeyDisabled   APIKeyStatus = "disabled"
)
