package public

import (
	"context"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublicSetupStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Check if initial setup is required
func NewGetPublicSetupStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublicSetupStatusLogic {
	return &GetPublicSetupStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPublicSetupStatusLogic) GetPublicSetupStatus() (resp *types.SetupStatusResponse, err error) {
	// Check if any admin user exists
	adminCount, err := l.svcCtx.DB.CountUsers(l.ctx, db.CountUsersParams{
		FilterRole:   "super_admin",
		FilterStatus: nil,
	})
	if err != nil {
		l.Error("failed to count admin users:", err)
		return nil, err
	}
	hasAdmin := adminCount > 0

	// Check for AWS configuration
	awsSettings, err := l.svcCtx.DB.GetPlatformSettingsByCategory(l.ctx, "aws")
	if err != nil {
		l.Error("failed to get AWS settings:", err)
		return nil, err
	}

	awsConfigured := make(map[string]bool)
	for _, setting := range awsSettings {
		hasValue := (setting.ValueText.Valid && setting.ValueText.String != "") ||
			(setting.ValueEncrypted.Valid && setting.ValueEncrypted.String != "")
		if hasValue {
			awsConfigured[setting.Key] = true
		}
	}
	hasAws := awsConfigured["aws_access_key"] && awsConfigured["aws_secret_key"]

	// Check for SMTP configuration (legacy)
	smtpRequiredSettings := []string{"smtp_host", "smtp_port", "smtp_user", "smtp_password", "from_email"}
	smtpSettings, err := l.svcCtx.DB.GetPlatformSettingsByCategory(l.ctx, "email")
	if err != nil {
		l.Error("failed to get SMTP settings:", err)
		return nil, err
	}

	smtpConfigured := make(map[string]bool)
	for _, setting := range smtpSettings {
		hasValue := (setting.ValueText.Valid && setting.ValueText.String != "") ||
			(setting.ValueEncrypted.Valid && setting.ValueEncrypted.String != "")
		if hasValue {
			smtpConfigured[setting.Key] = true
		}
	}

	var missingSettings []string
	for _, required := range smtpRequiredSettings {
		if !smtpConfigured[required] {
			missingSettings = append(missingSettings, required)
		}
	}
	hasSmtp := len(missingSettings) == 0

	// Platform is configured if we have either AWS or SMTP
	platformConfigured := hasAdmin && (hasAws || hasSmtp)

	// Setup is required if no admin exists
	setupRequired := !hasAdmin

	return &types.SetupStatusResponse{
		SetupRequired:      setupRequired,
		HasAdmin:           hasAdmin,
		HasAws:             hasAws,
		HasSmtp:            hasSmtp,
		PlatformConfigured: platformConfigured,
		MissingSettings:    missingSettings,
	}, nil
}
