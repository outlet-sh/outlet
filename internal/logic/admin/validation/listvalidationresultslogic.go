package validation

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListValidationResultsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListValidationResultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListValidationResultsLogic {
	return &ListValidationResultsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListValidationResultsLogic) ListValidationResults(req *types.ListValidationResultsRequest) (resp *types.ListValidationResultsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
