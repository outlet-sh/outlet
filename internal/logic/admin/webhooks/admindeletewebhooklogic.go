package webhooks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteWebhookLogic {
	return &AdminDeleteWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteWebhookLogic) AdminDeleteWebhook(req *types.DeleteWebhookRequest) (resp *types.Response, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("organization not found in context")
	}

	webhookID := req.Id

	// Get existing webhook to verify ownership (scoped to org)
	_, err = l.svcCtx.DB.GetWebhook(l.ctx, db.GetWebhookParams{
		ID:    webhookID,
		OrgID: orgID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("webhook not found")
	}
	if err != nil {
		l.Errorf("Failed to get webhook: %v", err)
		return nil, fmt.Errorf("failed to get webhook")
	}

	err = l.svcCtx.DB.DeleteWebhook(l.ctx, db.DeleteWebhookParams{
		ID:    webhookID,
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
