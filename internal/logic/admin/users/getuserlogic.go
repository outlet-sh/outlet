package users

import (
	"context"

	"outlet/internal/svc"
	"outlet/internal/types"
	"outlet/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.UserInfo, err error) {
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, req.Id)
	if err != nil {
		l.Errorf("Failed to get user: %v", err)
		return nil, err
	}

	return &types.UserInfo{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Active:    user.Status == "active",
		CreatedAt: utils.FormatNullString(user.CreatedAt),
	}, nil
}
