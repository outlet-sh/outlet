package lists

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"outlet/internal/db"
	"outlet/internal/svc"
	"outlet/internal/types"
	"outlet/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListLogic {
	return &UpdateListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateListLogic) UpdateList(req *types.UpdateListRequest) (resp *types.ListInfo, err error) {
	listID, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid list ID: %w", err)
	}

	params := db.UpdateEmailListParams{
		ID:   listID,
		Name: req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
		ConfirmationSubject: req.ConfirmationSubject,
		ConfirmationBody:    req.ConfirmationBody,
	}

	if req.DoubleOptin != nil {
		params.DoubleOptin = sql.NullInt64{
			Int64: boolToInt64(*req.DoubleOptin),
			Valid: true,
		}
	}

	list, err := l.svcCtx.DB.UpdateEmailList(l.ctx, params)
	if err != nil {
		l.Errorf("Failed to update list: %v", err)
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	subscriberCount, _ := l.svcCtx.DB.CountListSubscribers(l.ctx, db.CountListSubscribersParams{
		ListID:       list.ID,
		FilterStatus: "active",
	})

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
