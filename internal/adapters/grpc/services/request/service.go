package request

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/errors"
	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/services/request/v1"
	"github.com/MAD-py/pandora-core/internal/app/request"
)

type service struct {
	pb.RequestServiceServer

	updateStatusUC request.UpdateExecutionStatusUseCase
}

func (s *service) UpdateExecutionStatus(ctx context.Context, req *pb.UpdateExecutionStatusRequest) (*emptypb.Empty, error) {
	err := s.updateStatusUC.Execute(
		ctx, req.GetId(), updateExecutionStatusRequestToDomain(req),
	)
	if err != nil {
		return nil, status.Error(
			errors.CodeToGRPCCode(err.Code()),
			err.Error(),
		)
	}
	return &emptypb.Empty{}, nil
}

func RegisterService(s *grpc.Server, deps *bootstrap.Dependencies) {
	service := &service{
		updateStatusUC: request.NewUpdateExecutionStatusUseCase(
			deps.Validator,
			deps.Repositories.Request(),
		),
	}
	pb.RegisterRequestServiceServer(s, service)
}
