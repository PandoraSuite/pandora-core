package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var requestLogExecutionStatus = []string{
	"success", "failed", "unauthorized", "server error",
}

type RequestLog struct {
	ID              pgtype.Int4
	EnvironmentID   pgtype.Int4
	ServiceID       pgtype.Int4
	APIKey          pgtype.Text
	RequestTime     pgtype.Timestamptz
	ExecutionStatus pgtype.Text
	CreatedAt       pgtype.Timestamptz
}

func (p *RequestLog) ValidateModel() error {
	return p.validateExecutionStatus()
}

func (p *RequestLog) validateExecutionStatus() error {
	if executionStatus, _ := p.ExecutionStatus.Value(); executionStatus != nil {
		if slices.Contains(requestLogExecutionStatus, executionStatus.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid execution status: must be %s",
		strings.Join(requestLogExecutionStatus, ", "),
	)
}
