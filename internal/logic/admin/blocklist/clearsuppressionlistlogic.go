package blocklist

import (
	"context"
	"fmt"

	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearSuppressionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearSuppressionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearSuppressionListLogic {
	return &ClearSuppressionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearSuppressionListLogic) ClearSuppressionList() (resp *types.Response, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("organization ID not found")
	}

	err = l.svcCtx.DB.ClearSuppressionList(l.ctx, orgID)
	if err != nil {
		l.Errorf("Failed to clear suppression list: %v", err)
		return nil, err
	}

	return &types.Response{
		Success: true,
		Message: "Suppression list cleared",
	}, nil
}
