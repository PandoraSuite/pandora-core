package enums

import (
	"encoding/json"
	"fmt"
)

type ServiceStatus int

const (
	ServiceStatusNull ServiceStatus = iota
	ServiceActive
	ServiceDeactivated
	ServiceDeprecated
)

func (s ServiceStatus) String() string {
	switch s {
	case ServiceStatusNull:
		return ""
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

func (s *ServiceStatus) Scan(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid status: %v", v)
	}

	status, err := ParseServiceStatus(str)
	if err != nil {
		return err
	}

	*s = status
	return nil
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
	return []byte(fmt.Sprintf("\"%s\"", s.String())), nil
}

func ParseServiceStatus(s string) (ServiceStatus, error) {
	switch s {
	case "":
		return ServiceStatusNull, nil
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
