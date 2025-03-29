package enums

import (
	"encoding/json"
	"fmt"
)

type EnvironmentStatus int

const (
	EnvironmentStatusNull EnvironmentStatus = iota
	EnvironmentActive
	EnvironmentDeactivated
)

func (s EnvironmentStatus) String() string {
	switch s {
	case EnvironmentStatusNull:
		return ""
	case EnvironmentActive:
		return "active"
	case EnvironmentDeactivated:
		return "deactivated"
	default:
		panic("unknown EnvironmentStatus")
	}
}

func (s *EnvironmentStatus) Scan(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid status: %v", v)
	}

	status, err := ParseEnvironmentStatus(str)
	if err != nil {
		return err
	}

	*s = status
	return nil
}

func (s *EnvironmentStatus) UnmarshalJSON(b []byte) error {
	var ss string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}

	parsed, err := ParseEnvironmentStatus(ss)
	if err != nil {
		return err
	}

	*s = parsed
	return nil
}

func (s *EnvironmentStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func ParseEnvironmentStatus(s string) (EnvironmentStatus, error) {
	switch s {
	case "":
		return EnvironmentStatusNull, nil
	case "active":
		return EnvironmentActive, nil
	case "deactivated":
		return EnvironmentDeactivated, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", s)
	}
}
