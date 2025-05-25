package enums

type ClientType string

const (
	ClientTypeNull     ClientType = ""
	ClientDeveloper    ClientType = "developer"
	ClientOrganization ClientType = "organization"
)
