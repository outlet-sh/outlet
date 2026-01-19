package mcpctx

import (
	"context"
	"sync"
	"testing"

	"github.com/outlet-sh/outlet/internal/db"

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

// ========== NewToolContext Tests ==========

func TestNewToolContext_APIKeyMode(t *testing.T) {
	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Test Organization")

	tc := NewToolContext(nil, org, "req-123", "test-agent/1.0")

	assert.NotNil(t, tc, "ToolContext should not be nil")
	assert.Equal(t, AuthModeAPIKey, tc.AuthMode(), "AuthMode should be APIKey")
	assert.Equal(t, orgID, tc.OrgID(), "OrgID should match")
	assert.Equal(t, "req-123", tc.RequestID(), "RequestID should match")
	assert.Equal(t, "test-agent/1.0", tc.UserAgent(), "UserAgent should match")
	assert.True(t, tc.HasOrg(), "HasOrg should be true for API key mode")
}

func TestNewUserToolContext_OAuthMode(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")

	tc := NewUserToolContext(nil, user, "req-456", "test-agent/1.0", "session-123")

	assert.NotNil(t, tc, "ToolContext should not be nil")
	assert.Equal(t, AuthModeOAuth, tc.AuthMode(), "AuthMode should be OAuth")
	assert.Equal(t, userID, tc.UserID(), "UserID should match")
	assert.Equal(t, "", tc.OrgID(), "OrgID should be empty before selection")
	assert.False(t, tc.HasOrg(), "HasOrg should be false before org selection")
	assert.Equal(t, "session-123", tc.SessionID(), "SessionID should match")
}

// ========== Org Selection Tests ==========

func TestToolContext_SelectOrg(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")
	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Selected Organization")

	// Before selection
	assert.False(t, tc.HasOrg(), "HasOrg should be false before selection")
	assert.Equal(t, "", tc.OrgID(), "OrgID should be empty before selection")

	// Select org
	tc.SelectOrg(org)

	// After selection
	assert.True(t, tc.HasOrg(), "HasOrg should be true after selection")
	assert.Equal(t, orgID, tc.OrgID(), "OrgID should match selected org")
	assert.Equal(t, org.Name, tc.Org().Name, "Org name should match")
}

func TestToolContext_SelectOrg_CallsCallback(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")
	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Selected Organization")

	var callbackUserID string
	var callbackOrgID string
	callbackCalled := false

	tc.SetOrgSelectionCallback(func(uid, oid string) {
		callbackCalled = true
		callbackUserID = uid
		callbackOrgID = oid
	})

	tc.SelectOrg(org)

	assert.True(t, callbackCalled, "Callback should be called")
	assert.Equal(t, userID, callbackUserID, "Callback should receive user ID")
	assert.Equal(t, orgID, callbackOrgID, "Callback should receive org ID")
}

func TestToolContext_RestoreOrg_DoesNotCallCallback(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")
	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Restored Organization")

	callbackCalled := false
	tc.SetOrgSelectionCallback(func(uid, oid string) {
		callbackCalled = true
	})

	tc.RestoreOrg(org)

	assert.False(t, callbackCalled, "Callback should not be called for RestoreOrg")
	assert.True(t, tc.HasOrg(), "HasOrg should be true after restore")
	assert.Equal(t, orgID, tc.OrgID(), "OrgID should match restored org")
}

func TestToolContext_ClearOrg(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")
	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Organization to Clear")
	tc.SelectOrg(org)

	assert.True(t, tc.HasOrg(), "HasOrg should be true after selection")

	tc.ClearOrg()

	assert.False(t, tc.HasOrg(), "HasOrg should be false after clear")
	assert.Equal(t, "", tc.OrgID(), "OrgID should be empty after clear")
}

// ========== APIKey Mode Tests ==========

func TestToolContext_APIKeyMode_AlwaysHasOrg(t *testing.T) {
	orgID := uuid.New().String()
	org := createTestOrg(orgID, "API Key Organization")

	tc := NewToolContext(nil, org, "req-123", "test-agent/1.0")

	// In API key mode, org is always set
	assert.True(t, tc.HasOrg(), "HasOrg should always be true in API key mode")
	assert.Equal(t, orgID, tc.OrgID(), "OrgID should be set")
	assert.Equal(t, org.Name, tc.Org().Name, "Org should match")
}

func TestToolContext_APIKeyMode_NoUserID(t *testing.T) {
	orgID := uuid.New().String()
	org := createTestOrg(orgID, "API Key Organization")

	tc := NewToolContext(nil, org, "req-123", "test-agent/1.0")

	assert.Equal(t, "", tc.UserID(), "UserID should be empty in API key mode")
	assert.Nil(t, tc.User(), "User should be nil in API key mode")
	assert.False(t, tc.IsSuperAdmin(), "IsSuperAdmin should be false when no user")
}

// ========== User and Permissions Tests ==========

func TestToolContext_IsSuperAdmin(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		expectedResult bool
	}{
		{"Super admin role", "super_admin", true},
		{"Admin role", "admin", false},
		{"User role", "user", false},
		{"Empty role", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := uuid.New().String()
			user := createTestUser(userID, "test@example.com", tt.role)
			tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

			assert.Equal(t, tt.expectedResult, tc.IsSuperAdmin(), "IsSuperAdmin should return %v for role %s", tt.expectedResult, tt.role)
		})
	}
}

func TestToolContext_User(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")

	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	assert.NotNil(t, tc.User(), "User should not be nil")
	assert.Equal(t, userID, tc.User().ID, "User ID should match")
	assert.Equal(t, "test@example.com", tc.User().Email, "User email should match")
}

// ========== RequireOrg Tests ==========

func TestToolContext_RequireOrg_Success(t *testing.T) {
	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Test Organization")

	tc := NewToolContext(nil, org, "req-123", "test-agent/1.0")

	err := tc.RequireOrg()
	assert.NoError(t, err, "RequireOrg should not return error when org is set")
}

func TestToolContext_RequireOrg_Failure(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")

	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	err := tc.RequireOrg()
	require.Error(t, err, "RequireOrg should return error when no org selected")
	assert.Equal(t, ErrNoOrgSelected, err, "Error should be ErrNoOrgSelected")
}

// ========== ToolError Tests ==========

func TestToolError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ToolError
		expected string
	}{
		{
			name:     "Error without field",
			err:      &ToolError{Code: "not_found", Message: "Resource not found"},
			expected: "not_found: Resource not found",
		},
		{
			name:     "Error with field",
			err:      &ToolError{Code: "validation", Message: "Invalid value", Field: "email"},
			expected: "validation: Invalid value (field: email)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("Invalid email format", "email")

	assert.Equal(t, "validation", err.Code, "Code should be 'validation'")
	assert.Equal(t, "Invalid email format", err.Message, "Message should match")
	assert.Equal(t, "email", err.Field, "Field should be 'email'")
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("User not found")

	assert.Equal(t, "not_found", err.Code, "Code should be 'not_found'")
	assert.Equal(t, "User not found", err.Message, "Message should match")
	assert.Empty(t, err.Field, "Field should be empty")
}

func TestNewConflictError(t *testing.T) {
	err := NewConflictError("Email already exists")

	assert.Equal(t, "conflict", err.Code, "Code should be 'conflict'")
	assert.Equal(t, "Email already exists", err.Message, "Message should match")
	assert.Empty(t, err.Field, "Field should be empty")
}

// ========== Context Value Tests ==========

func TestWithToolContext_And_FromContext(t *testing.T) {
	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Test Organization")
	tc := NewToolContext(nil, org, "req-123", "test-agent/1.0")

	ctx := context.Background()
	ctx = WithToolContext(ctx, tc)

	retrieved := ToolContextFromContext(ctx)
	require.NotNil(t, retrieved, "ToolContext should be retrievable from context")
	assert.Equal(t, tc, retrieved, "Retrieved context should match original")
	assert.Equal(t, orgID, retrieved.OrgID(), "OrgID should match")
}

func TestToolContextFromContext_NilWhenNotSet(t *testing.T) {
	ctx := context.Background()

	retrieved := ToolContextFromContext(ctx)
	assert.Nil(t, retrieved, "Should return nil when ToolContext not set")
}

func TestToolContextFromContext_WrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), toolContextKey{}, "wrong type")

	retrieved := ToolContextFromContext(ctx)
	assert.Nil(t, retrieved, "Should return nil for wrong type")
}

// ========== Concurrent Access Tests ==========

func TestToolContext_ConcurrentOrgSelection(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")
	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	var wg sync.WaitGroup
	numGoroutines := 100

	// Create different orgs to select
	orgs := make([]db.Organization, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		orgs[i] = createTestOrg(uuid.New().String(), "Org "+string(rune('A'+i%26)))
	}

	// Concurrent org selection
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			tc.SelectOrg(orgs[idx])
			// Read immediately after write
			_ = tc.HasOrg()
			_ = tc.OrgID()
		}(i)
	}

	wg.Wait()

	// After all concurrent selections, should have an org
	assert.True(t, tc.HasOrg(), "Should have an org after concurrent selections")
}

func TestToolContext_ConcurrentReads(t *testing.T) {
	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Test Organization")
	tc := NewToolContext(nil, org, "req-123", "test-agent/1.0")

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = tc.HasOrg()
			_ = tc.OrgID()
			_ = tc.Org()
			_ = tc.AuthMode()
			_ = tc.RequestID()
			_ = tc.UserAgent()
		}()
	}

	wg.Wait()
	// Test passes if no race condition panic occurs
}

// ========== Empty/Nil Org Tests ==========

func TestToolContext_Org_ReturnsEmptyWhenNil(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")
	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-123")

	org := tc.Org()
	assert.Equal(t, db.Organization{}, org, "Org should return empty struct when not set")
}

// ========== Session ID Tests ==========

func TestToolContext_SessionID(t *testing.T) {
	userID := uuid.New().String()
	user := createTestUser(userID, "test@example.com", "admin")

	tc := NewUserToolContext(nil, user, "req-123", "test-agent/1.0", "session-abc-123")

	assert.Equal(t, "session-abc-123", tc.SessionID(), "SessionID should match")
}

func TestToolContext_SessionID_EmptyForAPIKeyMode(t *testing.T) {
	orgID := uuid.New().String()
	org := createTestOrg(orgID, "Test Organization")

	tc := NewToolContext(nil, org, "req-123", "test-agent/1.0")

	assert.Empty(t, tc.SessionID(), "SessionID should be empty for API key mode")
}
