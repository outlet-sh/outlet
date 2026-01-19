package contacts

import (
	"context"
	"database/sql"
	"strings"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/middleware"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddContactTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddContactTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddContactTagsLogic {
	return &AddContactTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddContactTagsLogic) AddContactTags(req *types.AddContactTagsRequest) (resp *types.Response, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return &types.Response{Success: false, Message: "unauthorized"}, nil
	}

	var contact db.Contact

	// Find contact by ID or email
	if strings.Contains(req.Id, "@") {
		contact, err = l.svcCtx.DB.GetContactByOrgAndEmail(l.ctx, db.GetContactByOrgAndEmailParams{
			OrgID: sql.NullString{String: orgID, Valid: true},
			Email: req.Id,
		})
	} else {
		contact, err = l.svcCtx.DB.GetContactByOrgID(l.ctx, db.GetContactByOrgIDParams{
			ID:    req.Id,
			OrgID: sql.NullString{String: orgID, Valid: true},
		})
	}

	if err != nil {
		l.Errorf("Failed to get contact: %v", err)
		return &types.Response{Success: false, Message: "contact not found"}, nil
	}

	// Add each tag
	for _, tag := range req.Tags {
		_, _ = l.svcCtx.DB.AddContactTag(l.ctx, db.AddContactTagParams{
			ContactID: sql.NullString{String: contact.ID, Valid: true},
			Tag:       tag,
		})
	}

	l.Infof("Added tags to contact: org=%s id=%s tags=%v", orgID, contact.ID, req.Tags)

	return &types.Response{Success: true, Message: "tags added"}, nil
}
