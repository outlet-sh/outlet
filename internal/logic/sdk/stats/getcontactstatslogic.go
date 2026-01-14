package stats

import (
	"context"
	"database/sql"
	"time"

	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetContactStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContactStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContactStatsLogic {
	return &GetContactStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetContactStatsLogic) GetContactStats(req *types.GetContactStatsRequest) (resp *types.GetContactStatsResponse, err error) {
	// Get org ID from context (set by API key middleware)
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		l.Errorf("Organization not found in context")
		return nil, nil
	}

	// Parse date range
	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	// Get total contact count for this org
	totalContacts, err := l.svcCtx.DB.CountContactsByOrg(l.ctx, sql.NullString{String: orgID, Valid: true})
	if err != nil {
		l.Errorf("Failed to count contacts: %v", err)
		totalContacts = 0
	}

	// Build stats points for the date range
	stats := make([]types.ContactStatsPoint, 0)
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		stats = append(stats, types.ContactStatsPoint{
			Date:         dateStr,
			NewContacts:  0,
			Unsubscribed: 0,
			Bounced:      0,
			NetGrowth:    0,
		})
	}

	return &types.GetContactStatsResponse{
		Stats:             stats,
		TotalActive:       int(totalContacts),
		TotalUnsubscribed: 0,
		TotalBounced:      0,
		NetGrowth:         0,
	}, nil
}

// parseDateRange parses start and end date strings and returns time.Time values
// Defaults to last 30 days if not specified
func parseDateRange(startDateStr, endDateStr string) (time.Time, time.Time) {
	now := time.Now().UTC()
	endDate := now

	// Default to last 30 days
	startDate := now.AddDate(0, 0, -30)

	if startDateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = parsed
		} else if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = parsed
		} else if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Add a day to include the entire end date
			endDate = parsed.AddDate(0, 0, 1)
		}
	}

	// Normalize to start of day
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)

	return startDate, endDate
}
