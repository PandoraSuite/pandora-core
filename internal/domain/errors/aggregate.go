package errors

import (
	"fmt"
	"strings"
)

var _ Error = (*AggregateError)(nil)

type AggregateError []Error

func (e AggregateError) Error() string {
	if len(e) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("<%s> %d errors:\n", e.Code(), len(e)))
	for i, err := range e {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(fmt.Sprintf("  %s", err.Error()))
	}
	return b.String()
}

func (e AggregateError) Code() ErrorCode {
	return CodeAggregate
}

func (e AggregateError) Unwrap() error {
	if len(e) == 0 {
		return nil
	}

	if len(e) == 1 {
		return e[0]
	}

	return e
}

func (e AggregateError) PriorityCode() ErrorCode {
	if len(e) == 0 {
		return CodeAggregate
	}

	best := e[0].Code()
	bestPriority := ErrorCodePriority[best]

	for _, err := range e[1:] {
		if priority, ok := ErrorCodePriority[err.Code()]; ok && priority < bestPriority {
			best = err.Code()
			bestPriority = priority
		}
	}

	return best
}

func NewAggregateError(errs ...Error) Error {
	if len(errs) == 0 {
		return nil
	}

	return (AggregateError)(errs)
}

func Aggregate(err, newErr Error) Error {
	if err == nil {
		return newErr
	}

	if newErr == nil {
		return err
	}

	if agg, ok := err.(AggregateError); ok {
		return append(agg, newErr)
	}

	if agg, ok := newErr.(AggregateError); ok {
		return append(AggregateError{err}, agg...)
	}

	return NewAggregateError(err, newErr)
}
