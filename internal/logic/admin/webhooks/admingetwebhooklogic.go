package webhooks

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"
	"outlet/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetWebhookLogic {
	return &AdminGetWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetWebhookLogic) AdminGetWebhook(req *types.GetWebhookRequest) (resp *types.WebhookInfo, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("organization not found in context")
	}

	webhookID := req.Id

	w, err := l.svcCtx.DB.GetWebhook(l.ctx, db.GetWebhookParams{
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

	// Decode events from JSON string
	var events []string
	if err := json.Unmarshal([]byte(w.Events), &events); err != nil {
		l.Errorf("Failed to decode events: %v", err)
		events = []string{}
	}

	info := &types.WebhookInfo{
		Id:                w.ID,
		Url:               w.Url,
		Events:            events,
		Active:            w.Active.Int64 == 1,
		CreatedAt:         utils.FormatNullString(w.CreatedAt),
		DeliveriesTotal:   int(w.DeliveriesTotal.Int64),
		DeliveriesSuccess: int(w.DeliveriesSuccess.Int64),
		DeliveriesFailed:  int(w.DeliveriesFailed.Int64),
	}
	if w.LastDeliveryAt.Valid {
		info.LastDeliveryAt = utils.FormatNullString(w.LastDeliveryAt)
	}
	if w.LastStatus.Valid {
		info.LastStatus = int(w.LastStatus.Int64)
	}

	return info, nil
}
