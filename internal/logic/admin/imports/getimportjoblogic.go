package imports

import (
	"context"
	"errors"
	"strconv"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/middleware"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetImportJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetImportJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImportJobLogic {
	return &GetImportJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetImportJobLogic) GetImportJob(req *types.GetImportJobRequest) (resp *types.ImportJobInfo, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, errors.New("org_id not found in context")
	}

	job, err := l.svcCtx.DB.GetImportJob(l.ctx, db.GetImportJobParams{
		ID:    req.Id,
		OrgID: orgID,
	})
	if err != nil {
		l.Errorf("Failed to get import job: %v", err)
		return nil, err
	}

	listID := ""
	if job.ListID.Valid {
		listID = strconv.FormatInt(job.ListID.Int64, 10)
	}

	status := "pending"
	if job.Status.Valid {
		status = job.Status.String
	}

	return &types.ImportJobInfo{
		Id:            job.ID,
		OrgId:         job.OrgID,
		ListId:        listID,
		Type:          job.Type,
		Status:        status,
		Filename:      job.Filename,
		TotalRows:     int(job.TotalRows.Int64),
		ProcessedRows: int(job.ProcessedRows.Int64),
		SuccessCount:  int(job.SuccessCount.Int64),
		ErrorCount:    int(job.ErrorCount.Int64),
	}, nil
}
