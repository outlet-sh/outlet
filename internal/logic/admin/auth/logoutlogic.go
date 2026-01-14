package auth

import (
	"context"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() (resp *types.AnalyticsResponse, err error) {
	// Since we're using stateless JWT tokens, logout is handled client-side
	// by deleting the token. The server doesn't need to do anything.
	// This endpoint exists for consistency and potential future enhancements
	// (like token blacklisting or session tracking).

	l.Info("User logged out")

	return &types.AnalyticsResponse{
		Success: true,
		Message: "Logged out successfully",
	}, nil
}
