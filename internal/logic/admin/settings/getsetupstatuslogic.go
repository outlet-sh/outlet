package settings

import (
	"context"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSetupStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSetupStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSetupStatusLogic {
	return &GetSetupStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSetupStatusLogic) GetSetupStatus() (resp *types.SetupStatusResponse, err error) {
	// Required email settings for platform to be fully configured
	requiredSettings := []string{"smtp_host", "smtp_port", "smtp_user", "smtp_password", "from_email"}

	// Get all email-related platform settings
	settings, err := l.svcCtx.DB.GetPlatformSettingsByCategory(l.ctx, "email")
	if err != nil {
		l.Error("failed to get platform settings:", err)
		return nil, err
	}

	// Build a map of configured settings
	configuredSettings := make(map[string]bool)
	for _, setting := range settings {
		// Check if the setting has a value (either text or encrypted)
		hasValue := (setting.ValueText.Valid && setting.ValueText.String != "") ||
			(setting.ValueEncrypted.Valid && setting.ValueEncrypted.String != "")
		if hasValue {
			configuredSettings[setting.Key] = true
		}
	}

	// Find missing settings
	var missingSettings []string
	for _, required := range requiredSettings {
		if !configuredSettings[required] {
			missingSettings = append(missingSettings, required)
		}
	}

	return &types.SetupStatusResponse{
		PlatformConfigured: len(missingSettings) == 0,
		MissingSettings:    missingSettings,
	}, nil
}
