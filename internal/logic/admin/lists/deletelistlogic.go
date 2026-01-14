package lists

import (
	"context"
	"fmt"
	"strconv"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteListLogic {
	return &DeleteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteListLogic) DeleteList(req *types.DeleteListRequest) (resp *types.AnalyticsResponse, err error) {
	listID, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid list ID: %w", err)
	}

	err = l.svcCtx.DB.DeleteEmailList(l.ctx, listID)
	if err != nil {
		l.Errorf("Failed to delete list: %v", err)
		return nil, fmt.Errorf("failed to delete list: %w", err)
	}

	return &types.AnalyticsResponse{
		Success: true,
		Message: "List deleted successfully",
	}, nil
}
