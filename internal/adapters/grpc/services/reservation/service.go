package reservation

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/services/reservation/v1"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/utils"
	"github.com/MAD-py/pandora-core/internal/app/reservation"
)

type service struct {
	pb.ReservationServiceServer

	commitUC   reservation.CommitUseCase
	rollbackUC reservation.RollbackUseCase
}

func (s *service) Commit(ctx context.Context, req *pb.CommitRequest) (*emptypb.Empty, error) {
	err := s.commitUC.Execute(ctx, req.GetParams().Id)
	if err != nil {
		return nil, status.Error(
			utils.GetDomainErrorStatusCode(err),
			err.Error(),
		)
	}
	return &emptypb.Empty{}, nil

}

func (s *service) Rollback(ctx context.Context, req *pb.RollbackRequest) (*emptypb.Empty, error) {
	err := s.rollbackUC.Execute(ctx, req.GetParams().Id)
	if err != nil {
		return nil, status.Error(
			utils.GetDomainErrorStatusCode(err),
			err.Error(),
		)
	}
	return &emptypb.Empty{}, nil
}

func RegisterService(s *grpc.Server, deps *bootstrap.Dependencies) {
	service := &service{
		commitUC: reservation.NewCommitUseCase(
			deps.Validator,
			deps.Repositories.Reservation(),
		),
		rollbackUC: reservation.NewRollbackUseCase(
			deps.Validator,
			deps.Repositories.Reservation(),
			deps.Repositories.Environment(),
		),
	}
	pb.RegisterReservationServiceServer(s, service)
}
