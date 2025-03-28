package enums

import (
	"encoding/json"
	"fmt"
)

type ProjectServiceResetFrequency int

const (
	ProjectServiceNull ProjectServiceResetFrequency = iota
	ProjectServiceDaily
	ProjectServiceWeekly
	ProjectServiceBiweekly
	ProjectServiceMonthly
)

func (rf ProjectServiceResetFrequency) String() string {
	switch rf {
	case ProjectServiceNull:
		return "null"
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

func (rf *ProjectServiceResetFrequency) Scan(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid reset frequency: %v", v)
	}

	restFrequency, err := ParseProjectServiceResetFrequency(str)
	if err != nil {
		return err
	}

	*rf = restFrequency
	return nil
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
	return []byte(rf.String()), nil
}

func ParseProjectServiceResetFrequency(rf string) (ProjectServiceResetFrequency, error) {
	switch rf {
	case "":
		return ProjectServiceNull, nil
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
