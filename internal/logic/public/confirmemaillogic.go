package public

import (
	"context"
	"database/sql"

	"outlet/internal/db"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmEmailLogic {
	return &ConfirmEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmEmailLogic) ConfirmEmail(req *types.ConfirmEmailRequest) (resp *types.ConfirmEmailResponse, err error) {
	if req.Token == "" {
		return &types.ConfirmEmailResponse{
			Success: false,
			Message: "Invalid confirmation link",
		}, nil
	}

	// Look up contact by verification token
	contact, err := l.svcCtx.DB.GetContactByVerificationToken(l.ctx, sql.NullString{String: req.Token, Valid: true})
	if err != nil {
		logx.Infof("Invalid or expired verification token: %s", req.Token)
		return &types.ConfirmEmailResponse{
			Success:  false,
			Message:  "This confirmation link is invalid or has expired.",
			Redirect: "/confirm-expired",
		}, nil
	}

	// Verify the email
	err = l.svcCtx.DB.VerifyContactEmail(l.ctx, contact.ID)
	if err != nil {
		logx.Errorf("Failed to verify email for contact %s: %v", contact.ID, err)
		return &types.ConfirmEmailResponse{
			Success: false,
			Message: "Failed to confirm email. Please try again.",
		}, nil
	}

	// Update contact status
	_ = l.svcCtx.DB.UpdateContactStatus(l.ctx, db.UpdateContactStatusParams{
		ID:     contact.ID,
		Status: sql.NullString{String: "verified", Valid: true},
	})

	// Add verified tag
	_, _ = l.svcCtx.DB.AddContactTag(l.ctx, db.AddContactTagParams{
		ContactID: sql.NullString{String: contact.ID, Valid: true},
		Tag:       "email_verified",
	})

	logx.Infof("Email verified for contact %s (%s)", contact.ID, contact.Email)

	return &types.ConfirmEmailResponse{
		Success:  true,
		Message:  "Email confirmed successfully!",
		Redirect: "/confirmed",
	}, nil
}
