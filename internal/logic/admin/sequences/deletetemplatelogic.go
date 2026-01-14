package sequences

import (
	"context"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTemplateLogic {
	return &DeleteTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTemplateLogic) DeleteTemplate(req *types.DeleteTemplateRequest) (resp *types.AnalyticsResponse, err error) {
	err = l.svcCtx.DB.DeleteTemplate(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.AnalyticsResponse{
		Success: true,
		Message: "Template deleted successfully",
	}, nil
}
