package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"outlet/internal/db"

	"github.com/zeromicro/go-zero/core/logx"
)

type contextKey string

const (
	OrgIDKey  contextKey = "org_id"
	OrgKey    contextKey = "org"
)

type APIKeyMiddleware struct {
	store *db.Store
	cache sync.Map // map[apiKey]db.Organization
}

func NewAPIKeyMiddleware(store *db.Store) *APIKeyMiddleware {
	return &APIKeyMiddleware{store: store}
}

// getOrg retrieves an organization by API key, using cache when available.
// Cache persists until InvalidateCache is called or service restarts.
func (m *APIKeyMiddleware) getOrg(ctx context.Context, apiKey string) (db.Organization, error) {
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
func (m *APIKeyMiddleware) InvalidateCache(apiKey string) {
	m.cache.Delete(apiKey)
}

func (m *APIKeyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
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
