package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/MAD-py/pandora-core/docs"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

type Server struct {
	addr string

	server *http.Server

	srvService    inbound.ServiceHTTPPort
	apiKeyService inbound.APIKeyHTTPPort
	// clientService inbound.ClientHTTPPort
	// environmentService inbound.EnvironmentHTTPPort
	// projectService inbound.ProjectHTTPPort
}

func (srv *Server) setupRoutes(router *gin.RouterGroup) {
	apiKeys := router.Group("/api-keys")
	{
		apiKeys.POST("/", handlers.CreateAPIKey(srv.apiKeyService))
	}

	environments := router.Group("/environments")
	{
		environments.GET(
			"/:environment_id/api-keys",
			handlers.GetAPIKeysByEnvironment(srv.apiKeyService),
		)
	}

	services := router.Group("/services")
	{
		services.POST("", handlers.CreateService(srv.srvService))
		services.GET("", handlers.GetAllServices(srv.srvService))
		services.GET("/active", handlers.GetActiveServices(srv.srvService))
	}
}

func (srv *Server) Run() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv.setupRoutes(router.Group("/api/v1"))

	srv.server = &http.Server{
		Addr:    srv.addr,
		Handler: router,
	}

	log.Printf("HTTP Server running on :%s\n", srv.addr)
	if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP Server error: %v\n", err)
	}
}

func NewServer(
	addr string,
	srvService inbound.ServiceHTTPPort,
	apiKeyService inbound.APIKeyHTTPPort,
) *Server {
	return &Server{
		addr:          addr,
		srvService:    srvService,
		apiKeyService: apiKeyService,
	}
}
