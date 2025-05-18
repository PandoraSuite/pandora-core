package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/errors"
	domainerr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		raw := c.Errors.Last().Err
		httpError := errors.MapToHTTPError(raw)

		code := httpError.Code
		if httpError.Code == domainerr.CodeAggregate {
			code = httpError.PriorityCode()
		}

		status := errors.CodeToStatusCode(code)

		c.JSON(status, httpError)
	}
}
