package organizations

import (
	"context"

	"outlet/internal/db"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrgEmailSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrgEmailSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrgEmailSettingsLogic {
	return &UpdateOrgEmailSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrgEmailSettingsLogic) UpdateOrgEmailSettings(req *types.UpdateOrgEmailSettingsRequest) (resp *types.OrgInfo, err error) {
	// Update email settings (from_name, from_email, reply_to)
	org, err := l.svcCtx.DB.UpdateOrgEmailSettings(l.ctx, db.UpdateOrgEmailSettingsParams{
		ID:        req.Id,
		FromName:  req.FromName,
		FromEmail: req.FromEmail,
		ReplyTo:   req.ReplyTo,
	})
	if err != nil {
		return nil, err
	}

	// Invalidate API key middleware cache so SDK gets updated org settings
	l.svcCtx.APIKeyMiddleware.InvalidateCache(org.ApiKey)

	return &types.OrgInfo{
		Id:               org.ID,
		Name:             org.Name,
		Slug:             org.Slug,
		ApiKey:           org.ApiKey,
		BillingStatus:    "trial",
		Plan:             "starter",
		MaxContacts:      int(org.MaxContacts.Int64),
		StripeConfigured: false,
		FromName:         org.FromName.String,
		FromEmail:        org.FromEmail.String,
		ReplyTo:          org.ReplyTo.String,
		AppUrl:           org.AppUrl.String,
		CreatedAt:        org.CreatedAt.String,
	}, nil
}
