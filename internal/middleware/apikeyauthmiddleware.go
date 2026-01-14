package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"outlet/internal/db"

	"github.com/zeromicro/go-zero/core/logx"
)

// APIKeyAuthMiddleware validates API keys from the X-API-Key header
// and sets the organization context for downstream handlers.
type APIKeyAuthMiddleware struct {
	store *db.Store
	cache sync.Map // map[apiKey]db.Organization
}

// NewAPIKeyAuthMiddleware creates a new API key authentication middleware.
// The store parameter is required for database lookups.
func NewAPIKeyAuthMiddleware(store *db.Store) *APIKeyAuthMiddleware {
	return &APIKeyAuthMiddleware{store: store}
}

// getOrg retrieves an organization by API key, using cache when available.
// Cache persists until InvalidateCache is called or service restarts.
func (m *APIKeyAuthMiddleware) getOrg(ctx context.Context, apiKey string) (db.Organization, error) {
	// Check cache first
	if cached, ok := m.cache.Load(apiKey); ok {
		return cached.(db.Organization), nil
	}

	// Lookup from database
	org, err := m.store.GetOrganizationByAPIKey(ctx, apiKey)
	if err != nil {
		return db.Organization{}, err
	}

	// Cache the result
	m.cache.Store(apiKey, org)

	return org, nil
}

// InvalidateCache removes an API key from the cache (call when org is updated)
func (m *APIKeyAuthMiddleware) InvalidateCache(apiKey string) {
	m.cache.Delete(apiKey)
}

// Handle validates the X-API-Key header and sets org context for downstream handlers.
// Returns 401 Unauthorized if the API key is missing or invalid.
// Returns 403 Forbidden if the organization account is suspended.
func (m *APIKeyAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
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

		org, err := m.getOrg(r.Context(), apiKey)
		if err != nil {
			logx.Errorf("API key lookup failed: %v", err)
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
