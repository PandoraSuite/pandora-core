package enums

import (
	"encoding/json"
	"fmt"
)

type ProjectServiceResetFrequency int

const (
	ProjectServiceDaily ProjectServiceResetFrequency = iota
	ProjectServiceWeekly
	ProjectServiceBiweekly
	ProjectServiceMonthly
)

func (rf ProjectServiceResetFrequency) String() string {
	switch rf {
	case ProjectServiceDaily:
		return "daily"
	case ProjectServiceWeekly:
		return "weekly"
	case ProjectServiceBiweekly:
		return "biweekly"
	case ProjectServiceMonthly:
		return "monthly"
	default:
		panic("unknown ProjectServiceResetFrequency")
	}
}

func (rf *ProjectServiceResetFrequency) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := ParseProjectServiceResetFrequency(s)
	if err != nil {
		return err
	}
	*rf = parsed
	return nil
}

func (rf *ProjectServiceResetFrequency) MarshalJSON() ([]byte, error) {
	return json.Marshal(rf.String())
}

func ParseProjectServiceResetFrequency(rf string) (ProjectServiceResetFrequency, error) {
	switch rf {
	case "daily":
		return ProjectServiceDaily, nil
	case "weekly":
		return ProjectServiceWeekly, nil
	case "biweekly":
		return ProjectServiceBiweekly, nil
	case "monthly":
		return ProjectServiceMonthly, nil
	default:
		return 0, fmt.Errorf("invalid reset frequency: %s", rf)
	}
}
