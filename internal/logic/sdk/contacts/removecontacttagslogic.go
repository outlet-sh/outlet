package contacts

import (
	"context"
	"database/sql"
	"strings"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveContactTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveContactTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveContactTagsLogic {
	return &RemoveContactTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveContactTagsLogic) RemoveContactTags(req *types.RemoveContactTagsRequest) (resp *types.Response, err error) {
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

	// Remove each tag
	for _, tag := range req.Tags {
		_ = l.svcCtx.DB.RemoveContactTag(l.ctx, db.RemoveContactTagParams{
			ContactID: sql.NullString{String: contact.ID, Valid: true},
			Tag:       tag,
		})
	}

	l.Infof("Removed tags from contact: org=%s id=%s tags=%v", orgID, contact.ID, req.Tags)

	return &types.Response{Success: true, Message: "tags removed"}, nil
}
