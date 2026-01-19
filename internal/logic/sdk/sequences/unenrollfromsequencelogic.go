package sequences

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/middleware"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnenrollFromSequenceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnenrollFromSequenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnenrollFromSequenceLogic {
	return &UnenrollFromSequenceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnenrollFromSequenceLogic) UnenrollFromSequence(req *types.UnenrollSequenceRequest) (resp *types.Response, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return &types.Response{Success: false, Message: "organization not found"}, fmt.Errorf("organization not found")
	}

	if req.Email == "" {
		return &types.Response{Success: false, Message: "email is required"}, fmt.Errorf("email is required")
	}

	contact, err := l.svcCtx.DB.GetContactByOrgAndEmail(l.ctx, db.GetContactByOrgAndEmailParams{
		OrgID: sql.NullString{String: orgID, Valid: true},
		Email: req.Email,
	})
	if err != nil {
		l.Errorf("Contact not found: %s, error: %v", req.Email, err)
		return &types.Response{Success: false, Message: "contact not found"}, fmt.Errorf("contact not found")
	}

	if req.SequenceSlug != "" {
		return l.unenrollFromSpecificSequence(contact.ID, orgID, req.SequenceSlug)
	}

	return l.unenrollFromAllSequences(contact.ID, orgID)
}

func (l *UnenrollFromSequenceLogic) unenrollFromSpecificSequence(contactID, orgID string, sequenceSlug string) (*types.Response, error) {
	sequence, err := l.svcCtx.DB.GetSequenceByOrgAndSlug(l.ctx, db.GetSequenceByOrgAndSlugParams{
		OrgID: sql.NullString{String: orgID, Valid: true},
		Slug:  sequenceSlug,
	})
	if err != nil {
		l.Errorf("Sequence not found: %s, error: %v", sequenceSlug, err)
		return &types.Response{Success: false, Message: "sequence not found"}, fmt.Errorf("sequence not found")
	}

	state, err := l.svcCtx.DB.GetContactSequenceState(l.ctx, db.GetContactSequenceStateParams{
		ContactID:  sql.NullString{String: contactID, Valid: true},
		SequenceID: sql.NullString{String: sequence.ID, Valid: true},
	})
	if err != nil {
		l.Errorf("Enrollment not found: sequence=%s error: %v", sequenceSlug, err)
		return &types.Response{Success: false, Message: "not enrolled in sequence"}, fmt.Errorf("not enrolled in sequence")
	}

	if state.UnsubscribedAt.Valid {
		return &types.Response{Success: true, Message: "already unenrolled from sequence"}, nil
	}

	err = l.svcCtx.DB.CancelContactSequence(l.ctx, db.CancelContactSequenceParams{
		ContactID:  sql.NullString{String: contactID, Valid: true},
		SequenceID: sql.NullString{String: sequence.ID, Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to cancel sequence: %v", err)
		return &types.Response{Success: false, Message: "failed to unenroll from sequence"}, fmt.Errorf("failed to unenroll from sequence")
	}

	err = l.svcCtx.DB.CancelPendingEmailsForContactSequence(l.ctx, db.CancelPendingEmailsForContactSequenceParams{
		ContactID:  sql.NullString{String: contactID, Valid: true},
		SequenceID: sql.NullString{String: sequence.ID, Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to cancel pending emails: %v", err)
	}

	l.Infof("Unenrolled contact from sequence: org=%s contact=%s sequence=%s", orgID, contactID, sequenceSlug)

	return &types.Response{Success: true, Message: "unenrolled from sequence"}, nil
}

func (l *UnenrollFromSequenceLogic) unenrollFromAllSequences(contactID, orgID string) (*types.Response, error) {
	err := l.svcCtx.DB.CancelAllContactSequences(l.ctx, db.CancelAllContactSequencesParams{
		ContactID: sql.NullString{String: contactID, Valid: true},
		OrgID:     sql.NullString{String: orgID, Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to cancel all sequences: %v", err)
		return &types.Response{Success: false, Message: "failed to unenroll from sequences"}, fmt.Errorf("failed to unenroll from sequences")
	}

	err = l.svcCtx.DB.CancelEmailsForContact(l.ctx, sql.NullString{String: contactID, Valid: true})
	if err != nil {
		l.Errorf("Failed to cancel pending emails: %v", err)
	}

	l.Infof("Unenrolled contact from all sequences: org=%s contact=%s", orgID, contactID)

	return &types.Response{Success: true, Message: "unenrolled from all sequences"}, nil
}
