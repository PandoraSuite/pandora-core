package enums

type ProjectServiceResetFrequency string

const (
	ProjectServiceResetFrequencyNull     ProjectServiceResetFrequency = ""
	ProjectServiceResetFrequencyDaily    ProjectServiceResetFrequency = "daily"
	ProjectServiceResetFrequencyWeekly   ProjectServiceResetFrequency = "weekly"
	ProjectServiceResetFrequencyBiweekly ProjectServiceResetFrequency = "biweekly"
	ProjectServiceResetFrequencyMonthly  ProjectServiceResetFrequency = "monthly"
)
