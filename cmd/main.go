package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc"
	grpcBootstrap "github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http"
	httpBootstrap "github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/security"
	"github.com/MAD-py/pandora-core/internal/adapters/taskengine"
	taskengineBootstrap "github.com/MAD-py/pandora-core/internal/adapters/taskengine/bootstrap"
	"github.com/MAD-py/pandora-core/internal/config"
	"github.com/MAD-py/pandora-core/internal/validator"
)

func main() {
	time.Local = time.UTC

	log.Println("[INFO] Starting Pandora Core (API RESTful + gRPC + TaskEngine)...")

	cfg := config.LoadConfig()

	log.Printf("[INFO] HTTP, gRPC and TaskEngine config loaded")

	validator := validator.NewValidator()
	log.Println("[INFO] Validator initialized")

	repositories := persistence.NewRepositories(
		persistence.PostgresDriver, cfg.DBDNS(),
	)
	log.Println("[INFO] Repositories initialized")

	jwtProvider := security.NewJWTProvider([]byte(cfg.HTTPConfig().JWTSecret()))
	log.Println("[INFO] JWT provider initialized")

	credentialsRepo := security.NewCredentialsRepository(cfg.HTTPConfig().CredentialsFile())
	log.Println("[INFO] Credentials repository initialized")

	gRPCDeps := grpcBootstrap.NewDependencies(validator, repositories)

	grpcSrv := grpc.NewServer(
		fmt.Sprintf(":%s", cfg.GRPCConfig().Port()),
		gRPCDeps,
	)

	httpDeps := httpBootstrap.NewDependencies(
		validator,
		repositories,
		jwtProvider,
		credentialsRepo,
	)

	httpSrv := http.NewServer(
		fmt.Sprintf(":%s", cfg.HTTPConfig().Port()),
		cfg.HTTPConfig().ExposeVersion(),
		httpDeps,
	)

	taskEngineDeps := taskengineBootstrap.NewDependencies(repositories)

	taskEngine, err := taskengine.NewEngine(
		cfg.TaskEngineConfig().DBDNS(), taskEngineDeps,
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create TaskEngine: %v", err)
	}

	var g errgroup.Group

	g.Go(grpcSrv.Run)
	g.Go(httpSrv.Run)
	g.Go(taskEngine.Run)

	if err := g.Wait(); err != nil {
		log.Fatalf("[FATAL] One of the services failed: %v", err)
		os.Exit(1)
	}
}
