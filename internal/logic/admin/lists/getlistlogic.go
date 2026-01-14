package lists

import (
	"context"
	"fmt"
	"strconv"

	"outlet/internal/db"
	"outlet/internal/svc"
	"outlet/internal/types"
	"outlet/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListLogic {
	return &GetListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListLogic) GetList(req *types.GetListRequest) (resp *types.ListInfo, err error) {
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid list ID: %w", err)
	}

	list, err := l.svcCtx.DB.GetEmailList(l.ctx, id)
	if err != nil {
		l.Errorf("Failed to get list %d: %v", id, err)
		return nil, err
	}

	subscriberCount, err := l.svcCtx.DB.CountListSubscribers(l.ctx, db.CountListSubscribersParams{
		ListID:       id,
		FilterStatus: "active",
	})
	if err != nil {
		l.Errorf("Failed to count subscribers for list %d: %v", id, err)
		return nil, err
	}

	return &types.ListInfo{
		Id:              strconv.FormatInt(list.ID, 10),
		OrgId:           list.OrgID,
		Name:            list.Name,
		Slug:            list.Slug,
		Description:     list.Description.String,
		DoubleOptin:     list.DoubleOptin.Int64 == 1,
		SubscriberCount: int(subscriberCount),
		CreatedAt:       utils.FormatNullString(list.CreatedAt),
		UpdatedAt:       utils.FormatNullString(list.UpdatedAt),
	}, nil
}
