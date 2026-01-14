package campaigns

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScheduleCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScheduleCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleCampaignLogic {
	return &ScheduleCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScheduleCampaignLogic) ScheduleCampaign(req *types.ScheduleCampaignRequest) (resp *types.CampaignInfo, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, errors.New("org_id not found in context")
	}

	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		return nil, errors.New("invalid scheduled_at format, use ISO8601")
	}

	campaign, err := l.svcCtx.DB.ScheduleCampaign(l.ctx, db.ScheduleCampaignParams{
		ID:          req.Id,
		OrgID:       orgID,
		ScheduledAt: sql.NullString{String: scheduledAt.Format(time.RFC3339), Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to schedule campaign: %v", err)
		return nil, err
	}

	info := campaignToInfo(campaign)
	return &info, nil
}
