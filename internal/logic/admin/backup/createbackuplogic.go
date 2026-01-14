package backup

import (
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"outlet/internal/db"
	backupService "outlet/internal/services/backup"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBackupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBackupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBackupLogic {
	return &CreateBackupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBackupLogic) CreateBackup(req *types.CreateBackupRequest) (resp *types.CreateBackupResponse, err error) {
	// Generate backup ID and filename
	backupID := uuid.New().String()
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("outlet-backup-%s.db", timestamp)
	if req.Compress {
		filename += ".gz"
	}

	// Determine backup directory
	backupDir := filepath.Join(".", "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupPath := filepath.Join(backupDir, filename)

	// Create backup record in database
	storageType := "local"
	if req.UploadToS3 {
		storageType = "s3"
	}

	record, err := l.svcCtx.DB.CreateBackupRecord(l.ctx, db.CreateBackupRecordParams{
		ID:          backupID,
		Filename:    filename,
		FilePath:    sql.NullString{String: backupPath, Valid: true},
		FileSize:    0,
		BackupType:  "manual",
		StorageType: storageType,
		S3Bucket:    sql.NullString{String: req.S3Bucket, Valid: req.S3Bucket != ""},
		S3Key:       sql.NullString{},
		Status:      "in_progress",
		CreatedBy:   sql.NullString{}, // TODO: get from auth context
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	// Get S3 config if uploading to S3
	var s3Config *backupService.S3Config
	if req.UploadToS3 {
		s3Config = l.getS3Config()
		if s3Config == nil {
			return nil, fmt.Errorf("S3 not configured or credentials missing")
		}
	}

	// Perform the backup using SQLite VACUUM INTO
	go func() {
		ctx := context.Background()
		var finalErr error

		defer func() {
			if finalErr != nil {
				l.svcCtx.DB.UpdateBackupStatus(ctx, db.UpdateBackupStatusParams{
					ID:           backupID,
					Status:       "failed",
					ErrorMessage: sql.NullString{String: finalErr.Error(), Valid: true},
				})
			}
		}()

		// Create backup using VACUUM INTO
		tempPath := backupPath + ".tmp"
		_, finalErr = l.svcCtx.DB.GetDB().ExecContext(ctx, fmt.Sprintf("VACUUM INTO '%s'", tempPath))
		if finalErr != nil {
			return
		}

		// Compress if requested
		if req.Compress {
			if finalErr = compressFile(tempPath, backupPath); finalErr != nil {
				os.Remove(tempPath)
				return
			}
			os.Remove(tempPath)
		} else {
			if finalErr = os.Rename(tempPath, backupPath); finalErr != nil {
				return
			}
		}

		// Get file size
		info, err := os.Stat(backupPath)
		if err != nil {
			finalErr = err
			return
		}

		// Upload to S3 if requested
		var s3Key string
		if s3Config != nil {
			s3Key, finalErr = backupService.UploadToS3(ctx, *s3Config, backupPath, filename)
			if finalErr != nil {
				return
			}
		}

		// Update backup record with completion
		_, finalErr = l.svcCtx.DB.UpdateBackupComplete(ctx, db.UpdateBackupCompleteParams{
			ID:       backupID,
			FileSize: info.Size(),
		})
		if finalErr != nil {
			return
		}

		// Update S3 key if uploaded
		if s3Key != "" {
			l.svcCtx.DB.UpdateBackupS3Key(ctx, db.UpdateBackupS3KeyParams{
				ID:    backupID,
				S3Key: sql.NullString{String: s3Key, Valid: true},
			})
		}
	}()

	return &types.CreateBackupResponse{
		Backup: dbBackupToType(record),
	}, nil
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func dbBackupToType(b db.BackupHistory) types.BackupInfo {
	return types.BackupInfo{
		Id:           b.ID,
		Filename:     b.Filename,
		FilePath:     b.FilePath.String,
		FileSize:     b.FileSize,
		BackupType:   b.BackupType,
		StorageType:  b.StorageType,
		S3Bucket:     b.S3Bucket.String,
		S3Key:        b.S3Key.String,
		Status:       b.Status,
		ErrorMessage: b.ErrorMessage.String,
		CreatedBy:    b.CreatedBy.String,
		StartedAt:    b.StartedAt.String,
		CompletedAt:  b.CompletedAt.String,
		CreatedAt:    b.CreatedAt.String,
	}
}

// getS3Config retrieves S3 configuration from platform settings
func (l *CreateBackupLogic) getS3Config() *backupService.S3Config {
	settings, err := l.svcCtx.DB.GetPlatformSettingsByCategory(l.ctx, "backup")
	if err != nil {
		l.Errorf("Failed to get backup settings: %v", err)
		return nil
	}

	cfg := &backupService.S3Config{}
	var accessKeyEncrypted, secretKeyEncrypted []byte

	for _, s := range settings {
		switch s.Key {
		case "backup.s3_enabled":
			if s.ValueText.String != "true" {
				return nil // S3 not enabled
			}
		case "backup.s3_bucket":
			cfg.Bucket = s.ValueText.String
		case "backup.s3_region":
			cfg.Region = s.ValueText.String
		case "backup.s3_prefix":
			cfg.Prefix = s.ValueText.String
		case "backup.s3_access_key":
			if s.ValueEncrypted.Valid {
				var err error
				accessKeyEncrypted, err = hex.DecodeString(s.ValueEncrypted.String)
				if err != nil {
					l.Errorf("Failed to decode S3 access key: %v", err)
				}
			}
		case "backup.s3_secret_key":
			if s.ValueEncrypted.Valid {
				var err error
				secretKeyEncrypted, err = hex.DecodeString(s.ValueEncrypted.String)
				if err != nil {
					l.Errorf("Failed to decode S3 secret key: %v", err)
				}
			}
		}
	}

	// Decrypt credentials if crypto service is available
	if l.svcCtx.CryptoService != nil {
		if len(accessKeyEncrypted) > 0 {
			decrypted, err := l.svcCtx.CryptoService.DecryptString(accessKeyEncrypted)
			if err == nil {
				cfg.AccessKey = decrypted
			} else {
				l.Errorf("Failed to decrypt S3 access key: %v", err)
			}
		}
		if len(secretKeyEncrypted) > 0 {
			decrypted, err := l.svcCtx.CryptoService.DecryptString(secretKeyEncrypted)
			if err == nil {
				cfg.SecretKey = decrypted
			} else {
				l.Errorf("Failed to decrypt S3 secret key: %v", err)
			}
		}
	}

	// Validate required fields
	if cfg.Bucket == "" {
		l.Errorf("S3 bucket not configured")
		return nil
	}

	return cfg
}
