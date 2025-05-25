package enums

type ProjectStatus string

const (
	ProjectStatusNull     ProjectStatus = ""
	ProjectStatusEnabled  ProjectStatus = "enabled"
	ProjectStatusDisabled ProjectStatus = "disabled"
)

func ParseProjectStatus(status string) (ProjectStatus, bool) {
	switch s := ProjectStatus(status); s {
	case ProjectStatusNull, ProjectStatusEnabled, ProjectStatusDisabled:
		return s, true
	default:
		return ProjectStatusNull, false
	}
}
