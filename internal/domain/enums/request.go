package enums

type RequestExecutionStatus string

const (
	RequestExecutionStatusNull         RequestExecutionStatus = ""
	RequestExecutionStatusSuccess      RequestExecutionStatus = "success"
	RequestExecutionStatusForwarded    RequestExecutionStatus = "forwarded"
	RequestExecutionStatusClientError  RequestExecutionStatus = "client_error"
	RequestExecutionStatusServiceError RequestExecutionStatus = "service_error"
	RequestExecutionStatusUnauthorized RequestExecutionStatus = "unauthorized"
)

func ParseRequestExecutionStatus(status string) (RequestExecutionStatus, bool) {
	switch s := RequestExecutionStatus(status); s {
	case RequestExecutionStatusNull,
		RequestExecutionStatusSuccess,
		RequestExecutionStatusForwarded,
		RequestExecutionStatusClientError,
		RequestExecutionStatusServiceError,
		RequestExecutionStatusUnauthorized:
		return s, true
	default:
		return RequestExecutionStatusNull, false
	}
}

type RequestBodyContentType string

const (
	RequestBodyContentTypeNull        RequestBodyContentType = ""
	RequestBodyContentTypeXML         RequestBodyContentType = "application/xml"
	RequestBodyContentTypeJSON        RequestBodyContentType = "application/json"
	RequestBodyContentTypeText        RequestBodyContentType = "text/plain"
	RequestBodyContentTypeHTML        RequestBodyContentType = "text/html"
	RequestBodyContentTypeForm        RequestBodyContentType = "multipart/form-data"
	RequestBodyContentTypeFormURL     RequestBodyContentType = "application/x-www-form-urlencoded"
	RequestBodyContentTypeOctetStream RequestBodyContentType = "application/octet-stream"
)

func ParseRequestBodyContentType(status string) (RequestBodyContentType, bool) {
	switch t := RequestBodyContentType(status); t {
	case RequestBodyContentTypeNull,
		RequestBodyContentTypeXML,
		RequestBodyContentTypeJSON,
		RequestBodyContentTypeText,
		RequestBodyContentTypeHTML,
		RequestBodyContentTypeForm,
		RequestBodyContentTypeFormURL,
		RequestBodyContentTypeOctetStream:
		return t, true
	default:
		return RequestBodyContentTypeNull, false
	}
}
