package enums

type HealthStatus string

const (
	HealthStatusNull     HealthStatus = ""
	HealthStatusOK       HealthStatus = "OK"
	HealthStatusDown     HealthStatus = "DOWN"
	HealthStatusDegraded HealthStatus = "DEGRADED"
)

func ParseHealthStatus(status string) (HealthStatus, bool) {
	switch s := HealthStatus(status); s {
	case HealthStatusNull,
		HealthStatusOK,
		HealthStatusDown,
		HealthStatusDegraded:
		return s, true
	default:
		return HealthStatusNull, false
	}
}
