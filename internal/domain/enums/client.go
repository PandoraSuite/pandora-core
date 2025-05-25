package enums

type ClientType string

const (
	ClientTypeNull         ClientType = ""
	ClientTypeDeveloper    ClientType = "developer"
	ClientTypeOrganization ClientType = "organization"
)
