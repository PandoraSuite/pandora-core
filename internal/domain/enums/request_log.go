package enums

import (
	"encoding/json"
	"fmt"
)

type RequestLogExecutionStatus int

const (
	RequestLogSuccess RequestLogExecutionStatus = iota
	RequestLogFailed
	RequestLogUnauthorized
	RequestLogServerError
)

func (es RequestLogExecutionStatus) String() string {
	switch es {
	case RequestLogSuccess:
		return "success"
	case RequestLogFailed:
		return "failed"
	case RequestLogUnauthorized:
		return "unauthorized"
	case RequestLogServerError:
		return "server error"
	default:
		panic("unknown RequestLogExecutionStatus")
	}
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
	return json.Marshal(es.String())
}

func ParseRequestLogExecutionStatus(es string) (RequestLogExecutionStatus, error) {
	switch es {
	case "success":
		return RequestLogSuccess, nil
	case "failed":
		return RequestLogFailed, nil
	case "unauthorized":
		return RequestLogUnauthorized, nil
	case "server error":
		return RequestLogServerError, nil
	default:
		return 0, fmt.Errorf("invalid execution status: %s", es)
	}
}
