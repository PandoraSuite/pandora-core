package errors

import "errors"

var (
	ErrNotFound             = errors.New("record not found")
	ErrPersistence          = errors.New("persistence error")
	ErrUndefinedEntity      = errors.New("undefined entity")
	ErrUniqueViolation      = errors.New("unique key violation")
	ErrNotNullViolation     = errors.New("an attempt was made to insert a NULL value in a mandatory column")
	ErrForeignKeyViolation  = errors.New("foreign key violation")
	ErrRestrictionViolation = errors.New("restriction violation")
)
