package webhooks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWebhookLogic {
	return &UpdateWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWebhookLogic) UpdateWebhook(req *types.UpdateWebhookRequest) (resp *types.WebhookInfo, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("organization not found in context")
	}

	if req.Id == "" {
		return nil, fmt.Errorf("invalid webhook ID")
	}

	events := req.Events
	if events == nil {
		events = []string{}
	}

	active := int64(0)
	if req.Active {
		active = 1
	}

	webhook, err := l.svcCtx.DB.UpdateWebhook(l.ctx, db.UpdateWebhookParams{
		ID:     req.Id,
		OrgID:  orgID,
		Url:    req.Url,
		Events: strings.Join(events, ","),
		Active: sql.NullInt64{Int64: active, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("webhook not found")
		}
		l.Errorf("Failed to update webhook: %v", err)
		return nil, fmt.Errorf("failed to update webhook")
	}

	info := &types.WebhookInfo{
		Id:                webhook.ID,
		Url:               webhook.Url,
		Events:            strings.Split(webhook.Events, ","),
		Active:            webhook.Active.Valid && webhook.Active.Int64 == 1,
		CreatedAt:         webhook.CreatedAt.String,
		DeliveriesTotal:   int(webhook.DeliveriesTotal.Int64),
		DeliveriesSuccess: int(webhook.DeliveriesSuccess.Int64),
		DeliveriesFailed:  int(webhook.DeliveriesFailed.Int64),
	}
	if webhook.LastDeliveryAt.Valid {
		info.LastDeliveryAt = webhook.LastDeliveryAt.String
	}
	if webhook.LastStatus.Valid {
		info.LastStatus = int(webhook.LastStatus.Int64)
	}

	return info, nil
}
