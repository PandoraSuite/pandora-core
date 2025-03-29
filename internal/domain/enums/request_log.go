package enums

import (
	"encoding/json"
	"fmt"
)

type RequestLogExecutionStatus int

const (
	RequestLogExecutionStatusNull RequestLogExecutionStatus = iota
	RequestLogSuccess
	RequestLogFailed
	RequestLogPending
	RequestLogUnauthorized
	RequestLogServerError
)

func (es RequestLogExecutionStatus) String() string {
	switch es {
	case RequestLogExecutionStatusNull:
		return ""
	case RequestLogSuccess:
		return "success"
	case RequestLogFailed:
		return "failed"
	case RequestLogPending:
		return "pending"
	case RequestLogUnauthorized:
		return "unauthorized"
	case RequestLogServerError:
		return "server error"
	default:
		panic("unknown RequestLogExecutionStatus")
	}
}

func (s *RequestLogExecutionStatus) Scan(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid execution status: %v", v)
	}

	status, err := ParseRequestLogExecutionStatus(str)
	if err != nil {
		return err
	}

	*s = status
	return nil
}

func (es *RequestLogExecutionStatus) UnmarshalJSON(b []byte) error {
	var ss string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}
	parsed, err := ParseRequestLogExecutionStatus(ss)
	if err != nil {
		return err
	}
	*es = parsed
	return nil
}

func (es *RequestLogExecutionStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", es.String())), nil
}

func ParseRequestLogExecutionStatus(es string) (RequestLogExecutionStatus, error) {
	switch es {
	case "":
		return RequestLogExecutionStatusNull, nil
	case "success":
		return RequestLogSuccess, nil
	case "failed":
		return RequestLogFailed, nil
	case "pending":
		return RequestLogPending, nil
	case "unauthorized":
		return RequestLogUnauthorized, nil
	case "server error":
		return RequestLogServerError, nil
	default:
		return 0, fmt.Errorf("invalid execution status: %s", es)
	}
}
