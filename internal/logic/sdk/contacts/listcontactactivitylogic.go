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

type ListContactActivityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListContactActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContactActivityLogic {
	return &ListContactActivityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListContactActivityLogic) ListContactActivity(req *types.ListContactActivityRequest) (resp *types.ListContactActivityResponse, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return &types.ListContactActivityResponse{Activities: []types.ContactActivityInfo{}}, nil
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
		return &types.ListContactActivityResponse{Activities: []types.ContactActivityInfo{}}, nil
	}

	// Build activity list from contact events
	var activities []types.ContactActivityInfo

	// Add subscription event
	if contact.CreatedAt.Valid {
		activities = append(activities, types.ContactActivityInfo{
			Event:     "subscribed",
			Timestamp: contact.CreatedAt.String,
			Details:   "Contact created",
		})
	}

	// Add verification event
	if contact.VerifiedAt.Valid {
		activities = append(activities, types.ContactActivityInfo{
			Event:     "email_verified",
			Timestamp: contact.VerifiedAt.String,
			Details:   "Email address verified",
		})
	}

	// Add unsubscribe event
	if contact.UnsubscribedAt.Valid {
		activities = append(activities, types.ContactActivityInfo{
			Event:     "unsubscribed",
			Timestamp: contact.UnsubscribedAt.String,
			Details:   "Contact unsubscribed",
		})
	}

	// Limit results if requested
	if req.Limit > 0 && len(activities) > req.Limit {
		activities = activities[:req.Limit]
	}

	return &types.ListContactActivityResponse{Activities: activities}, nil
}
