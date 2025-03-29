package enums

import (
	"encoding/json"
	"fmt"
)

type APIKeyStatus int

const (
	APIKeyStatusNull APIKeyStatus = iota
	APIKeyActive
	APIKeyDeactivated
)

func (s APIKeyStatus) String() string {
	switch s {
	case APIKeyStatusNull:
		return ""
	case APIKeyActive:
		return "active"
	case APIKeyDeactivated:
		return "deactivated"
	default:
		panic("unknown APIKeyStatus")
	}
}

func (s *APIKeyStatus) Scan(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid status: %v", v)
	}

	status, err := ParseAPIKeyStatus(str)
	if err != nil {
		return err
	}

	*s = status
	return nil
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
	case "":
		return APIKeyStatusNull, nil
	case "active":
		return APIKeyActive, nil
	case "deactivated":
		return APIKeyDeactivated, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", s)
	}
}
