package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/MAD-py/pandora-core/internal/adapters/http/docs"
)

func setupSwagger(engine *gin.Engine) {
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
