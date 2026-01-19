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

type GetSequenceEnrollmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSequenceEnrollmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSequenceEnrollmentsLogic {
	return &GetSequenceEnrollmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSequenceEnrollmentsLogic) GetSequenceEnrollments(req *types.GetSequenceEnrollmentRequest) (resp *types.ListSequenceEnrollmentsResponse, err error) {
	orgID, ok := l.ctx.Value(middleware.OrgIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("organization not found")
	}

	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	contact, err := l.svcCtx.DB.GetContactByOrgAndEmail(l.ctx, db.GetContactByOrgAndEmailParams{
		OrgID: sql.NullString{String: orgID, Valid: true},
		Email: req.Email,
	})
	if err != nil {
		l.Errorf("Contact not found: %s, error: %v", req.Email, err)
		return nil, fmt.Errorf("contact not found")
	}

	if req.SequenceSlug != "" {
		state, err := l.svcCtx.DB.GetContactSequenceStateWithDetails(l.ctx, db.GetContactSequenceStateWithDetailsParams{
			ContactID: sql.NullString{String: contact.ID, Valid: true},
			OrgID:     sql.NullString{String: orgID, Valid: true},
			Slug:      req.SequenceSlug,
		})
		if err != nil {
			l.Errorf("Enrollment not found: email=%s sequence=%s error: %v", req.Email, req.SequenceSlug, err)
			return nil, fmt.Errorf("enrollment not found")
		}

		enrollment := l.stateToEnrollmentInfo(state)
		return &types.ListSequenceEnrollmentsResponse{
			Enrollments: []types.SequenceEnrollmentInfo{enrollment},
		}, nil
	}

	states, err := l.svcCtx.DB.ListContactSequenceStatesWithDetails(l.ctx, db.ListContactSequenceStatesWithDetailsParams{
		ContactID: sql.NullString{String: contact.ID, Valid: true},
		OrgID:     sql.NullString{String: orgID, Valid: true},
	})
	if err != nil {
		l.Errorf("Failed to list enrollments: %v", err)
		return nil, fmt.Errorf("failed to retrieve enrollments")
	}

	enrollments := make([]types.SequenceEnrollmentInfo, 0, len(states))
	for _, state := range states {
		enrollment := l.listStateToEnrollmentInfo(state)
		enrollments = append(enrollments, enrollment)
	}

	return &types.ListSequenceEnrollmentsResponse{
		Enrollments: enrollments,
	}, nil
}

func (l *GetSequenceEnrollmentsLogic) stateToEnrollmentInfo(state db.GetContactSequenceStateWithDetailsRow) types.SequenceEnrollmentInfo {
	status := determineEnrollmentStatus(state.IsActive, state.CompletedAt, state.UnsubscribedAt, state.PausedAt)

	info := types.SequenceEnrollmentInfo{
		EnrollmentId:  state.ID,
		SequenceSlug:  state.SequenceSlug,
		SequenceName:  state.SequenceName,
		Status:        status,
		CurrentStep:   int(state.CurrentPosition.Int64),
		TotalSteps:    int(state.TotalSteps),
		EmailsSent:    int(state.EmailsSent),
		EmailsOpened:  int(state.EmailsOpened),
		EmailsClicked: int(state.EmailsClicked),
	}

	if state.StartedAt.Valid {
		info.EnrolledAt = state.StartedAt.String
	}
	if state.CompletedAt.Valid {
		info.CompletedAt = state.CompletedAt.String
	}
	if state.NextEmailAt != "" {
		info.NextEmailAt = state.NextEmailAt
	}

	return info
}

func (l *GetSequenceEnrollmentsLogic) listStateToEnrollmentInfo(state db.ListContactSequenceStatesWithDetailsRow) types.SequenceEnrollmentInfo {
	status := determineEnrollmentStatus(state.IsActive, state.CompletedAt, state.UnsubscribedAt, state.PausedAt)

	info := types.SequenceEnrollmentInfo{
		EnrollmentId:  state.ID,
		SequenceSlug:  state.SequenceSlug,
		SequenceName:  state.SequenceName,
		Status:        status,
		CurrentStep:   int(state.CurrentPosition.Int64),
		TotalSteps:    int(state.TotalSteps),
		EmailsSent:    int(state.EmailsSent),
		EmailsOpened:  int(state.EmailsOpened),
		EmailsClicked: int(state.EmailsClicked),
	}

	if state.StartedAt.Valid {
		info.EnrolledAt = state.StartedAt.String
	}
	if state.CompletedAt.Valid {
		info.CompletedAt = state.CompletedAt.String
	}
	if state.NextEmailAt != "" {
		info.NextEmailAt = state.NextEmailAt
	}

	return info
}

func determineEnrollmentStatus(isActive sql.NullInt64, completedAt, unsubscribedAt, pausedAt sql.NullString) string {
	if unsubscribedAt.Valid {
		return "cancelled"
	}

	if completedAt.Valid {
		return "completed"
	}

	if pausedAt.Valid {
		return "paused"
	}

	if isActive.Valid && isActive.Int64 != 1 {
		return "paused"
	}

	return "active"
}
