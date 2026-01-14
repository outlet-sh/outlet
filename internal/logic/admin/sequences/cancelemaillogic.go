package sequences

import (
	"context"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelEmailLogic {
	return &CancelEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelEmailLogic) CancelEmail(req *types.CancelEmailRequest) (resp *types.AnalyticsResponse, err error) {
	err = l.svcCtx.DB.CancelEmail(l.ctx, req.Id)
	if err != nil {
		l.Errorf("Failed to cancel email %s: %v", req.Id, err)
		return nil, err
	}

	return &types.AnalyticsResponse{
		Success: true,
		Message: "Email cancelled successfully",
	}, nil
}
