package validation

import (
	"context"

	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminValidateEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminValidateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminValidateEmailLogic {
	return &AdminValidateEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminValidateEmailLogic) AdminValidateEmail(req *types.ValidateEmailRequest) (resp *types.ValidationResult, err error) {
	// todo: add your logic here and delete this line

	return
}
