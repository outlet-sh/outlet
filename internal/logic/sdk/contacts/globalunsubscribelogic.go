package contacts

import (
	"context"
	"database/sql"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GlobalUnsubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGlobalUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GlobalUnsubscribeLogic {
	return &GlobalUnsubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GlobalUnsubscribeLogic) GlobalUnsubscribe(req *types.GlobalUnsubscribeRequest) (resp *types.Response, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return &types.Response{Success: false, Message: "unauthorized"}, nil
	}

	if req.Email == "" {
		return &types.Response{Success: false, Message: "email required"}, nil
	}

	err = l.svcCtx.DB.GlobalUnsubscribeByOrgAndEmail(l.ctx, db.GlobalUnsubscribeByOrgAndEmailParams{
		OrgID: sql.NullString{String: orgID, Valid: true},
		Email: req.Email,
	})
	if err != nil {
		l.Errorf("Failed to global unsubscribe: %v", err)
		return &types.Response{Success: false, Message: "unsubscribe failed"}, nil
	}

	l.Infof("Global unsubscribe: org=%s email=%s reason=%s", orgID, req.Email, req.Reason)

	return &types.Response{Success: true, Message: "unsubscribed"}, nil
}
