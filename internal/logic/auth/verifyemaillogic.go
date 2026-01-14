package auth

import (
	"context"
	"database/sql"
	"fmt"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Verify email address
func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyEmailLogic) VerifyEmail(req *types.VerifyEmailRequest) (resp *types.VerifyEmailResponse, err error) {
	// Validate input
	if req.Token == "" {
		return nil, fmt.Errorf("verification token is required")
	}

	// Validate verification token
	authToken, err := l.svcCtx.DB.GetAuthTokenByToken(l.ctx, req.Token)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired verification token")
	}
	if err != nil {
		l.Errorf("Error validating verification token: %v", err)
		return nil, fmt.Errorf("failed to validate verification token")
	}

	// Verify token type
	if authToken.TokenType != "email_verification" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Update user's email_verified status and activate account
	err = l.svcCtx.DB.SetUserEmailVerified(l.ctx, authToken.UserID)
	if err != nil {
		l.Errorf("Error verifying email: %v", err)
		return nil, fmt.Errorf("failed to verify email")
	}

	// Mark the token as used
	err = l.svcCtx.DB.MarkAuthTokenUsed(l.ctx, authToken.ID)
	if err != nil {
		l.Errorf("Error marking token as used: %v", err)
		// Continue - email was verified successfully
	}

	l.Infof("Email verified successfully for user: %s", authToken.UserID)

	return &types.VerifyEmailResponse{
		Message: "Email verified successfully. Your account is now active.",
	}, nil
}
