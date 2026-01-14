package config

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Platform struct {
		Name string // Platform name for branding (e.g., "Outlet.sh")
	}
	App struct {
		BaseURL        string // Base URL for the app (e.g., https://outlet.sh)
		Domain         string // Production domain for HTTPS (e.g., outlet.sh)
		ProductionMode bool   // Enable production mode (HTTPS, Let's Encrypt)
	}
	Database struct {
		Path string `json:",default=./data/outlet.db"` // SQLite database file path
	}
	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
	}
	Admin struct {
		Email    string // Super admin email for initial seeding
		Password string // Super admin password for initial seeding
	}
	Email struct {
		// High-volume dispatcher settings (per-org SMTP config is in database)
		WorkerCount int     `json:",default=10"`  // Concurrent email workers
		RateLimit   float64 `json:",default=14"`  // Emails per second (SES limit)
		RateBurst   int     `json:",default=50"`  // Max burst size
		BatchSize   int     `json:",default=100"` // Emails per batch fetch
	}
	SMTP SMTPConfig
	Encryption struct {
		Key string // 32-byte hex-encoded key for AES-256 encryption
	}
}

// Validate checks that all required configuration values are set
// Returns a list of missing/invalid configurations
func (c *Config) Validate() []string {
	var errors []string

	// Required for all environments
	if c.App.BaseURL == "" {
		errors = append(errors, "APP_BASE_URL is required")
	}

	// Validate URLs have proper format (not empty after env expansion)
	if c.App.BaseURL != "" && !strings.HasPrefix(c.App.BaseURL, "http") {
		errors = append(errors, fmt.Sprintf("APP_BASE_URL must start with http:// or https://, got: %s", c.App.BaseURL))
	}

	return errors
}

// MustValidate panics if configuration is invalid
func (c *Config) MustValidate() {
	errors := c.Validate()
	if len(errors) > 0 {
		panic(fmt.Sprintf("Configuration errors:\n  - %s", strings.Join(errors, "\n  - ")))
	}
}

// ValidateAndWarn logs warnings for missing optional config but doesn't fail
func (c *Config) ValidateAndWarn() []string {
	var warnings []string

	// Optional but recommended
	if c.Auth.AccessSecret == "" {
		warnings = append(warnings, "JWT_SECRET not set - authentication disabled")
	}

	return warnings
}

// SMTPConfig configures the SMTP ingress server
type SMTPConfig struct {
	Enabled           bool   `json:",optional"`           // Enable SMTP server (default: false)
	Port              int    `json:",default=587"`        // SMTP port (default: 587)
	Domain            string `json:",optional"`           // Server domain (uses App.Domain if empty)
	TLSCert           string `json:",optional"`           // Path to TLS certificate
	TLSKey            string `json:",optional"`           // Path to TLS private key
	MaxMessageBytes   int    `json:",default=26214400"`   // Max message size (default: 25MB)
	MaxRecipients     int    `json:",default=100"`        // Max recipients per message
	AllowInsecureAuth bool   `json:",optional"`           // Allow auth without TLS (default: false)
}
