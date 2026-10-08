package validation

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidateBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateBatchLogic {
	return &ValidateBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidateBatchLogic) ValidateBatch(req *types.ValidateBatchRequest) (resp *types.ValidateBatchResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
