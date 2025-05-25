package enums

type EnvironmentStatus string

const (
	EnvironmentStatusNull EnvironmentStatus = ""
	EnvironmentEnabled    EnvironmentStatus = "enabled"
	EnvironmentDisabled   EnvironmentStatus = "disabled"
)
