package mcpctx

import (
	"context"
	"sync"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/svc"
)

// AuthMode indicates how the MCP session was authenticated.
type AuthMode int

const (
	// AuthModeAPIKey means authenticated via X-API-Key header (org-scoped).
	AuthModeAPIKey AuthMode = iota
	// AuthModeOAuth means authenticated via Bearer token (user-scoped, needs org selection).
	AuthModeOAuth
)

// OrgSelectionCallback is called when org selection changes.
// The handler uses this to persist org selection for session recovery.
// userID is included for DB persistence.
type OrgSelectionCallback func(userID, orgID string)

// ToolContext carries context for all MCP tools.
// It supports both API key auth (org-scoped) and OAuth auth (user-scoped with org selection).
type ToolContext struct {
	svc       *svc.ServiceContext
	requestID string
	userAgent string
	sessionID string

	// Auth mode
	authMode AuthMode

	// API key auth: org is set directly
	org *db.Organization

	// OAuth auth: user is set, org must be selected
	user        *db.User
	selectedOrg *db.Organization
	mu          sync.RWMutex // protects selectedOrg

	// Callback for persisting org selection (set by Handler)
	onOrgSelect OrgSelectionCallback
}

// NewToolContext creates a new org-scoped tool context (API key auth).
func NewToolContext(svc *svc.ServiceContext, org db.Organization, requestID, userAgent string) *ToolContext {
	return &ToolContext{
		svc:       svc,
		org:       &org,
		requestID: requestID,
		userAgent: userAgent,
		authMode:  AuthModeAPIKey,
	}
}

// NewUserToolContext creates a new user-scoped tool context (OAuth auth).
// The user must select an org using org_select before using other tools.
func NewUserToolContext(svc *svc.ServiceContext, user db.User, requestID, userAgent, sessionID string) *ToolContext {
	return &ToolContext{
		svc:       svc,
		user:      &user,
		requestID: requestID,
		userAgent: userAgent,
		sessionID: sessionID,
		authMode:  AuthModeOAuth,
	}
}

// SetOrgSelectionCallback sets the callback for persisting org selection.
func (t *ToolContext) SetOrgSelectionCallback(cb OrgSelectionCallback) {
	t.onOrgSelect = cb
}

// SessionID returns the MCP session ID.
func (t *ToolContext) SessionID() string {
	return t.sessionID
}

// OrgID returns the organization ID for scoping queries.
// Returns empty string if no org is available (OAuth without selection).
func (t *ToolContext) OrgID() string {
	org := t.currentOrg()
	if org == nil {
		return ""
	}
	return org.ID
}

// Org returns the full organization record.
// Returns empty Organization if no org is available.
func (t *ToolContext) Org() db.Organization {
	org := t.currentOrg()
	if org == nil {
		return db.Organization{}
	}
	return *org
}

// HasOrg returns true if an organization is available for operations.
func (t *ToolContext) HasOrg() bool {
	return t.currentOrg() != nil
}

// currentOrg returns the current org based on auth mode.
func (t *ToolContext) currentOrg() *db.Organization {
	if t.authMode == AuthModeAPIKey {
		return t.org
	}
	// OAuth mode: use selected org
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.selectedOrg
}

// AuthMode returns the authentication mode.
func (t *ToolContext) AuthMode() AuthMode {
	return t.authMode
}

// User returns the authenticated user (OAuth mode only).
func (t *ToolContext) User() *db.User {
	return t.user
}

// UserID returns the authenticated user's ID (OAuth mode only).
// Returns empty string if not in OAuth mode.
func (t *ToolContext) UserID() string {
	if t.user == nil {
		return ""
	}
	return t.user.ID
}

// IsSuperAdmin returns true if the user has the super_admin role.
func (t *ToolContext) IsSuperAdmin() bool {
	if t.user == nil {
		return false
	}
	return t.user.Role == "super_admin"
}

// SelectOrg sets the current organization for OAuth sessions.
// Also calls the persistence callback if set.
func (t *ToolContext) SelectOrg(org db.Organization) {
	t.mu.Lock()
	t.selectedOrg = &org
	cb := t.onOrgSelect
	userID := t.UserID()
	t.mu.Unlock()

	// Call callback outside the lock to persist the selection
	if cb != nil {
		cb(userID, org.ID)
	}
}

// RestoreOrg sets the current organization without triggering the callback.
// Used when restoring org selection from persistent storage.
func (t *ToolContext) RestoreOrg(org db.Organization) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.selectedOrg = &org
}

// ClearOrg clears the selected organization (OAuth mode).
func (t *ToolContext) ClearOrg() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.selectedOrg = nil
}

// DB returns the database store for queries.
func (t *ToolContext) DB() *db.Store {
	return t.svc.DB
}

// Svc returns the full service context for advanced operations.
func (t *ToolContext) Svc() *svc.ServiceContext {
	return t.svc
}

// RequestID returns the request ID for tracing.
func (t *ToolContext) RequestID() string {
	return t.requestID
}

// UserAgent returns the client's user agent string.
func (t *ToolContext) UserAgent() string {
	return t.userAgent
}

// ToolError represents a structured error for MCP tool responses.
type ToolError struct {
	Code    string `json:"code"`    // "not_found", "validation", "conflict", "rate_limit"
	Message string `json:"message"` // Human-readable description
	Field   string `json:"field"`   // For validation errors
}

func (e *ToolError) Error() string {
	if e.Field != "" {
		return e.Code + ": " + e.Message + " (field: " + e.Field + ")"
	}
	return e.Code + ": " + e.Message
}

// NewValidationError creates a validation error for a specific field.
func NewValidationError(message, field string) *ToolError {
	return &ToolError{Code: "validation", Message: message, Field: field}
}

// NewNotFoundError creates a not found error.
func NewNotFoundError(message string) *ToolError {
	return &ToolError{Code: "not_found", Message: message}
}

// NewConflictError creates a conflict error (duplicate, already exists).
func NewConflictError(message string) *ToolError {
	return &ToolError{Code: "conflict", Message: message}
}

// ErrNoOrgSelected is returned when an org-scoped operation is attempted without an org selected.
var ErrNoOrgSelected = &ToolError{
	Code:    "no_org_selected",
	Message: "No organization selected. Use org_list to see available organizations and org_select to choose one.",
}

// RequireOrg returns an error if no org is available.
// Use this at the start of tools that need org context.
func (t *ToolContext) RequireOrg() error {
	if !t.HasOrg() {
		return ErrNoOrgSelected
	}
	return nil
}

// toolContextKey is used to store ToolContext in context.Context
type toolContextKey struct{}

// WithToolContext adds ToolContext to a context.
func WithToolContext(ctx context.Context, tc *ToolContext) context.Context {
	return context.WithValue(ctx, toolContextKey{}, tc)
}

// ToolContextFromContext retrieves ToolContext from a context.
func ToolContextFromContext(ctx context.Context) *ToolContext {
	if tc, ok := ctx.Value(toolContextKey{}).(*ToolContext); ok {
		return tc
	}
	return nil
}
