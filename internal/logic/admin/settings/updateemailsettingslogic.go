package settings

import (
	"context"
	"database/sql"
	"fmt"

	"outlet/internal/db"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEmailSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEmailSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmailSettingsLogic {
	return &UpdateEmailSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEmailSettingsLogic) UpdateEmailSettings(req *types.UpdateEmailSettingsRequest) (resp *types.UpdateSettingsResponse, err error) {
	const category = "email"

	// Save SMTP host (non-sensitive)
	if req.SmtpHost != "" {
		_, err := l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:         "smtp_host",
			ValueText:   sql.NullString{String: req.SmtpHost, Valid: true},
			Description: sql.NullString{String: "SMTP server hostname", Valid: true},
			Category:    category,
			IsSensitive: sql.NullInt64{Int64: 0, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save smtp_host: %w", err)
		}
	}

	// Save SMTP port (non-sensitive)
	if req.SmtpPort > 0 {
		_, err := l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:         "smtp_port",
			ValueText:   sql.NullString{String: fmt.Sprintf("%d", req.SmtpPort), Valid: true},
			Description: sql.NullString{String: "SMTP server port", Valid: true},
			Category:    category,
			IsSensitive: sql.NullInt64{Int64: 0, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save smtp_port: %w", err)
		}
	}

	// Save SMTP user (non-sensitive - just the username)
	if req.SmtpUser != "" {
		_, err := l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:         "smtp_user",
			ValueText:   sql.NullString{String: req.SmtpUser, Valid: true},
			Description: sql.NullString{String: "SMTP username", Valid: true},
			Category:    category,
			IsSensitive: sql.NullInt64{Int64: 0, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save smtp_user: %w", err)
		}
	}

	// Save SMTP password (sensitive - encrypted)
	if req.SmtpPassword != "" {
		if l.svcCtx.CryptoService == nil {
			return nil, fmt.Errorf("encryption service not configured - cannot store sensitive data")
		}

		encrypted, err := l.svcCtx.CryptoService.EncryptString(req.SmtpPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt smtp_password: %w", err)
		}

		_, err = l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:            "smtp_password",
			ValueEncrypted: sql.NullString{String: string(encrypted), Valid: true},
			Description:    sql.NullString{String: "SMTP password", Valid: true},
			Category:       category,
			IsSensitive:    sql.NullInt64{Int64: 1, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save smtp_password: %w", err)
		}
	}

	// Save From Email (non-sensitive)
	if req.FromEmail != "" {
		_, err := l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:         "from_email",
			ValueText:   sql.NullString{String: req.FromEmail, Valid: true},
			Description: sql.NullString{String: "Default from email address", Valid: true},
			Category:    category,
			IsSensitive: sql.NullInt64{Int64: 0, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save from_email: %w", err)
		}
	}

	// Save From Name (non-sensitive)
	if req.FromName != "" {
		_, err := l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:         "from_name",
			ValueText:   sql.NullString{String: req.FromName, Valid: true},
			Description: sql.NullString{String: "Default from name", Valid: true},
			Category:    category,
			IsSensitive: sql.NullInt64{Int64: 0, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save from_name: %w", err)
		}
	}

	// Save Reply-To (non-sensitive)
	if req.ReplyTo != "" {
		_, err := l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:         "reply_to",
			ValueText:   sql.NullString{String: req.ReplyTo, Valid: true},
			Description: sql.NullString{String: "Default reply-to email address", Valid: true},
			Category:    category,
			IsSensitive: sql.NullInt64{Int64: 0, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save reply_to: %w", err)
		}
	}

	return &types.UpdateSettingsResponse{
		Success: true,
		Message: "Email settings saved successfully",
	}, nil
}
