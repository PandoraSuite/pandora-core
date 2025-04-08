package enums

import (
	"encoding/json"
	"fmt"
)

type ReserveExecutionStatusCode int

const (
	ReserveExecutionStatusNull ReserveExecutionStatusCode = iota
	ReserveExecutionStatusOk
	ReserveExecutionStatusKeyNotFound
	ReserveExecutionStatusInvalidKey
	ReserveExecutionStatusDeactivatedKey
	ReserveExecutionStatusExpiredKey
	ReserveExecutionStatusServiceNotFound
	ReserveExecutionStatusDeactivatedService
	ReserveExecutionStatusDeprecatedService
	ReserveExecutionStatusExceededRequests
	ReserveExecutionStatusActiveReservations
)

func (es ReserveExecutionStatusCode) String() string {
	switch es {
	case ReserveExecutionStatusOk:
		return "OK"
	case ReserveExecutionStatusKeyNotFound:
		return "KEY_NOT_FOUND"
	case ReserveExecutionStatusInvalidKey:
		return "INVALID_KEY"
	case ReserveExecutionStatusDeactivatedKey:
		return "DEACTIVATED_KEY"
	case ReserveExecutionStatusExpiredKey:
		return "EXPIRED_KEY"
	case ReserveExecutionStatusServiceNotFound:
		return "SERVICE_NOT_FOUND"
	case ReserveExecutionStatusDeactivatedService:
		return "DEACTIVATED_SERVICE"
	case ReserveExecutionStatusDeprecatedService:
		return "DEPRECATED_SERVICE"
	case ReserveExecutionStatusExceededRequests:
		return "EXCEEDED_AVAILABLE_REQUEST"
	case ReserveExecutionStatusActiveReservations:
		return "ACTIVE_RESERVATIONS"
	default:
		panic("unknown ReserveExecutionStatusCode")
	}
}

func (s *ReserveExecutionStatusCode) Scan(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid reserve execution status code: %v", v)
	}

	status, err := ParseReserveExecutionStatusCode(str)
	if err != nil {
		return err
	}

	*s = status
	return nil
}

func (es *ReserveExecutionStatusCode) UnmarshalJSON(b []byte) error {
	var ss string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}

	parsed, err := ParseReserveExecutionStatusCode(ss)
	if err != nil {
		return err
	}

	*es = parsed
	return nil
}

func (es *ReserveExecutionStatusCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(es.String())
}

func ParseReserveExecutionStatusCode(es string) (ReserveExecutionStatusCode, error) {
	switch es {
	case "OK":
		return ReserveExecutionStatusOk, nil
	case "KEY_NOT_FOUND":
		return ReserveExecutionStatusKeyNotFound, nil
	case "INVALID_KEY":
		return ReserveExecutionStatusInvalidKey, nil
	case "DEACTIVATED_KEY":
		return ReserveExecutionStatusDeactivatedKey, nil
	case "EXPIRED_KEY":
		return ReserveExecutionStatusExpiredKey, nil
	case "SERVICE_NOT_FOUND":
		return ReserveExecutionStatusServiceNotFound, nil
	case "DEACTIVATED_SERVICE":
		return ReserveExecutionStatusDeactivatedService, nil
	case "DEPRECATED_SERVICE":
		return ReserveExecutionStatusDeprecatedService, nil
	case "EXCEEDED_AVAILABLE_REQUEST":
		return ReserveExecutionStatusExceededRequests, nil
	case "ACTIVE_RESERVATIONS":
		return ReserveExecutionStatusActiveReservations, nil
	default:
		return 0, fmt.Errorf("invalid reserve execution status code: %s", es)
	}
}
