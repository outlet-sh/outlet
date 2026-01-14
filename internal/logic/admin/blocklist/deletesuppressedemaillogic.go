package blocklist

import (
	"context"
	"fmt"
	"strconv"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSuppressedEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSuppressedEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSuppressedEmailLogic {
	return &DeleteSuppressedEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSuppressedEmailLogic) DeleteSuppressedEmail(req *types.DeleteSuppressedEmailRequest) (resp *types.Response, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("organization ID not found")
	}

	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid suppression ID: %v", err)
	}

	err = l.svcCtx.DB.DeleteSuppressionByID(l.ctx, db.DeleteSuppressionByIDParams{
		ID:    id,
		OrgID: orgID,
	})
	if err != nil {
		l.Errorf("Failed to delete suppressed email: %v", err)
		return nil, err
	}

	return &types.Response{
		Success: true,
		Message: "Email removed from suppression list",
	}, nil
}
