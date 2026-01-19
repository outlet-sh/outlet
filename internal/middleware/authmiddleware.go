package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct {
	// Config will be injected by go-zero when middleware is used
	// For now, we'll use a simple approach
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Try to get from cookie
			cookie, err := r.Cookie("auth_token")
			if err == nil && cookie.Value != "" {
				authHeader = "Bearer " + cookie.Value
			}
		}

		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "No authentication token provided",
			})
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]

		// Parse token to extract claims without full validation
		// The actual validation will happen in the logic layer with the proper secret
		token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			logx.Errorf("Token parsing error: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid token format",
			})
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Add user info to context for use in handlers
			if userID, exists := claims["user_id"]; exists {
				ctx := context.WithValue(r.Context(), "userId", userID)
				if email, exists := claims["email"]; exists {
					ctx = context.WithValue(ctx, "email", email)
				}
				if role, exists := claims["role"]; exists {
					ctx = context.WithValue(ctx, "role", role)
				}
				// Store the raw token for validation in logic layer
				ctx = context.WithValue(ctx, "token", tokenString)

				// Read X-Org-Id header for org-scoped admin operations
				orgIDStr := r.Header.Get("X-Org-Id")
				if orgIDStr != "" {
					// Validate it's a proper UUID, but store as string (logic files expect string)
					if _, err := uuid.Parse(orgIDStr); err == nil {
						ctx = context.WithValue(ctx, OrgIDKey, orgIDStr)
					}
				}

				next(w, r.WithContext(ctx))
				return
			}
		}

		// If we get here, token was invalid
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid token claims",
		})
	}
}
