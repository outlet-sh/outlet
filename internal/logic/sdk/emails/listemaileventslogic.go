package emails

import (
	"context"
	"database/sql"
	"sort"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEmailEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListEmailEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEmailEventsLogic {
	return &ListEmailEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListEmailEventsLogic) ListEmailEvents(req *types.ListEmailEventsRequest) (resp *types.ListEmailEventsResponse, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, nil
	}

	if req.MessageId == "" {
		return nil, nil
	}

	send, err := l.svcCtx.DB.GetTransactionalSendByTrackingAndOrg(l.ctx, db.GetTransactionalSendByTrackingAndOrgParams{
		TrackingToken: sql.NullString{String: req.MessageId, Valid: true},
		OrgID:         orgID,
	})
	if err != nil {
		l.Errorf("Failed to get email for events: %v", err)
		return nil, err
	}

	events := make([]types.EmailEventInfo, 0)

	if send.CreatedAt.Valid {
		events = append(events, types.EmailEventInfo{
			Event:     "queued",
			Timestamp: send.CreatedAt.String,
		})
	}

	if send.SentAt.Valid {
		events = append(events, types.EmailEventInfo{
			Event:     "sent",
			Timestamp: send.SentAt.String,
		})
	}

	if send.DeliveredAt.Valid {
		events = append(events, types.EmailEventInfo{
			Event:     "delivered",
			Timestamp: send.DeliveredAt.String,
		})
	}

	if send.OpenedAt.Valid {
		events = append(events, types.EmailEventInfo{
			Event:     "opened",
			Timestamp: send.OpenedAt.String,
		})
	}

	if send.ClickedAt.Valid {
		events = append(events, types.EmailEventInfo{
			Event:     "clicked",
			Timestamp: send.ClickedAt.String,
		})
	}

	if send.Status.Valid && send.Status.String == "bounced" {
		data := ""
		if send.ErrorMessage.Valid {
			data = send.ErrorMessage.String
		}
		events = append(events, types.EmailEventInfo{
			Event:     "bounced",
			Timestamp: send.CreatedAt.String,
			Data:      data,
		})
	}

	if send.Status.Valid && send.Status.String == "failed" {
		data := ""
		if send.ErrorMessage.Valid {
			data = send.ErrorMessage.String
		}
		events = append(events, types.EmailEventInfo{
			Event:     "failed",
			Timestamp: send.CreatedAt.String,
			Data:      data,
		})
	}

	// Sort events by timestamp (oldest first)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp < events[j].Timestamp
	})

	l.Infof("ListEmailEvents: org=%s messageId=%s events=%d", orgID, req.MessageId, len(events))

	return &types.ListEmailEventsResponse{
		Events: events,
	}, nil
}
