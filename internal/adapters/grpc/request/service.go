package request

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/request/v1"
)

type service struct {
	pb.RequestServiceServer
}

func (s *service) SetRequestStatus(ctx context.Context, req *pb.SetRequestStatusRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func RegisterService(server *grpc.Server) {
	service := &service{}
	pb.RegisterRequestServiceServer(server, service)
}
