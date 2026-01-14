package sdk

import (
	"context"
	"database/sql"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsubscribeFromListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsubscribeFromListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsubscribeFromListLogic {
	return &UnsubscribeFromListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsubscribeFromListLogic) UnsubscribeFromList(req *types.SubscribeRequest) (resp *types.Response, err error) {
	// Get org ID from context (set by API key middleware)
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return &types.Response{Success: false, Message: "Organization not found"}, nil
	}

	// Validate required fields
	if req.Email == "" {
		return &types.Response{Success: false, Message: "Email is required"}, nil
	}
	if req.Slug == "" {
		return &types.Response{Success: false, Message: "List slug is required"}, nil
	}

	// Get the email list by slug
	list, err := l.svcCtx.DB.GetEmailListByOrgAndSlug(l.ctx, db.GetEmailListByOrgAndSlugParams{
		OrgID: orgID,
		Slug:  req.Slug,
	})
	if err != nil {
		l.Errorf("List not found: %s, error: %v", req.Slug, err)
		return &types.Response{Success: false, Message: "List not found"}, nil
	}

	// Get contact by email
	contact, err := l.svcCtx.DB.GetContactByOrgAndEmail(l.ctx, db.GetContactByOrgAndEmailParams{
		OrgID: sql.NullString{String: orgID, Valid: true},
		Email: req.Email,
	})
	if err != nil {
		l.Infof("Contact not found for unsubscribe: %s", req.Email)
		return &types.Response{Success: true, Message: "Unsubscribed"}, nil
	}

	// Unsubscribe from list
	err = l.svcCtx.DB.UnsubscribeFromList(l.ctx, db.UnsubscribeFromListParams{
		ListID:    list.ID,
		ContactID: contact.ID,
	})
	if err != nil {
		l.Errorf("Failed to unsubscribe: %v", err)
		return &types.Response{Success: false, Message: "Failed to unsubscribe"}, nil
	}

	l.Infof("Unsubscribed: org=%s email=%s list=%s", orgID, req.Email, req.Slug)
	return &types.Response{Success: true, Message: "Successfully unsubscribed"}, nil
}
