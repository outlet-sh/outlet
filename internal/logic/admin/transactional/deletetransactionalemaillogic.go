package transactional

import (
	"context"
	"errors"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/middleware"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTransactionalEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTransactionalEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTransactionalEmailLogic {
	return &DeleteTransactionalEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTransactionalEmailLogic) DeleteTransactionalEmail(req *types.DeleteTransactionalEmailRequest) (resp *types.AnalyticsResponse, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, errors.New("org_id not found in context")
	}

	err = l.svcCtx.DB.DeleteTransactionalEmail(l.ctx, db.DeleteTransactionalEmailParams{
		ID:    req.Id,
		OrgID: orgID,
	})
	if err != nil {
		l.Errorf("Failed to delete transactional email: %v", err)
		return nil, err
	}

	return &types.AnalyticsResponse{
		Success: true,
		Message: "Transactional email deleted successfully",
	}, nil
}
