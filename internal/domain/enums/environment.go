package enums

type EnvironmentStatus string

const (
	EnvironmentStatusNull     EnvironmentStatus = ""
	EnvironmentStatusEnabled  EnvironmentStatus = "enabled"
	EnvironmentStatusDisabled EnvironmentStatus = "disabled"
)

func ParseEnvironmentStatus(status string) (EnvironmentStatus, bool) {
	switch s := EnvironmentStatus(status); s {
	case EnvironmentStatusNull, EnvironmentStatusEnabled, EnvironmentStatusDisabled:
		return s, true
	default:
		return EnvironmentStatusNull, false
	}
}
