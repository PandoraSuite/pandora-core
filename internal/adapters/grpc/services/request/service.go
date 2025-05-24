package request

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/services/request/v1"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/utils"
	"github.com/MAD-py/pandora-core/internal/app/request"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type service struct {
	pb.RequestServiceServer

	updateStatusUC request.UpdateStatusUseCase
}

func (s *service) SetRequestStatus(ctx context.Context, req *pb.SetRequestStatusRequest) (*emptypb.Empty, error) {
	err := s.updateStatusUC.Execute(
		ctx, req.GetId(), enums.RequestExecutionStatus(req.GetStatus()),
	)
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
		updateStatusUC: request.NewUpdateStatusUseCase(
			deps.Validator,
			deps.Repositories.Request(),
		),
	}
	pb.RegisterRequestServiceServer(s, service)
}
