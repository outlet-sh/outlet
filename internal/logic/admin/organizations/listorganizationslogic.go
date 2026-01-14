package organizations

import (
	"context"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrganizationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrganizationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrganizationsLogic {
	return &ListOrganizationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOrganizationsLogic) ListOrganizations() (resp *types.OrgListResponse, err error) {
	orgs, err := l.svcCtx.DB.ListOrganizations(l.ctx)
	if err != nil {
		return nil, err
	}

	orgInfos := make([]types.OrgInfo, len(orgs))
	for i, org := range orgs {
		orgInfos[i] = types.OrgInfo{
			Id:               org.ID,
			Name:             org.Name,
			Slug:             org.Slug,
			ApiKey:           org.ApiKey,
			BillingStatus:    "trial",
			Plan:             "starter",
			MaxContacts:      int(org.MaxContacts.Int64),
			StripeConfigured: false,
			AppUrl:           org.AppUrl.String,
			CreatedAt:        org.CreatedAt.String,
		}
	}

	return &types.OrgListResponse{
		Organizations: orgInfos,
		Total:         len(orgInfos),
	}, nil
}
