package organizations

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"

	"outlet/internal/db"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrganizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrganizationLogic {
	return &CreateOrganizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrganizationLogic) CreateOrganization(req *types.CreateOrgRequest) (resp *types.OrgInfo, err error) {
	// Generate API key
	apiKeyBytes := make([]byte, 32)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return nil, err
	}
	apiKey := hex.EncodeToString(apiKeyBytes)

	// Set defaults
	maxContacts := int64(1000)
	if req.MaxContacts > 0 {
		maxContacts = int64(req.MaxContacts)
	}

	// Generate org ID
	orgID := uuid.New().String()

	org, err := l.svcCtx.DB.CreateOrganization(l.ctx, db.CreateOrganizationParams{
		ID:          orgID,
		Name:        req.Name,
		Slug:        req.Slug,
		ApiKey:      apiKey,
		MaxContacts: sql.NullInt64{Int64: maxContacts, Valid: true},
		Settings:    sql.NullString{String: "{}", Valid: true},
		AppUrl:      sql.NullString{String: req.AppUrl, Valid: req.AppUrl != ""},
	})
	if err != nil {
		return nil, err
	}

	// Seed default email designs for the org
	l.seedDefaultDesigns(org.ID, req.Name)

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

// seedDefaultDesigns creates starter email designs for a new organization
func (l *CreateOrganizationLogic) seedDefaultDesigns(orgID string, orgName string) {
	// Simple design - minimal footer only
	simpleHTML := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f9fafb;">
  <div style="max-width: 600px; margin: 0 auto; padding: 40px 20px;">
    {{content}}

    <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #e5e7eb; text-align: center; color: #6b7280; font-size: 12px;">
      <p style="margin: 0;">` + orgName + `</p>
      <p style="margin: 8px 0 0 0;">
        <a href="{{unsubscribe_url}}" style="color: #6b7280;">Unsubscribe</a>
      </p>
    </div>
  </div>
</body>
</html>`

	// Branded design - header + footer with branding
	brandedHTML := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f9fafb;">
  <div style="max-width: 600px; margin: 0 auto;">
    <!-- Header -->
    <div style="background-color: #1f2937; padding: 24px 20px; text-align: center;">
      <h1 style="margin: 0; color: #ffffff; font-size: 24px; font-weight: 600;">` + orgName + `</h1>
    </div>

    <!-- Content -->
    <div style="background-color: #ffffff; padding: 40px 20px;">
      {{content}}
    </div>

    <!-- Footer -->
    <div style="background-color: #f3f4f6; padding: 24px 20px; text-align: center; color: #6b7280; font-size: 12px;">
      <p style="margin: 0; font-weight: 500;">` + orgName + `</p>
      <p style="margin: 12px 0 0 0;">
        <a href="{{unsubscribe_url}}" style="color: #6b7280;">Unsubscribe</a> ·
        <a href="{{preferences_url}}" style="color: #6b7280;">Email Preferences</a>
      </p>
    </div>
  </div>
</body>
</html>`

	// Create Simple design
	_, err := l.svcCtx.DB.CreateEmailDesign(l.ctx, db.CreateEmailDesignParams{
		OrgID:       orgID,
		Name:        "Simple",
		Slug:        "simple",
		Description: sql.NullString{String: "Minimal design with footer only", Valid: true},
		Category:    sql.NullString{String: "general", Valid: true},
		HtmlBody:    simpleHTML,
		PlainText:   sql.NullString{String: "{{content}}\n\n---\n" + orgName + "\nUnsubscribe: {{unsubscribe_url}}", Valid: true},
		IsActive:    sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to seed simple design: %v", err)
	}

	// Create Branded design
	_, err = l.svcCtx.DB.CreateEmailDesign(l.ctx, db.CreateEmailDesignParams{
		OrgID:       orgID,
		Name:        "Branded",
		Slug:        "branded",
		Description: sql.NullString{String: "Full branded design with header and footer", Valid: true},
		Category:    sql.NullString{String: "general", Valid: true},
		HtmlBody:    brandedHTML,
		PlainText:   sql.NullString{String: "=== " + orgName + " ===\n\n{{content}}\n\n---\n" + orgName + "\nUnsubscribe: {{unsubscribe_url}}", Valid: true},
		IsActive:    sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to seed branded design: %v", err)
	}
}

