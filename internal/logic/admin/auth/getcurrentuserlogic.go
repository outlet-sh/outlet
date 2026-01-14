package auth

import (
	"context"
	"database/sql"
	"fmt"

	"outlet/internal/svc"
	"outlet/internal/types"
	"outlet/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserLogic) GetCurrentUser() (resp *types.UserInfo, err error) {
	userIDValue := l.ctx.Value("userId")
	if userIDValue == nil {
		return nil, fmt.Errorf("user ID not found in context")
	}

	userID, ok := userIDValue.(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in context")
	}

	user, err := l.svcCtx.DB.GetUserByID(l.ctx, userID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		l.Errorf("Database error getting user: %v", err)
		return nil, fmt.Errorf("database error")
	}

	return &types.UserInfo{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: utils.FormatNullString(user.CreatedAt),
	}, nil
}
