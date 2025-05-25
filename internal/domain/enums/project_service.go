package enums

type ProjectServiceResetFrequency string

const (
	ProjectServiceResetFrequencyNull     ProjectServiceResetFrequency = ""
	ProjectServiceResetFrequencyDaily    ProjectServiceResetFrequency = "daily"
	ProjectServiceResetFrequencyWeekly   ProjectServiceResetFrequency = "weekly"
	ProjectServiceResetFrequencyBiweekly ProjectServiceResetFrequency = "biweekly"
	ProjectServiceResetFrequencyMonthly  ProjectServiceResetFrequency = "monthly"
)

func ParseProjectServiceResetFrequency(status string) (ProjectServiceResetFrequency, bool) {
	switch rf := ProjectServiceResetFrequency(status); rf {
	case ProjectServiceResetFrequencyNull,
		ProjectServiceResetFrequencyDaily,
		ProjectServiceResetFrequencyWeekly,
		ProjectServiceResetFrequencyBiweekly,
		ProjectServiceResetFrequencyMonthly:
		return rf, true
	default:
		return ProjectServiceResetFrequencyNull, false
	}
}
