package backup

import (
	"context"
	"database/sql"
	"fmt"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBackupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBackupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBackupLogic {
	return &GetBackupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBackupLogic) GetBackup(req *types.GetBackupRequest) (resp *types.BackupInfo, err error) {
	backup, err := l.svcCtx.DB.GetBackup(l.ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("backup not found")
		}
		return nil, err
	}

	result := dbBackupToType(backup)
	return &result, nil
}
