package stats

import (
	"context"

	"github.com/outlet-sh/outlet/internal/middleware"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmailStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmailStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailStatsLogic {
	return &GetEmailStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmailStatsLogic) GetEmailStats(req *types.GetEmailStatsRequest) (resp *types.GetEmailStatsResponse, err error) {
	// Get org ID from context (set by API key middleware)
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		l.Errorf("Organization not found in context")
		return nil, nil
	}

	_ = orgID // Mark as used

	// Parse date range
	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	// Build stats points for the date range
	stats := make([]types.EmailStatsPoint, 0)
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		stats = append(stats, types.EmailStatsPoint{
			Date:      dateStr,
			Sent:      0,
			Delivered: 0,
			Opened:    0,
			Clicked:   0,
			Bounced:   0,
			OpenRate:  0,
			ClickRate: 0,
		})
	}

	return &types.GetEmailStatsResponse{
		Stats:          stats,
		TotalSent:      0,
		TotalDelivered: 0,
		TotalOpened:    0,
		TotalClicked:   0,
		TotalBounced:   0,
		AvgOpenRate:    0,
		AvgClickRate:   0,
	}, nil
}
