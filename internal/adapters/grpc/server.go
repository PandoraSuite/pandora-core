package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	protovalidator "github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	apikey "github.com/MAD-py/pandora-core/internal/adapters/grpc/services/api_key"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/services/request"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/services/reservation"
)

type Server struct {
	addr string

	server *grpc.Server

	deps *bootstrap.Dependencies
}

func (s *Server) setupServices() {
	// Register standard gRPC health service for compatibility with grpc_health_probe
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s.server, healthServer)

	// Register our application services
	apikey.RegisterService(s.server, s.deps)
	request.RegisterService(s.server, s.deps)
	reservation.RegisterService(s.server, s.deps)
}

func (s *Server) Run() error {
	validator, err := protovalidator.New()
	if err != nil {
		panic("failed to create protovalidate validator")
	}

	logger := interceptorLogger()

	s.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(
				logger,
				logging.WithLogOnEvents(
					logging.StartCall,
					logging.FinishCall,
				),
			),
			protovalidate.UnaryServerInterceptor(validator),
		),
	)
	s.setupServices()

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("[ERROR] Failed to listen: %v\n", err)
		return err
	}

	log.Printf("[INFO] gRPC server is running on port: %s\n", s.addr)
	log.Printf("[INFO] Pandora Core is fully initialized and ready to accept requests.\n\n")
	if err := s.server.Serve(listener); err != nil {
		log.Printf("[ERROR] Failed to serve: %v\n", err)
		return err
	}
	return nil
}

func NewServer(addr string, deps *bootstrap.Dependencies) *Server {
	return &Server{addr: addr, deps: deps}
}

func interceptorLogger() logging.Logger {
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			msg = fmt.Sprintf("DEBUG :%v", msg)
		case logging.LevelInfo:
			msg = fmt.Sprintf("INFO :%v", msg)
		case logging.LevelWarn:
			msg = fmt.Sprintf("WARN :%v", msg)
		case logging.LevelError:
			msg = fmt.Sprintf("ERROR :%v", msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
		logger.Println(append([]any{"msg", msg}, fields...))
	})
}
