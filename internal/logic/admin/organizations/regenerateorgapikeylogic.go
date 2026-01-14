package organizations

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"outlet/internal/db"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegenerateOrgApiKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegenerateOrgApiKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegenerateOrgApiKeyLogic {
	return &RegenerateOrgApiKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegenerateOrgApiKeyLogic) RegenerateOrgApiKey(req *types.RegenerateApiKeyRequest) (resp *types.OrgInfo, err error) {
	// Get current org to retrieve old API key for cache invalidation
	currentOrg, err := l.svcCtx.DB.GetOrganizationByID(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	oldApiKey := currentOrg.ApiKey

	// Generate new API key
	apiKeyBytes := make([]byte, 32)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return nil, err
	}
	newApiKey := hex.EncodeToString(apiKeyBytes)

	// Update the API key
	org, err := l.svcCtx.DB.RegenerateAPIKey(l.ctx, db.RegenerateAPIKeyParams{
		ID:     req.Id,
		ApiKey: newApiKey,
	})
	if err != nil {
		return nil, err
	}

	// Invalidate old API key from cache so it stops working immediately
	l.svcCtx.APIKeyMiddleware.InvalidateCache(oldApiKey)

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
