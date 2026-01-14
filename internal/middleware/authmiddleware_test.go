package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestToken creates a JWT token for testing purposes
func createTestToken(claims jwt.MapClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// createUnsignedToken creates a JWT token string without signing (for testing ParseUnverified)
func createUnsignedToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign with a dummy secret since the middleware uses ParseUnverified
	tokenString, _ := token.SignedString([]byte("any-secret"))
	return tokenString
}

func TestNewAuthMiddleware(t *testing.T) {
	middleware := NewAuthMiddleware()
	assert.NotNil(t, middleware, "NewAuthMiddleware should return a non-nil middleware")
}

func TestAuthMiddleware_Handle_NoAuthorizationHeader(t *testing.T) {
	middleware := NewAuthMiddleware()
	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called when no auth token is provided")
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Should return 401 Unauthorized")

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "No authentication token provided", response["error"])
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestAuthMiddleware_Handle_ValidTokenFromCookie(t *testing.T) {
	middleware := NewAuthMiddleware()
	handlerCalled := false
	var capturedUserID interface{}
	var capturedEmail interface{}
	var capturedRole interface{}
	var capturedToken interface{}

	userID := uuid.New().String()
	email := "test@example.com"
	role := "admin"

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		capturedUserID = r.Context().Value("userId")
		capturedEmail = r.Context().Value("email")
		capturedRole = r.Context().Value("role")
		capturedToken = r.Context().Value("token")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called with valid cookie token")
	assert.Equal(t, userID, capturedUserID)
	assert.Equal(t, email, capturedEmail)
	assert.Equal(t, role, capturedRole)
	assert.Equal(t, tokenString, capturedToken)
}

func TestAuthMiddleware_Handle_ValidTokenFromHeader(t *testing.T) {
	middleware := NewAuthMiddleware()
	handlerCalled := false
	var capturedUserID interface{}
	var capturedEmail interface{}
	var capturedRole interface{}

	userID := uuid.New().String()
	email := "user@example.com"
	role := "user"

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		capturedUserID = r.Context().Value("userId")
		capturedEmail = r.Context().Value("email")
		capturedRole = r.Context().Value("role")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called with valid Bearer token")
	assert.Equal(t, userID, capturedUserID)
	assert.Equal(t, email, capturedEmail)
	assert.Equal(t, role, capturedRole)
}

func TestAuthMiddleware_Handle_InvalidAuthorizationHeaderFormat(t *testing.T) {
	middleware := NewAuthMiddleware()
	handlerCalled := false

	tests := []struct {
		name          string
		authHeader    string
		expectedError string
	}{
		{"Only token no prefix", "some-token-value", "Invalid authorization header format"},
		{"Wrong prefix", "Basic some-token-value", "Invalid authorization header format"},
		{"Bearer with no token", "Bearer", "Invalid authorization header format"},
		{"Empty bearer", "Bearer ", "Invalid token format"}, // Empty token after Bearer passes split but fails parse
		{"Too many parts", "Bearer token extra", "Invalid authorization header format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rr := httptest.NewRecorder()

			handler(rr, req)

			assert.False(t, handlerCalled, "Handler should not be called with invalid header format")
			assert.Equal(t, http.StatusUnauthorized, rr.Code)

			var response map[string]string
			err := json.NewDecoder(rr.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedError, response["error"])
		})
	}
}

func TestAuthMiddleware_Handle_MalformedToken(t *testing.T) {
	middleware := NewAuthMiddleware()
	handlerCalled := false

	tests := []struct {
		name  string
		token string
	}{
		{"Completely invalid", "not-a-jwt-token"},
		{"Invalid base64", "eyJ.eyJ.sig"},
		{"Missing parts", "eyJhbGciOiJIUzI1NiJ9"},
		{"Empty token", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			rr := httptest.NewRecorder()

			handler(rr, req)

			assert.False(t, handlerCalled, "Handler should not be called with malformed token")
			assert.Equal(t, http.StatusUnauthorized, rr.Code)

			var response map[string]string
			err := json.NewDecoder(rr.Body).Decode(&response)
			require.NoError(t, err)
			assert.Contains(t, response["error"], "Invalid")
		})
	}
}

func TestAuthMiddleware_Handle_MissingUserIDClaim(t *testing.T) {
	middleware := NewAuthMiddleware()
	handlerCalled := false

	// Token without user_id claim
	claims := jwt.MapClaims{
		"email": "test@example.com",
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called without user_id claim")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid token claims", response["error"])
}

func TestAuthMiddleware_Handle_PartialClaims(t *testing.T) {
	middleware := NewAuthMiddleware()
	handlerCalled := false
	var capturedUserID interface{}
	var capturedEmail interface{}
	var capturedRole interface{}

	// Token with only user_id (no email or role)
	userID := uuid.New().String()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		capturedUserID = r.Context().Value("userId")
		capturedEmail = r.Context().Value("email")
		capturedRole = r.Context().Value("role")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler(rr, req)

	// Should succeed because user_id is present
	assert.True(t, handlerCalled, "Handler should be called when user_id is present")
	assert.Equal(t, userID, capturedUserID)
	assert.Nil(t, capturedEmail, "Email should be nil when not in claims")
	assert.Nil(t, capturedRole, "Role should be nil when not in claims")
}

func TestAuthMiddleware_Handle_HeaderTakesPrecedenceOverCookie(t *testing.T) {
	middleware := NewAuthMiddleware()
	var capturedUserID interface{}

	headerUserID := uuid.New().String()
	cookieUserID := uuid.New().String()

	headerClaims := jwt.MapClaims{
		"user_id": headerUserID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	headerToken := createUnsignedToken(headerClaims)

	cookieClaims := jwt.MapClaims{
		"user_id": cookieUserID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	cookieToken := createUnsignedToken(cookieClaims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = r.Context().Value("userId")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+headerToken)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: cookieToken,
	})
	rr := httptest.NewRecorder()

	handler(rr, req)

	// Header should take precedence
	assert.Equal(t, headerUserID, capturedUserID, "Authorization header should take precedence over cookie")
}

func TestAuthMiddleware_Handle_WithOrgIDHeader(t *testing.T) {
	middleware := NewAuthMiddleware()
	var capturedOrgID interface{}

	userID := uuid.New().String()
	orgID := uuid.New()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		capturedOrgID = r.Context().Value(OrgIDKey)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("X-Org-Id", orgID.String())
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, orgID, capturedOrgID, "OrgIDKey should be set in context from X-Org-Id header")
}

func TestAuthMiddleware_Handle_InvalidOrgIDHeader(t *testing.T) {
	middleware := NewAuthMiddleware()
	var capturedOrgID interface{}

	userID := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		capturedOrgID = r.Context().Value(OrgIDKey)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("X-Org-Id", "not-a-valid-uuid")
	rr := httptest.NewRecorder()

	handler(rr, req)

	// Handler should still be called, just without OrgIDKey set
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Nil(t, capturedOrgID, "OrgIDKey should not be set with invalid UUID")
}

func TestAuthMiddleware_Handle_EmptyOrgIDHeader(t *testing.T) {
	middleware := NewAuthMiddleware()
	var capturedOrgID interface{}

	userID := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		capturedOrgID = r.Context().Value(OrgIDKey)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	// Empty X-Org-Id header
	req.Header.Set("X-Org-Id", "")
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Nil(t, capturedOrgID, "OrgIDKey should not be set with empty header")
}

func TestAuthMiddleware_Handle_EmptyCookie(t *testing.T) {
	middleware := NewAuthMiddleware()
	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: "",
	})
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called with empty cookie value")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "No authentication token provided", response["error"])
}

func TestAuthMiddleware_Handle_ContentTypeHeader(t *testing.T) {
	middleware := NewAuthMiddleware()

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		// Should not be called
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"),
		"Error responses should have Content-Type: application/json")
}

func TestAuthMiddleware_Handle_DifferentHTTPMethods(t *testing.T) {
	middleware := NewAuthMiddleware()

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	}

	userID := uuid.New().String()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			handlerCalled := false

			handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
			})

			req := httptest.NewRequest(method, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+tokenString)
			rr := httptest.NewRecorder()

			handler(rr, req)

			assert.True(t, handlerCalled, "Handler should be called for %s method", method)
		})
	}
}

func TestAuthMiddleware_Handle_NumericClaims(t *testing.T) {
	middleware := NewAuthMiddleware()
	var capturedUserID interface{}

	// Test with numeric user_id (some systems use integers)
	claims := jwt.MapClaims{
		"user_id": float64(12345), // JSON numbers are float64
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString := createUnsignedToken(claims)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = r.Context().Value("userId")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, float64(12345), capturedUserID, "Should handle numeric user_id claims")
}
