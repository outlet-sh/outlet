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

type GetTransactionalEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransactionalEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionalEmailLogic {
	return &GetTransactionalEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransactionalEmailLogic) GetTransactionalEmail(req *types.GetTransactionalEmailRequest) (resp *types.TransactionalEmailInfo, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, errors.New("org_id not found in context")
	}

	email, err := l.svcCtx.DB.GetTransactionalEmail(l.ctx, db.GetTransactionalEmailParams{
		ID:    req.Id,
		OrgID: orgID,
	})
	if err != nil {
		l.Errorf("Failed to get transactional email: %v", err)
		return nil, err
	}

	info := transactionalEmailToInfo(email)
	return &info, nil
}
