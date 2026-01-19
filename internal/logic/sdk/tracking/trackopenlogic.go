package tracking

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrackOpenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrackOpenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrackOpenLogic {
	return &TrackOpenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrackOpenLogic) TrackOpen(req *types.TrackOpenRequest) (resp *types.Response, err error) {
	if err := l.svcCtx.Tracking.RecordOpen(l.ctx, req.Token); err != nil {
		return &types.Response{Success: false, Message: "invalid token"}, nil
	}

	return &types.Response{Success: true}, nil
}
