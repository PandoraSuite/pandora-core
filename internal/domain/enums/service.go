package enums

type ServiceStatus string

const (
	ServiceStatusNull       ServiceStatus = ""
	ServiceStatusEnabled    ServiceStatus = "enabled"
	ServiceStatusDisabled   ServiceStatus = "disabled"
	ServiceStatusDeprecated ServiceStatus = "deprecated"
)

func ParseServiceStatus(status string) (ServiceStatus, bool) {
	switch s := ServiceStatus(status); s {
	case ServiceStatusNull,
		ServiceStatusEnabled,
		ServiceStatusDisabled,
		ServiceStatusDeprecated:
		return s, true
	default:
		return ServiceStatusNull, false
	}
}
