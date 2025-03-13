package enums

import (
	"encoding/json"
	"fmt"
)

type ProjectStatus int

const (
	ProjectStatusNull ProjectStatus = iota
	ProjectInProduction
	ProjectInDevelopment
	ProjectDeactivated
)

func (s ProjectStatus) String() string {
	switch s {
	case ProjectStatusNull:
		return ""
	case ProjectInProduction:
		return "in_production"
	case ProjectInDevelopment:
		return "in_development"
	case ProjectDeactivated:
		return "deactivated"
	default:
		panic("unknown ProjectStatus")
	}
}

func (s *ProjectStatus) UnmarshalJSON(b []byte) error {
	var ss string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}
	parsed, err := ParseProjectStatus(ss)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

func (s *ProjectStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func ParseProjectStatus(s string) (ProjectStatus, error) {
	switch s {
	case "":
		return ProjectStatusNull, nil
	case "in_production":
		return ProjectInProduction, nil
	case "in_development":
		return ProjectInDevelopment, nil
	case "deactivated":
		return ProjectDeactivated, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", s)
	}
}
