package organizations

import (
	"context"
	"database/sql"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrganizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrganizationLogic {
	return &UpdateOrganizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrganizationLogic) UpdateOrganization(req *types.UpdateOrgRequest) (resp *types.OrgInfo, err error) {
	// Build update params with optional fields
	params := db.UpdateOrganizationParams{
		ID:   req.Id,
		Name: req.Name, // Empty string becomes NULL via NULLIF in query
	}

	if req.MaxContacts > 0 {
		params.MaxContacts = sql.NullInt64{Int64: int64(req.MaxContacts), Valid: true}
	}
	// AppUrl uses COALESCE/NULLIF so empty string preserves existing value
	params.AppUrl = req.AppUrl

	org, err := l.svcCtx.DB.UpdateOrganization(l.ctx, params)
	if err != nil {
		return nil, err
	}

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
