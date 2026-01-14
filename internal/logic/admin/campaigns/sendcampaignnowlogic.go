package campaigns

import (
	"context"
	"database/sql"
	"errors"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCampaignNowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendCampaignNowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCampaignNowLogic {
	return &SendCampaignNowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCampaignNowLogic) SendCampaignNow(req *types.SendCampaignNowRequest) (resp *types.CampaignInfo, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, errors.New("org_id not found in context")
	}

	campaign, err := l.svcCtx.DB.UpdateCampaignStatus(l.ctx, db.UpdateCampaignStatusParams{
		ID:     req.Id,
		OrgID:  orgID,
		Status: sql.NullString{String: "sending", Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to start campaign: %v", err)
		return nil, err
	}

	info := campaignToInfo(campaign)
	return &info, nil
}
