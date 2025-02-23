package enums

import (
	"encoding/json"
	"fmt"
)

type EnvironmentStatus int

const (
	EnvironmentActive EnvironmentStatus = iota
	EnvironmentDeactivated
)

func (s EnvironmentStatus) String() string {
	switch s {
	case EnvironmentActive:
		return "active"
	case EnvironmentDeactivated:
		return "deactivated"
	default:
		panic("unknown EnvironmentStatus")
	}
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
	case "active":
		return EnvironmentActive, nil
	case "deactivated":
		return EnvironmentDeactivated, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", s)
	}
}
