package request

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/request/v1"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/utils"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

type service struct {
	pb.RequestServiceServer

	requestLogService inbound.RequestLogGRPCPort
}

func (s *service) SetRequestStatus(ctx context.Context, req *pb.SetRequestStatusRequest) (*emptypb.Empty, error) {
	err := s.requestLogService.UpdateExecutionStatus(
		ctx, req.GetId(), enums.RequestLogExecutionStatus(req.GetStatus()),
	)
	if err != nil {
		return nil, status.Error(
			utils.GetDomainErrorStatusCode(err),
			err.Message,
		)
	}
	return &emptypb.Empty{}, nil
}

func RegisterService(server *grpc.Server, requestLogService inbound.RequestLogGRPCPort) {
	service := &service{requestLogService: requestLogService}
	pb.RegisterRequestServiceServer(server, service)
}
