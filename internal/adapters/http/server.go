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

func (s *Server) setupRoutes(router *gin.RouterGroup) {

	auth := router.Group("/auth")
	{
		auth.POST("/login", handlers.Authenticate(s.authService))
	}

	protected := router.Group("")
	protected.Use(middleware.ValidateToken(s.authService))
	{
		auth := protected.Group("/auth")
		{
			auth.POST(
				"/change-password", handlers.ChangePassword(s.authService),
			)
		}

		protected.Use(middleware.ForcePasswordReset(s.authService))
		apiKeys := protected.Group("/api-keys")
		{
			apiKeys.POST("", handlers.CreateAPIKey(s.apiKeyService))
		}

		environments := protected.Group("/environments")
		{
			environments.POST(
				"", handlers.CreateEnvironment(s.environmentService),
			)
			environments.GET(
				"/:id/api-keys",
				handlers.GetAPIKeysByEnvironment(s.apiKeyService),
			)
			environments.POST(
				"/:id/services",
				handlers.AssignServiceToEnvironment(s.environmentService),
			)
		}

		services := protected.Group("/services")
		{
			services.POST("", handlers.CreateService(s.srvService))
			services.GET("", handlers.GetAllServices(s.srvService))
		}

		projects := protected.Group("/projects")
		{
			projects.POST("", handlers.CreateProject(s.projectService))
			projects.GET(
				"/:id/environments",
				handlers.GetEnvironmentsByProject(s.projectService),
			)
			projects.POST(
				"/:id/services",
				handlers.AssignServiceToProject(s.projectService),
			)
		}

		clients := protected.Group("/clients")
		{
			clients.POST("", handlers.CreateClient(s.clientService))
			clients.GET("", handlers.GetAllClients(s.clientService))
			clients.GET(
				":id/projects",
				handlers.GetProjectsByClient(s.clientService),
			)
		}
	}
}

func (s *Server) Run(exposeVersion bool) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	if exposeVersion {
		router.Use(middleware.VersionHeader())
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.setupRoutes(router.Group("/api/v1"))

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	log.Printf("[INFO] API is running on port: %s\n", s.addr)
	log.Printf("[INFO] Pandora Core is fully initialized and ready to accept requests.\n\n")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
