package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/moha/kaafipay-backend/internal/utils"
)

// AuthMiddleware handles authentication for protected routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header", nil)
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization format", nil)
			return
		}

		token := parts[1]
		if token == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token", nil)
			return
		}

		// TODO: Validate token and get user ID
		// This is a placeholder. You should implement proper token validation
		// and user ID extraction from your authentication service
		userID, err := uuid.Parse("00000000-0000-0000-0000-000000000000") // Replace with actual user ID from token
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token", nil)
			return
		}

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
