package tracking

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/events"
)

var (
	ErrInvalidToken = errors.New("invalid or missing token")
	ErrNotFound     = errors.New("record not found")
)

// Service handles email tracking operations
type Service struct {
	db  *db.Queries
	bus *events.Subject // optional; opens and clicks are published here
}

// New creates a new tracking service. bus may be nil.
func New(db *db.Queries, bus *events.Subject) *Service {
	return &Service{db: db, bus: bus}
}

// publish emits an open or click for the email on the event bus.
func (s *Service) publish(topic, status string, e db.GetEmailByTrackingTokenRow) {
	if s.bus == nil {
		return
	}
	_ = events.Emit(s.bus, topic, events.EmailEvent{
		EmailID:   e.ID,
		ContactID: e.ContactID.String,
		Email:     e.Email,
		Status:    status,
		Timestamp: time.Now(),
	})
}

// RecordOpen records an email open event by tracking token
func (s *Service) RecordOpen(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	emailRecord, err := s.db.GetEmailByTrackingToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		return ErrNotFound
	}

	if err := s.db.RecordEmailOpen(ctx, emailRecord.ID); err != nil {
		return err
	}
	s.publish(events.TopicEmailOpened, "opened", emailRecord)
	return nil
}

// RecordClick records an email click event by tracking token
func (s *Service) RecordClick(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	emailRecord, err := s.db.GetEmailByTrackingToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		return ErrNotFound
	}

	if err := s.db.RecordEmailClick(ctx, emailRecord.ID); err != nil {
		return err
	}
	s.publish(events.TopicEmailClicked, "clicked", emailRecord)
	return nil
}

// Unsubscribe unsubscribes a contact by tracking token and cancels pending emails
func (s *Service) Unsubscribe(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	contact, err := s.db.GetContactByTrackingToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		return ErrNotFound
	}

	if err := s.db.UnsubscribeContact(ctx, contact.ID); err != nil {
		return err
	}

	return s.db.CancelEmailsForContact(ctx, sql.NullString{String: contact.ID, Valid: true})
}
