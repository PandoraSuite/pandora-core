package enums

import (
	"encoding/json"
	"fmt"
)

type ReserveExecutionStatusCode int

const (
	ReserveExecutionStatusNull ReserveExecutionStatusCode = iota
	ReserveExecutionStatusKeyNotFound
	ReserveExecutionStatusInvalidKey
	ReserveExecutionStatusDeactivatedKey
	ReserveExecutionStatusExpiredKey
	ReserveExecutionStatusServiceNotFound
	ReserveExecutionStatusDeactivatedService
	ReserveExecutionStatusDeprecatedService
	ReserveExecutionStatusEnvironmentNotFound
	ReserveExecutionStatusDeactivatedEnvironment
	ReserveExecutionStatusExceededRequests
	ReserveExecutionStatusActiveReservations

	ReservationExecutionStatusNotFound
	ReservationExecutionStatusInvalidService
	ReservationExecutionStatusInvalidServiceVersion
	ReservationExecutionStatusServiceNotActive
	ReservationExecutionStatusInvalidEnvironment
	ReservationExecutionStatusEnvironmentNotActive

	ValidateStatusInvalidEnvironmentKey
	ValidateStatusEnvironmentServiceInvalid
)

func (es ReserveExecutionStatusCode) String() string {
	switch es {
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
	case ReserveExecutionStatusEnvironmentNotFound:
		return "ENVIRONMENT_NOT_FOUND"
	case ReserveExecutionStatusDeactivatedEnvironment:
		return "DEACTIVATED_ENVIRONMENT"
	case ReserveExecutionStatusExceededRequests:
		return "EXCEEDED_AVAILABLE_REQUEST"
	case ReserveExecutionStatusActiveReservations:
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
	case "ENVIRONMENT_NOT_FOUND":
		return ReserveExecutionStatusEnvironmentNotFound, nil
	case "DEACTIVATED_ENVIRONMENT":
		return ReserveExecutionStatusDeactivatedEnvironment, nil
	case "EXCEEDED_AVAILABLE_REQUEST":
		return ReserveExecutionStatusExceededRequests, nil
	case "ACTIVE_RESERVATIONS":
		return ReserveExecutionStatusActiveReservations, nil
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
