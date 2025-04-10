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

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/api_key"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/request"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/reservation"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

type Server struct {
	addr string

	server *grpc.Server

	apiKeyService      inbound.APIKeyGRPCPort
	reservationService inbound.ReservationGRPCPort
}

func (s *Server) setupServices() {
	api_key.RegisterService(s.server, s.apiKeyService)
	request.RegisterService(s.server)
	reservation.RegisterService(s.server, s.reservationService)
}

func (s *Server) Run() {
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
		log.Fatalf("[ERROR] Failed to listen: %v", err)
	}

	log.Printf("[INFO] gRPC server is running on port: %s\n", s.addr)
	log.Printf("[INFO] Pandora Core is fully initialized and ready to accept requests.\n\n")
	if err := s.server.Serve(listener); err != nil {
		log.Fatalf("[ERROR] Failed to serve: %v", err)
	}
}

func NewServer(
	addr string,
	apiKeyService inbound.APIKeyGRPCPort,
	reservationService inbound.ReservationGRPCPort,
) *Server {
	return &Server{
		addr:               addr,
		apiKeyService:      apiKeyService,
		reservationService: reservationService,
	}
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
