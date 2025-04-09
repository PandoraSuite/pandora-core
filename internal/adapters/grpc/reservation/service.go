package reservation

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/reservation/v1"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

type service struct {
	pb.ReservationServiceServer

	reservationService inbound.ReservationGRPCPort
}

func (s *service) Commit(ctx context.Context, req *pb.CommitRequest) (*emptypb.Empty, error) {
	params := req.GetParams()
	err := s.reservationService.Commit(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil

}

func (s *service) Rollback(ctx context.Context, req *pb.RollbackRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func RegisterService(server *grpc.Server, reservationService inbound.ReservationGRPCPort) {
	service := &service{reservationService: reservationService}
	pb.RegisterReservationServiceServer(server, service)
}
