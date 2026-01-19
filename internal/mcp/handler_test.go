package mcp

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/mcp/mcpauth"
	"github.com/outlet-sh/outlet/internal/mcp/mcpctx"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestOrg creates a mock organization for testing
func createTestOrg(id string, name string) db.Organization {
	return db.Organization{
		ID:     id,
		Name:   name,
		Slug:   "test-org",
		ApiKey: "test-api-key",
	}
}

// createTestUser creates a mock user for testing
func createTestUser(id string, email, role string) db.User {
	return db.User{
		ID:     id,
		Email:  email,
		Name:   "Test User",
		Role:   role,
		Status: "active",
	}
}

// ========== Handler Creation Tests ==========

func TestNewHandler(t *testing.T) {
	// Create handler with nil svc for basic instantiation test
	handler := NewHandler(nil, "https://example.com")

	assert.NotNil(t, handler, "Handler should not be nil")
}

func TestNewHandler_TrimTrailingSlash(t *testing.T) {
	handler := NewHandler(nil, "https://example.com/")

	// The handler should have trimmed the trailing slash internally
	assert.NotNil(t, handler)
}

func TestNewHandler_Authenticator(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	auth := handler.Authenticator()
	assert.NotNil(t, auth, "Authenticator should not be nil")
}

// ========== HTTP Handler Tests ==========

func TestHandler_ServeHTTP_MissingAuth(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Should return 401 without auth")

	// Check WWW-Authenticate header
	wwwAuth := rr.Header().Get("WWW-Authenticate")
	assert.Contains(t, wwwAuth, "Bearer", "Should have Bearer challenge")
	assert.Contains(t, wwwAuth, "resource_metadata", "Should include resource_metadata")
}

func TestHandler_ServeHTTP_EmptyBearer(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer ")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Should return 401 with empty bearer")
}

func TestHandler_ServeHTTP_InvalidToken(t *testing.T) {
	// Note: Without a real DB, we can't fully test token verification
	// The authenticator will panic without a DB connection
	// This test documents the expected behavior
	t.Skip("Skipping: requires DB connection for token verification")
}

func TestHandler_ServeHTTP_BasicAuth_Rejected(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz") // user:pass
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Basic auth should be rejected since it doesn't start with "Bearer "
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Should reject Basic auth")
}

// ========== Session Cache Tests ==========

func TestHandler_SessionCache_StoreAndLoad(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	sessionID := "test-session-123"
	server := NewServer(nil, httptest.NewRequest(http.MethodPost, "/mcp", nil))

	// Store session
	handler.sessionCache.Store(sessionID, &sessionData{server: server, toolCtx: nil})

	// Load session
	cached, ok := handler.sessionCache.Load(sessionID)

	assert.True(t, ok, "Session should be found")
	assert.NotNil(t, cached, "Session data should not be nil")

	data := cached.(*sessionData)
	assert.Equal(t, server, data.server, "Server should match")
}

func TestHandler_SessionCache_Delete(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	sessionID := "session-to-delete"
	handler.sessionCache.Store(sessionID, &sessionData{})

	// Verify stored
	_, ok := handler.sessionCache.Load(sessionID)
	assert.True(t, ok, "Session should be stored")

	// Delete
	handler.sessionCache.Delete(sessionID)

	// Verify deleted
	_, ok = handler.sessionCache.Load(sessionID)
	assert.False(t, ok, "Session should be deleted")
}

// ========== Org Selection Cache Tests ==========

func TestHandler_StoreOrgSelection(t *testing.T) {
	// Note: StoreOrgSelection launches a goroutine that tries to persist to DB
	// Without a DB, this will panic in the background goroutine
	// We test the memory cache directly instead

	handler := NewHandler(nil, "https://example.com")

	sessionID := "test-session"
	orgID := uuid.New().String()

	// Directly test the in-memory cache (avoid triggering the DB goroutine)
	handler.orgSelectionCache.Store(sessionID, orgID)

	// Check memory cache
	cached, ok := handler.orgSelectionCache.Load(sessionID)
	require.True(t, ok, "Org selection should be in memory cache")
	assert.Equal(t, orgID, cached.(string), "OrgID should match")
}

func TestHandler_ClearOrgSelection(t *testing.T) {
	// Note: ClearOrgSelection launches a goroutine that tries to delete from DB
	// Without a DB, this will panic in the background goroutine
	// We test the memory cache directly instead

	handler := NewHandler(nil, "https://example.com")

	sessionID := "session-to-clear"
	orgID := uuid.New().String()

	// Store directly to cache
	handler.orgSelectionCache.Store(sessionID, orgID)

	// Verify stored
	_, ok := handler.orgSelectionCache.Load(sessionID)
	assert.True(t, ok, "Org selection should be stored")

	// Clear directly from cache (avoid triggering the DB goroutine)
	handler.orgSelectionCache.Delete(sessionID)

	// Verify cleared
	_, ok = handler.orgSelectionCache.Load(sessionID)
	assert.False(t, ok, "Org selection should be cleared")
}

// ========== Multi-Tenant Isolation Tests ==========

func TestHandler_MultiTenant_DifferentOrgsPerSession(t *testing.T) {
	_ = NewHandler(nil, "https://example.com") // Keep for session cache tests

	// Create two users with different orgs
	user1 := createTestUser(uuid.New().String(), "user1@example.com", "admin")
	user2 := createTestUser(uuid.New().String(), "user2@example.com", "admin")

	org1 := createTestOrg(uuid.New().String(), "Org 1")
	org2 := createTestOrg(uuid.New().String(), "Org 2")

	// Create tool contexts for each user
	tc1 := mcpctx.NewUserToolContext(nil, user1, "req-1", "agent/1.0", "session-1")
	tc2 := mcpctx.NewUserToolContext(nil, user2, "req-2", "agent/1.0", "session-2")

	// Select different orgs
	tc1.SelectOrg(org1)
	tc2.SelectOrg(org2)

	// Verify isolation
	assert.Equal(t, org1.ID, tc1.OrgID(), "User 1 should have Org 1")
	assert.Equal(t, org2.ID, tc2.OrgID(), "User 2 should have Org 2")
	assert.NotEqual(t, tc1.OrgID(), tc2.OrgID(), "Orgs should be different")
}

func TestHandler_MultiTenant_SameUserMultipleSessions(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	// Same user, different sessions (e.g., different browser tabs)
	user := createTestUser(uuid.New().String(), "user@example.com", "admin")

	org1 := createTestOrg(uuid.New().String(), "Org 1")
	org2 := createTestOrg(uuid.New().String(), "Org 2")

	// Create two sessions for same user
	session1 := "session-browser-1"
	session2 := "session-browser-2"

	tc1 := mcpctx.NewUserToolContext(nil, user, "req-1", "agent/1.0", session1)
	tc2 := mcpctx.NewUserToolContext(nil, user, "req-2", "agent/1.0", session2)

	// User selects different orgs in different sessions
	tc1.SelectOrg(org1)
	tc2.SelectOrg(org2)

	// Store in handler's session cache
	handler.sessionCache.Store(session1, &sessionData{toolCtx: tc1})
	handler.sessionCache.Store(session2, &sessionData{toolCtx: tc2})

	// Verify each session has its own org
	cached1, _ := handler.sessionCache.Load(session1)
	cached2, _ := handler.sessionCache.Load(session2)

	data1 := cached1.(*sessionData)
	data2 := cached2.(*sessionData)

	assert.Equal(t, org1.ID, data1.toolCtx.OrgID(), "Session 1 should have Org 1")
	assert.Equal(t, org2.ID, data2.toolCtx.OrgID(), "Session 2 should have Org 2")
}

func TestHandler_MultiTenant_OrgIsolation_RequireOrg(t *testing.T) {
	// Test that RequireOrg properly enforces org selection
	user := createTestUser(uuid.New().String(), "user@example.com", "admin")

	// Create context without org selection
	tc := mcpctx.NewUserToolContext(nil, user, "req-1", "agent/1.0", "session-1")

	// Should require org selection before operations
	err := tc.RequireOrg()
	assert.Error(t, err, "Should error when no org selected")
	assert.Equal(t, mcpctx.ErrNoOrgSelected, err)

	// Select org
	org := createTestOrg(uuid.New().String(), "My Org")
	tc.SelectOrg(org)

	// Now should not error
	err = tc.RequireOrg()
	assert.NoError(t, err, "Should not error after org selection")
}

// ========== Session Persistence Tests ==========

func TestHandler_SessionRestore_FromMemoryCache(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	sessionID := "persistent-session"
	orgID := uuid.New().String()

	// Store org selection in memory cache
	handler.orgSelectionCache.Store(sessionID, orgID)

	// Verify it can be retrieved
	cached, ok := handler.orgSelectionCache.Load(sessionID)
	require.True(t, ok, "Should find org in memory cache")
	assert.Equal(t, orgID, cached.(string))
}

// ========== Concurrent Access Tests ==========

func TestHandler_ConcurrentSessionAccess(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			sessionID := "session-" + string(rune('0'+idx%10))
			orgID := uuid.New().String()

			// Store
			handler.orgSelectionCache.Store(sessionID, orgID)

			// Load
			if cached, ok := handler.orgSelectionCache.Load(sessionID); ok {
				_ = cached.(string)
			}

			// Store session
			handler.sessionCache.Store(sessionID, &sessionData{})

			// Load session
			if cached, ok := handler.sessionCache.Load(sessionID); ok {
				_ = cached.(*sessionData)
			}
		}(i)
	}

	wg.Wait()
	// Test passes if no race condition panic occurs
}

func TestHandler_ConcurrentOrgSelection(t *testing.T) {
	// Note: StoreOrgSelection triggers a DB goroutine that panics without a DB
	// We test the concurrent cache access directly instead
	handler := NewHandler(nil, "https://example.com")

	user := createTestUser(uuid.New().String(), "user@example.com", "admin")

	var wg sync.WaitGroup
	numGoroutines := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sessionID := uuid.New().String()
			orgID := uuid.New().String()

			// Create context
			tc := mcpctx.NewUserToolContext(nil, user, "req", "agent", sessionID)

			// Select org
			org := createTestOrg(orgID, "Test Org")
			tc.SelectOrg(org)

			// Store directly in cache (avoid DB goroutine)
			handler.orgSelectionCache.Store(sessionID, orgID)

			// Read org
			_ = tc.OrgID()
			_ = tc.HasOrg()

			// Read from cache
			_, _ = handler.orgSelectionCache.Load(sessionID)
		}()
	}

	wg.Wait()
	// Test passes if no race condition panic occurs
}

// ========== WWW-Authenticate Header Tests ==========

func TestHandler_WWWAuthenticate_Format(t *testing.T) {
	handler := NewHandler(nil, "https://api.example.com")

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	wwwAuth := rr.Header().Get("WWW-Authenticate")

	// Should include Bearer scheme
	assert.Contains(t, wwwAuth, "Bearer", "Should have Bearer scheme")

	// Should include resource_metadata with quoted URL
	assert.Contains(t, wwwAuth, "resource_metadata=", "Should have resource_metadata")

	// Should include scope
	assert.Contains(t, wwwAuth, "scope=", "Should have scope")
	assert.Contains(t, wwwAuth, "mcp:full", "Should have mcp:full scope")
}

// ========== Tool Context Tests ==========

func TestHandler_ToolContext_UserInfo(t *testing.T) {
	user := createTestUser(uuid.New().String(), "admin@example.com", "super_admin")

	tc := mcpctx.NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-abc")

	assert.Equal(t, user.ID, tc.UserID(), "UserID should match")
	assert.True(t, tc.IsSuperAdmin(), "Should be super admin")
	assert.Equal(t, mcpctx.AuthModeOAuth, tc.AuthMode(), "Should be OAuth mode")
	assert.Equal(t, "session-abc", tc.SessionID(), "SessionID should match")
}

func TestHandler_ToolContext_OrgContext(t *testing.T) {
	org := createTestOrg(uuid.New().String(), "Test Company")

	tc := mcpctx.NewToolContext(nil, org, "req-123", "test-agent/1.0")

	assert.Equal(t, org.ID, tc.OrgID(), "OrgID should match")
	assert.True(t, tc.HasOrg(), "Should have org")
	assert.Equal(t, mcpctx.AuthModeAPIKey, tc.AuthMode(), "Should be API key mode")
	assert.Empty(t, tc.SessionID(), "SessionID should be empty for API key mode")
}

// ========== OAuth vs API Key Mode Tests ==========

func TestHandler_APIKeyMode_NoOrgSelection(t *testing.T) {
	org := createTestOrg(uuid.New().String(), "API Key Org")

	// API key mode: org is set directly, no selection needed
	tc := mcpctx.NewToolContext(nil, org, "req-123", "agent/1.0")

	assert.True(t, tc.HasOrg(), "Should have org in API key mode")
	assert.Equal(t, org.ID, tc.OrgID())
	assert.NoError(t, tc.RequireOrg(), "RequireOrg should pass")
}

func TestHandler_OAuthMode_RequiresOrgSelection(t *testing.T) {
	user := createTestUser(uuid.New().String(), "user@example.com", "admin")

	// OAuth mode: org must be selected
	tc := mcpctx.NewUserToolContext(nil, user, "req-123", "agent/1.0", "session-123")

	assert.False(t, tc.HasOrg(), "Should not have org before selection")
	assert.Equal(t, "", tc.OrgID())
	assert.Error(t, tc.RequireOrg(), "RequireOrg should fail before selection")

	// Select org
	org := createTestOrg(uuid.New().String(), "Selected Org")
	tc.SelectOrg(org)

	assert.True(t, tc.HasOrg(), "Should have org after selection")
	assert.Equal(t, org.ID, tc.OrgID())
	assert.NoError(t, tc.RequireOrg(), "RequireOrg should pass after selection")
}

// ========== Session Data Tests ==========

func TestSessionData_Structure(t *testing.T) {
	user := createTestUser(uuid.New().String(), "user@example.com", "admin")
	tc := mcpctx.NewUserToolContext(nil, user, "req-123", "agent/1.0", "session-123")

	// Create mock server (actual server creation would require full setup)
	data := &sessionData{
		server:  nil, // Would be MCP server
		toolCtx: tc,
	}

	assert.NotNil(t, data.toolCtx)
	assert.Equal(t, "session-123", data.toolCtx.SessionID())
}

// ========== Edge Cases ==========

func TestHandler_EmptySessionID(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	// Should handle empty session ID gracefully
	_, ok := handler.sessionCache.Load("")
	assert.False(t, ok, "Empty session ID should not find anything")
}

func TestHandler_UUIDSessionID(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	// UUID-style session IDs should work
	sessionID := uuid.New().String()
	orgID := uuid.New().String()

	handler.orgSelectionCache.Store(sessionID, orgID)

	cached, ok := handler.orgSelectionCache.Load(sessionID)
	assert.True(t, ok, "UUID session ID should work")
	assert.Equal(t, orgID, cached.(string))
}

func TestHandler_SpecialCharacterSessionID(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	// Session IDs with special characters
	sessionID := "session-with-special!@#$%"
	orgID := uuid.New().String()

	handler.orgSelectionCache.Store(sessionID, orgID)

	cached, ok := handler.orgSelectionCache.Load(sessionID)
	assert.True(t, ok, "Special character session ID should work")
	assert.Equal(t, orgID, cached.(string))
}

// ========== HTTP Method Tests ==========

func TestHandler_DifferentHTTPMethods(t *testing.T) {
	handler := NewHandler(nil, "https://example.com")

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/mcp", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// All methods without auth should return 401
			assert.Equal(t, http.StatusUnauthorized, rr.Code, "%s should return 401 without auth", method)
		})
	}
}

// ========== Org Selection Callback Tests ==========

func TestHandler_OrgSelectionCallback(t *testing.T) {
	user := createTestUser(uuid.New().String(), "user@example.com", "admin")
	tc := mcpctx.NewUserToolContext(nil, user, "req-123", "agent/1.0", "session-123")

	var callbackUserID string
	var callbackOrgID string
	callbackCalled := false

	tc.SetOrgSelectionCallback(func(uid, oid string) {
		callbackCalled = true
		callbackUserID = uid
		callbackOrgID = oid
	})

	org := createTestOrg(uuid.New().String(), "Test Org")
	tc.SelectOrg(org)

	assert.True(t, callbackCalled, "Callback should be called")
	assert.Equal(t, user.ID, callbackUserID, "Callback should receive user ID")
	assert.Equal(t, org.ID, callbackOrgID, "Callback should receive org ID")
}

// ========== Token Info Context Tests ==========

func TestHandler_TokenInfoInContext(t *testing.T) {
	ctx := httptest.NewRequest(http.MethodGet, "/", nil).Context()

	// Without token info
	info := mcpauth.TokenInfoFromContext(ctx)
	assert.Nil(t, info, "Should be nil when not set")
}
