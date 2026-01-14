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

type UpdateContactLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateContactLogic {
	return &UpdateContactLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateContactLogic) UpdateContact(req *types.UpdateContactRequest) (resp *types.SDKContactInfo, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, nil
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
		return nil, err
	}

	// Update contact name if provided
	nameToUpdate := req.Name
	if nameToUpdate == "" {
		nameToUpdate = contact.Name
	}

	updated, err := l.svcCtx.DB.UpdateSDKContact(l.ctx, db.UpdateSDKContactParams{
		ID:    contact.ID,
		OrgID: sql.NullString{String: orgID, Valid: true},
		Name:  nameToUpdate,
	})
	if err != nil {
		l.Errorf("Failed to update contact: %v", err)
		return nil, err
	}

	// Get contact tags
	tags, _ := l.svcCtx.DB.GetContactTags(l.ctx, sql.NullString{String: updated.ID, Valid: true})
	var tagList []string
	for _, t := range tags {
		tagList = append(tagList, t.Tag)
	}

	// Determine status
	status := "active"
	if updated.Status.Valid && updated.Status.String != "" {
		status = updated.Status.String
	}
	if updated.UnsubscribedAt.Valid {
		status = "unsubscribed"
	}
	if updated.BlockedAt.Valid {
		status = "blocked"
	}

	l.Infof("Updated contact: org=%s id=%s", orgID, updated.ID)

	return &types.SDKContactInfo{
		Id:            updated.ID,
		Email:         updated.Email,
		Name:          updated.Name,
		Status:        status,
		EmailVerified: updated.EmailVerified == 1,
		Tags:          tagList,
		Lists:         []string{},
		Source:        updated.Source.String,
		CreatedAt:     updated.CreatedAt.String,
		UpdatedAt:     updated.UpdatedAt.String,
	}, nil
}
