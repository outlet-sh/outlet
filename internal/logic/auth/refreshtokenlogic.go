package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"outlet/internal/svc"
	"outlet/internal/types"
	"outlet/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Refresh access token using refresh token
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest) (resp *types.LoginResponse, err error) {
	if req.RefreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	// Validate the refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, l.svcCtx.Config.Auth.RefreshSecret)
	if err != nil {
		l.Errorf("Invalid refresh token: %v", err)
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Get user from database to ensure they still exist and are active
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, claims.UserID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		l.Errorf("Database error getting user: %v", err)
		return nil, fmt.Errorf("failed to verify user")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, fmt.Errorf("account is not active")
	}

	// Check if account is locked
	if user.LockedUntil.Valid {
		lockedUntil, err := time.Parse(time.RFC3339, user.LockedUntil.String)
		if err == nil && time.Now().Before(lockedUntil) {
			return nil, fmt.Errorf("account is locked until %s", user.LockedUntil.String)
		}
	}

	// Generate new access token
	accessExpire := time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second
	accessToken, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
		l.svcCtx.Config.Auth.AccessSecret,
		accessExpire,
	)
	if err != nil {
		l.Errorf("Error generating access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token")
	}

	// Generate new refresh token (rotate refresh tokens for security)
	refreshExpire := time.Duration(l.svcCtx.Config.Auth.RefreshExpire) * time.Second
	newRefreshToken, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
		l.svcCtx.Config.Auth.RefreshSecret,
		refreshExpire,
	)
	if err != nil {
		l.Errorf("Error generating refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	l.Infof("Token refreshed for user: %s", user.Email)

	expiresAt := time.Now().Add(accessExpire)

	return &types.LoginResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
		User: types.UserInfo{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: utils.FormatNullString(user.CreatedAt),
		},
	}, nil
}
