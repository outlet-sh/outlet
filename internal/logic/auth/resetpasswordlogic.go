package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Reset password with token
func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordRequest) (resp *types.ResetPasswordResponse, err error) {
	// Validate input
	if req.Token == "" {
		return nil, fmt.Errorf("reset token is required")
	}
	if req.NewPassword == "" {
		return nil, fmt.Errorf("new password is required")
	}
	if len(req.NewPassword) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}
	if req.NewPassword != req.ConfirmPassword {
		return nil, fmt.Errorf("passwords do not match")
	}

	// Validate reset token
	authToken, err := l.svcCtx.DB.GetAuthTokenByToken(l.ctx, req.Token)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired reset token")
	}
	if err != nil {
		l.Errorf("Error validating reset token: %v", err)
		return nil, fmt.Errorf("failed to validate reset token")
	}

	// Verify token type
	if authToken.TokenType != "password_reset" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("Error hashing password: %v", err)
		return nil, fmt.Errorf("failed to process password")
	}

	// Update user's password
	err = l.svcCtx.DB.UpdateUserPassword(l.ctx, db.UpdateUserPasswordParams{
		ID:           authToken.UserID,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		l.Errorf("Error updating password: %v", err)
		return nil, fmt.Errorf("failed to update password")
	}

	// Mark the token as used
	err = l.svcCtx.DB.MarkAuthTokenUsed(l.ctx, authToken.ID)
	if err != nil {
		l.Errorf("Error marking token as used: %v", err)
		// Continue - password was updated successfully
	}

	// Reset failed login attempts if any
	err = l.svcCtx.DB.ResetFailedLogins(l.ctx, authToken.UserID)
	if err != nil {
		l.Errorf("Error resetting failed logins: %v", err)
		// Continue - password was updated successfully
	}

	l.Infof("Password reset successfully for user: %s", authToken.UserID)

	return &types.ResetPasswordResponse{
		Message: "Password has been reset successfully. You can now login with your new password.",
	}, nil
}
