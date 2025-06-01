package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
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

// @securitydefinitions.apikey ScopedToken
// @in header
// @name Authorization

type Server struct {
	addr string

	exposeVersion bool

	server *http.Server

	deps *bootstrap.Dependencies
}

func (s *Server) Run() error {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	setupSwagger(engine)

	if s.exposeVersion {
		engine.Use(middlewares.VersionHeader())
	}

	engine.Use(middlewares.ErrorHandler())

	v1 := engine.Group("/api/v1")

	{
		routes.RegisterLoginRoutes(v1, s.deps)
	}

	v1Protected := v1.Group("")
	v1Protected.Use(
		middlewares.ValidateAccessToken(
			auth.NewTokenValidationUseCase(
				s.deps.Validator, s.deps.TokenProvider,
			),
		),
	)

	{
		routes.RegisterAuthRoutes(v1Protected, s.deps)
	}

	v1Protected.Use(
		middlewares.ForcePasswordReset(
			auth.NewResetPasswordUseCase(
				s.deps.Validator, s.deps.CredentialsRepo,
			),
		),
	)

	{
		routes.RegisterServiceRoutes(v1Protected, s.deps)
		routes.RegisterClientRoutes(v1Protected, s.deps)
		routes.RegisterProjectRoutes(v1Protected, s.deps)
		routes.RegisterEnvironmentRoutes(v1Protected, s.deps)
		routes.RegisterAPIKeyRoutes(v1Protected, s.deps)
	}

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: engine,
	}

	log.Printf("[INFO] HTTP API is running on port: %s\n", s.addr)
	log.Printf("[INFO] Pandora Core is fully initialized and ready to accept requests.\n\n")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("[ERROR] Failed to start server: %v\n", err)
		return err
	}

	return nil
}

func NewServer(
	addr string,
	exposeVersion bool,
	deps *bootstrap.Dependencies,
) *Server {
	return &Server{
		addr:          addr,
		deps:          deps,
		exposeVersion: exposeVersion,
	}
}
