package tracking

import (
	"context"
	"database/sql"
	"errors"

	"github.com/outlet-sh/outlet/internal/db"
)

var (
	ErrInvalidToken = errors.New("invalid or missing token")
	ErrNotFound     = errors.New("record not found")
)

// Service handles email tracking operations
type Service struct {
	db *db.Queries
}

// New creates a new tracking service
func New(db *db.Queries) *Service {
	return &Service{db: db}
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

	return s.db.RecordEmailOpen(ctx, emailRecord.ID)
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

	return s.db.RecordEmailClick(ctx, emailRecord.ID)
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
