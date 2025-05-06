package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/middlewares"
	"github.com/MAD-py/pandora-core/internal/adapters/http/routes"
	"github.com/MAD-py/pandora-core/internal/app/auth"
)

// @title Pandora Core
// @version 1.0
// @description API for centralized API key management and service access control.
// @termsOfService http://example.com/terms/

// @tag.name Authentication
// @tag.name Services
// @tag.name Clients
// @tag.name Projects
// @tag.name Environments
// @tag.name API Keys

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
}

func (s *Server) setupRoutes(router *gin.RouterGroup) {

	// auth := router.Group("/auth")
	// {
	// 	auth.POST("/login", handlers.Authenticate(s.authService))
	// }

	// protected := router.Group("")
	// protected.Use(middleware.ValidateToken(s.authService))
	// {
	// 	auth := protected.Group("/auth")
	// 	{
	// 		auth.POST(
	// 			"/change-password", handlers.ChangePassword(s.authService),
	// 		)
	// 	}

	// 	protected.Use(middleware.ForcePasswordReset(s.authService))

	// 	services := protected.Group("/services")
	// 	{
	// 		services.GET("", handlers.GetAllServices(s.srvService))
	// 		services.POST("", handlers.CreateService(s.srvService))
	// 		services.DELETE("/:id", handlers.DeleteService(s.srvService))
	// 		services.PATCH(
	// 			"/:id/status",
	// 			handlers.UpdateStatusService(s.srvService),
	// 		)
	// 	}

	// 	clients := protected.Group("/clients")
	// 	{
	// 		clients.GET("", handlers.GetAllClients(s.clientService))
	// 		clients.POST("", handlers.CreateClient(s.clientService))
	// 		clients.GET("/:id", handlers.GetClient(s.clientService))
	// 		clients.PATCH("/:id", handlers.UpdateClient(s.clientService))
	// 		clients.GET(
	// 			"/:id/projects",
	// 			handlers.GetProjectsByClient(s.clientService),
	// 		)
	// 	}

	// 	projects := protected.Group("/projects")
	// 	{
	// 		projects.POST("", handlers.CreateProject(s.projectService))
	// 		projects.GET("/:id", handlers.GetProject(s.projectService))
	// 		projects.PATCH("/:id", handlers.UpdateProject(s.projectService))
	// 		projects.GET(
	// 			"/:id/environments",
	// 			handlers.GetEnvironmentsByProject(s.projectService),
	// 		)
	// 		projects.POST(
	// 			"/:id/services",
	// 			handlers.AssignServiceToProject(s.projectService),
	// 		)
	// 		projects.DELETE(
	// 			"/:id/services/:service_id",
	// 			handlers.RemoveServiceFromProject(s.projectService),
	// 		)
	// 		projects.PATCH(
	// 			"/:id/services/:service_id",
	// 			handlers.UpdateProjectService(s.projectService),
	// 		)
	// 		projects.POST(
	// 			"/:id/services/:service_id/reset-requests",
	// 			handlers.ResetServiceAvailableRequests(s.projectService),
	// 		)
	// 	}

	// 	environments := protected.Group("/environments")
	// 	{
	// 		environments.POST(
	// 			"", handlers.CreateEnvironment(s.environmentService),
	// 		)
	// 		environments.GET(
	// 			"/:id", handlers.GetEnvironment(s.environmentService),
	// 		)
	// 		environments.PATCH(
	// 			"/:id", handlers.UpdateEnvironment(s.environmentService),
	// 		)
	// 		environments.GET(
	// 			"/:id/api-keys",
	// 			handlers.GetAPIKeysByEnvironment(s.apiKeyService),
	// 		)
	// 		environments.POST(
	// 			"/:id/services",
	// 			handlers.AssignServiceToEnvironment(s.environmentService),
	// 		)
	// 		environments.DELETE(
	// 			"/:id/services/:service_id",
	// 			handlers.RemoveServiceFromEnvironment(s.environmentService),
	// 		)
	// 		environments.PATCH(
	// 			"/:id/services/:service_id",
	// 			handlers.UpdateEnvironmentService(s.environmentService),
	// 		)
	// 		environments.PATCH(
	// 			"/:id/services/:service_id/reset-requests",
	// 			handlers.ResetServiceRequestsFromEnvironment(s.environmentService),
	// 		)
	// 	}

	// 	apiKeys := protected.Group("/api-keys")
	// 	{
	// 		apiKeys.POST("", handlers.CreateAPIKey(s.apiKeyService))
	// 		apiKeys.PATCH("/:id", handlers.UpdateAPIKey(s.apiKeyService))
	// 	}

	// }
}

func (s *Server) Run(exposeVersion bool) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	if exposeVersion {
		router.Use(middlewares.VersionHeader())
	}

	setupSwagger(router)

	v1 := router.Group("/api/v1")

	{
		routes.RegisterAuthRoutes(v1)
	}

	v1Protected := v1.Group("")
	v1Protected.Use(
		middlewares.ValidateToken(
			auth.NewTokenValidationUseCase(nil, nil),
		),
	)

	{
		routes.RegisterProtectedAuthRoutes(v1Protected)
	}

	v1Protected.Use(
		middlewares.ForcePasswordReset(
			auth.NewResetPasswordUseCase(nil, nil),
		),
	)

	{
		routes.RegisterProtectedServiceRoutes(v1Protected)
		routes.RegisterProtectedAPIKeyRoutes(v1Protected)
	}

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

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}
