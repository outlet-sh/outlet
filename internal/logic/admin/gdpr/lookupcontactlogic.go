package gdpr

import (
	"context"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LookupContactLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLookupContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LookupContactLogic {
	return &LookupContactLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LookupContactLogic) LookupContact(req *types.GDPRLookupRequest) (resp *types.GDPRLookupResponse, err error) {
	if req.Email == "" {
		return &types.GDPRLookupResponse{
			Found: false,
		}, nil
	}

	// Look up contact by email
	contact, err := l.svcCtx.DB.GetContactByEmail(l.ctx, req.Email)
	if err != nil {
		// Not found is not an error, just return found=false
		return &types.GDPRLookupResponse{
			Found: false,
		}, nil
	}

	return &types.GDPRLookupResponse{
		Found:     true,
		ContactId: contact.ID,
		Email:     contact.Email,
	}, nil
}
