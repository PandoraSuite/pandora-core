package errors

import (
	"fmt"
	"strings"
)

func newErrorWithAttrs(
	code ErrorCode,
	entity, reason string,
	attrs map[string]any,
) Error {
	message := fmt.Sprintf("%s %s", entity, reason)

	kv := make([]string, 0, len(attrs))
	for k, v := range attrs {
		kv = append(kv, fmt.Sprintf("%s=%v", k, v))
	}

	internalMsg := fmt.Sprintf(
		"%s: %s; %s",
		entity, reason,
		strings.Join(kv, ", "),
	)

	return newDomainError(code, message, internalMsg)
}

func NewNotFound(entity string, idOrAttrs any) Error {
	return newErrorWithAttrs(
		CodeNotFound, entity, "not found", normalizeAttrs(idOrAttrs),
	)
}

func NewAlreadyExists(entity string, idOrAttrs any) Error {
	return newErrorWithAttrs(
		CodeAlreadyExists, entity, "already exists", normalizeAttrs(idOrAttrs),
	)
}

func NewValidationFailed(entity string, field string, detail string) Error {
	attrs := map[string]any{"field": field, "detail": detail}
	return newErrorWithAttrs(
		CodeValidationFailed, entity, "validation failed", attrs,
	)
}

func NewForbidden(entity, permission string, attrs map[string]any) Error {
	if attrs == nil {
		attrs = map[string]any{}
	}
	attrs["permission"] = permission
	return newErrorWithAttrs(
		CodeForbidden, entity, "forbidden", attrs,
	)
}

func NewUnauthorized(reason string) Error {
	return newDomainError(
		CodeUnauthorized, "Unauthorized", fmt.Sprintf("auth: %s", reason),
	)
}

func NewInternal(reason string, err error) Error {
	return wrap(CodeInternal, reason, err.Error(), err)
}

func normalizeAttrs(idOrAttrs any) map[string]any {
	switch v := idOrAttrs.(type) {
	case map[string]any:
		return v
	default:
		return map[string]any{"id": v}
	}
}
