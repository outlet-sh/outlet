package stats

import (
	"context"
	"database/sql"

	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStatsOverviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStatsOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStatsOverviewLogic {
	return &GetStatsOverviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStatsOverviewLogic) GetStatsOverview(req *types.GetStatsOverviewRequest) (resp *types.StatsOverviewResponse, err error) {
	// Get org ID from context (set by API key middleware)
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		l.Errorf("Organization not found in context")
		return nil, nil
	}

	_ = req // Mark as used

	// Get total contact count for this org
	totalContacts, err := l.svcCtx.DB.CountContactsByOrg(l.ctx, sql.NullString{String: orgID, Valid: true})
	if err != nil {
		l.Errorf("Failed to count contacts: %v", err)
		totalContacts = 0
	}

	return &types.StatsOverviewResponse{
		TotalContacts:   int(totalContacts),
		NewContacts:     0,
		ActiveContacts:  0,
		EmailsSent:      0,
		EmailsDelivered: 0,
		EmailsOpened:    0,
		EmailsClicked:   0,
		EmailsBounced:   0,
		OpenRate:        0,
		ClickRate:       0,
		BounceRate:      0,
	}, nil
}
