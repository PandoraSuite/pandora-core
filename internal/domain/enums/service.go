package enums

type ServiceStatus string

const (
	ServiceStatusNull       ServiceStatus = ""
	ServiceStatusEnabled    ServiceStatus = "enabled"
	ServiceStatusDisabled   ServiceStatus = "disabled"
	ServiceStatusDeprecated ServiceStatus = "deprecated"
)
