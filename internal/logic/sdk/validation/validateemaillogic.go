package validation

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateEmailLogic {
	return &ValidateEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidateEmailLogic) ValidateEmail(req *types.ValidateEmailRequest) (resp *types.ValidationResult, err error) {
	// todo: add your logic here and delete this line

	return
}
