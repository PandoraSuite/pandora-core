package enums

import (
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
	parsed, err := ParseClientType(str)
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

func (t *ClientType) UnmarshalJSON(b []byte) error {
	parsed, err := ParseClientType(string(b))
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

func (t *ClientType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.String())), nil
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
