package enums

import (
	"encoding/json"
	"fmt"
)

type ClientType int

const (
	ClientTypeNull ClientType = iota
	ClientDeveloper
	ClientOrganization
)

func (t ClientType) String() string {
	switch t {
	case ClientTypeNull:
		return ""
	case ClientDeveloper:
		return "developer"
	case ClientOrganization:
		return "organization"
	default:
		panic("unknown ClientType")
	}
}

func (t *ClientType) Scan(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid type: %v", v)
	}

	ct, err := ParseClientType(str)
	if err != nil {
		return err
	}

	*t = ct
	return nil
}

func (t *ClientType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := ParseClientType(s)
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

func (t *ClientType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func ParseClientType(t string) (ClientType, error) {
	switch t {
	case "":
		return ClientTypeNull, nil
	case "developer":
		return ClientDeveloper, nil
	case "organization":
		return ClientOrganization, nil
	default:
		return 0, fmt.Errorf("invalid type: %s", t)
	}
}
