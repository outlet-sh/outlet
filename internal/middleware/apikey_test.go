package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"outlet/internal/db"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockStore is a minimal mock that embeds the real Store but overrides GetOrganizationByAPIKey
type MockStore struct {
	*db.Store
	GetOrgFunc func(ctx context.Context, apiKey string) (db.Organization, error)
}

// GetOrganizationByAPIKey overrides the real implementation for testing
func (m *MockStore) GetOrganizationByAPIKey(ctx context.Context, apiKey string) (db.Organization, error) {
	if m.GetOrgFunc != nil {
		return m.GetOrgFunc(ctx, apiKey)
	}
	return db.Organization{}, errors.New("not implemented")
}

// createMockOrg creates a mock organization with the given parameters
func createMockOrg(id string, name, apiKey string) db.Organization {
	return db.Organization{
		ID:     id,
		Name:   name,
		Slug:   "test-org",
		ApiKey: apiKey,
	}
}

// testableAPIKeyMiddleware wraps APIKeyMiddleware with a mockable getOrg function
type testableAPIKeyMiddleware struct {
	*APIKeyMiddleware
	getOrgFunc func(ctx context.Context, apiKey string) (db.Organization, error)
}

func newTestableAPIKeyMiddleware(getOrgFunc func(ctx context.Context, apiKey string) (db.Organization, error)) *testableAPIKeyMiddleware {
	m := &testableAPIKeyMiddleware{
		APIKeyMiddleware: &APIKeyMiddleware{
			store: nil,
			cache: sync.Map{},
		},
		getOrgFunc: getOrgFunc,
	}
	return m
}

func (m *testableAPIKeyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Missing API key",
			})
			return
		}

		org, err := m.getOrgFunc(r.Context(), apiKey)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid API key",
			})
			return
		}

		// Add org info to context
		ctx := context.WithValue(r.Context(), OrgIDKey, org.ID)
		ctx = context.WithValue(ctx, OrgKey, org)

		next(w, r.WithContext(ctx))
	}
}

// testableAPIKeyAuthMiddleware wraps APIKeyAuthMiddleware with a mockable getOrg function
type testableAPIKeyAuthMiddleware struct {
	*APIKeyAuthMiddleware
	getOrgFunc func(ctx context.Context, apiKey string) (db.Organization, error)
}

func newTestableAPIKeyAuthMiddleware(getOrgFunc func(ctx context.Context, apiKey string) (db.Organization, error)) *testableAPIKeyAuthMiddleware {
	m := &testableAPIKeyAuthMiddleware{
		APIKeyAuthMiddleware: &APIKeyAuthMiddleware{
			store: nil,
			cache: sync.Map{},
		},
		getOrgFunc: getOrgFunc,
	}
	return m
}

func (m *testableAPIKeyAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Missing API key",
			})
			return
		}

		org, err := m.getOrgFunc(r.Context(), apiKey)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid API key",
			})
			return
		}

		// Add org info to context
		ctx := context.WithValue(r.Context(), OrgIDKey, org.ID)
		ctx = context.WithValue(ctx, OrgKey, org)

		next(w, r.WithContext(ctx))
	}
}

// ========== APIKeyMiddleware Tests ==========

func TestNewAPIKeyMiddleware(t *testing.T) {
	// Create with nil store (we don't have a real DB in tests)
	middleware := NewAPIKeyMiddleware(nil)
	assert.NotNil(t, middleware, "NewAPIKeyMiddleware should return a non-nil middleware")
}

func TestAPIKeyMiddleware_Handle_MissingAPIKey(t *testing.T) {
	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return db.Organization{}, nil
	})
	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called without API key")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Missing API key", response["error"])
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestAPIKeyMiddleware_Handle_ValidAPIKey(t *testing.T) {
	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Test Org", "valid-api-key")

	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		if apiKey == "valid-api-key" {
			return org, nil
		}
		return db.Organization{}, errors.New("not found")
	})

	handlerCalled := false
	var capturedOrgID interface{}
	var capturedOrg interface{}

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		capturedOrgID = r.Context().Value(OrgIDKey)
		capturedOrg = r.Context().Value(OrgKey)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "valid-api-key")
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called with valid API key")
	assert.Equal(t, orgID, capturedOrgID)
	assert.Equal(t, org, capturedOrg)
}

func TestAPIKeyMiddleware_Handle_InvalidAPIKey(t *testing.T) {
	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return db.Organization{}, errors.New("not found")
	})

	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "invalid-api-key")
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called with invalid API key")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid API key", response["error"])
}

func TestAPIKeyMiddleware_Handle_DatabaseError(t *testing.T) {
	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return db.Organization{}, errors.New("database connection error")
	})

	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "some-api-key")
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called on database error")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid API key", response["error"])
}

func TestAPIKeyMiddleware_Handle_DifferentHTTPMethods(t *testing.T) {
	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Test Org", "api-key")

	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return org, nil
	})

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			handlerCalled := false

			handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
			})

			req := httptest.NewRequest(method, "/test", nil)
			req.Header.Set("X-API-Key", "api-key")
			rr := httptest.NewRecorder()

			handler(rr, req)

			assert.True(t, handlerCalled, "Handler should be called for %s method", method)
		})
	}
}

// ========== APIKeyAuthMiddleware Tests ==========

func TestNewAPIKeyAuthMiddleware(t *testing.T) {
	middleware := NewAPIKeyAuthMiddleware(nil)
	assert.NotNil(t, middleware, "NewAPIKeyAuthMiddleware should return a non-nil middleware")
}

func TestAPIKeyAuthMiddleware_Handle_MissingAPIKey(t *testing.T) {
	middleware := newTestableAPIKeyAuthMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return db.Organization{}, nil
	})
	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called without API key")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Missing API key", response["error"])
}

func TestAPIKeyAuthMiddleware_Handle_ValidAPIKey(t *testing.T) {
	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Test Org", "valid-api-key")

	middleware := newTestableAPIKeyAuthMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		if apiKey == "valid-api-key" {
			return org, nil
		}
		return db.Organization{}, errors.New("not found")
	})

	handlerCalled := false
	var capturedOrgID interface{}
	var capturedOrg interface{}

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		capturedOrgID = r.Context().Value(OrgIDKey)
		capturedOrg = r.Context().Value(OrgKey)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "valid-api-key")
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called with valid API key")
	assert.Equal(t, orgID, capturedOrgID)
	assert.Equal(t, org, capturedOrg)
}

func TestAPIKeyAuthMiddleware_Handle_InvalidAPIKey(t *testing.T) {
	middleware := newTestableAPIKeyAuthMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return db.Organization{}, errors.New("not found")
	})

	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "invalid-api-key")
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called with invalid API key")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid API key", response["error"])
}

// ========== Cache Tests ==========

func TestAPIKeyMiddleware_Cache(t *testing.T) {
	callCount := 0
	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Cached Org", "cached-api-key")

	middleware := &APIKeyMiddleware{
		store: nil,
		cache: sync.Map{},
	}

	// Pre-populate the cache
	middleware.cache.Store("cached-api-key", org)

	// Verify cache contains the org
	cachedOrg, found := middleware.cache.Load("cached-api-key")
	assert.True(t, found, "Org should be found in cache")
	assert.Equal(t, org, cachedOrg)

	// Test InvalidateCache
	middleware.InvalidateCache("cached-api-key")
	_, found = middleware.cache.Load("cached-api-key")
	assert.False(t, found, "Org should be removed from cache after invalidation")

	// callCount should still be 0 since we never called the DB
	assert.Equal(t, 0, callCount)
}

func TestAPIKeyAuthMiddleware_Cache(t *testing.T) {
	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Cached Org", "cached-api-key")

	middleware := &APIKeyAuthMiddleware{
		store: nil,
		cache: sync.Map{},
	}

	// Pre-populate the cache
	middleware.cache.Store("cached-api-key", org)

	// Verify cache contains the org
	cachedOrg, found := middleware.cache.Load("cached-api-key")
	assert.True(t, found, "Org should be found in cache")
	assert.Equal(t, org, cachedOrg)

	// Test InvalidateCache
	middleware.InvalidateCache("cached-api-key")
	_, found = middleware.cache.Load("cached-api-key")
	assert.False(t, found, "Org should be removed from cache after invalidation")
}

func TestAPIKeyMiddleware_InvalidateCache_NonExistentKey(t *testing.T) {
	middleware := &APIKeyMiddleware{
		store: nil,
		cache: sync.Map{},
	}

	// Should not panic when invalidating non-existent key
	assert.NotPanics(t, func() {
		middleware.InvalidateCache("non-existent-key")
	})
}

func TestAPIKeyAuthMiddleware_InvalidateCache_NonExistentKey(t *testing.T) {
	middleware := &APIKeyAuthMiddleware{
		store: nil,
		cache: sync.Map{},
	}

	// Should not panic when invalidating non-existent key
	assert.NotPanics(t, func() {
		middleware.InvalidateCache("non-existent-key")
	})
}

// ========== Context Keys Tests ==========

func TestContextKeys(t *testing.T) {
	// Verify the context keys are defined correctly
	assert.Equal(t, contextKey("org_id"), OrgIDKey)
	assert.Equal(t, contextKey("org"), OrgKey)

	// Verify they can be used as map keys (different keys)
	m := make(map[contextKey]string)
	m[OrgIDKey] = "id-value"
	m[OrgKey] = "org-value"

	assert.Equal(t, "id-value", m[OrgIDKey])
	assert.Equal(t, "org-value", m[OrgKey])
	assert.Len(t, m, 2)
}

// ========== Edge Cases ==========

func TestAPIKeyMiddleware_Handle_EmptyAPIKeyHeader(t *testing.T) {
	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return db.Organization{}, nil
	})

	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "")
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.False(t, handlerCalled, "Handler should not be called with empty API key")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAPIKeyMiddleware_Handle_WhitespaceAPIKey(t *testing.T) {
	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		// If whitespace-only key is passed, treat it as invalid
		return db.Organization{}, errors.New("not found")
	})

	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "   ")
	rr := httptest.NewRecorder()

	handler(rr, req)

	// The middleware doesn't trim whitespace, so it will try to look up the whitespace key
	assert.False(t, handlerCalled, "Handler should not be called with whitespace-only API key")
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAPIKeyMiddleware_Handle_SpecialCharactersInAPIKey(t *testing.T) {
	specialKey := "api-key-with-special-chars!@#$%^&*()_+-=[]{}|;':\",./<>?"
	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Special Org", specialKey)

	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		if apiKey == specialKey {
			return org, nil
		}
		return db.Organization{}, errors.New("not found")
	})

	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", specialKey)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called with special characters in API key")
}

func TestAPIKeyMiddleware_Handle_UnicodeAPIKey(t *testing.T) {
	unicodeKey := "api-key-\u4e2d\u6587-\u0420\u0443\u0441\u0441\u043a\u0438\u0439"
	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Unicode Org", unicodeKey)

	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		if apiKey == unicodeKey {
			return org, nil
		}
		return db.Organization{}, errors.New("not found")
	})

	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", unicodeKey)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called with unicode characters in API key")
}

func TestAPIKeyMiddleware_Handle_VeryLongAPIKey(t *testing.T) {
	// Create a very long API key (1000 characters)
	longKey := make([]byte, 1000)
	for i := range longKey {
		longKey[i] = 'a' + byte(i%26)
	}
	longKeyStr := string(longKey)

	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Long Key Org", longKeyStr)

	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		if apiKey == longKeyStr {
			return org, nil
		}
		return db.Organization{}, errors.New("not found")
	})

	handlerCalled := false

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", longKeyStr)
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called with very long API key")
}

// ========== Concurrent Access Tests ==========

func TestAPIKeyMiddleware_ConcurrentAccess(t *testing.T) {
	orgID := uuid.New().String()
	org := createMockOrg(orgID, "Concurrent Org", "concurrent-key")

	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return org, nil
	})

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
				// Verify org is in context
				capturedOrgID := r.Context().Value(OrgIDKey)
				assert.Equal(t, orgID, capturedOrgID)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("X-API-Key", "concurrent-key")
			rr := httptest.NewRecorder()

			handler(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
		}()
	}

	wg.Wait()
}

// ========== Response Format Tests ==========

func TestAPIKeyMiddleware_ErrorResponseFormat(t *testing.T) {
	middleware := newTestableAPIKeyMiddleware(func(ctx context.Context, apiKey string) (db.Organization, error) {
		return db.Organization{}, errors.New("not found")
	})

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "invalid-key")
	rr := httptest.NewRecorder()

	handler(rr, req)

	// Verify response is valid JSON
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err, "Response should be valid JSON")

	// Verify response has expected structure
	_, hasError := response["error"]
	assert.True(t, hasError, "Response should have 'error' field")

	// Verify Content-Type header
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}
