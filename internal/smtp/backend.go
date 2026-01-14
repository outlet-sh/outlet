package smtp

import (
	"context"
	"errors"
	"io"

	"outlet/internal/db"
	"outlet/internal/svc"

	"github.com/emersion/go-smtp"
	"github.com/zeromicro/go-zero/core/logx"
)

// Backend implements smtp.Backend for authenticating SMTP connections
type Backend struct {
	svcCtx *svc.ServiceContext
}

// NewBackend creates a new SMTP backend
func NewBackend(svcCtx *svc.ServiceContext) *Backend {
	return &Backend{svcCtx: svcCtx}
}

// NewSession is called for anonymous connections (not allowed)
func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		svcCtx: b.svcCtx,
		conn:   c,
	}, nil
}

// Session implements smtp.Session for handling individual SMTP connections
type Session struct {
	svcCtx     *svc.ServiceContext
	conn       *smtp.Conn
	org        db.Organization
	authed     bool
	from       string
	recipients []string
}

// AuthPlain handles PLAIN authentication
// Username can be "api" or org slug (flexible), password is the API key
func (s *Session) AuthPlain(username, password string) error {
	if password == "" {
		logx.Infof("SMTP: Auth failed - empty password from %s", s.conn.Conn().RemoteAddr())
		return errors.New("invalid credentials")
	}

	// Use the API key middleware's getOrg method for validation
	org, err := s.svcCtx.DB.GetOrganizationByAPIKey(context.Background(), password)
	if err != nil {
		logx.Infof("SMTP: Auth failed - invalid API key from %s (user=%s)", s.conn.Conn().RemoteAddr(), username)
		return errors.New("invalid credentials")
	}

	s.org = org
	s.authed = true
	logx.Infof("SMTP: Authenticated as org %s (%s) from %s", org.Slug, org.ID, s.conn.Conn().RemoteAddr())
	return nil
}

// Mail is called for MAIL FROM command
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.authed {
		return errors.New("authentication required")
	}
	s.from = from
	logx.Debugf("SMTP: MAIL FROM: %s (org=%s)", from, s.org.Slug)
	return nil
}

// Rcpt is called for RCPT TO command
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.authed {
		return errors.New("authentication required")
	}
	s.recipients = append(s.recipients, to)
	logx.Debugf("SMTP: RCPT TO: %s (org=%s)", to, s.org.Slug)
	return nil
}

// Data is called when email data is received
func (s *Session) Data(r io.Reader) error {
	if !s.authed {
		return errors.New("authentication required")
	}
	if len(s.recipients) == 0 {
		return errors.New("no recipients specified")
	}

	// Process the email
	processor := NewEmailProcessor(s.svcCtx, s.org, s.from, s.recipients)
	messageID, err := processor.Process(r)
	if err != nil {
		logx.Errorf("SMTP: Failed to process email from %s to %v: %v", s.from, s.recipients, err)
		return err
	}

	logx.Infof("SMTP: Message accepted from=%s to=%v msgId=%s org=%s", s.from, s.recipients, messageID, s.org.Slug)
	return nil
}

// Reset clears session state between messages
func (s *Session) Reset() {
	s.from = ""
	s.recipients = nil
}

// Logout is called when the connection closes
func (s *Session) Logout() error {
	logx.Debugf("SMTP: Connection closed from %s", s.conn.Conn().RemoteAddr())
	return nil
}
