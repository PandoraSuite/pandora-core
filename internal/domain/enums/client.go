package enums

type ClientType string

const (
	ClientTypeNull         ClientType = ""
	ClientTypeDeveloper    ClientType = "developer"
	ClientTypeOrganization ClientType = "organization"
)

func ParseClientType(status string) (ClientType, bool) {
	switch t := ClientType(status); t {
	case ClientTypeNull, ClientTypeDeveloper, ClientTypeOrganization:
		return t, true
	default:
		return ClientTypeNull, false
	}
}
