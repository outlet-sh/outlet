package validation

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyValidationJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyValidationJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyValidationJobLogic {
	return &ApplyValidationJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyValidationJobLogic) ApplyValidationJob(req *types.ApplyValidationJobRequest) (resp *types.ApplyValidationJobResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
