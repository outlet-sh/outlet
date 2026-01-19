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

type ResumeSequenceEnrollmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResumeSequenceEnrollmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResumeSequenceEnrollmentLogic {
	return &ResumeSequenceEnrollmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResumeSequenceEnrollmentLogic) ResumeSequenceEnrollment(req *types.ResumeSequenceRequest) (resp *types.Response, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return &types.Response{Success: false, Message: "organization not found"}, fmt.Errorf("organization not found")
	}

	if req.Email == "" {
		return &types.Response{Success: false, Message: "email is required"}, fmt.Errorf("email is required")
	}
	if req.SequenceSlug == "" {
		return &types.Response{Success: false, Message: "sequence_slug is required"}, fmt.Errorf("sequence_slug is required")
	}

	sequence, err := l.svcCtx.DB.GetSequenceByOrgAndSlug(l.ctx, db.GetSequenceByOrgAndSlugParams{
		OrgID: sql.NullString{String: orgID, Valid: true},
		Slug:  req.SequenceSlug,
	})
	if err != nil {
		l.Errorf("Sequence not found: %s, error: %v", req.SequenceSlug, err)
		return &types.Response{Success: false, Message: "sequence not found"}, fmt.Errorf("sequence not found")
	}

	contact, err := l.svcCtx.DB.GetContactByOrgAndEmail(l.ctx, db.GetContactByOrgAndEmailParams{
		OrgID: sql.NullString{String: orgID, Valid: true},
		Email: req.Email,
	})
	if err != nil {
		l.Errorf("Contact not found: %s, error: %v", req.Email, err)
		return &types.Response{Success: false, Message: "contact not found"}, fmt.Errorf("contact not found")
	}

	state, err := l.svcCtx.DB.GetContactSequenceState(l.ctx, db.GetContactSequenceStateParams{
		ContactID:  sql.NullString{String: contact.ID, Valid: true},
		SequenceID: sql.NullString{String: sequence.ID, Valid: true},
	})
	if err != nil {
		l.Errorf("Enrollment not found: email=%s sequence=%s error: %v", req.Email, req.SequenceSlug, err)
		return &types.Response{Success: false, Message: "not enrolled in sequence"}, fmt.Errorf("not enrolled in sequence")
	}

	if state.CompletedAt.Valid {
		return &types.Response{Success: false, Message: "sequence already completed"}, fmt.Errorf("sequence already completed")
	}
	if state.UnsubscribedAt.Valid {
		return &types.Response{Success: false, Message: "sequence was cancelled"}, fmt.Errorf("sequence was cancelled")
	}

	if !state.PausedAt.Valid && state.IsActive.Valid && state.IsActive.Int64 == 1 {
		return &types.Response{Success: true, Message: "sequence already active"}, nil
	}

	err = l.svcCtx.DB.ResumeContactSequence(l.ctx, db.ResumeContactSequenceParams{
		ContactID:  sql.NullString{String: contact.ID, Valid: true},
		SequenceID: sql.NullString{String: sequence.ID, Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to resume sequence: %v", err)
		return &types.Response{Success: false, Message: "failed to resume sequence"}, fmt.Errorf("failed to resume sequence")
	}

	l.Infof("Resumed sequence enrollment: org=%s email=%s sequence=%s", orgID, req.Email, req.SequenceSlug)

	return &types.Response{Success: true, Message: "sequence resumed"}, nil
}
