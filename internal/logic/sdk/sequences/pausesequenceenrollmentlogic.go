package sequences

import (
	"context"
	"database/sql"
	"fmt"

	"outlet/internal/db"
	"outlet/internal/middleware"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PauseSequenceEnrollmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPauseSequenceEnrollmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PauseSequenceEnrollmentLogic {
	return &PauseSequenceEnrollmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PauseSequenceEnrollmentLogic) PauseSequenceEnrollment(req *types.PauseSequenceRequest) (resp *types.Response, err error) {
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

	if state.PausedAt.Valid {
		return &types.Response{Success: true, Message: "sequence already paused"}, nil
	}

	if state.CompletedAt.Valid {
		return &types.Response{Success: false, Message: "sequence already completed"}, fmt.Errorf("sequence already completed")
	}
	if state.UnsubscribedAt.Valid {
		return &types.Response{Success: false, Message: "sequence already cancelled"}, fmt.Errorf("sequence already cancelled")
	}

	err = l.svcCtx.DB.PauseContactSequence(l.ctx, db.PauseContactSequenceParams{
		ContactID:  sql.NullString{String: contact.ID, Valid: true},
		SequenceID: sql.NullString{String: sequence.ID, Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to pause sequence: %v", err)
		return &types.Response{Success: false, Message: "failed to pause sequence"}, fmt.Errorf("failed to pause sequence")
	}

	l.Infof("Paused sequence enrollment: org=%s email=%s sequence=%s", orgID, req.Email, req.SequenceSlug)

	return &types.Response{Success: true, Message: "sequence paused"}, nil
}
