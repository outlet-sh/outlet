package tracking

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrackUnsubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrackUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrackUnsubscribeLogic {
	return &TrackUnsubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrackUnsubscribeLogic) TrackUnsubscribe(req *types.TrackUnsubscribeRequest) (resp *types.Response, err error) {
	if err := l.svcCtx.Tracking.Unsubscribe(l.ctx, req.Token); err != nil {
		return &types.Response{Success: false, Message: "invalid token"}, nil
	}

	return &types.Response{Success: true}, nil
}
