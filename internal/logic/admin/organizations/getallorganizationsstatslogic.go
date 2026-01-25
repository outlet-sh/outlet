package organizations

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllOrganizationsStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllOrganizationsStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllOrganizationsStatsLogic {
	return &GetAllOrganizationsStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllOrganizationsStatsLogic) GetAllOrganizationsStats() (resp *types.OrgBulkStatsResponse, err error) {
	resp = &types.OrgBulkStatsResponse{
		Stats: make(map[string]types.OrgStats),
	}

	// Get contact counts for all orgs
	contactStats, err := l.svcCtx.DB.GetAllOrgsContactStats(l.ctx)
	if err != nil {
		l.Logger.Errorf("Failed to get contact stats: %v", err)
	} else {
		for _, stat := range contactStats {
			if !stat.OrgID.Valid {
				continue
			}
			orgID := stat.OrgID.String
			orgStats := resp.Stats[orgID]
			orgStats.TotalContacts = stat.TotalContacts
			resp.Stats[orgID] = orgStats
		}
	}

	// Get email stats for last 30 days for all orgs
	emailStats, err := l.svcCtx.DB.GetAllOrgsEmailStats30Days(l.ctx)
	if err != nil {
		l.Logger.Errorf("Failed to get email stats: %v", err)
	} else {
		for _, stat := range emailStats {
			orgStats := resp.Stats[stat.OrgID]
			orgStats.EmailsSent30d = toInt64(stat.EmailsSent)
			resp.Stats[stat.OrgID] = orgStats
		}
	}

	// Get list counts for all orgs
	listStats, err := l.svcCtx.DB.GetAllOrgsListCounts(l.ctx)
	if err != nil {
		l.Logger.Errorf("Failed to get list stats: %v", err)
	} else {
		for _, stat := range listStats {
			orgStats := resp.Stats[stat.OrgID]
			orgStats.ListCount = stat.ListCount
			resp.Stats[stat.OrgID] = orgStats
		}
	}

	return resp, nil
}
