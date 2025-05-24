package enums

type RequestExecutionStatus string

const (
	RequestExecutionStatusNull          RequestExecutionStatus = ""
	RequestExecutionStatusSuccess       RequestExecutionStatus = "success"
	RequestExecutionStatusForwarded     RequestExecutionStatus = "forwarded"
	RequestExecutionStatusClientError   RequestExecutionStatus = "client_error"
	RequestExecutionStatusServiceError  RequestExecutionStatus = "service_error"
	RequestExecutionStatusUnauthorized  RequestExecutionStatus = "unauthorized"
	RequestExecutionStatusQuotaExceeded RequestExecutionStatus = "quota_exceeded"
)

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
