package validation

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelValidationJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelValidationJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelValidationJobLogic {
	return &CancelValidationJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelValidationJobLogic) CancelValidationJob(req *types.ValidationJobIdRequest) (resp *types.ValidationJobInfo, err error) {
	// todo: add your logic here and delete this line

	return
}
