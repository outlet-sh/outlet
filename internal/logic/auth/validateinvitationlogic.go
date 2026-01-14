package auth

import (
	"context"
	"time"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateInvitationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Validate invitation token and get user details
func NewValidateInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateInvitationLogic {
	return &ValidateInvitationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidateInvitationLogic) ValidateInvitation(req *types.ValidateInvitationRequest) (resp *types.ValidateInvitationResponse, err error) {
	// Find the invitation token
	var email, role string
	var expiresAt, usedAt *time.Time
	err = l.svcCtx.DB.GetDB().QueryRowContext(l.ctx, `
		SELECT email, role, expires_at, used_at
		FROM invitation_tokens
		WHERE token = $1
	`, req.Token).Scan(&email, &role, &expiresAt, &usedAt)
	if err != nil {
		return &types.ValidateInvitationResponse{
			Valid:   false,
			Message: "Invalid invitation token",
		}, nil
	}

	// Check if token has been used
	if usedAt != nil {
		return &types.ValidateInvitationResponse{
			Valid:   false,
			Message: "This invitation has already been used",
		}, nil
	}

	// Check if token has expired
	if expiresAt != nil && time.Now().After(*expiresAt) {
		return &types.ValidateInvitationResponse{
			Valid:   false,
			Message: "This invitation has expired",
		}, nil
	}

	// Get user details from the email
	var firstName, lastName, phone string
	err = l.svcCtx.DB.GetDB().QueryRowContext(l.ctx, `
		SELECT first_name, last_name, COALESCE(phone, '') as phone
		FROM users
		WHERE email = $1 AND role = $2
	`, email, role).Scan(&firstName, &lastName, &phone)
	if err != nil {
		// If user doesn't exist yet, return just the email and role
		// (for invitations that don't pre-create the user)
		return &types.ValidateInvitationResponse{
			Email:     email,
			FirstName: "",
			LastName:  "",
			Phone:     "",
			Role:      role,
			Valid:     true,
		}, nil
	}

	return &types.ValidateInvitationResponse{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Role:      role,
		Valid:     true,
	}, nil
}
