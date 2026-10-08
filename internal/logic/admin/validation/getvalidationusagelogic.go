package validation

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetValidationUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetValidationUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetValidationUsageLogic {
	return &GetValidationUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetValidationUsageLogic) GetValidationUsage(req *types.ValidationUsageRequest) (resp *types.ValidationUsageResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
