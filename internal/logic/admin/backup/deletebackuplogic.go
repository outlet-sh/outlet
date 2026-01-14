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

type DeleteBackupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBackupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBackupLogic {
	return &DeleteBackupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBackupLogic) DeleteBackup(req *types.DeleteBackupRequest) (resp *types.Response, err error) {
	// Get backup to find file path
	backup, err := l.svcCtx.DB.GetBackup(l.ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("backup not found")
		}
		return nil, err
	}

	// Delete local file if it exists
	if backup.FilePath.Valid && backup.FilePath.String != "" {
		if err := os.Remove(backup.FilePath.String); err != nil && !os.IsNotExist(err) {
			l.Errorf("Failed to delete backup file %s: %v", backup.FilePath.String, err)
		}
	}

	// TODO: Delete from S3 if stored there
	if backup.StorageType == "s3" && backup.S3Bucket.Valid && backup.S3Key.Valid {
		// Will implement S3 delete when S3 service is available
		l.Infof("S3 backup deletion not yet implemented: %s/%s", backup.S3Bucket.String, backup.S3Key.String)
	}

	// Delete from database
	if err := l.svcCtx.DB.DeleteBackup(l.ctx, req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete backup record: %w", err)
	}

	return &types.Response{
		Success: true,
		Message: "Backup deleted successfully",
	}, nil
}
