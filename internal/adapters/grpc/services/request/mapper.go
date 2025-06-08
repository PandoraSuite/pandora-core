package request

import (
	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/services/request/v1"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

func updateExecutionStatusRequestToDomain(
	req *pb.UpdateExecutionStatusRequest,
) *dto.RequestExecutionStatusUpdate {
	return &dto.RequestExecutionStatusUpdate{
		Detail:          req.Detail,
		StatusCode:      int(req.StatusCode),
		ExecutionStatus: enums.RequestExecutionStatus(req.ExecutionStatus),
	}
}
