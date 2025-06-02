package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/version"
)

func VersionHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", fmt.Sprintf("pandora/%s", version.Version))
		c.Next()
	}
}
