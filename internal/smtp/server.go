package smtp

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/outlet-sh/outlet/internal/config"
	"github.com/outlet-sh/outlet/internal/svc"

	"github.com/emersion/go-smtp"
	"github.com/zeromicro/go-zero/core/logx"
)

// Server wraps the go-smtp server with Outlet.sh integration
type Server struct {
	config  config.SMTPConfig
	svcCtx  *svc.ServiceContext
	server  *smtp.Server
	started bool
}

// NewServer creates a new SMTP ingress server
func NewServer(svcCtx *svc.ServiceContext, cfg config.SMTPConfig) *Server {
	return &Server{
		config: cfg,
		svcCtx: svcCtx,
	}
}

// Start begins listening for SMTP connections
func (s *Server) Start() error {
	if s.started {
		return fmt.Errorf("SMTP server already started")
	}

	backend := NewBackend(s.svcCtx)

	s.server = smtp.NewServer(backend)
	s.server.Addr = fmt.Sprintf(":%d", s.config.GetPort())
	s.server.Domain = s.config.Domain
	s.server.ReadTimeout = 60 * time.Second
	s.server.WriteTimeout = 60 * time.Second
	s.server.MaxMessageBytes = int64(s.config.MaxMessageBytes)
	s.server.MaxRecipients = s.config.MaxRecipients
	s.server.AllowInsecureAuth = s.config.IsAllowInsecureAuth()
	logx.Infof("SMTP: AllowInsecureAuth config value: %q, parsed: %v", s.config.AllowInsecureAuth, s.server.AllowInsecureAuth)

	// Configure TLS if certificates provided
	if s.config.TLSCert != "" && s.config.TLSKey != "" {
		cert, err := tls.LoadX509KeyPair(s.config.TLSCert, s.config.TLSKey)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}
		s.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		logx.Infof("SMTP: TLS configured with certificate from %s", s.config.TLSCert)
	} else {
		logx.Info("SMTP: Running without TLS (STARTTLS not available)")
	}

	s.started = true

	// Start server in goroutine
	go func() {
		logx.Infof("SMTP: Server starting on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil {
			logx.Errorf("SMTP: Server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the SMTP server
func (s *Server) Stop() error {
	if !s.started || s.server == nil {
		return nil
	}

	logx.Info("SMTP: Shutting down server...")
	if err := s.server.Close(); err != nil {
		return fmt.Errorf("failed to close SMTP server: %w", err)
	}

	s.started = false
	logx.Info("SMTP: Server stopped")
	return nil
}

// Addr returns the server's listen address
func (s *Server) Addr() string {
	if s.server != nil {
		return s.server.Addr
	}
	return ""
}
