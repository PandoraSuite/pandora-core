package reservation

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/reservation/v1"
)

type service struct {
	pb.ReservationServiceServer
}

func (s *service) Commit(ctx context.Context, req *pb.CommitRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil

}

func (s *service) Rollback(ctx context.Context, req *pb.RollbackRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func RegisterService(server *grpc.Server) {
	service := &service{}
	pb.RegisterReservationServiceServer(server, service)
}
