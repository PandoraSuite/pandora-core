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
	return ErrorCodeAggregate
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
