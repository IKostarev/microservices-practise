package rest

import (
	"gateway/pkg/ctxutil"
	"gateway/pkg/jwtutil"
	"gateway/pkg/redis"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func ValidateTokenMiddleware(
	jwtUtil *jwtutil.JWTUtil,
	excludedPaths []string,
	manager *redis.RedisManager,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a UUID for each request
			ctx := ctxutil.SetRequestIDToContext(r.Context(), uuid.New().String())

			// Check if the route is excluded from validation
			for _, path := range excludedPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Validate the token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, exp, err := jwtUtil.VerifyToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			if !manager.IsTokenValid(ctx, userID, exp) {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add the user ID to the context
			ctx = ctxutil.SetUserIDToContext(ctx, userID)

			// Token is valid, proceed with the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
