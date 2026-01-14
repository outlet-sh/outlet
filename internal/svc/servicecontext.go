package svc

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	"outlet/internal/config"
	"outlet/internal/db"
	"outlet/internal/db/migrations"
	"outlet/internal/events"
	"outlet/internal/middleware"
	"outlet/internal/services/crypto"
	"outlet/internal/services/email"
	"outlet/internal/services/tracking"
	"outlet/internal/services/webhook"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/rest"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type ServiceContext struct {
	Config            config.Config
	DB                *db.Store
	Auth              rest.Middleware
	APIKeyAuth        rest.Middleware
	APIKeyMiddleware  *middleware.APIKeyMiddleware
	AuthRateLimit     rest.Middleware
	EmailService      *email.Service
	CryptoService     *crypto.Service
	Tracking          *tracking.Service
	Events            *events.Subject
	WebhookDispatcher *webhook.Dispatcher
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Initialize SQLite database connection
	dbPath := c.Database.Path
	if dbPath == "" {
		dbPath = "./data/outlet.db"
	}

	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Open SQLite database with WAL mode for better concurrency
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable WAL mode and other SQLite optimizations
	if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
		log.Printf("Warning: Failed to enable WAL mode: %v", err)
	}
	if _, err := conn.Exec("PRAGMA foreign_keys=ON"); err != nil {
		log.Printf("Warning: Failed to enable foreign keys: %v", err)
	}
	if _, err := conn.Exec("PRAGMA busy_timeout=5000"); err != nil {
		log.Printf("Warning: Failed to set busy timeout: %v", err)
	}

	// Test database connection
	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("SQLite database connected: %s", dbPath)

	// Run database migrations (non-destructive - only applies pending migrations)
	log.Printf("Running database migrations...")
	if err := migrations.Run(conn); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Printf("Database migrations completed")

	// Create database store
	store := db.NewStore(conn)

	// Seed super admin user if configured and not already exists
	log.Printf("DEBUG: Admin config - email=%q, password_len=%d", c.Admin.Email, len(c.Admin.Password))
	if c.Admin.Email != "" && c.Admin.Password != "" {
		log.Printf("Admin seeding configured for: %s", c.Admin.Email)
		seedSuperAdmin(context.Background(), store, c.Admin.Email, c.Admin.Password)
	} else {
		log.Printf("Admin seeding skipped - ADMIN_EMAIL or ADMIN_PASSWORD not set (email=%q)", c.Admin.Email)
	}

	// Initialize Email service (SMTP config loaded from platform_settings in database)
	emailService := email.NewService(store)
	if c.App.BaseURL != "" {
		emailService.SetBaseURL(c.App.BaseURL)
	}

	// Initialize Crypto service for org credential encryption
	var cryptoService *crypto.Service
	if c.Encryption.Key != "" {
		var err error
		cryptoService, err = crypto.NewService(c.Encryption.Key)
		if err != nil {
			log.Printf("Warning: Failed to initialize crypto service: %v", err)
		} else {
			log.Printf("Crypto service initialized successfully")
		}
	} else {
		log.Printf("Warning: ENCRYPTION_KEY not set - credential storage (Stripe, AI keys) will be disabled")
	}

	// Initialize Tracking service
	trackingService := tracking.New(store.Queries)

	// Initialize rate limit middleware for auth endpoints
	authRateLimiter := middleware.NewRateLimitMiddleware(middleware.DefaultAuthRateLimitConfig())

	// Initialize API key middleware with cache
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(store)

	// Initialize Event bus for webhooks and automation
	eventSubject := events.NewSubject(
		events.WithBufferSize(1024),
		events.WithReplay(100),
	)
	log.Printf("Event bus initialized")

	// Initialize and start Webhook Dispatcher for outbound webhook delivery
	webhookDispatcher := webhook.NewDispatcher(store.Queries, eventSubject)
	if err := webhookDispatcher.Start(context.Background()); err != nil {
		log.Printf("Warning: Failed to start webhook dispatcher: %v", err)
	} else {
		log.Printf("Webhook dispatcher started")
	}

	return &ServiceContext{
		Config:            c,
		DB:                store,
		Auth:              middleware.NewAuthMiddleware().Handle,
		APIKeyAuth:        apiKeyMiddleware.Handle,
		APIKeyMiddleware:  apiKeyMiddleware,
		AuthRateLimit:     authRateLimiter.Handle,
		EmailService:      emailService,
		CryptoService:     cryptoService,
		Tracking:          trackingService,
		Events:            eventSubject,
		WebhookDispatcher: webhookDispatcher,
	}
}

// seedSuperAdmin creates the super admin user if it doesn't already exist
func seedSuperAdmin(ctx context.Context, store *db.Store, email, password string) {
	email = strings.ToLower(strings.TrimSpace(email))

	// Check if user already exists
	exists, err := store.CheckEmailExists(ctx, email)
	if err != nil {
		log.Printf("Warning: Failed to check if super admin exists: %v", err)
		return
	}

	if exists > 0 {
		log.Printf("Super admin user already exists: %s", email)
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Failed to hash super admin password: %v", err)
		return
	}

	// Create super admin user
	_, err = store.CreateUser(ctx, db.CreateUserParams{
		ID:            uuid.NewString(),
		Email:         email,
		PasswordHash:  string(passwordHash),
		Name:          "Super Admin",
		Role:          "super_admin",
		Status:        "active",
		EmailVerified: 1,
	})
	if err != nil {
		log.Printf("Warning: Failed to create super admin user: %v", err)
		return
	}

	log.Printf("Super admin user created successfully: %s", email)
}
