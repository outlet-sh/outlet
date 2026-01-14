package sequences

import (
	"context"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEntryRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteEntryRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEntryRuleLogic {
	return &DeleteEntryRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEntryRuleLogic) DeleteEntryRule(req *types.DeleteEntryRuleRequest) (resp *types.AnalyticsResponse, err error) {
	err = l.svcCtx.DB.DeleteEntryRule(l.ctx, req.Id)
	if err != nil {
		l.Errorf("Failed to delete entry rule: %v", err)
		return nil, err
	}

	return &types.AnalyticsResponse{
		Success: true,
		Message: "Entry rule deleted successfully",
	}, nil
}
