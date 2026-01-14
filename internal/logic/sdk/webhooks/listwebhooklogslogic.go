package webhooks

import (
	"context"
	"database/sql"
	"fmt"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWebhookLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWebhookLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWebhookLogsLogic {
	return &ListWebhookLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWebhookLogsLogic) ListWebhookLogs(req *types.ListWebhookLogsRequest) (resp *types.ListWebhookLogsResponse, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("organization not found in context")
	}

	if req.Id == "" {
		return nil, fmt.Errorf("invalid webhook ID")
	}

	_, err = l.svcCtx.DB.GetWebhook(l.ctx, db.GetWebhookParams{
		ID:    req.Id,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("webhook not found")
		}
		l.Errorf("Failed to get webhook: %v", err)
		return nil, fmt.Errorf("failed to get webhook")
	}

	limit := int64(50)
	if req.Limit > 0 && req.Limit <= 100 {
		limit = int64(req.Limit)
	}

	logs, err := l.svcCtx.DB.ListWebhookLogs(l.ctx, db.ListWebhookLogsParams{
		WebhookID:  req.Id,
		LimitCount: limit,
	})
	if err != nil {
		l.Errorf("Failed to list webhook logs: %v", err)
		return nil, fmt.Errorf("failed to list webhook logs")
	}

	result := make([]types.WebhookLogInfo, 0, len(logs))
	for _, log := range logs {
		info := types.WebhookLogInfo{
			Id:          log.ID,
			Event:       log.Event,
			Payload:     log.Payload,
			StatusCode:  int(log.StatusCode.Int64),
			Duration:    int(log.DurationMs.Int64),
			DeliveredAt: log.DeliveredAt.String,
		}
		if log.Response.Valid {
			info.Response = log.Response.String
		}
		if log.Error.Valid {
			info.Error = log.Error.String
		}
		result = append(result, info)
	}

	return &types.ListWebhookLogsResponse{
		Logs: result,
	}, nil
}
