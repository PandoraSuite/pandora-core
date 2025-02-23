package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func IntToPgtypeInt4(value int) pgtype.Int4 {
	if value == 0 {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(value), Valid: true}
}

func StringToPgtypeText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: value, Valid: true}
}

func TimeToPgtypeTimestamptz(value time.Time) pgtype.Timestamptz {
	if value.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: value, Valid: true}
}

func PgtypeInt4ToInt(field pgtype.Int4) int {
	if field.Valid {
		return int(field.Int32)
	}
	return 0
}

func PgtypeTextToString(field pgtype.Text) string {
	if field.Valid {
		return field.String
	}
	return ""
}

func PgtypeTimestamptzToTime(field pgtype.Timestamptz) time.Time {
	if field.Valid {
		return field.Time
	}
	return time.Time{}
}
