package campaigns

import (
	"context"
	"errors"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCampaignStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCampaignStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCampaignStatsLogic {
	return &GetCampaignStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCampaignStatsLogic) GetCampaignStats(req *types.GetCampaignRequest) (resp *types.CampaignStatsResponse, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, errors.New("org_id not found in context")
	}

	campaign, err := l.svcCtx.DB.GetCampaign(l.ctx, db.GetCampaignParams{
		ID:    req.Id,
		OrgID: orgID,
	})
	if err != nil {
		l.Errorf("Failed to get campaign: %v", err)
		return nil, err
	}

	links, err := l.svcCtx.DB.GetCampaignLinkStats(l.ctx, req.Id)
	if err != nil {
		l.Errorf("Failed to get campaign link stats: %v", err)
	}

	linkStats := make([]types.CampaignLinkStat, 0)
	for _, link := range links {
		linkStats = append(linkStats, types.CampaignLinkStat{
			Url:        link.LinkUrl,
			Name:       link.LinkName.String,
			ClickCount: int(link.ClickCount),
		})
	}

	info := campaignToInfo(campaign)
	return &types.CampaignStatsResponse{
		Campaign: info,
		Links:    linkStats,
	}, nil
}
