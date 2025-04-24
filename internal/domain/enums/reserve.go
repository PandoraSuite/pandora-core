package enums

import (
	"encoding/json"
	"fmt"
)

type ValidateStatusCode int

const (
	ValidateStatusNull ValidateStatusCode = iota
	ValidateStatusDeactivatedKey
	ValidateStatusExpiredKey
	ValidateStatusServiceNotFound
	ValidateStatusDeactivatedEnvironment
	ValidateStatusExceededRequests
	ValidateStatusActiveReservations
	ValidateStatusInvalidEnvironmentKey
	ValidateStatusEnvironmentServiceInvalid

	ReservationExecutionStatusNotFound
	ReservationExecutionStatusInvalidService
	ReservationExecutionStatusInvalidServiceVersion
	ReservationExecutionStatusServiceNotActive
	ReservationExecutionStatusInvalidEnvironment
	ReservationExecutionStatusEnvironmentNotActive

	ValidateStatusKeyNotFound
	ValidateStatusInvalidKey
)

func (es ValidateStatusCode) String() string {
	switch es {
	case ValidateStatusKeyNotFound:
		return "KEY_NOT_FOUND"
	case ValidateStatusInvalidKey:
		return "INVALID_KEY"
	case ValidateStatusDeactivatedKey:
		return "DEACTIVATED_KEY"
	case ValidateStatusExpiredKey:
		return "EXPIRED_KEY"
	case ValidateStatusServiceNotFound:
		return "SERVICE_NOT_FOUND"
	case ValidateStatusDeactivatedEnvironment:
		return "DEACTIVATED_ENVIRONMENT"
	case ValidateStatusExceededRequests:
		return "EXCEEDED_AVAILABLE_REQUEST"
	case ValidateStatusActiveReservations:
		return "ACTIVE_RESERVATIONS"
	case ReservationExecutionStatusNotFound:
		return "RESERVATION_NOT_FOUND"
	case ReservationExecutionStatusInvalidService:
		return "INVALID_SERVICE"
	case ReservationExecutionStatusInvalidServiceVersion:
		return "INVALID_SERVICE_VERSION"
	case ReservationExecutionStatusServiceNotActive:
		return "SERVICE_NOT_ACTIVE"
	case ReservationExecutionStatusInvalidEnvironment:
		return "INVALID_ENVIRONMENT"
	case ReservationExecutionStatusEnvironmentNotActive:
		return "ENVIRONMENT_NOT_ACTIVE"
	case ValidateStatusInvalidEnvironmentKey:
		return "INVALID_ENVIRONMENT_KEY"
	case ValidateStatusEnvironmentServiceInvalid:
		return "ENVIRONMENT_SERVICE_INVALID"
	default:
		panic("unknown ValidateStatusCode")
	}
}

func (s *ValidateStatusCode) Scan(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid reserve execution status code: %v", v)
	}

	status, err := ParseValidateStatusCode(str)
	if err != nil {
		return err
	}

	*s = status
	return nil
}

func (es *ValidateStatusCode) UnmarshalJSON(b []byte) error {
	var ss string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}

	parsed, err := ParseValidateStatusCode(ss)
	if err != nil {
		return err
	}

	*es = parsed
	return nil
}

func (es *ValidateStatusCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(es.String())
}

func ParseValidateStatusCode(es string) (ValidateStatusCode, error) {
	switch es {
	case "KEY_NOT_FOUND":
		return ValidateStatusKeyNotFound, nil
	case "INVALID_KEY":
		return ValidateStatusInvalidKey, nil
	case "DEACTIVATED_KEY":
		return ValidateStatusDeactivatedKey, nil
	case "EXPIRED_KEY":
		return ValidateStatusExpiredKey, nil
	case "SERVICE_NOT_FOUND":
		return ValidateStatusServiceNotFound, nil
	case "DEACTIVATED_ENVIRONMENT":
		return ValidateStatusDeactivatedEnvironment, nil
	case "EXCEEDED_AVAILABLE_REQUEST":
		return ValidateStatusExceededRequests, nil
	case "ACTIVE_RESERVATIONS":
		return ValidateStatusActiveReservations, nil
	case "RESERVATION_NOT_FOUND":
		return ReservationExecutionStatusNotFound, nil
	case "INVALID_SERVICE":
		return ReservationExecutionStatusInvalidService, nil
	case "INVALID_SERVICE_VERSION":
		return ReservationExecutionStatusInvalidServiceVersion, nil
	case "SERVICE_NOT_ACTIVE":
		return ReservationExecutionStatusServiceNotActive, nil
	case "INVALID_ENVIRONMENT":
		return ReservationExecutionStatusInvalidEnvironment, nil
	case "ENVIRONMENT_NOT_ACTIVE":
		return ReservationExecutionStatusEnvironmentNotActive, nil
	case "INVALID_ENVIRONMENT_KEY":
		return ValidateStatusInvalidEnvironmentKey, nil
	case "ENVIRONMENT_SERVICE_INVALID":
		return ValidateStatusEnvironmentServiceInvalid, nil
	default:
		return 0, fmt.Errorf("invalid reserve execution status code: %s", es)
	}
}
