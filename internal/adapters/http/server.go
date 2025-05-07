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

func (s *Server) Run(exposeVersion bool) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	if exposeVersion {
		router.Use(middlewares.VersionHeader())
	}

	setupSwagger(router)

	v1 := router.Group("/api/v1")

	{
		routes.RegisterLoginRoutes(v1)
	}

	v1Protected := v1.Group("")
	v1Protected.Use(
		middlewares.ValidateToken(
			auth.NewTokenValidationUseCase(nil, nil),
		),
	)

	{
		routes.RegisterAuthRoutes(v1Protected)
	}

	v1Protected.Use(
		middlewares.ForcePasswordReset(
			auth.NewResetPasswordUseCase(nil, nil),
		),
	)

	{
		routes.RegisterServiceRoutes(v1Protected)
		routes.RegisterClientRoutes(v1Protected)
		routes.RegisterProjectRoutes(v1Protected)
		routes.RegisterEnvironmentRoutes(v1Protected)
		routes.RegisterAPIKeyRoutes(v1Protected)
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
