package backup

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadBackupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadBackupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadBackupLogic {
	return &DownloadBackupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DownloadBackup returns the backup file info. The actual download is handled by a custom handler.
// This method validates the backup exists and is ready for download.
func (l *DownloadBackupLogic) DownloadBackup(req *types.DownloadBackupRequest) error {
	// Get backup record
	backup, err := l.svcCtx.DB.GetBackup(l.ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("backup not found")
		}
		return err
	}

	// Check if backup is completed
	if backup.Status != "completed" {
		return fmt.Errorf("backup is not ready for download (status: %s)", backup.Status)
	}

	// Check if file exists locally
	if !backup.FilePath.Valid || backup.FilePath.String == "" {
		return fmt.Errorf("backup file not available")
	}

	// Verify file exists on disk
	if _, err := os.Stat(backup.FilePath.String); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found on disk")
	}

	// Success - the handler will stream the file
	return nil
}

// GetBackupFilePath returns the file path for a backup (used by custom download handler)
func (l *DownloadBackupLogic) GetBackupFilePath(id string) (string, string, error) {
	backup, err := l.svcCtx.DB.GetBackup(l.ctx, id)
	if err != nil {
		return "", "", err
	}

	if backup.Status != "completed" {
		return "", "", fmt.Errorf("backup not ready")
	}

	if !backup.FilePath.Valid {
		return "", "", fmt.Errorf("no file path")
	}

	return backup.FilePath.String, backup.Filename, nil
}
