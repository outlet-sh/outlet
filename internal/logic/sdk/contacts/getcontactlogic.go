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

type GetContactLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContactLogic {
	return &GetContactLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetContactLogic) GetContact(req *types.GetContactRequest) (resp *types.SDKContactInfo, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, nil
	}

	var contact db.Contact

	// Try to find by ID first, then by email
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
		return nil, err
	}

	// Get contact tags
	tags, _ := l.svcCtx.DB.GetContactTags(l.ctx, sql.NullString{String: contact.ID, Valid: true})
	var tagList []string
	for _, t := range tags {
		tagList = append(tagList, t.Tag)
	}

	// Determine status
	status := "active"
	if contact.Status.Valid && contact.Status.String != "" {
		status = contact.Status.String
	}
	if contact.UnsubscribedAt.Valid {
		status = "unsubscribed"
	}
	if contact.BlockedAt.Valid {
		status = "blocked"
	}

	return &types.SDKContactInfo{
		Id:            contact.ID,
		Email:         contact.Email,
		Name:          contact.Name,
		Status:        status,
		EmailVerified: contact.EmailVerified == 1,
		Tags:          tagList,
		Lists:         []string{},
		Source:        contact.Source.String,
		CreatedAt:     contact.CreatedAt.String,
		UpdatedAt:     contact.UpdatedAt.String,
	}, nil
}
