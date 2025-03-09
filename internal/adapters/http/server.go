package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/MAD-py/pandora-core/docs"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/middleware"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

// @title Pandora Core
// @version 1.0
// @description API for centralized API key management and service access control.
// @termsOfService http://example.com/terms/

// @contact.name Pandora Core Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl /api/v1/auth/login

type Server struct {
	addr string

	server *http.Server

	srvService         inbound.ServiceHTTPPort
	authService        inbound.AuthHTTPPort
	apiKeyService      inbound.APIKeyHTTPPort
	clientService      inbound.ClientHTTPPort
	projectService     inbound.ProjectHTTPPort
	environmentService inbound.EnvironmentHTTPPort
}

func (srv *Server) setupRoutes(router *gin.RouterGroup) {

	auth := router.Group("/auth")
	{
		auth.POST("/login", handlers.Authenticate(srv.authService))
	}

	protected := router.Group("")
	protected.Use(middleware.ValidateToken(srv.authService))
	{
		auth := protected.Group("/auth")
		{
			auth.POST("/change-password", handlers.ChangePassword(srv.authService))
		}

		protected.Use(middleware.ForcePasswordReset(srv.authService))
		apiKeys := protected.Group("/api-keys")
		{
			apiKeys.POST("", handlers.CreateAPIKey(srv.apiKeyService))
		}

		environments := protected.Group("/environments")
		{
			environments.GET(
				"/:id/api-keys",
				handlers.GetAPIKeysByEnvironment(srv.apiKeyService),
			)
		}

		services := protected.Group("/services")
		{
			services.POST("", handlers.CreateService(srv.srvService))
			services.GET("", handlers.GetAllServices(srv.srvService))
			services.GET("/active", handlers.GetActiveServices(srv.srvService))
		}

		projects := protected.Group("/projects")
		{
			projects.POST("", handlers.CreateProject(srv.projectService))
			projects.POST(
				"/:project_id/services/:service_id/assign",
				handlers.AssignServiceToProject(srv.projectService),
			)
		}

		clients := protected.Group("/clients")
		{
			clients.POST("", handlers.CreateClient(srv.clientService))
			clients.GET("", handlers.GetAllClients(srv.clientService))
			clients.GET(
				":id/projects",
				handlers.GetProjectsByClient(srv.projectService),
			)
		}
	}
}

func (srv *Server) Run() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv.setupRoutes(router.Group("/api/v1"))

	srv.server = &http.Server{
		Addr:    srv.addr,
		Handler: router,
	}

	log.Printf("[INFO] API is running on port: %s\n", srv.addr)
	log.Printf("[INFO] Pandora Core is fully initialized and ready to accept requests.\n\n")
	if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[ERROR] Failed to start server: %v", err)
	}
}

func NewServer(
	addr string,
	srvService inbound.ServiceHTTPPort,
	authService inbound.AuthHTTPPort,
	apiKeyService inbound.APIKeyHTTPPort,
	clientService inbound.ClientHTTPPort,
	projectService inbound.ProjectHTTPPort,
	environmentService inbound.EnvironmentHTTPPort,
) *Server {
	return &Server{
		addr:               addr,
		srvService:         srvService,
		authService:        authService,
		apiKeyService:      apiKeyService,
		clientService:      clientService,
		projectService:     projectService,
		environmentService: environmentService,
	}
}
