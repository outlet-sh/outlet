package transactional

import (
	"context"
	"errors"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionalEmailStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransactionalEmailStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionalEmailStatsLogic {
	return &GetTransactionalEmailStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransactionalEmailStatsLogic) GetTransactionalEmailStats(req *types.GetTransactionalEmailRequest) (resp *types.TransactionalStatsResponse, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, errors.New("org_id not found in context")
	}

	_, err = l.svcCtx.DB.GetTransactionalEmail(l.ctx, db.GetTransactionalEmailParams{
		ID:    req.Id,
		OrgID: orgID,
	})
	if err != nil {
		l.Errorf("Failed to get transactional email: %v", err)
		return nil, err
	}

	stats, err := l.svcCtx.DB.GetTransactionalStats(l.ctx, req.Id)
	if err != nil {
		l.Errorf("Failed to get transactional email stats: %v", err)
		return &types.TransactionalStatsResponse{
			Total:     0,
			Sent:      0,
			Delivered: 0,
			Opened:    0,
			Clicked:   0,
			Bounced:   0,
			Failed:    0,
		}, nil
	}

	return &types.TransactionalStatsResponse{
		Total:     int(stats.Total),
		Sent:      int(stats.Sent.Float64),
		Delivered: int(stats.Delivered.Float64),
		Opened:    int(stats.Opened.Float64),
		Clicked:   int(stats.Clicked.Float64),
		Bounced:   int(stats.Bounced.Float64),
		Failed:    int(stats.Failed.Float64),
	}, nil
}
