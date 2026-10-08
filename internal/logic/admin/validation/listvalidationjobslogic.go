package validation

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListValidationJobsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListValidationJobsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListValidationJobsLogic {
	return &ListValidationJobsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListValidationJobsLogic) ListValidationJobs(req *types.ListValidationJobsRequest) (resp *types.ListValidationJobsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
