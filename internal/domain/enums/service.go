package enums

import (
	"encoding/json"
	"fmt"
)

type ServiceStatus int

const (
	ServiceActive ServiceStatus = iota
	ServiceDeactivated
	ServiceDeprecated
)

func (s ServiceStatus) String() string {
	switch s {
	case ServiceActive:
		return "active"
	case ServiceDeactivated:
		return "deactivated"
	case ServiceDeprecated:
		return "deprecated"
	default:
		panic("unknown ServiceStatus")
	}
}

func (s *ServiceStatus) UnmarshalJSON(b []byte) error {
	var ss string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}
	parsed, err := ParseServiceStatus(ss)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

func (s *ServiceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func ParseServiceStatus(s string) (ServiceStatus, error) {
	switch s {
	case "active":
		return ServiceActive, nil
	case "deactivated":
		return ServiceDeactivated, nil
	case "deprecated":
		return ServiceDeprecated, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", s)
	}
}
