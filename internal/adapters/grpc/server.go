package grpc

import (
	"log"
	"net"

	"github.com/bufbuild/protovalidate-go"
	interceptors "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/api_key"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/request"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/reservation"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

type Server struct {
	addr string

	server *grpc.Server

	apiKeyService inbound.APIKeyGRPCPort
}

func (s *Server) setupServices() {
	api_key.RegisterService(s.server, s.apiKeyService)
	request.RegisterService(s.server)
	reservation.RegisterService(s.server)
}

func (s *Server) Run() {
	validator, err := protovalidate.New()
	if err != nil {
		panic("failed to create protovalidate validator")
	}

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptors.UnaryServerInterceptor(validator),
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
) *Server {
	return &Server{
		addr:          addr,
		apiKeyService: apiKeyService,
	}
}
