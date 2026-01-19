package webhooks

import (
	"context"
	"fmt"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/middleware"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWebhookLogic {
	return &DeleteWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWebhookLogic) DeleteWebhook(req *types.DeleteWebhookRequest) (resp *types.Response, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("organization not found in context")
	}

	if req.Id == "" {
		return nil, fmt.Errorf("invalid webhook ID")
	}

	err = l.svcCtx.DB.DeleteWebhook(l.ctx, db.DeleteWebhookParams{
		ID:    req.Id,
		OrgID: orgID,
	})
	if err != nil {
		l.Errorf("Failed to delete webhook: %v", err)
		return nil, fmt.Errorf("failed to delete webhook")
	}

	return &types.Response{
		Success: true,
		Message: "Webhook deleted successfully",
	}, nil
}
