package enums

type EnvironmentStatus string

const (
	EnvironmentStatusNull     EnvironmentStatus = ""
	EnvironmentStatusEnabled  EnvironmentStatus = "enabled"
	EnvironmentStatusDisabled EnvironmentStatus = "disabled"
)
