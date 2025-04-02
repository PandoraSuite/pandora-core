package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	addr string

	server *grpc.Server
}

func (srv *Server) Run() {
	srv.server = grpc.NewServer()

	listener, err := net.Listen("tcp", srv.addr)
	if err != nil {
		log.Fatalf("[ERROR] Failed to listen: %v", err)
	}

	fmt.Printf("[INFO] gRPC server is running on port: %s\n", srv.addr)
	log.Printf("[INFO] Pandora Core is fully initialized and ready to accept requests.\n\n")
	if err := srv.server.Serve(listener); err != nil {
		log.Fatalf("[ERROR] Failed to serve: %v", err)
	}
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}
