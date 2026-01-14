package backup

import (
	"context"
	"database/sql"
	"strconv"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBackupSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBackupSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBackupSettingsLogic {
	return &GetBackupSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBackupSettingsLogic) GetBackupSettings() (resp *types.BackupSettingsResponse, err error) {
	// Get settings from platform_settings table
	settings, err := l.svcCtx.DB.GetPlatformSettingsByCategory(l.ctx, "backup")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Build response from settings
	resp = &types.BackupSettingsResponse{
		S3Enabled:       false,
		S3Bucket:        "",
		S3Region:        "",
		S3Prefix:        "",
		HasS3Creds:      false,
		ScheduleEnabled: false,
		ScheduleCron:    "0 3 * * *", // Default: 3 AM daily
		RetentionDays:   30,          // Default: 30 days
	}

	for _, s := range settings {
		value := s.ValueText.String
		switch s.Key {
		case "backup.s3_enabled":
			resp.S3Enabled = value == "true"
		case "backup.s3_bucket":
			resp.S3Bucket = value
		case "backup.s3_region":
			resp.S3Region = value
		case "backup.s3_prefix":
			resp.S3Prefix = value
		case "backup.s3_access_key":
			// Just check if credentials exist (don't return actual value)
			resp.HasS3Creds = s.ValueEncrypted.Valid && s.ValueEncrypted.String != ""
		case "backup.schedule_enabled":
			resp.ScheduleEnabled = value == "true"
		case "backup.schedule_cron":
			if value != "" {
				resp.ScheduleCron = value
			}
		case "backup.retention_days":
			if days, err := strconv.Atoi(value); err == nil {
				resp.RetentionDays = days
			}
		}
	}

	// Get last backup time
	lastBackup, err := l.svcCtx.DB.GetLatestBackup(l.ctx)
	if err == nil {
		resp.LastBackupAt = lastBackup.CompletedAt.String
	}

	// TODO: Calculate next backup time from cron expression

	return resp, nil
}
