package enums

type ProjectStatus string

const (
	ProjectStatusNull     ProjectStatus = ""
	ProjectStatusEnabled  ProjectStatus = "enabled"
	ProjectStatusDisabled ProjectStatus = "disabled"
)
