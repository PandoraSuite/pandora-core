package enums

import (
	"encoding/json"
	"fmt"
)

type APIKeyStatus int

const (
	APIKeyActive APIKeyStatus = iota
	APIKeyDeactivated
)

func (s APIKeyStatus) String() string {
	switch s {
	case APIKeyActive:
		return "active"
	case APIKeyDeactivated:
		return "deactivated"
	default:
		panic("unknown APIKeyStatus")
	}
}

func (s *APIKeyStatus) UnmarshalJSON(b []byte) error {
	var ss string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}
	parsed, err := ParseAPIKeyStatus(ss)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

func (s *APIKeyStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func ParseAPIKeyStatus(s string) (APIKeyStatus, error) {
	switch s {
	case "active":
		return APIKeyActive, nil
	case "deactivated":
		return APIKeyDeactivated, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", s)
	}
}
